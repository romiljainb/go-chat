package main

import (
	"net"
	"fmt"
	"bufio"
)

type Message struct {
	msg    string
	sender User
}

type Server struct {
	ip net.IP
	port int
	serverAddr string
	srvType string
	srvLevel string
}
type SrvInterface interface {
	Start() error
	Stop() error
	AcceptConns() error
	GetSrvInfo() error
} 

func (server *Server) acceptConn(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		conns <- conn
	}
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