package policy

import "packet-verify/packet_verify"

type IPCAddr struct {
	localApp string
	localPort int
	localHost string
	remotePort int
	remoteHost string
	remoteApp string
}

type IPCLink struct {
	cpRecv packet_verify.CounterProtocol
	cpSend packet_verify.CounterProtocol
	IPCAddr
}


