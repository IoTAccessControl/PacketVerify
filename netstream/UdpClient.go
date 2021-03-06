package netstream

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "os"
    "packet-verify/utils"
)

// const PORT = 8080

func RunUDPClient(serverAddr string) {
    conn, err := net.Dial("udp", serverAddr)
    
    
    if err != nil {
        fmt.Println("ERROR", err)
        os.Exit(1)
        //return
    }
  //  defer conn.Close()
    fmt.Printf("UDPClient Connect to: " + serverAddr)

    userInput := bufio.NewReader(os.Stdin)
    response := bufio.NewReader(conn)
    for {
        userLine, err := userInput.ReadBytes(byte('\n'))
        switch err {
        case nil:
            bs := make([]byte, 4)
            ul := uint32 (len(userLine))
            utils.Statistic("Send Pkt Len:%d\n", ul)
            binary.BigEndian.PutUint32(bs, ul)
            conn.Write(bs)
            conn.Write(userLine)
        case io.EOF:
            os.Exit(0)
        default:
            fmt.Println("ERROR", err)
            os.Exit(1)
        }

        serverLine, err := response.ReadBytes(byte('\n'))
        switch err {
        case nil:
            fmt.Print(string(serverLine))
        case io.EOF:
            os.Exit(0)
        default:
            fmt.Println("ERROR", err)
            os.Exit(2)
        }
    }
}