package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
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
	value, ok := groups[i]
	if ok {
		if existIn(value, conn) == false {
			groups[i] = append(groups[i], conn)
			fmt.Println("Joined a group ", groups)
		} else {
			fmt.Println("already member of: ", groups)
		}

	} else {
		newGrp = append(newGrp, conn)
		groups[i] = newGrp
		for _, value = range groups {
			fmt.Println(value)
		}
		fmt.Println("you have created new group ", groups)
	}
}

func existIn(memList []net.Conn, conn net.Conn) bool {
	for _, val := range memList {
		if val == conn {
			return true
		}

	}
	return false
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

func handlePeer(data []string, info []string, sender int) {
	rec, err := strconv.Atoi(info[1])
	if err != nil {
		fmt.Println(err)
	}

	if rec <= 0 || rec > len(peers) {
		peers[sender].Write([]byte("Reciever is out of range\n"))
	} else {

		msg := "Client " + strconv.Itoa(sender) + " : " + data[1] + "\n"
		peers[rec].Write([]byte(msg))
	}

}

func handleBroadcast(data []string, sender int) {
	msg := "Client " + strconv.Itoa(sender) + " : " + data[1] + "\n"
	for conn := range clients {
		conn.Write([]byte(msg))
	}

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

			data := strings.Split(strings.TrimSpace(message.msg), ":") //j 2
			info := strings.Split(data[0], " ")                        //[j,2]

			if info[0] == "p" {
				handlePeer(data, info, message.sender)
			} else if info[0] == "b" {
				handleBroadcast(data, message.sender)
			} else if info[0] == "g" {

			} else if info[0] == "j" {
				groupID, err := strconv.Atoi(info[1])
				if err != nil {
					fmt.Println("error occured", err)
				}
				joinGroup(message.connection, groupID)

				//fmt.Println("error occured", err)

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
