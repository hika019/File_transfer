package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
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

	messageBuf := make([]byte, SocketByte)
	//fmt.Println(messageBuf)

	_, err := conn.Read(messageBuf)
	CheckError(err)
	fmt.Println(messageBuf)
	fileNameLen := ByteToInt(messageBuf)

	if fileNameLen == uint16(0) {
		fmt.Println("ファイル名が不正")
		os.Exit(1)
	}

	fileName := string(messageBuf[:fileNameLen])
	fmt.Println("filename: ", fileName)

	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err)
	defer fp.Close()

	receiveCount := 0

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		messageBuf = make([]byte, SocketByte)
		_, err := conn.Read(messageBuf)

		//clientが先に切断した際err=EOFになる
		//err.Error()=EOFの場合messageLen=0  (エラー落ち回避)
		dataLen := ByteToInt(messageBuf)

		if dataLen == 0 {
			fmt.Println("download file")
			break
		}
		CheckError(err)

		//表示しないとデータが変になる
		fmt.Println(messageBuf)

		//時々変なデータがあるから回避
		if messageBuf[DataSizeBytePos0] != uint8(1) {
			fmt.Println("変なデータのため切断")
			break
		}

		if uint16(SocketDataByte) < dataLen {
			fmt.Println(messageBuf)
			fmt.Println("data_size が無効")
			dataLen = 0
			continue
		}

		fp.Write(messageBuf[:dataLen])

		//ファイルに書き込み
		receiveCount++
	}
	fmt.Println(receiveCount)
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
