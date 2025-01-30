package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "net/url"
    "strconv"
    "time"

    "golang.org/x/net/proxy"
)

// handleConnection handles a single client connection
func handleConnection(clientConn net.Conn, upstreamProxy string) {
    defer clientConn.Close()

    // Set up a SOCKS5 dialer to connect to the upstream proxy
    dialer, err := proxy.SOCKS5("tcp", upstreamProxy, nil, proxy.Direct)
    if err != nil {
        log.Printf("Failed to create SOCKS5 dialer: %v", err)
        return
    }

    // Read the SOCKS version and number of authentication methods
    buf := make([]byte, 256)
    _, err = io.ReadFull(clientConn, buf[:2])
    if err != nil {
        log.Printf("Failed to read SOCKS version and auth methods: %v", err)
        return
    }
    version := buf[0]
    nmethods := int(buf[1])

    // Read the authentication methods
    _, err = io.ReadFull(clientConn, buf[:nmethods])
    if err != nil {
        log.Printf("Failed to read authentication methods: %v", err)
        return
    }

    // For simplicity, we assume no authentication is required
    clientConn.Write([]byte{0x05, 0x00})

    // Read the SOCKS request
    _, err = io.ReadFull(clientConn, buf[:4])
    if err != nil {
        log.Printf("Failed to read SOCKS request: %v", err)
        return
    }
    cmd := buf[1]
    atyp := buf[3]

    var targetAddr string
    switch atyp {
    case 0x01: // IPv4 address
        _, err = io.ReadFull(clientConn, buf[:4])
        if err != nil {
            log.Printf("Failed to read IPv4 address: %v", err)
            return
        }
        ip := net.IPv4(buf[0], buf[1], buf[2], buf[3])
        _, err = io.ReadFull(clientConn, buf[:2])
        if err != nil {
            log.Printf("Failed to read port: %v", err)
            return
        }
        port := (int(buf[0]) << 8) | int(buf[1])
        targetAddr = fmt.Sprintf("%s:%d", ip.String(), port)
    case 0x03: // Domain name
        _, err = io.ReadFull(clientConn, buf[:1])
        if err != nil {
            log.Printf("Failed to read domain name length: %v", err)
            return
        }
        domainLen := int(buf[0])
        _, err = io.ReadFull(clientConn, buf[:domainLen])
        if err != nil {
            log.Printf("Failed to read domain name: %v", err)
            return
        }
        domain := string(buf[:domainLen])
        _, err = io.ReadFull(clientConn, buf[:2])
        if err != nil {
            log.Printf("Failed to read port: %v", err)
            return
        }
        port := (int(buf[0]) << 8) | int(buf[1])
        targetAddr = fmt.Sprintf("%s:%d", domain, port)
    case 0x04: // IPv6 address
        _, err = io.ReadFull(clientConn, buf[:16])
        if err != nil {
            log.Printf("Failed to read IPv6 address: %v", err)
            return
        }
        ip := net.IP(buf[:16])
        _, err = io.ReadFull(clientConn, buf[:2])
        if err != nil {
            log.Printf("Failed to read port: %v", err)
            return
        }
        port := (int(buf[0]) << 8) | int(buf[1])
        targetAddr = fmt.Sprintf("%s:%d", ip.String(), port)
    default:
        log.Printf("Unsupported address type: %d", atyp)
        return
    }

    if cmd == 0x01 { // CONNECT command
        // Connect to the target server through the upstream proxy
        targetConn, err := dialer.Dial("tcp", targetAddr)
        if err != nil {
            log.Printf("Failed to connect to target server: %v", err)
            return
        }
        defer targetConn.Close()

        // Send the success response to the client
        clientConn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

        // Start forwarding data between the client and the target server
        go func() {
            _, err := io.Copy(targetConn, clientConn)
            if err != nil {
                log.Printf("Error forwarding data from client to target: %v", err)
            }
        }()
        go func() {
            _, err := io.Copy(clientConn, targetConn)
            if err != nil {
                log.Printf("Error forwarding data from target to client: %v", err)
            }
        }()

        // Wait for the connections to close
        select {
        case <-time.After(30 * time.Minute):
            log.Println("Connection timed out")
        }
    }
}

func main() {
    listenAddr := ":1080"
    upstreamProxy := "127.0.0.1:9050" // Replace with your actual upstream SOCKS5 proxy address

    ln, err := net.Listen("tcp", listenAddr)
    if err != nil {
        log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
    }
    log.Printf("SOCKS proxy server listening on %s", listenAddr)

    for {
        clientConn, err := ln.Accept()
        if err != nil {
            log.Printf("Failed to accept client connection: %v", err)
            continue
        }
        go handleConnection(clientConn, upstreamProxy)
    }
}    