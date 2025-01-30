package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	// "time"
)

func main() {
	client1()
	client2()
	client3()
}

func client1() {
	fmt.Println("client1() ...")
	// Specify the local address and port
	localAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0")

	// Specify the remote address and port
	remoteAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")

	// Create a TCP connection
	conn, cErr := net.DialTCP("tcp", localAddr, remoteAddr)

	if cErr != nil {
		fmt.Println(cErr)
		return
	}

	// Set the local address for the connection
	fmt.Println(conn.LocalAddr())

	// Read and write data
	_, err := conn.Write([]byte("Hello, server!"))
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make([]byte, 1024)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	// Close the connection
	conn.Close()
}

func client2() {
	fmt.Println("client2() ...")

	// Specify the local address and port
	//localAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0")

	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0},
		//LocalAddr: localAddr,
	}

	conn, err := dialer.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Set the local address for the connection
	fmt.Println(conn.LocalAddr())

	// Read and write data
	_, err = conn.Write([]byte("Hello, server!"))
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make([]byte, 1024)
	_, err = conn.Read(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	// Close the connection
	conn.Close()
}

func client3() {
	fmt.Println("client3() ...")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// local address for the connection
	fmt.Println(conn.LocalAddr())

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter text to send to server. Type 'exit' to quit.")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
		conn.Write([]byte(input))
	}
}
