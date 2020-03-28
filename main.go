package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	msgs    = make(chan string)
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

		//msgs <- fmt.Sprintf("Client %v: %v", i, m)
		msgs <- fmt.Sprintf("%v", m)
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
		case msg := <-msgs:

			data := strings.Split(strings.TrimSpace(msg), ":")
			info := strings.Split(data[0], " ")

			fmt.Println(data)
			fmt.Println(info)

			if info[0] == "p" {
				fmt.Println(data[1])
				rec, err := strconv.Atoi(info[1])
				if err != nil {
					fmt.Println(err)
				}
				peers[rec].Write([]byte(data[1]))
			} else if info[0] == "b" {
				for conn := range clients {
					conn.Write([]byte(msg))
				}
			} else if info[0] == "g" {

			} else if info[0] == "j" {

				joinGroup(conn, i)
				fmt.Println("created a group ", groups)

			} else {
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
