package chat

import (
	"net"
	"bufio"
	"strings"
	"github.com/romiljainb/lets-go/connections"
)

type Message struct {
	msg    string
	sender User
}


type User struct {
	name  string
	ID    int
	uconn net.Conn
	//uconn ConnHandler
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

