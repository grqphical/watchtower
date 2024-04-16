package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "watchtower",
	Short: "A programmable TCP server to debug incoming network requests.",
	Long: `A programmable TCP server to debug incoming network requests.
Watchtower can be configured to watch for specific strings, regex expressions, and much more.
	
Example usage (Host a TCP server on port 2000 looking for any request with 'foo' in it):
watchtower -p 2000 -s foo`,
	Run:     runServer,
	Version: "0.1.1",
}

func handleConnection(conn net.Conn, terms []string, file string, responseFile string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// Peek so we can read the size of the buffer
	_, err := reader.Peek(1)
	if err != nil {
		slog.Error("Could not read incoming data")
		return
	}

	output, err := reader.Peek(reader.Buffered())
	if err != nil {
		slog.Error("Could not read incoming data")
		conn.Close()
		return
	}

	packet := string(output)
	slog.Info(fmt.Sprintf("Connection recieved from %s ", conn.RemoteAddr()))

	colouredPacket, foundMatch := replaceTerms(packet, terms)

	if file != "" && foundMatch {
		file, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("Failed to open output file")
		}

		file.Write([]byte(packet))

		file.Close()
	}

	if len(terms) != 0 {
		if foundMatch {
			fmt.Println(colouredPacket)
		} else {
			slog.Warn("NO MATCH")
		}
	}

	if responseFile == "" {
		conn.Write([]byte(""))
		return
	}

	if _, err := os.Stat(responseFile); err == nil {
		data, err := os.ReadFile(responseFile)
		if err != nil {
			slog.Error("Could not read output file")
			conn.Write([]byte(""))
			return
		}

		conn.Write(data)
	}

}

func runServer(cmd *cobra.Command, args []string) {
	addr, err := cmd.Flags().GetIP("address")
	if err != nil {
		slog.Error("Invalid IP Address")
		return
	}

	port, _ := cmd.Flags().GetInt("port")

	terms, _ := cmd.Flags().GetStringSlice("search-terms")
	file, _ := cmd.Flags().GetString("output-file")
	responseFile, _ := cmd.Flags().GetString("response-file")

	listener, err := net.Listen("tcp", addr.String()+":"+fmt.Sprintf("%d", port))
	if err != nil {
		slog.Error("Failed to start watchtower: " + err.Error())
		return
	}

	slog.Info(fmt.Sprintf("Listening on %s:%d", addr, port))
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error accepting connection: " + err.Error())
			break
		}
		go handleConnection(conn, terms, file, responseFile)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IPP("address", "a", net.ParseIP("127.0.0.1"), "IP To host server on")
	rootCmd.Flags().IntP("port", "p", 8000, "Port to host server on")
	rootCmd.Flags().StringSliceP("search-terms", "s", []string{}, `Term(s) to search incoming requests for. Set the environment variable
WATCHTOWER_USE_REGEX = 1 to use regex searches
	`)
	rootCmd.Flags().StringP("output-file", "f", "", "File to output matched values to")
	rootCmd.Flags().StringP("response-file", "r", "", "Data from a file to respond to incoming connection with")
}
