package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/proxy"
)

// $ go run client.go -proxy 127.0.0.1:1080 -user test -pass test -target https://www.baidu.com
// $ curl --socks5 127.0.0.1:1080 https://www.baidu.com --proxy-user test:test
func main() {
	target := flag.String("target", "http://example.org", "URL to get")
	proxyAddr := flag.String("proxy", "localhost:1080", "SOCKS5 proxy address to use")
	username := flag.String("user", "", "username for SOCKS5 proxy")
	password := flag.String("pass", "", "password for SOCKS5 proxy")
	flag.Parse()

	auth := proxy.Auth{
		User:     *username,
		Password: *password,
	}
	// Basic Auth
	dialer, err := proxy.SOCKS5("tcp", *proxyAddr, &auth, nil)
	// No Auth
	//dialer, err := proxy.SOCKS5("tcp", *proxyAddr, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	// zhou: TCP using proxy.
	// dialer.Dial("tcp", "192.168.0.1:1234")

	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	r, err := client.Get(*target)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	server.ListenAndServe()
}
