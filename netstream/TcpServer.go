package netstream

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"packet-verify/utils"
)

type TcpServer struct {

}

const PORT = 8080


func RunTCPServer(listenAddr string) {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fmt.Println("TCPServer Listent to: " + listenAddr)
	conns := clientConns(l)

	for {
		go handleTcpConn(<-conns)
	}

}

func clientConns(listener net.Listener) chan net.Conn {
    ch := make(chan net.Conn)
    i := 0
    go func() {
        for {
            client, err := listener.Accept()
            if client == nil {
                fmt.Printf("couldn't accept: " + err.Error())
                continue
            }
            i++
            fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
            ch <- client
        }
    }()
    return ch
}

func handleTcpConn(client net.Conn) {
    buf := bufio.NewReader(client)
    packet := make([]byte, 0, 4*1024)
    for {
        //line, err := b.ReadBytes('\n')
        n, err := io.ReadFull(buf, packet[:4])
        header := packet[:n]
        if err != nil {
            if err == io.EOF {
                break
            }
        }
        dataLen := binary.BigEndian.Uint32(header)
        utils.Statistic("%s Receive Packet Len: %d\n", PeerTag, dataLen)
        if uint64(dataLen) > uint64(cap(packet)) {
            packet = make([]byte, 0, dataLen)
        }
        n, err = io.ReadFull(buf, packet[:dataLen])
        packet = packet[:n]
        processPacket(packet, client)
    }
}

func processPacket(data []byte, client net.Conn)  {
    client.Write([]byte("ECHO:"))
    client.Write(data)
}