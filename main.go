package main

import (
	"flag"
	"fmt"
)

func main() {

	port := flag.Int("port", 8080, "a port number")
	ip := flag.String("ip", "127.0.0.1", "a ip string")
	serverType := flag.String("type", "tcp", "a server type string")

/*
	server := Server{
		ip: net.ParseIP(*ip),
		port: *port,
		srvType: *serverType,
		srvLevel: "simple",
	}
*/
    server, err := getServer(port, ip, serverType)
    if err != nil {
		fmt.Println("Error starting server", err.Error())
    }
	//server.serverAddr = *ip + ":" + strconv.Itoa(*port)
    server.setSrvAddr(ip, port)

	srvHandler, err := server.Start()
	if err != nil {
		fmt.Println("Error starting server", err.Error())
	}

	fmt.Println("Server Starting!!!")

	go server.AcceptConns(srvHandler)
	handleConns()

}

