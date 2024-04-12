package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	Run: runServer,
}

func replaceTerms(packet string, terms []string) string {
	foundMatch := false
	for _, term := range terms {
		if strings.Contains(packet, term) {
			// Colour the output to highlight the found value
			packet = strings.ReplaceAll(packet, term, color.New(color.FgCyan).Sprint(term))
			foundMatch = true
		}
	}

	if !foundMatch {
		return "NO MATCH\n"
	}

	return packet
}

func handleConnection(conn net.Conn, bufferSize int, terms []string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	output := make([]byte, bufferSize)

	_, err := reader.Read(output)
	if err != nil {
		conn.Close()
		return
	}

	packet := string(output)
	fmt.Printf("\nConnection recieved from %s ", conn.RemoteAddr())

	if len(terms) == 0 {
		fmt.Printf("\n%s", packet)
	} else {
		fmt.Print(replaceTerms(packet, terms))
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

	buffer_size, _ := cmd.Flags().GetInt("buffer-size")
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
		go handleConnection(conn, buffer_size, terms)
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
	rootCmd.Flags().Int("buffer-size", 1024, "Size of input buffer from network connection")
	rootCmd.Flags().StringP("network", "n", "tcp", "Network protocol to use. Either tcp or udp")
	rootCmd.Flags().StringSliceP("search-terms", "s", []string{}, "Term(s) to search incoming requests for")
}