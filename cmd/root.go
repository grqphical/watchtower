/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"

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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	fmt.Printf("\nConnection recieved from %s\n:", conn.RemoteAddr())
	output := make([]byte, 1024)
	_, err := reader.Read(output)
	if err != nil {
		conn.Close()
		return
	}
	fmt.Printf("%s", string(output))
	conn.Write([]byte(""))
}

func runServer(cmd *cobra.Command, args []string) {
	addr, err := cmd.Flags().GetIP("address")
	if err != nil {
		fmt.Println("Invalid IP address")
		return
	}

	port, _ := cmd.Flags().GetInt("port")

	listener, err := net.Listen("tcp", addr.String()+":"+fmt.Sprintf("%d", port))
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
		go handleConnection(conn)
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
}
