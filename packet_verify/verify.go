package packet_verify

type CounterProtocol struct {
	tokenRecv int32
	tokenSend int32
	sendRetry int32
	//recvRetry int32
}

func (proto *CounterProtocol) Init(v int32) {
	proto.tokenSend = v
	proto.tokenRecv = v
}

// TODO: implement hmac
func (proto *CounterProtocol) hmac(counter int32, data []byte) int32 {
	sign := int32(0)
	for _, b := range data {
		sign += counter * int32(b)
	}
	return sign
}

func (proto *CounterProtocol) SignPkt(pkt []byte) int32 {
	proto.tokenSend += 1
	return proto.hmac(proto.tokenSend, pkt)
}

func (proto *CounterProtocol) VerifyPkt(sign int32, pkt []byte) bool {
	proto.tokenRecv += 1
	curSign := proto.hmac(proto.tokenRecv, pkt)
	return curSign == sign
}

func (proto *CounterProtocol) SendCounterUpdate(pkt []byte) int32 {
	// +100
	ts := proto.tokenSend
	if proto.sendRetry < 5 {
		ts++
		proto.sendRetry++
	} else {
		ts = (ts / 100 + 1) * 100
	}
	proto.tokenSend = ts
	return proto.hmac(proto.tokenSend, pkt)
}

func (proto *CounterProtocol) SendCounterOk() {
	proto.sendRetry = 0
}

//
func (proto *CounterProtocol) TryRecvCounter(sign int32, pkt []byte) bool {
	retries := 5
	for retries > 0 {
		retries--
		proto.tokenRecv++
		if sign == proto.hmac(proto.tokenRecv, pkt) {
			return true
		}
	}
	ts := proto.tokenRecv
	ts = (ts / 100 + 1) * 100
	proto.tokenRecv = ts
	return sign == proto.hmac(proto.tokenRecv, pkt)
}

/*

https://www.zhihu.com/question/29383090/answer/1649199245
1. diffie hellman
*/
