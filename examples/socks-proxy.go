package main

import (
	"flag"
	"fmt"
	"go-socks5"
)

// zhou: implement "type CredentialStore interface{}".
//	     Why not use "socks5.StaticCredentials" directly ???
/*
type myCredentialStore struct {
	user     string
	password string
}

func (cs *myCredentialStore) Valid(user, password string) bool {
	return user == cs.user && password == cs.password
}
*/

// $ go run server.go -h ":1080" -u test -p test
func main() {
	listen := flag.String("h", ":1080", "IP:Port for SOCKS5 proxy")
	username := flag.String("u", "", "username for SOCKS5 proxy")
	password := flag.String("p", "", "password for SOCKS5 proxy")
	flag.Parse()

	var conf *socks5.Config

	if *username != "" && *password != "" {
		// basic Auth
		// curl --socks5 127.0.0.1:1090 https://www.baidu.com --proxy-user test:test
		myCredentialStore := make(socks5.StaticCredentials)
		myCredentialStore[*username] = *password

		auth := socks5.UserPassAuthenticator{
			//Credentials: &myCredentialStore{user: *username, password: *password},
			Credentials: &myCredentialStore,
		}

		conf = &socks5.Config{
			AuthMethods: []socks5.Authenticator{auth},
		}
		fmt.Println("Socks5 Proxy", *listen, "with Basic Auth ", *username, *password)
	} else {
		// No Auth
		// curl --socks5 127.0.0.1:1080 https://www.baidu.com
		conf = &socks5.Config{}
		fmt.Println("Socks5 Proxy", *listen, "with No Auth")
	}

	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	if err := server.ListenAndServe("tcp", *listen); err != nil {
		panic(err)
	}
}
