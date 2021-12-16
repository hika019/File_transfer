package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	protocol := "tcp"
	port := ":55555"

	tcpAddr, err := net.ResolveTCPAddr(protocol, port)
	CheckError(err)

	listner, err := net.ListenTCP(protocol, tcpAddr)
	CheckError(err)

	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)

	}
}

func handleClient(conn net.Conn) {

	addr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		return
	}

	fmt.Println(addr.IP.String())

	defer conn.Close()
	messageBuf := make([]byte, SocketSize)
	//fmt.Println(messageBuf)

	_, err := conn.Read(messageBuf)
	fileName := string(messageBuf)
	fileName = fileName[:(strings.Index(fileName, ":"))]
	fmt.Println("filename: ", fileName)

	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err)

	defer fp.Close()

	receiveCount := 0
	conn.SetReadDeadline(time.Now().Add(50 * time.Second))
	for {
		messageBuf := make([]byte, SocketSize)
		messageLen, err := conn.Read(messageBuf)

		//clientが先に切断した際err=EOFになる
		//err=EOFの場合messageLen=0で終了(エラー落ち回避)
		if messageLen == 0 {
			break
		}
		CheckError(err)

		var dataLen uint16 = 0

		fmt.Println(messageBuf)

		dataLen = ByteToInt(messageBuf)
		//fmt.Println(dataLen)

		if uint16(SocketDataSize) < dataLen {
			//そこそこ大きい画像を送るときに最初に謎データがくっつくでスキップ
			fmt.Println("data_size が無効")
			dataLen = 0
			continue
		}

		if dataLen == 0 {
			fmt.Println("Downloaded file data")
			fmt.Println(receiveCount)
			break
		}

		fmt.Println(receiveCount)

		fmt.Printf("%v byte\n", dataLen)
		fmt.Fprintf(fp, "%v", string(messageBuf[:dataLen]))
		fmt.Println("######")
		//ファイルに書き込み
		receiveCount++
	}
	/*
		hash := CreateFileHash(fileName)
		fmt.Println(hash)
		fmt.Println(len(hash))
		//conn.Write(hash)
	*/
}
