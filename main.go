package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type message struct {
	msg        string
	sender     int
	connection net.Conn
}

var (
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	msgs    = make(chan message)
	clients = make(map[net.Conn]int)
	peers   = make(map[int]net.Conn)
	groups  = make(map[int][]net.Conn)
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

func joinGroup(conn net.Conn, i int) {
	var newGrp []net.Conn
	if newGrp == nil && groups == nil {
		newGrp = append(newGrp, conn)
		groups[i] = newGrp
	}
	//return groups
}

func readConn(conn net.Conn, i int) {
	rd := bufio.NewReader(conn)
	for {
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		mdata := message{msg: m, sender: i, connection: conn}
		msgs <- mdata
	}
	dconns <- conn

}

// p 1: hi there
// p 3: hi there back
// b : megs to all
// g 2: hi all in group 2
// j 2:

func handleConns() {
	i := 0

	for {
		select {
		// read the incoming messages
		case conn := <-conns:
			clients[conn] = i
			i++
			peers[i] = conn

			go readConn(conn, i)

		// msg must be broadcast to everyone
		case message := <-msgs:

			data := strings.Split(strings.TrimSpace(message.msg), ":")
			info := strings.Split(data[0], " ")

			if info[0] == "p" {
				handlePeer(data, info, message.sender)
			} else if info[0] == "b" {
				for conn := range clients {
					conn.Write([]byte(data[1]))
				}
			} else if info[0] == "g" {

			} else if info[0] == "j" {

				joinGroup(message.connection, i)
				fmt.Println("created a group ", groups)

			} else {
				peers[message.sender].Write([]byte("Error parsing message info\n"))
				fmt.Println("Error parsing message info")
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
