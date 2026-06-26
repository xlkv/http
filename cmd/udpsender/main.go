package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr := "localhost:42069"
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing server:", err)
		os.Exit(1)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddr)

	for true {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", input)
	}
	// buffer := make([]byte, 1024)
	// n, err := conn.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error reading response:", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Server response: %s\n", string(buffer[:n]))
	// for true {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	_, err := fmt.Scan(">")
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	size, _ := reader.Read(buffer)
	// 	str, _ := reader.ReadString(buffer[size])
	// 	conn.Write([]byte(str))
	// 	fmt.Println(str)
	// }
}
