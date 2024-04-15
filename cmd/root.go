package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "watchtower",
	Short: "A programmable TCP/UDP server to debug incoming network requests.",
	Long: `A programmable TCP/UDP server to debug incoming network requests.
Watchtower can be configured to watch for specific strings, regex expressions, and much more.
	
Example usage (Host a TCP server on port 2000 looking for any request with 'foo' in it):
watchtower -p 2000 -t tcp -s foo`,
	Run:     runServer,
	Version: "0.1.1",
}

func replaceStringTerms(packet string, terms []string) (string, bool) {
	foundMatch := false
	for _, term := range terms {
		if strings.Contains(packet, term) {
			// Colour the output to highlight the found value
			packet = strings.ReplaceAll(packet, term, color.CyanString(term))
			foundMatch = true
		}
	}

	return packet, foundMatch
}

func replaceTerms(packet string, terms []string) (string, bool) {
	if os.Getenv("WATCHTOWER_USE_REGEX") != "1" {
		return replaceStringTerms(packet, terms)
	}

	foundMatch := false
	for _, term := range terms {
		regex := regexp.MustCompile(term)
		if regex.MatchString(packet) {
			packet = regex.ReplaceAllStringFunc(packet, func(s string) string {
				return color.CyanString(s)
			})
			foundMatch = true
		}
	}

	return packet, foundMatch

}

func handleConnection(conn net.Conn, terms []string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// Peek so we can read the size of the buffer
	_, err := reader.Peek(1)
	if err != nil {
		panic(err)
	}

	output, err := reader.Peek(reader.Buffered())
	if err != nil {
		conn.Close()
		return
	}

	packet := string(output)
	fmt.Printf("\nConnection recieved from %s ", conn.RemoteAddr())
	packet, foundMatch := replaceTerms(packet, terms)

	if len(terms) != 0 {
		if foundMatch {
			fmt.Printf("\n%s", packet)
		} else {
			fmt.Println("NO MATCH")
		}
	}

	conn.Write([]byte(""))
}

func runServer(cmd *cobra.Command, args []string) {
	addr, err := cmd.Flags().GetIP("address")
	if err != nil {
		fmt.Println("Invalid IP address")
		return
	}

	port, _ := cmd.Flags().GetInt("port")

	network, _ := cmd.Flags().GetString("network")
	if network != "tcp" && network != "udp" {
		fmt.Println("Invalid network. Only 'tcp' and 'udp' are supported")
		return
	}

	terms, _ := cmd.Flags().GetStringSlice("search-terms")

	listener, err := net.Listen(network, addr.String()+":"+fmt.Sprintf("%d", port))
	if err != nil {
		fmt.Println("Failed to start watchtower: " + err.Error())
		return
	}

	fmt.Printf("Listening on %s:%d\n", addr, port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go handleConnection(conn, terms)
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
	rootCmd.Flags().StringP("network", "n", "tcp", "Network protocol to use. Either tcp or udp")
	rootCmd.Flags().StringSliceP("search-terms", "s", []string{}, `Term(s) to search incoming requests for. Set the environment variable
WATCHTOWER_USE_REGEX = 1 to use regex searches
	`)
}
