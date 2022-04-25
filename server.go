package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/hika019/File_transfer/lib"
)

func main() {
	protocol := "tcp"
	port := ":55555"

	tcpAddr, err := net.ResolveTCPAddr(protocol, port)
	lib.CheckErrorExit(err)

	listner, err := net.ListenTCP(protocol, tcpAddr)
	lib.CheckErrorExit(err)

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
		if lib.CheckError(err) {
			return
		}
	}

	messageBuf := make([]byte, lib.SocketByte)

	messageLen, err := conn.Read(messageBuf)
	//EOFエラー回避
	if messageLen == 0 {
		return
	}

	if lib.CheckError(err) {
		return
	}
	fmt.Println(messageBuf)

	fileName, hash := lib.ByteToFileName(messageBuf[:messageLen])

	fileName = senderIP + "/" + fileName
	//fmt.Println("filename: ", fileName, hash)

	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if lib.CheckError(err) {
		return
	}
	defer fp.Close()

	receiveCount := 0

	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		messageBuf = make([]byte, lib.SocketByte)
		messageLen, err = conn.Read(messageBuf)

		//EOFエラー回避
		if messageLen == 0 {
			fmt.Println("download file")
			break
		}

		lib.CheckErrorExit(err)

		_, err = fp.Write(messageBuf[:messageLen])
		lib.CheckError(err)
		receiveCount++
	}

	if reflect.DeepEqual(hash, lib.CreateSHA256(fileName)) {
		fmt.Println("Consistency: Yes")
		conn.Write([]byte{0})
	} else {
		fmt.Println("Consistency: No")
		conn.Write([]byte{1})
	}

}

func Exists(path string) bool {
	f, err := os.Stat(path)
	return !(os.IsNotExist(err) || !f.IsDir())
}
