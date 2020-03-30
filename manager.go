package main

import (
	"net"
	"fmt"
	"strings"
	"strconv"
)

var (
	clients = make(map[net.Conn]int)
	peers   = make(map[int]net.Conn)
	users   = make(map[User]bool)
	groups  = make(map[int][]net.Conn)
)

type UserMgr struct {
    users map[string]User
    groups map[string][]User
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



func handleConns() {
	i := 1

    mgr := UserMgr{users:make(map[string]User), groups:make(map[string][]User), conlist:make(map[net.Conn]int),}

	for {
		select {
		// read the incoming Messages
		case conn := <-conns:
			_, exist := clients[conn]
			if !exist {
				createUser(conn, i, &mgr)
				i++
			}


		// msg must be broadcast to everyone
		case message := <-msgs:
			if len(strings.TrimSpace(message.msg)) == 0 {
				continue
			}

            client := message.sender

			data := strings.Split(strings.TrimSpace(message.msg), ":")
			info := strings.Split(data[0], " ")

			if info[0] == "p" {
                client.sendToPeer(data, info, &mgr)
			} else if info[0] == "b" {
                client.broadcast(data, &mgr)
			} else if info[0] == "g" {
				sendToGrp(data, info, message.sender)

			} else if info[0] == "j" {
				groupID, err := strconv.Atoi(info[1])
				if err != nil {
					fmt.Println("error occured", err)
				}
				joinGroup(message.sender.uconn, groupID)
			} else if info[0] == "l" {
				groupID, err := strconv.Atoi(info[1])
				if err != nil {
					fmt.Println("error occured", err)
				}
				leaveGroup(message.sender.uconn, groupID)
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

