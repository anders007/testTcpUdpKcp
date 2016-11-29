package main

import (
	"net"
	//"example"
	//"fmt"

	//"github.com/golang/protobuf/proto"
	"message"
)

const (
	HANDLER_TEST = 0
)

type HandlerFunc func(conn net.Conn, data []byte)

func HandleTest(conn net.Conn, data []byte) {
	// 进行解码
	//newTest := &example.Test{}
	//err := proto.Unmarshal(data, newTest)
	//if err != nil {
	//	fmt.Println("unmarshaling error: ", err)
	//	return
	//}
	//fmt.Println(newTest)

	msghead := message.NewHead(uint16(0), int32(len(data)))
	conn.Write(msghead.Encode())
	conn.Write(data)
}