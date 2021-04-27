package netstream


import (
	"fmt"
	"net"
	"strconv"
	"os"
)

// CheckError checks for errors
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func RunUDPServer() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(PORT))
	CheckError(err)
	fmt.Println("listening on :" + strconv.Itoa(PORT))
	
	ch := make(chan *net.UDPConn)

	go func() {
		client, err := net.ListenUDP("udp", ServerAddr)
		CheckError(err)
		ch <- client
	}()
	for {
		/* Now listen at selected port */
	
		// CheckError(err)
		// defer client.Close()
		go handleUdpConn(<-ch)
	}
}

func handleUdpConn(client *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := client.ReadFromUDP(buf)
		fmt.Printf("received: %s from: %s\n", string(buf[0:n]), addr)
	
		if err != nil {
			fmt.Println("error: ", err)
		}
		client.WriteTo([]byte("ECHO:"), addr)
		client.WriteTo(buf[0:n], addr)
	}
}