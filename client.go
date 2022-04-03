package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s message", os.Args[0])
		os.Exit(1)
	}

	protocol := "tcp"
	//serverIP := "192.168.11.50"
	serverIP := os.Args[1]
	serverPort := "55555"
	myIP := MyAddr()
	myPort := 55556

	fileName := os.Args[2]

	tcpAddr, err := net.ResolveTCPAddr(protocol, serverIP+":"+serverPort)
	CheckErrorExit(err)

	myAddr := new(net.TCPAddr)
	myAddr.IP = net.ParseIP(myIP)
	myAddr.Port = myPort

	conn, err := net.DialTCP(protocol, myAddr, tcpAddr)
	CheckErrorExit(err)

	messageBuf := fileNameToByte(fileName)
	fmt.Println(messageBuf)
	send(conn, fileName)

}

func send(conn net.Conn, fileName string) bool {
	defer conn.Close()
	fp, err := os.Open(fileName)
	CheckError(err)

	defer fp.Close()

	messageBuf := fileNameToByte(fileName)

	conn.Write(messageBuf)
	fmt.Println("Sent the file name")
	conn.SetDeadline(time.Now().Add(50 * time.Second))
	for {

		messageBuf = make([]byte, SocketByte)
		messageLen, err := fp.Read(messageBuf[:SocketByte])
		messageBuf = messageBuf[:messageLen]

		if messageLen == 0 {
			fmt.Println("hogee")
			conn.Write([]byte{})
			break
		}
		if CheckError(err) == true {
			return false
		}

		_, err = conn.Write(messageBuf)
		if CheckError(err) == true {
			return false
		}

	}
	fmt.Println("sent the file data")
	conn.SetDeadline(time.Now().Add(1 * time.Second))
	return DownloadHashAndSendStatus(conn, fileName)
}

func DownloadHash(conn net.Conn) []byte {
	DownloadHash := make([]byte, SHA256ByteLen)
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err := conn.Read(DownloadHash)
	CheckError(err)
	return DownloadHash
}

func DownloadHashAndSendStatus(conn net.Conn, fileName string) bool {
	fmt.Println("DownloadHashAndSendStatus")
	hash := CreateSHA256(fileName)
	fmt.Println(hash)

	//hashをダウンロード
	downloadHash := DownloadHash(conn)

	//ステータスの送信
	if reflect.DeepEqual(hash, downloadHash) {
		fmt.Println("Complete File Transefer")
		conn.Write([]byte{0})
		return true
	} else {
		fmt.Println("NOT Complete File Transefer!!")
		conn.Write([]byte{1})
		return false
	}
}
