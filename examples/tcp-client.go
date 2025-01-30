package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	target := flag.String("target", "localhost:8080", "URL to get")
	flag.Parse()

	//conn, err := net.Dial("tcp", "localhost:8080")
	conn, err := net.Dial("tcp", *target)

	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	/*
		file, ferr := conn.(*net.TCPConn).File()
		if ferr != nil {
			fmt.Println("Can not get file:", ferr.Error())
			return
		}

			// zhou: doesn't work
			serr := syscall.BindToDevice(int(file.Fd()), "eth0")
			if err != nil {
				fmt.Println("Error setting socket fd:", serr.Error())
				return
			}
	*/
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
