package main

import (
	"net"
)

var (
	conns   = make(chan net.Conn)
	connhandles   = make(chan ConnHandler)
	dconns  = make(chan net.Conn)
	msgs    = make(chan Message)
)
