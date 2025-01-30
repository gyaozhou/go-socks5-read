package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// $ go run normal_proxy.go -target "45.113.192.102:443"
func main() {
	target := flag.String("target", "http://example.org", "URL to get")
	proxyAddr := flag.String("proxy", ":1080", "normal proxy address to use")
	flag.Parse()

	listener, err := net.Listen("tcp", *proxyAddr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Println("Server is listening on", *proxyAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		fmt.Println("Connection from", conn.RemoteAddr())

		go handleRequest(conn, *target)
	}
}

func handleRequest(sConn net.Conn, target string) error {
	defer sConn.Close()

	tConn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return err
	}
	defer tConn.Close()

	fmt.Println("Connected to", target)

	/*
		go copyConn(sConn, tConn)

		buf := make([]byte, 1024)
		for {
			reqLen, err := tConn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				return
			}
			fmt.Println("[Received data] ", string(buf[:reqLen]))
			sConn.Write(buf[:reqLen])
		}

	*/

	// Start proxying
	errCh := make(chan error, 2)
	go iocopy(sConn, tConn, errCh)
	go iocopy(tConn, sConn, errCh)

	// Wait
	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			// return from this function closes target (and conn).
			return e
		}
	}
	return nil
}

/*
	func copyConn(sConn net.Conn, tConn net.Conn) {
		buf := make([]byte, 1024)
		for {
			reqLen, err := sConn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				return
			}
			fmt.Println("[Sent data]", string(buf[:reqLen]))
			tConn.Write(buf[:reqLen])
		}
	}
*/
func iocopy(dst io.Writer, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)

	//	if tcpConn, ok := dst.(closeWriter); ok {
	//		tcpConn.CloseWrite()
	//	}
	errCh <- err
}

/*
func getIP(url string) string {
	// Resolve the hostname to an IP address
	u, err := url.Parse(url)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return ""
	}
	addrs, err := net.LookupHost(u.Host)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return ""
	}

	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip != nil {
			fmt.Println(ip)
			// Use the IP address with net.Dial
			return fmt.Sprintf("%s:443", ip)
		}
	}
	return ""
}
*/
