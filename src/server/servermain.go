package main

import (
	"fmt"
	"net"
	"io"

	"message"
	"time"
	kcp "github.com/xtaci/kcp-go"
	"log"
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
//服务器端接收数据线程
//参数：
//      数据连接 conn
//      通讯通道 messages
//
////////////////////////////////////////////////////////
func Handler(conn net.Conn,messages chan string){

	headbuf := make([]byte, MESSAGE_HEAD_LEN)

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

		//fmt.Println("Rec[",conn.RemoteAddr().String(),"] Say :" ,string(buf[0:lenght]))
		go handlerMap[msghead.GetMsgId()](conn, data)

	}

}

////////////////////////////////////////////////////////
//
//服务器发送数据的线程
//
//参数
//      连接字典 conns
//      数据通道 messages
//
////////////////////////////////////////////////////////
func echoHandler(conns *map[string]net.Conn,messages chan string){


	for{
		msg:= <- messages
		fmt.Println(msg)

		for key,value := range *conns {

			fmt.Println("connection is connected from ...",key)
			_,err :=value.Write([]byte(msg))
			if(err != nil){
				fmt.Println(err.Error())
				delete(*conns,key)
			}

		}

		// 进行解码
		//newTest := &example.Test{}
		//err := proto.Unmarshal([]byte(msg), newTest)
		//if err != nil {
		//	log.Fatal("unmarshaling error: ", err)
		//}
		//fmt.Println(newTest)
	}

}


////////////////////////////////////////////////////////
//
//启动服务器
//参数
//  端口 port
//
////////////////////////////////////////////////////////
func StartServer(port string) {

	service := ":" + port //strconv.Itoa(port);
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	conns := make(map[string]net.Conn)
	messages := make(chan string, 10)
	//启动服务器广播线程
	//go echoHandler(&conns,messages)

	go func() {
		for {
			fmt.Println("Listening ...")
			conn, err := l.Accept()
			checkError(err, "Accept")
			fmt.Println("Accepting ...")
			conns[conn.RemoteAddr().String()] = conn
			//启动一个新线程
			go Handler(conn, messages)
		}
	}()

}

//缓冲区大小
const BUFFSIZE = 1024

//消息处理
func HandleUDPMsg(udpListener *net.UDPConn) {
	var udpbuff = make([]byte, BUFFSIZE)

	for {
		n, addr, err := udpListener.ReadFromUDP(udpbuff)
		checkError(err, "ReadFromUDP")

		//fmt.Println("rev udp msg, len = ", n)
		if n > 0 {
			msghead := message.MessageHead{}
			msghead.Decode(udpbuff[:MESSAGE_HEAD_LEN])
			udpListener.WriteToUDP(udpbuff[:n], addr)
		}
	}
}

func StartUDPServer(port string) {
	//监听地址
	udpAddr, err := net.ResolveUDPAddr("udp4", ":" + port)
	checkError(err, "ResolveUDPAddr")
	//监听连接
	udpListener, err := net.ListenUDP("udp4", udpAddr)
	checkError(err, "LestenUDP")

	//消息处理
	go HandleUDPMsg(udpListener)

}

func handleKCPMsg(session *kcp.UDPSession) {

	var kcpbuff = make([]byte, BUFFSIZE)

	for {
		n, err := session.Read(kcpbuff)
		checkError(err, "KCPRead")

		if n > 0 {
			msghead := message.MessageHead{}
			msghead.Decode(kcpbuff[:MESSAGE_HEAD_LEN])
			session.Write(kcpbuff[:n])
		}
	}
}

func StartKCPServer(port string) {
	lis, err := kcp.ListenWithOptions(":" + port, nil, 10, 3)
	checkError(err, "")

	for {
		if conn, err := lis.AcceptKCP(); err == nil {
			log.Println("remote address:", conn.RemoteAddr())
			conn.SetStreamMode(true)
			conn.SetNoDelay(1, 20, 2, 1)

			go handleKCPMsg(conn)
		} else {
			log.Printf("%+v", err)
		}
	}
}

type Data struct {
	Name string
}

func main() {
	initHandlerFunc()
	StartServer("3000")
	StartUDPServer("3001")
	StartKCPServer("3002")

	for {
		time.Sleep(3 * time.Millisecond)
	}
}