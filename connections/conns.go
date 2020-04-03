package connections

import (
)

type ConnInterface interface {
	Send() error
}

type ConnHandler struct {
    connType string
    connInf ConnInterface
}

type netConn struct {

}

