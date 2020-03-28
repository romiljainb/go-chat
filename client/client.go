package client

import (
	"net"
	"fmt"
	"encoding/gob"
)

/*
1. encoding/decoding
2. compression/decompress
3. serialization/deserization
4. use flags to lauch server tcp/http/ws
5. add group leader feature
6. add simple security
7. TLS
8. TCP network configuration
9. Cli for client app
10. 

*/

func main() {

	runClient()
}

func runClient() {
    // Connects to server
    con, error := net.Dial("tcp", "127.0.0.1:8272")

    // Handles eventual errors
    if error != nil {
        fmt.Println(error)
        return
    }

    fmt.Println("Connected to 127.0.0.1:8080.")

    // Sends a message
    message := "Hello world"
    encoder := gob.NewEncoder(con)
    error = encoder.Encode(message)

    // Checks for errors
    if error != nil {
        fmt.Println(error)
    }

    con.Close()

    fmt.Println("Message sent. Connection closed.")
}