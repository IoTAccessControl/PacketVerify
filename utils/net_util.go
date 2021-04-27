package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)
import "encoding/gob"

func Int2Bytes(v int32) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(v))
	return bs
}

func Bytes2Int(bys []byte) int {
	return int(binary.BigEndian.Uint32(bys))
}

func Struct2Bytes(st interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(st)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

/*
https://medium.com/hackernoon/today-i-learned-pass-by-reference-on-interface-parameter-in-golang-35ee8d8a848e
 */
func StructFromBytes(bys []byte, st interface{}) error  {
	var buf bytes.Buffer
	n, err := buf.Read(bys)
	if len(bys) != n || err != nil {
		return err
	}
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(st)
	return err
}

func SendPkt(w *bufio.Writer, cmd uint8,  buf []byte) {
	//w := bufio.NewWriter(conn)
	pkSize := int32(len(buf))
	w.Write(Int2Bytes(pkSize))
	w.WriteByte(byte(cmd))
	w.Write(buf)
	w.Flush()
}

func RecvPkt(r *bufio.Reader) (uint8, []byte) {
	pkt := make([]byte, 0, 4)
	//pkt := make([]byte, 0, 4*1024)
	io.ReadFull(r, pkt[:4])
	pktLen := Bytes2Int(pkt[:4])
	io.ReadFull(r, pkt[:1])
	cmd := pkt[:1][0]
	pkt = make([]byte, 0, pktLen)
	io.ReadFull(r, pkt[:pktLen])
	//fmt.Printf("RecvPkt: cmd=%d len=%d pkt=%d\n", cmd, pktLen, len(pkt[:pktLen]))
	return cmd, pkt[:pktLen]
}
