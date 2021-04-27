package proxy

import (
	"encoding/binary"
)


const (
	CMD_DH = 0x10
	CMD_DH1 = 0x11
	CMD_DH2 = 0x12
	CMD_DH3 = 0x13

	CMD_PKT = 0x20
	CMD_PKT_SYNC = 0x21
	CMD_PKT_SYNC_ACK = 0x22
)

type MarshalAble interface {
	Marshal() []byte
	Unmarshal(data []byte)
}

type PoorDHMsg struct {
	S int32
	G int32
}

func (msg *PoorDHMsg) Marshal() []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint32(bs[:4], uint32(msg.S))
	binary.BigEndian.PutUint32(bs[4:], uint32(msg.G))
	//fmt.Printf("%v %d %d %v\n", msg, uint32(msg.S), msg.G, bs)
	return bs
}

func (msg *PoorDHMsg) Unmarshal(data []byte) {
	msg.S = int32(binary.BigEndian.Uint32(data[:4]))
	msg.G = int32(binary.BigEndian.Uint32(data[4:]))
}


func MarshalMessage(st interface{}) {
	switch st.(type) {
	case *PoorDHMsg:
		msg := st.(*PoorDHMsg)
		msg.S = 3
		msg.G = 5
		break
	}
}

