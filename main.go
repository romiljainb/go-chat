package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type client struct {
	name  string
	ID    int
	uconn net.Conn
}
type message struct {
	msg    string
	sender client
}

type ConnectionManager interface {
}

type Manager struct {
}

var (
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	msgs    = make(chan message)
	clients = make(map[net.Conn]int)
	peers   = make(map[int]net.Conn)
	users   = make(map[client]bool)
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

func readConn(conn net.Conn, user client) {
	rd := bufio.NewReader(conn)
	for {
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		mdata := message{msg: m, sender: user}
		msgs <- mdata
	}
	dconns <- conn

}

func handlePeer(data []string, info []string, sender client) {
	rec, err := strconv.Atoi(info[1])
	if err != nil {
		fmt.Println(err)
	}

	if rec <= 0 || rec > len(peers) {
		peers[sender.ID].Write([]byte("Reciever is out of range\n"))
	} else {

		msg := sender.name + ": " + data[1] + "\n"
		peers[rec].Write([]byte(msg))
	}

}

func handleBroadcast(data []string, sender client) {
	msg := sender.name + ": " + data[1] + "\n"
	for conn := range clients {
		conn.Write([]byte(msg))
	}

}

func getUserDetails(conn net.Conn, id int) (client, bool) {
	rd := bufio.NewReader(conn)
	var user client
	for i := 0; i < 3; i++ {
		conn.Write([]byte("Enter an username: \n"))
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		pname := strings.TrimSpace(m)

		if len(pname) > 0 {
			user = client{name: pname, ID: id, uconn: conn}
			users[user] = true
			return user, true
		}

	}
	return user, false

}

func createUser(conn net.Conn, id int) {
	user, created := getUserDetails(conn, id)
	if created {
		clients[conn] = id
		peers[id] = conn
		go readConn(conn, user)
	}

}
func handleConns() {
	i := 1

	for {
		select {
		// read the incoming messages
		case conn := <-conns:
			_, exist := clients[conn]
			if !exist {
				createUser(conn, i)
				i++
			}

		// msg must be broadcast to everyone
		case message := <-msgs:
			if len(strings.TrimSpace(message.msg)) == 0 {
				continue
			}
			data := strings.Split(strings.TrimSpace(message.msg), ":")
			info := strings.Split(data[0], " ")

			if info[0] == "p" {
				handlePeer(data, info, message.sender)
			} else if info[0] == "b" {
				handleBroadcast(data, message.sender)
			} else if info[0] == "g" {

			} else {
				peers[message.sender.ID].Write([]byte("Error parsing message info\n"))
				fmt.Println("Error parsing message info")
			}
		case dconn := <-dconns:
			fmt.Println("Clinet %v logged off", clients[dconn])
			delete(clients, dconn)
		}
	}
}

func main() {

	port := flag.Int("port", 8080, "a port number")
	ip := flag.String("ip", "127.0.0.1", "a ip string")
	serverType := flag.String("type", "tcp", "a server type string")

	ipPort := *ip + ":" + strconv.Itoa(*port)

	ln, err := net.Listen(*serverType, ipPort)
	if err != nil {
		fmt.Println("Error starting server", err.Error())
	}
	fmt.Println("Server Starting!!!")

	go acceptConn(ln)
	handleConns()

}
