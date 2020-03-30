package main

import (
	"net"
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


func (client *User) broadcast(data []string, mgr *UserMgr ) {
	msg := client.name + ": " + data[1] + "\n"
	for _, user := range mgr.users{
		user.uconn.Write([]byte(msg))
	}
}

func (client *User) sendToPeer(data []string, info []string, mgr *UserMgr) {
	rec := info[1]

    _, exists := mgr.users[rec]
    if !exists {
        client.uconn.Write([]byte("User doesnt exist\n"))
    } else {
		msg := client.name + ": " + data[1] + "\n"
        (mgr.users[rec]).uconn.Write([]byte(msg))
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

