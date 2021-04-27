package netstream

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "os"
)

// const PORT = 8080

func RunTCPClient(serverAddr string) {
    conn, err := net.Dial("tcp", serverAddr)
    fmt.Println("TCPClient Connect to: " + serverAddr)
    
    if err != nil {
            fmt.Println("ERROR", err)
            os.Exit(1)
    }

    userInput := bufio.NewReader(os.Stdin)
    response := bufio.NewReader(conn)
    for {
        userLine, err := userInput.ReadBytes(byte('\n'))
        switch err {
        case nil:
            bs := make([]byte, 4)
            ul := uint32 (len(userLine))
            fmt.Printf("Send Pkt Len:%d\n", ul)
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