package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s message", os.Args[0])
		os.Exit(1)
	}

	protocol := "tcp"
	serverIP := "192.168.11.50"
	serverPort := "55555"
	myIP := "192.168.11.30"
	myPort := 55556

	fileName := os.Args[1]

	fp, err := os.Open(fileName)
	CheckError(err)

	tcpAddr, err := net.ResolveTCPAddr(protocol, serverIP+":"+serverPort)
	CheckError(err)

	myAddr := new(net.TCPAddr)
	myAddr.IP = net.ParseIP(myIP)
	myAddr.Port = myPort
	conn, err := net.DialTCP(protocol, myAddr, tcpAddr)
	CheckError(err)

	defer conn.Close()

	defer fp.Close()
	messageBuf := make([]byte, SocketSize)
	tmp := 0

	conn.SetDeadline(time.Now().Add(50 * time.Second))
	fmt.Println(fileName)
	conn.Write([]byte(fileName + ":"))
	fmt.Println("Sent the file name")

	for {
		messageLen, err := fp.Read(messageBuf[:SocketDataSize])
		fmt.Println(messageLen)
		messageBuf = IntToByte(messageBuf, uint16(messageLen))
		tmp++

		if messageLen == 0 {
			break
		}
		CheckError(err)

		fmt.Println(tmp)
		//fmt.Println(messageBuf)
		conn.Write(messageBuf)
	}
	fmt.Println("sent the file data")
	fmt.Println(tmp)

	hash := CreateSHA256(fileName)
	fmt.Println(hash)

	DownloadHash := make([]byte, SHA256ByteLen)

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Read(DownloadHash)
	CheckError(err)

	if reflect.DeepEqual(hash, DownloadHash) {
		fmt.Println("Complete File Transefer")
		conn.Write([]byte{0})
	} else {
		fmt.Println("NOT Complete File Transefer!!")
		conn.Write([]byte{1})
	}

}
