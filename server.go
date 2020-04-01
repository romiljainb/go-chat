package main

import (
	"net"
	"fmt"
	"bufio"
	"errors"
	"strconv"
)
type Message struct {
	msg    string
	sender User
}

type DataPkt struct {
	msg string
}

type ConnHandler struct {

}

type ConnInterface interface {
	Send() error
}

type Server struct {
	ip net.IP
	port int
	serverAddr string
	srvType string
	srvLevel string
}

type TCPServer struct {
	ip net.IP
	port int
	serverAddr string
	srvType string
	srvLevel string
}

type SrvInterface interface {
	Start() (interface{}, error)
	Stop() error
	AcceptConns(srv interface{}) error
    setSrvAddr(ip *string, port *int)
	//GetSrvInfo() error
}

type SrvFactory struct {
    srvType *string
}


func getServer( port *int, ip *string, serverType *string) (SrvInterface,error) {
    switch *serverType{
    case "tcp":
        s := new(TCPServer)
        s.ip = net.ParseIP(*ip)
        s.port = *port
        s.srvType = *serverType
        s.srvLevel = "simple"
        return s, nil
    default:
        return nil, errors.New("invalid server type")
    }
    /*
        server = Server{
            ip: net.ParseIP(*ip),
            port: *port,
            srvType: *serverType,
            srvLevel: "simple",
        }
    */

}

// TODO: change net.Listener to something else like interface{} for other servers
// 1. cheapt solution
// 2. overloading
// 3. reflection (maybe)
func (server *Server) Start() (interface{}, error) {

	var ln interface{}
	var err error

	if server.srvType == "tcp" {
		ln, err = net.Listen(server.srvType, server.serverAddr)
		if err != nil {
			fmt.Println("Error starting server", err.Error())
			return nil, err
		}
	} else if server.srvType == "http" {

	} else if server.srvType == "grpc" {

	} else { //default

	}

	return ln, nil
}

func (server *TCPServer) Start() (interface{}, error) {

	var ln interface{}
	var err error
    ln, err = net.Listen(server.srvType, server.serverAddr)
    if err != nil {
        fmt.Println("Error starting server", err.Error())
        return nil, err
    }
	return ln, nil
}

func (server *Server) Stop() error {
	return nil
}

func (server *TCPServer) Stop() error {
	return nil
}


func (server *Server) AcceptConns(ln interface{}) error {

	for {
		conn, err := ln.(net.Listener).Accept()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		conns <- conn
	}
}

func (server *TCPServer) AcceptConns(ln interface{}) error {

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

func (server *Server) setSrvAddr(ip *string, port *int){
	server.serverAddr = *ip + ":" + strconv.Itoa(*port)
}

func (server *TCPServer) setSrvAddr(ip *string, port *int){
	server.serverAddr = *ip + ":" + strconv.Itoa(*port)
}
