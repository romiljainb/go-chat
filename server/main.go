package main

import (
	"flag"
	"fmt"
    "github.com/romiljainb/lets-go/chat"
    "github.com/romiljainb/lets-go/connections"
)

func main() {

	port := flag.Int("port", 8080, "a port number")
	ip := flag.String("ip", "127.0.0.1", "a ip string")
	serverType := flag.String("type", "http", "a server type string")

    server, err := getServer(port, ip, serverType)
    if err != nil {
		fmt.Println("Error starting server", err.Error())
    }

    server.setSrvAddr(ip, port)

	srvHandler, err := server.Start()
	if err != nil {
		fmt.Println("Error starting server", err.Error())
	}

	fmt.Println("Server Starting!!!")

	go server.AcceptConns(srvHandler)
	handleConns()

}

