package main

import (
	"bufio"
	"fmt"
	"net"
)

var (
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	msgs    = make(chan string)
	clients = make(map[net.Conn]int)
	//peers = make(map[int]net.Conn)
	//groups = make(map[int][]net.Conn)
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

func readConn(conn net.Conn, i int) {
	rd := bufio.NewReader(conn)
	for {
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		msgs <- fmt.Sprintf("Client %v: %v", i, m)
	}
	dconns <- conn

}

func handleConns() {
	i := 0

	for {
		select {
		// read the incoming messages
		case conn := <-conns:
			clients[conn] = i
			i++
			go readConn(conn, i)

		// msg must be broadcast to everyone
		case msg := <-msgs:
			for conn := range clients {
				conn.Write([]byte(msg))
			}
		case dconn := <-dconns:
			fmt.Println("Clinet %v logged off", clients[dconn])
			delete(clients, dconn)
		}

	}

}

func main() {

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error starting server", err.Error())
	}
	fmt.Println("Server Starting!!!")

	go acceptConn(ln)
	handleConns()

}
