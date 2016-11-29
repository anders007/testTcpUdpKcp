package main

import (
	"net"
	"log"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"example"
)

const (
	HANDLER_TEST = 0
)

const (
	TEST_TIMES = 100
)

type HandlerFunc func(conn net.Conn, data []byte)

var responsflag = make([]int, TEST_TIMES)

func HandleTest(conn net.Conn, data []byte) {
	// 进行解码
	checkMsg := &example.CheckRtt{}
	err := proto.Unmarshal(data, checkMsg)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	fmt.Println(checkMsg)
	current := time.Now().UnixNano() / int64(time.Millisecond)
	if checkMsg.GetId()>=0 && checkMsg.GetId()<TEST_TIMES {
		responsflag[checkMsg.GetId()] = int(current - checkMsg.GetTime())
	}
}