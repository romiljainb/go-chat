package connections

import (
)

type ConnHandler struct {
    connType string
    connInf interface{}
}

type netConn struct {

}

type ConnInterface interface {
	Send() error
}