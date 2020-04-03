package chat

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"bufio"
	connH "github.com/romiljainb/lets-go/connections"
)

var (
	clients = make(map[net.Conn]int)
	peers   = make(map[int]net.Conn)
	users   = make(map[User]bool)
	groups  = make(map[int][]net.Conn)
)

type UserMgr struct {
	users   map[string]User
	groups  map[string][]User
	conlist map[net.Conn]int
}

type UserMgrInterface interface {
	addUser(user User)
	removeUser(user User)
	getUser(name string) User
}

func createUser(conn net.Conn, id int, mgr *UserMgr) {
	user, created := getUserDetails(conn, id)
	if created {
		//old
		clients[conn] = id
		peers[id] = conn

		//new
		mgr.users[user.name] = user

		go readConn(conn, user)
	}

}

func HandleConns() {
	i := 1

	mgr := UserMgr{users: make(map[string]User), groups: make(map[string][]User), conlist: make(map[net.Conn]int)}

	for {
		select {
		// read the incoming Messages
		case conn := <- connH.Conns:
			_, exist := clients[conn]
			if !exist {
				createUser(conn, i, &mgr)
				i++
			}

		// msg must be broadcast to everyone
		case dataPkt := <-msgs:
			if len(strings.TrimSpace(dataPkt.msg)) == 0 {
				continue
			}

			client := dataPkt.sender

			data := strings.Split(strings.TrimSpace(dataPkt.msg), ":")
			info := strings.Split(data[0], " ")

			if info[0] == "p" {
            client.sendToPeer(data, info, &mgr)
			} else if info[0] == "b" {
				client.broadcast(data, &mgr)
			} else if info[0] == "g" {
				client.sendToGrp(data, info, &mgr)

			} else if info[0] == "j" {
				groupID, err := strconv.Atoi(info[1])
				if err != nil {
					fmt.Println("error occured", err)
				}
				joinGroup(dataPkt.sender.uconn, groupID)
			} else if info[0] == "l" {
				groupID, err := strconv.Atoi(info[1])
				if err != nil {
					fmt.Println("error occured", err)
				}
				leaveGroup(dataPkt.sender.uconn, groupID)
			} else {
				peers[dataPkt.sender.ID].Write([]byte("Error parsing message info\n"))
				fmt.Println("Error parsing message info")
			}
		case dconn := <-connH.Dconns:
			fmt.Println("Clinet %v logged off", clients[dconn])
			delete(clients, dconn)
		}
	}
}

func readConn(conn net.Conn, user User) {
	rd := bufio.NewReader(conn)
	for {
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		dataPkt := DataPkt{msg: m, sender: user}
		msgs <- dataPkt
	}
	connH.Dconns <- conn

}
