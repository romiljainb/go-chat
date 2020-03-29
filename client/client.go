package main

import (
	"net"
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
    "flag"
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


    port := flag.Int("port", 8080, "a port number")
	ip := flag.String("ip", "127.0.0.1", "a ip string")
	serverType := flag.String("type", "tcp", "a server type string")
	

    serverAddr := *ip + ":" + strconv.Itoa(*port)

    con, error := net.Dial(*serverType, serverAddr)
    if error != nil {
        fmt.Println(error)
        return
    }
    
    fmt.Println("Connected to %s", serverAddr)


    for {
        reader := bufio.NewReader(os.Stdin)
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error Reading stdin line, please try again")
        }

        if strings.HasPrefix(message,"quit()") {
            break
        }

        clientUser := Client{Name: "User", Msg: message}
        con.Write([]byte(clientUser.Msg))
    }

    con.Close()
    fmt.Printf("Connection to server %s closed.\n", serverAddr)
}