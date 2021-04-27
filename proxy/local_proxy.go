package proxy

import (
	"bufio"
	"io"
	"net"
	"packet-verify/netstream"
	"packet-verify/packet_verify"
	"packet-verify/utils"
)

func StartLocalProxy(listenAddr string, remoteAddr string) {
	listen, err := net.Listen("tcp", listenAddr)
	utils.LogInfo("%s Start to Listen: %s\n", netstream.PeerTag, listenAddr)
	if listen == nil {
		utils.Fatal("Listen port failed: %v", err)
	}

	for {
		local, err := listen.Accept()
		utils.LogInfo("%s accept connect: %v\n", netstream.PeerTag, local)
		utils.LogInfo("%s: %v <-> %v\n", netstream.PeerTag, local.LocalAddr(), local.RemoteAddr())
		if local == nil {
			utils.Fatal("accept failed: %v", err)
		}

		var COUNTER packet_verify.CounterProtocol

		//var ipc policy.IPCLink
		remote, err := net.Dial("tcp", remoteAddr)
		utils.LogInfo("%s Connect to: %v\n", netstream.PeerTag, remoteAddr)
		if remote == nil {
			utils.Fatal("%s remote dial failed: %v\n", netstream.PeerTag, err)
			return
		}

		localReader := bufio.NewReader(local)
		remoteReader := bufio.NewReader(remote)
		localWriter := bufio.NewWriter(local)
		remoteWriter := bufio.NewWriter(remote)

		peer := PeerConn{local, remote, localReader, remoteReader,
			localWriter, remoteWriter}

		// init counter
		if netstream.PeerTag == netstream.PeerLocalProxy {
			go func() {
				counter := localTokenNegotiation(peer)
				COUNTER.Init(counter)
				utils.LogInfo("%s Start forwarding. Counter: %d \n", netstream.PeerTag, counter)
				localForward(peer, COUNTER)
			}()

		}
	}
}

func localTokenNegotiation(peer PeerConn) int32  {
	var v int32 = 52
	dh := PoorDHMsg{}
	dh.S = v
	// TODO diffie hellman
	peer.SendPktRemote(CMD_DH, dh.Marshal())

	//cmd, bys := peer.RecvPktRemote()

	//if cmd == CMD_DH {
	//	//var dh PoorDHMsg
	//	dh.Unmarshal(bys)
	//	// counter 2
	//}

	// send counter

	//var rdh PoorDHMsg
	//rdh.Unmarshal(bys)

	//utils.LogInfo("Start localTokenNegotiation=========: %v %v %v\n", bys, dh, rdh)

	//w := bufio.NewWriter(remote)
	//
	//w.Write(utils.Int2Bytes(pkSize))
	////w.Write(utils.Int2Bytes(cmd))
	//w.WriteByte(byte(cmd))
	//w.Write(bys)
	//
	//w.Flush()
	//remote.W

	return v
}

func localForward(peer PeerConn, counter packet_verify.CounterProtocol) {
	done := make(chan int, 2)
	go func() {
		// // client -> [local proxy] -> remote proxy
		// outgress
		buf := make([]byte, 4 * 1024)
		for {
			// receive raw data
			//pkt, err := io.ReadAll(peer.localReader)
			n, err := peer.localReader.Read(buf)
			pkt := buf[:n]
			//fmt.Printf("%s read all------> %v\n", netstream.PeerTag, pkt)
			if n == 0 || err == io.EOF {
				// finish
				//done<-1
				break
			}

			if !LocalPktLoss() {
				sign := counter.SignPkt(pkt)
				bys1 := utils.Int2Bytes(sign)
				data := append(bys1, pkt...)
				// forward
				peer.SendPktRemote(CMD_PKT, data)
				utils.Statistic("[%d] %s forward to server: %d sign: %d\n", len(pkt), netstream.PeerTag, len(data), sign)
			}

		}

		done <- 1
	}()

	go func() {
		// remote proxy -> [local proxy] -> client
		// ingress from remote proxy
		for {
			cmd, pkt := peer.RecvPktRemote()
			if cmd == 0 && len(pkt) == 0 {
				break
			}
			if (cmd & CMD_PKT) == CMD_PKT {
				sign := int32(utils.Bytes2Int(pkt[:4]))
				payload := pkt[4:]
				if cmd == CMD_PKT {
					if counter.TryRecvCounter(sign, payload) {
						peer.ForwardToLocal(payload)
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
						peer.ForwardToLocal(payload)
					} else {
						utils.Fatal("BUG: Counter sync failed.")
					}
				}
			}
		}
		done <- 2
	}()

	defer peer.Close()
	<-done
	<-done
	utils.LogInfo("%s forwarding finish!", netstream.PeerTag)
}

