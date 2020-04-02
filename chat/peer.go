package chat

import (
	"github.com/romiljainb/lets-go/connections"
)

func (client *User) broadcast(data []string, mgr *UserMgr ) {
	msg := client.name + ": " + data[1] + "\n"
	//dataPkt := DataPkt{msg: msg}
	for _, user := range mgr.users{
		user.uconn.Write([]byte(msg))
		// if srvType == "simple" 
		//     sendData = dataPkt.msg.([]byte)
		// else 
		//     sendData = dataPkt.([]byte)
		//user.uconn.Write(sendData)
	}
}

func (client *User) sendToPeer(data []string, info []string, mgr *UserMgr) {
	rec := info[1]

    _, exists := mgr.users[rec]
    if !exists {
        client.uconn.Write([]byte("User doesnt exist\n"))
        // user.connHandler.Write(struct)
    } else {
		msg := client.name + ": " + data[1] + "\n"
        (mgr.users[rec]).uconn.Write([]byte(msg))
    }
}
