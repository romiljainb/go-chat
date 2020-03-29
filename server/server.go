package main

import (
	"net"
)
type Server struct {
	ip net.IpAddr
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

/*
func (server Server) acceptConn(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		conns <- conn
	}
}
*/