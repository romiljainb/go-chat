package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type User struct {
	name  string
	ID    int
	uconn net.Conn
}

type Message struct {
	msg    string
	sender User
}

type UserMgr struct {
    users map[string]User
    groups map[string][]User

}

type UserInterface interface {
    broadcast()
    sendToPeer(peer string)
    leaveGroup(group string)
    joinGroup(group string)
}

type UserMgrInterface interface {
    addUser(user User)
    removeUser(user User)
    getUser(name string) User
}

type ConnectionManager interface {
}

type Manager struct {
}

var (
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	msgs    = make(chan Message)
	clients = make(map[net.Conn]int)
	peers   = make(map[int]net.Conn)
	users   = make(map[User]bool)
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

func leaveGroup(conn net.Conn, i int) {
	value, ok := groups[i]
	if ok {
		for index, val := range value {
			if val == conn {

				value[len(value)-1], value[index] = value[index], value[len(value)-1]
				value = value[:len(value)-1]
				fmt.Println(value)
				conn.Write([]byte("you left group" + "\n"))
				sendMsgToAll(value, conn, clients[conn])
				break

			} else {
				fmt.Println("not a member of this group" + "\n")
			}

		}

	} else {
		fmt.Println("group doesnt exist" + "\n")
	}

}

func sendMsgToAll(clients []net.Conn, client net.Conn, clientID int) {
	for _, conn := range clients {
		conn.Write([]byte(strconv.Itoa(clientID) + " has left this group" + "\n"))
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

func readConn(conn net.Conn, user User) {
	rd := bufio.NewReader(conn)
	for {
		m, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		mdata := Message{msg: m, sender: user}
		msgs <- mdata
	}
	dconns <- conn

}

func handlePeer(data []string, info []string, sender User) {
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

func handleBroadcast(data []string, sender User) {
	msg := sender.name + ": " + data[1] + "\n"
	for conn := range clients {
		conn.Write([]byte(msg))
	}

}

func sendToGrp(data []string, info []string, sender User) {
	rec, err := strconv.Atoi(info[1])
	if err != nil {
		fmt.Println(err)
	}
	if rec <= 0 {
		fmt.Println("group doesn't exist" + "\n")
		sender.uconn.Write([]byte("group doesn't exist" + "\n"))
	} else {
		value, ok := groups[rec]
		if ok {
			if existIn(value, sender.uconn) == true {
				msg := sender.name + " : " + data[1] + "\n"
				for _, receiver := range groups[rec] {
					receiver.Write([]byte(msg))
				}
			} else {
				sender.uconn.Write([]byte("not member of group" + "\n"))
			}

		} else {
			//sender.uconn.Write([]byte("group doesn't exist" + "\n"))
		}

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
		// read the incoming Messages
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

func main() {

	port := flag.Int("port", 8080, "a port number")
	ip := flag.String("ip", "127.0.0.1", "a ip string")
	serverType := flag.String("type", "tcp", "a server type string")

	serverAddr := *ip + ":" + strconv.Itoa(*port)

	ln, err := net.Listen(*serverType, serverAddr)

	if err != nil {
		fmt.Println("Error starting server", err.Error())
	}
	fmt.Println("Server Starting!!!")

	go acceptConn(ln)
	handleConns()

}
