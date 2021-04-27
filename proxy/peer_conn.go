package proxy

import (
	"bufio"
	"math/rand"
	"net"
	"packet-verify/utils"
)

type PeerConn struct {
	localConn net.Conn
	remoteConn net.Conn
	localReader *bufio.Reader
	remoteReader *bufio.Reader
	localWriter *bufio.Writer
	remoteWriter *bufio.Writer
}

func (peer *PeerConn) SendPktLocal(cmd uint8, data []byte)  {
	utils.SendPkt(peer.localWriter, cmd, data)
}

func (peer *PeerConn) RecvPktLocal() (uint8, []byte)  {
	return utils.RecvPkt(peer.localReader)
}

func (peer *PeerConn) SendPktRemote(cmd uint8, data []byte)  {
	utils.SendPkt(peer.remoteWriter, cmd, data)
}

func (peer *PeerConn) RecvPktRemote() (uint8, []byte)  {
	return utils.RecvPkt(peer.remoteReader)
}

func (peer *PeerConn) ForwardToRemote(data []byte) {
	peer.remoteWriter.Write(data)
	peer.remoteWriter.Flush()
}

func (peer *PeerConn) ForwardToLocal(data []byte) {
	peer.localWriter.Write(data)
	peer.localWriter.Flush()
}



func (peer *PeerConn) Close() {
	peer.localWriter.Flush()
	peer.localConn.Close()
	peer.remoteWriter.Flush()
	peer.remoteConn.Close()
}

func LocalPktLoss() bool {
	lossRate := 0
	return (rand.Int() % 100) < lossRate
}

func RemotePktLoss() bool {
	lossRate := 0
	return (rand.Int() % 100) < lossRate
}

