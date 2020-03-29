package main

import (
	"net"
    "fmt"
    "os"
    "bufio"
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
10. client name 
11. ci/cd -> travis
12. memcache

*/

type Client struct {
    Name string
    Msg string
}

func main() {

	runClient()
}

func runClient() {
    // Connects to server
    con, error := net.Dial("tcp", "127.0.0.1:8080")
    if error != nil {
        fmt.Println(error)
        return
    }
    fmt.Println("Connected to 127.0.0.1:8080.")


    reader := bufio.NewReader(os.Stdin)
    message, _ := reader.ReadString('\n')

    message = "b :" + message
    clientUser := Client{Name: "User", Msg: message}

    con.Write([]byte(clientUser.Msg))

    con.Close()

    fmt.Println("Message sent. Connection closed.")
}