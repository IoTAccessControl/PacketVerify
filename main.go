package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"packet-verify/netstream"
	"packet-verify/proxy"
	"packet-verify/utils"
	"syscall"
)

/*
localClient -----> localProxy -----> remoteProxy -----> remoteServer

 */

const HOST = "127.0.0.1"

const (
	LocalProxy = "LocalProxy"
	RemoteProxy = "RemoteProxy"
	TCPClient = "TCPClient"
	TCPServer = "TCPServer"
	UDPClient = "UDPClient"
	UDPServer = "UDPServer"
)

var localProxyAddr = HOST + ":8082"
var remoteProxyAddr = HOST + ":8083"
var serverAddr = HOST + ":8084"

func main() {

	SetupCloseHandler()

	mode := flag.String("mode", "TCPServer", "set request mode")
	//sever := flag.String("server", "xxx", "set sever addr")

	flag.Parse()

	//fmt.Println(*mode, *sever)

	netstream.PeerTag = *mode

 	utils.InitLogger(netstream.PeerTag)

	if *mode == TCPServer {
		netstream.RunTCPServer(serverAddr)
	} else if *mode == TCPClient {
		netstream.RunTCPClient(localProxyAddr)
	} else if *mode == UDPServer {
		netstream.RunUDPServer()
	} else if *mode == LocalProxy {
		proxy.StartLocalProxy(localProxyAddr, remoteProxyAddr)
	} else if *mode == RemoteProxy {
		proxy.StartRemoteProxy(remoteProxyAddr, serverAddr)
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		// kill all process
		exec.Command("pkill -f packet-verify")
		os.Exit(0)
	}()
}