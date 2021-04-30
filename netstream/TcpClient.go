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

func RunTCPClient(serverAddr string) {
    conn, err := net.Dial("tcp", serverAddr)
    utils.LogInfo("TCPClient Connect to: " + serverAddr + "\n")
    
    if err != nil {
            fmt.Println("ERROR", err)
            os.Exit(1)
    }

 //   userInput := bufio.NewReader(os.Stdin)
    response := bufio.NewReader(conn)
    userInput :="1234567890123456789\n"
    //for {
   //     userLine, err := userInput.ReadBytes(byte('\n'))
        userLine := []byte(userInput)
        // switch err {
        // case nil:
            bs := make([]byte, 4)
            ul := uint32 (len(userLine))
            utils.Statistic("TCPClient Send Pkt [%d]\n", ul)
         //   fmt.Printf("Send Pkt Len:%d\n", ul)
            binary.BigEndian.PutUint32(bs, ul)
            conn.Write(bs)
            conn.Write(userLine)
        // case io.EOF:
        //     os.Exit(0)
        // default:
        //     fmt.Println("ERROR", err)
        //     os.Exit(1)
        // }

        serverLine, err := response.ReadBytes(byte('\n'))
        switch err {
        case nil:
            fmt.Print(string(serverLine))
            utils.Statistic("Received ECHO.\n")
        case io.EOF:
            os.Exit(0)
        default:
            fmt.Println("ERROR", err)
            os.Exit(2)
        }
    //}
}