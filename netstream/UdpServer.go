package netstream


import (
	"fmt"
	"net"
	"os"
	"packet-verify/utils"
)

// CheckError checks for errors
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	//	return
		os.Exit(0)
	}
}

func RunUDPServer(udpAddr string) {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", udpAddr)
	CheckError(err)
	fmt.Printf("UDPServer Listen to:" + udpAddr)
	
	ch := make(chan *net.UDPConn)

	go func() {
		client, err := net.ListenUDP("udp", ServerAddr)
		CheckError(err)
		ch <- client
	//	defer client.Close()
	
	for {
		/* Now listen at selected port */
		go handleUdpConn(<-ch)
	}
	}()
}

func handleUdpConn(client *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := client.ReadFromUDP(buf)
	
		if err != nil {
			CheckError(err)
			//fmt.Println("error: ", err)
		}
		utils.Statistic("received: %s from: %s\n", string(buf[0:n]), addr)
		client.WriteTo([]byte("ECHO:"), addr)
		client.WriteTo(buf[0:n], addr)
	}
}