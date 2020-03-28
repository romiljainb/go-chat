package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"strconv"
)

type message struct{
    msg string
    sender int
}

var (
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	msgs    = make(chan message)
	clients = make(map[net.Conn]int)
	peers = make(map[int]net.Conn)
	groups = make(map[int][]net.Conn)
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

        mdata := message{msg: m, sender: i}
        msgs <- mdata
	}
	dconns <- conn

}


func handlePeer(data []string, info []string) {
    rec, err := strconv.Atoi(info[1])
    if err != nil {
        fmt.Println(err)
    }
    peers[rec].Write([]byte(data[1]))

}

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
                handlePeer(data, info)
			} else if info[0] == "b" {
				for conn := range clients {
					conn.Write([]byte(data[1]))
				}
			} else if info[0] == "g" {

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
