
package main
import (
	"net"
	"fmt"
	"strconv"
)
func sendToGrp(data []string, info []string, sender User) {
	rec, err := strconv.Atoi(info[1])
	if err != nil {
		fmt.Println(err)
	}
	if rec <= 0  {
		fmt.Println("group doesn't exist")
		sender.uconn.Write([]byte("group doesn't exist"))
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
			sender.uconn.Write([]byte("not member of group"))
		}

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

