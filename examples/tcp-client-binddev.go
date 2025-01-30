package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
)

func main() {
	// Specify the interface and port to dial
	interfaceName := "ens35"
	port := 8080

	// Get the interface's IP address
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Println("Error getting interface:", err)
		return
	}

	// Create a dialer with the interface
	dialer := &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := syscall.BindToDevice(int(fd), iface.Name)
				if err != nil {
					fmt.Println("Error binding to interface:", err)
				}
			})
		},
	}

	// Dial the specific interface and port
	conn, err := dialer.Dial("tcp", fmt.Sprintf("20.11.22.33:%d", port))
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected from", conn.LocalAddr())
	fmt.Println("Connected to", conn.RemoteAddr())

	// Send data
	// data := []byte("Hello, world!")
	// _, err = conn.Write(data)
	// if err != nil {
	// 	fmt.Println("Error sending data:", err)
	// 	return
	// }

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

	fmt.Println("Data sent successfully")
}
