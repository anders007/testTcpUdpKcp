package main

import (
	"fmt"
	"net"
	"time"
	"io"

	"message"
	"example"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
)

const (
	TCP_TIMEOUT           	= 60 // tcp read timeout
	MESSAGE_HEAD_LEN	= 7
)

var handlerMap = make(map[uint16]HandlerFunc)

func checkError(err error,info string) (res bool) {

	if(err != nil){
		fmt.Println(info+"  " + err.Error())
		return false
	}
	return true
}

func initHandlerFunc() {
	handlerMap[HANDLER_TEST] = HandleTest
}

////////////////////////////////////////////////////////
//
//客户端发送线程
//参数
//      发送连接 conn
//
////////////////////////////////////////////////////////
func chatSend(conn net.Conn){

	var times int32
	for {
		if times >= TEST_TIMES {
			break
		}

		msgcheck := &example.CheckRtt{
			Id:	proto.Int32(times),
			Time:	proto.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
		}

		//data := make([]byte, 40960)
		data, err := proto.Marshal(msgcheck)
		msghead := message.NewHead(uint16(0), int32(len(data)))
		_,err =conn.Write(msghead.Encode())
		_,err = conn.Write(data)
		//fmt.Println(lens)
		if(err != nil){
			fmt.Println(err.Error())
			conn.Close()
			break
		}
		times++

		time.Sleep(30 * time.Millisecond)
	}

}

////////////////////////////////////////////////////////
//
//客户端启动函数
//参数
//      远程ip地址和端口 tcpaddr
//
////////////////////////////////////////////////////////
func StartClient(tcpaddr string){

	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)
	checkError(err,"ResolveTCPAddr")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err,"DialTCP")
	//启动客户端发送线程
	go chatSend(conn)
	go analysis()

	headbuf := make([]byte, MESSAGE_HEAD_LEN)
	//开始客户端轮训
	for{
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err := io.ReadFull(conn, headbuf)
		//fmt.Println("Read msg len %d", lenght)
		if(checkError(err,"Connection")==false){
			conn.Close()
			break
		}
		msghead := message.MessageHead{}
		msghead.Decode(headbuf)

		data := make([]byte, msghead.GetDataLen())
		_, err = io.ReadFull(conn, data)
		if err != nil {
			conn.Close()
			break
		}

		go handlerMap[msghead.GetMsgId()](conn, data)
	}
}

func analysis() {
	time.Sleep(20 * time.Second)
	miss, resp, totalelasp := 0, 0, 0
	for _, elapse := range responsflag {
		if elapse < 0 {
			miss++
		} else {
			totalelasp += elapse
			resp++
		}
	}

	//fmt.Println(responsflag)
	missrate := float32(miss) / float32(TEST_TIMES)
	avgrtt := float32(totalelasp) / float32(resp)
	fmt.Println("missrate = ", missrate, "; average rtt = ", avgrtt)
}

var sendbuff = make([]byte, 1024)
func sendUDPCheckMsg(conn *net.UDPConn){

	var times int32
	for {
		if times >= TEST_TIMES {
			break
		}

		msgcheck := &example.CheckRtt{
			Id:	proto.Int32(times),
			Time:	proto.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
		}

		//data := make([]byte, 40960)
		data, err := proto.Marshal(msgcheck)
		msghead := message.NewHead(uint16(0), int32(len(data)))

		//_,err =conn.Write(msghead.Encode())
		//_,err = conn.Write(data)

		len := copy(sendbuff, msghead.Encode())
		len += copy(sendbuff[len:], data)

		_, writeerr := conn.Write(sendbuff[:len])
		checkError(writeerr, "UDPWrite")

		//fmt.Println("udp write len = ", writelen)

		//fmt.Println(lens)
		if(err != nil){
			fmt.Println(err.Error())
			conn.Close()
			break
		}
		times++

		time.Sleep(300 * time.Millisecond)
	}

}

//缓冲区
const BUFFSIZE = 1024

func startUDPClient(addr string) {
	//udp地址
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	checkError(err, "ResolveUDPAddr")

	//udp连接
	udpConn, err := net.DialUDP("udp4", nil, udpAddr)
	checkError(err, "DialUDP")

	go sendUDPCheckMsg(udpConn)
	go analysis()

	for {
		var udpbuff = make([]byte, BUFFSIZE)
		n, _, err := udpConn.ReadFromUDP(udpbuff)
		checkError(err, "ReadFromUDP")

		if n > 0 {
			msghead := message.MessageHead{}
			msghead.Decode(udpbuff[:MESSAGE_HEAD_LEN])

			go handlerMap[msghead.GetMsgId()](udpConn, udpbuff[MESSAGE_HEAD_LEN:n])
		}
	}
}

func sendKCPCheckMsg(conn *kcp.UDPSession){

	var times int32
	for {
		if times >= TEST_TIMES {
			break
		}

		msgcheck := &example.CheckRtt{
			Id:	proto.Int32(times),
			Time:	proto.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
		}

		//data := make([]byte, 40960)
		data, err := proto.Marshal(msgcheck)
		msghead := message.NewHead(uint16(0), int32(len(data)))

		//_,err =conn.Write(msghead.Encode())
		//_,err = conn.Write(data)

		len := copy(sendbuff, msghead.Encode())
		len += copy(sendbuff[len:], data)
		conn.Write(sendbuff[:len])

		//fmt.Println(lens)
		if(err != nil){
			fmt.Println(err.Error())
			conn.Close()
			break
		}
		times++

		time.Sleep(300 * time.Millisecond)
	}

}

func startKCPClient(addr string) {
	kcpconn, err := kcp.DialWithOptions(addr, nil, 10, 3)
	checkError(err, "kcp.DialWithOptions")

	kcpconn.SetStreamMode(true)
	kcpconn.SetNoDelay(1, 20, 2, 1)

	go sendKCPCheckMsg(kcpconn)
	go analysis()

	for {
		var udpbuff = make([]byte, BUFFSIZE)
		n, err := kcpconn.Read(udpbuff)
		checkError(err, "KCP Read")

		if n > 0 {
			msghead := message.MessageHead{}
			msghead.Decode(udpbuff[:MESSAGE_HEAD_LEN])

			go handlerMap[msghead.GetMsgId()](kcpconn, udpbuff[MESSAGE_HEAD_LEN:n])
		}
	}
}

func main() {
	for i, _ := range responsflag {
		responsflag[i] = -1
	}

	initHandlerFunc()
	//StartClient("54.238.231.182:3000")
	//StartClient("127.0.0.1:3000")

	//startUDPClient("127.0.0.1:3001")
	//startUDPClient("45.32.255.89:3001")

	//startKCPClient("127.0.0.1:3002")
	startKCPClient("45.32.255.89:3002")
}
