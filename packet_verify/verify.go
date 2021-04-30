package packet_verify

import (
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"packet-verify/utils"
)

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
func (proto *CounterProtocol) HmacSha256(counter int32, data []byte) int32 {
	sign := int32(0)
//	message := string(data[:])
	//secret := "0933e54e76b24731a2d84b6b463ec04c"
	//key := []byte(secret)
	key := utils.Int2Bytes(counter)
	h := hmac.New(sha256.New, key)
	h.Write(data)
	sha := hex.EncodeToString(h.Sum(nil))
	result := base64.StdEncoding.EncodeToString([]byte(sha))
	sign = int32(utils.Bytes2Int([]byte(result)))
	return sign
}
// func (proto *CounterProtocol) HmacSha256(counter int32, data []byte) int32 {
// 	sign := int32(0)
// 	for _, b := range data {
// 		sign += counter * int32(b)
// 	}
// 	return sign
// }


func (proto *CounterProtocol) SignPkt(pkt []byte) int32 {
	proto.tokenSend += 1
	fmt.Printf("current counter: %d\n", proto.tokenSend)
	return proto.HmacSha256(proto.tokenSend, pkt)
}

func (proto *CounterProtocol) VerifyPkt(sign int32, pkt []byte) bool {
	proto.tokenRecv += 1
	curSign := proto.HmacSha256(proto.tokenRecv, pkt)
	return curSign == sign
}

func (proto *CounterProtocol) SendCounterUpdate(pkt []byte) int32 {
	// +100
	ts := proto.tokenSend
	if proto.sendRetry < 5 { //这里为啥要重复五次
		ts++
		proto.sendRetry++
	} else {
		ts = (ts / 100 + 1) * 100
	}
	proto.tokenSend = ts
	return proto.HmacSha256(proto.tokenSend, pkt)
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
		if sign == proto.HmacSha256(proto.tokenRecv, pkt) {
			return true
		}
	}
	ts := proto.tokenRecv
	fmt.Printf("Old counter: %d\n", ts)
	ts = (ts / 100 + 1) * 100
	fmt.Printf("New counter: %d\n", ts)
	proto.tokenRecv = ts
	return sign == proto.HmacSha256(proto.tokenRecv, pkt)
}

/*

https://www.zhihu.com/question/29383090/answer/1649199245
1. diffie hellman
*/
