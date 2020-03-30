package main

import (
	"net"
	"strconv"
	"fmt"
	"bufio"
	"strings"
)

type Message struct {
	msg    string
	sender User
}
type User struct {
	name  string
	ID    int
	uconn net.Conn
}


type UserInterface interface {
    broadcast()
    sendToPeer(peer string)
    leaveGroup(group string)
    joinGroup(group string)
}


func (client User) broadcast(data []string, mgr *UserMgr ) {
	msg := client.name + ": " + data[1] + "\n"
	for _, user := range mgr.users{
		user.uconn.Write([]byte(msg))
	}
}

func (sender User) sendToPeer(data []string, info []string ) {
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


func getUserDetails(conn net.Conn, id int) (User, bool) {
	rd := bufio.NewReader(conn)
	var user User
	for i := 0; i < 3; i++ {
		conn.Write([]byte("Enter an username: \n"))
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		pname := strings.TrimSpace(m)

		if len(pname) > 0 {
			user = User{name: pname, ID: id, uconn: conn}
			users[user] = true
			return user, true
		}

	}
	return user, false

}

