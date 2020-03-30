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
	Start() (net.Listener, error)
	Stop() error
	AcceptConns(srv interface{}) error
	GetSrvInfo() error
} 

// TODO: change net.Listener to something else like interface{} for other servers
func (server *Server) Start() (net.Listener, error) {
	ln, err := net.Listen(server.srvType, server.serverAddr)
	if err != nil {
		fmt.Println("Error starting server", err.Error())
		return nil, err
	}
	return ln, nil
}

func (server *Server) Stop() error {
	return nil
}

func (server *Server) AcceptConns(ln net.Listener) error {

	for {
		conn, err := ln.(net.Listener).Accept()
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