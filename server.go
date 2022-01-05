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
	CheckErrorExit(err)

	listner, err := net.ListenTCP(protocol, tcpAddr)
	CheckErrorExit(err)

	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn)

	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	addr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		return
	}

	senderIP := addr.IP.String()

	fmt.Println(senderIP)

	//dirの作成
	if !Exists(senderIP) {
		err := os.Mkdir(senderIP, 0777)
		if CheckError(err) {
			return
		}
	}

	messageBuf := make([]byte, SocketByte)

	messageLen, err := conn.Read(messageBuf)
	//EOFエラー回避
	if messageLen == 0 {
		return
	}

	if CheckError(err) {
		return
	}
	fmt.Println(messageBuf)
	fileNameLen := ByteToInt(messageBuf)

	if fileNameLen == uint16(0) {
		fmt.Println("ファイル名が不正")
		return
	}

	fileName := string(messageBuf[:fileNameLen])
	fmt.Println("filename: ", fileName)

	fp, err := os.OpenFile(senderIP+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if CheckError(err) {
		return
	}
	defer fp.Close()

	receiveCount := 0

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		messageBuf = make([]byte, SocketByte)
		_, err := conn.Read(messageBuf)

		//clientが先に切断した際err=EOFになる
		//err.Error()=EOFの場合messageLen=0  (エラー落ち回避)
		dataLen := ByteToInt(messageBuf)

		if messageBuf[DataSizeBytePos0] != uint8(1) {
			fmt.Println("download file")
			break
		}
		CheckErrorExit(err)

		//ファイルに書き込み
		fp.Write(messageBuf[:dataLen])
		receiveCount++
	}
	fmt.Println(receiveCount)
	hash := CreateSHA256(fileName)
	fmt.Println(hash)

	//hashを送る
	conn.Write(hash)
	fmt.Println("Send File hash")

	//ステータスのダウンロード
	if DownloadFileStatus(conn) {
		fmt.Println("Complete File Transefer")
	} else {
		fmt.Println("NOT Complete File Transefer!!")
	}

}

func DownloadFileStatus(conn net.Conn) bool {
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	status := []byte{1}
	_, err := conn.Read(status)
	if CheckError(err) == false {
		return false
	}
	return reflect.DeepEqual(status, []byte{0})
}

func Exists(path string) bool {
	f, err := os.Stat(path)
	return !(os.IsNotExist(err) || !f.IsDir())
}
