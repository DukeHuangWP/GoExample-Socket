package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const buffSize = 1024

func main(){

	socketAddr := "localhost:51680"
	var err error
	var socketServer *net.Conn //socket server公用控制器
	go func(addr string) {
		var conn net.Conn
		buff := make([]byte, buffSize)
		for {
			if socketServer == nil {
				conn, err = net.Dial("tcp", addr) //與server進行連線
				if err != nil {
					log.Printf("Reconnect error:%v\n", err)
					socketServer = nil
					time.Sleep(1 * time.Second)
					continue
				} else {
					socketServer = &conn
					defer conn.Close()
					log.Printf("Suecess to connect socket > '%v'\n", addr)
				}
			}

			endChr, err := conn.Read(buff) //接收server訊息
			if err != nil {
				log.Println("Receive fail!")
				socketServer = nil
				continue
			}
			log.Printf("Received from server : %s", buff[:endChr])
		}
	}(socketAddr)

	var sendMsg string
	for {
		fmt.Scanln(&sendMsg)
		conn := *socketServer //指向公用socket server
		_,err :=conn.Write([]byte(sendMsg))
		if err !=nil{
			log.Println("Sent fail !")
			socketServer = nil
			continue //若發送失敗接收者負責重連
		}
		log.Printf("Send: %s", sendMsg)
		time.Sleep(1 * time.Second)
	}
}
