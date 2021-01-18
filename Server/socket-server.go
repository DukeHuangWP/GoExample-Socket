package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const buffSize = 1024

func main() {

	//建立socket，監聽埠
	listenPort := ":51680"
	netListen, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Printf("listen tcp '%v' fail > '%v'\n", listenPort, err)
		os.Exit(1)
	}
	defer netListen.Close()

	log.Println("Waiting for clients")
	var clientList = make(map[*net.Conn]struct{})
	go func() { //監聽採用無限制迴圈
		for {
			conn, err := netListen.Accept()
			if err != nil {
				continue
			}

			clientIP := conn.RemoteAddr().String()
			log.Println(clientIP, " tcp connect success")
			buffer := make([]byte, buffSize) //超過緩存範圍訊息將被斷行打印
			go func() {                      //每個客戶端開goroutine
				for {
					clientList[&conn] = struct{}{}   //struct不占記憶體空間
					endChr, err := conn.Read(buffer) //阻塞式，當無client連接時會進行等待
					if err != nil {
						log.Println(clientIP, " connection error: ", err)
						delete(clientList, &conn)
						break
					}

					msg := string(buffer[:endChr])
					switch msg {
					case "bye": //關閉客戶端
						conn.Write([]byte("Bye bye!"))
						delete(clientList, &conn) //移除廣播清單
						return                    //關閉goroutine
					default: //接收客戶端訊息
						log.Println(clientIP, "receive data string:", msg)
					}
				}
				return
			}()

		}
	}()

	var sendMsg string
	for { //廣播給所有清單內的客戶端
		fmt.Scanln(&sendMsg)
		for client, _ := range clientList {
			conn := *client
			_, err := conn.Write([]byte(sendMsg))
			if err != nil {
				log.Println("廣播 fail : ", conn.RemoteAddr())
				continue //監聽程序會負責移除斷線或異常的客戶端
			}
		}
		log.Printf("廣播 : %s", sendMsg)
		time.Sleep(1 * time.Second)
	}
}
