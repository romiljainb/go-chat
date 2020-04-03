package chat

import (
	"net"
	"bufio"
	"strings"
	"time"
	connH "github.com/romiljainb/lets-go/connections"
)

var msgs = make(chan DataPkt)

type DataPkt struct {
	msg string
	sender User
	msgSent time.Time
	msgOrigin string
}

type User struct {
	name  string
	ID    int
	uconn net.Conn
	uconnHandler connH.ConnHandler
}


type UserInterface interface {
    broadcast()
    sendToPeer(peer string)
    leaveGroup(group string)
    joinGroup(group string)
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

