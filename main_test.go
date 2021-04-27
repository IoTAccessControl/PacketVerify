package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"packet-verify/packet_verify"
	"testing"
)

/*

 */
type User struct {
	Name string
}

func Test1(t *testing.T) {
	u := &User{Name: "Leto"}
	println(u.Name)
	Modify(u)
	println(u.Name)
}

func Modify(u *User) {
	u = &User{Name: "Paul"}
}

/*
Test serialize
 */
func TestMarshall(t *testing.T) {
	dh := packet_verify.PoorDH{}
	dh.G = 122
	dh.S = 211

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&dh)
	println(err)
	bys := buf.Bytes()

	var buf2 bytes.Buffer
	dh2 := packet_verify.PoorDH{}

	n, err := buf2.Read(bys)
	fmt.Println("read bys", n, err)
	dec := gob.NewDecoder(&buf2)
	err = dec.Decode(&dh2)

	//fmt.Println(bys)
	fmt.Println(dh, dh2)
	//testing.Ass
	t.Log("Logging", dh2.S)
}