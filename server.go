package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
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
		fmt.Println(messageLen)
		CheckError(err)

		//clientが先に切断した際err=EOFになる
		//err.Error()=EOFの場合messageLen=0  (エラー落ち回避)
		if messageLen == 0 {
			break
		}
		CheckError(err)

		var dataLen uint16 = 0

		dataLen = ByteToInt(messageBuf)
		fmt.Println(receiveCount)
		//fmt.Println(messageBuf)

		//そこそこ大きい画像を送るときに謎データがくっつくのでスキップ
		if uint16(SocketDataSize) < dataLen {
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
		fmt.Println(messageBuf)
		fmt.Printf("%v byte\n", dataLen)
		fmt.Fprintf(fp, "%v", string(messageBuf))
		fmt.Println("######")
		//ファイルに書き込み
		receiveCount++
	}

	hash := CreateSHA256(fileName)
	fmt.Println(hash)

	//FIXME:ハッシュの確認の部分を実行するとファイルが送れない
	/*
		conn.Write(hash)

		fmt.Println("DownloadFileStatus")
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		status := []byte{1}
		_, err = conn.Read(status)
		if !reflect.DeepEqual(status, []byte{0}) {
			fmt.Println("NOT Complete File Transefer!!")
		} else {
			fmt.Println("Complete File Transefer")
		}
	*/
}

func DownloadFileStatus(conn net.Conn) bool {
	fmt.Println("DownloadFileStatus")

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	status := make([]byte, 1)
	_, err := conn.Read(status)
	CheckError(err)
	return reflect.DeepEqual(status, []byte{2})
}
