package proxy

import (
	"bufio"
	"io"
	"net"
	"packet-verify/netstream"
	"packet-verify/packet_verify"
	"packet-verify/utils"
	"fmt"
)

func StartRemoteProxy(listenAddr string, remoteAddr string) {
	listen, err := net.Listen("tcp", listenAddr)
	utils.LogInfo("%s Start to Listen: %s\n", netstream.PeerTag, listenAddr)
	if listen == nil {
		utils.Fatal("Listen port failed: %v\n", err)
	}

	for {
		local, err := listen.Accept()
	//	utils.LogInfo("%s accept connect: %v\n", netstream.PeerTag, local.LocalAddr())
	//	utils.LogInfo("%s accept connect: %v <-> %v\n", netstream.PeerTag, local.LocalAddr(), local.RemoteAddr())
		if local == nil {
			utils.Fatal("accept failed: %v", err)
		}

		remote, err := net.Dial("tcp", remoteAddr)
		utils.LogInfo("%s Connect to: %v\n", netstream.PeerTag, remoteAddr)
		if remote == nil {
			utils.Fatal("%s remote dial failed: %v\n", netstream.PeerTag, err)
			return
		}

		var COUNTER packet_verify.CounterProtocol
		localReader := bufio.NewReader(local)
		remoteReader := bufio.NewReader(remote)
		localWriter := bufio.NewWriter(local)
		remoteWriter := bufio.NewWriter(remote)

		peer := PeerConn{local, remote, localReader, remoteReader,
			localWriter, remoteWriter}

		if netstream.PeerTag == netstream.PeerRemoteProxy {
			go func() {
				counter := remoteTokenNegotiation(peer)
				COUNTER.Init(counter)
				utils.LogInfo("%s Start forwarding.\n", netstream.PeerTag)
				fmt.Printf("Counter: %d\n", counter)
				remoteForward(peer, COUNTER)
			}()
		}
	}
}


func remoteTokenNegotiation(peer PeerConn) int32 {
	// receive counter
	cmd, bys := peer.RecvPktLocal()
	if cmd == CMD_DH {
		var dh PoorDHMsg
		dh.Unmarshal(bys)
		return dh.S
	}

	return 0
}

func remoteForward(peer PeerConn, counter packet_verify.CounterProtocol) {
	done := make(chan int, 2)
	go func() {
		// local proxy -> [remote proxy] -> server
		// ingress forward
		//buf := make([]byte, 4 * 1024)
		for {
			// receive raw data
			cmd, pkt := peer.RecvPktLocal()
			if cmd == 0 && len(pkt) == 0 {
				break
			}
			if (cmd & CMD_PKT) == CMD_PKT {
				sign := int32(utils.Bytes2Int(pkt[:4]))
				payload := pkt[4:]

				fmt.Printf("%s recv from local: %d %d %d\n", netstream.PeerTag, len(pkt), len(payload), sign)
				//utils.LogInfo("%s recv from local: %d %d %d\n", netstream.PeerTag, len(pkt), len(payload), sign)
				if cmd == CMD_PKT {
					if counter.TryRecvCounter(sign, payload) {
						peer.ForwardToRemote(payload)
						utils.Statistic("[%d] %s forwarding to server.\n", len(payload), netstream.PeerTag)
						//utils.Statistic("[%d] %s forward to server: %d sign: %d\n", len(payload), netstream.PeerTag, len(payload), sign)
					} else {
						// notify packet loss
						sign = counter.SendCounterUpdate(payload)
						bys1 := utils.Int2Bytes(sign)
						data := append(bys1, payload...)
						peer.SendPktRemote(CMD_PKT_SYNC, data)
					}
				} else if cmd == CMD_PKT_SYNC_ACK {
					counter.SendCounterOk()
				} else if cmd == CMD_PKT_SYNC {
					if counter.TryRecvCounter(sign, payload) {
						peer.ForwardToRemote(payload)
					} else {
						utils.Fatal("BUG: Counter sync failed.")
					}
				}
			}
		}

		done<-1
	}()

	go func() {
		// server -> [remote proxy] -> local proxy
		buf := make([]byte, 4 * 1024)
		for {
			n, err := peer.remoteReader.Read(buf)
			pkt := buf[:n]
			//pkt, err := io.ReadAll(peer.localReader)
			//fmt.Printf("%s read all------> %d %v\n", netstream.PeerTag, len(pkt), pkt)
			if n == 0 || err == io.EOF {
				// finish
				break
			}

			if !RemotePktLoss() {
				sign := counter.SignPkt(pkt)
				bys1 := utils.Int2Bytes(sign)
				data := append(bys1, pkt...)
				// forward
				peer.SendPktLocal(CMD_PKT, data)
				//utils.Statistic("%s forwarding to local: %d sign: %d\n", netstream.PeerTag, len(pkt), sign)
				utils.LogInfo("%s forwarding to local: %d\n", netstream.PeerTag, len(pkt))
			}
		}

		done<-2
	}()
	defer peer.Close()
	<-done
	<-done
	utils.LogInfo("%s forwarding finish!", netstream.PeerTag)
}