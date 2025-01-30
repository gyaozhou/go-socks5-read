package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		/*
			file, ferr := conn.(*net.TCPConn).File()
			if ferr != nil {
				fmt.Println("Can not get file:", ferr.Error())
				return
			}

				fd, err := syscall.GetsockoptInt(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_TYPE)
				if err != nil {
					fmt.Println("Error getting socket fd:", err.Error())
					continue
				}

					if err := syscall.SetNonblock(fd, true); err != nil {
						fmt.Println("Error setting non-blocking mode:", err.Error())
						continue
					}
		*/
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		fmt.Println("Received data:", string(buf[:reqLen]))
		conn.Write(buf[:reqLen])
	}
}
