package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
)


var (
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	msgs    = make(chan Message)
)

func acceptConn(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		conns <- conn
	}
}



func readConn(conn net.Conn, user User) {
	rd := bufio.NewReader(conn)
	for {
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		mdata := Message{msg: m, sender: user}
		msgs <- mdata
	}
	dconns <- conn

}



func main() {

	port := flag.Int("port", 8080, "a port number")
	ip := flag.String("ip", "127.0.0.1", "a ip string")
	serverType := flag.String("type", "tcp", "a server type string")

	serverAddr := *ip + ":" + strconv.Itoa(*port)

	ln, err := net.Listen(*serverType, serverAddr)

	if err != nil {
		fmt.Println("Error starting server", err.Error())
	}
	fmt.Println("Server Starting!!!")

	go acceptConn(ln)
	handleConns()

}
