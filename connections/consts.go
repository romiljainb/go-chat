package connections

import (
	"net"
)

var (
	Conns   	= make(chan net.Conn)
	Connhandles = make(chan ConnHandler)
	Dconns  	= make(chan net.Conn)
)
