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
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s message\n", os.Args[0])
		os.Exit(1)
	}
	//TODO args[1]dir/fileの場合エラー

	protocol := "tcp"
	//serverIP := "192.168.11.50"
	serverIP := os.Args[1]
	serverPort := "55555"
	myIP := lib.MyAddr()
	myPort := 55556

	fileName := os.Args[2]

	tcpAddr, err := net.ResolveTCPAddr(protocol, serverIP+":"+serverPort)
	lib.CheckErrorExit(err)

	myAddr := new(net.TCPAddr)
	myAddr.IP = net.ParseIP(myIP)
	myAddr.Port = myPort

	conn, err := net.DialTCP(protocol, myAddr, tcpAddr)
	lib.CheckErrorExit(err)

	defer conn.Close()

	messageBuf := lib.FileNameToByte(fileName)
	fmt.Println(messageBuf)
	send(conn, fileName)

}

func send(conn net.Conn, fileName string) bool {
	defer conn.Close()
	fp, err := os.Open(fileName)
	lib.CheckError(err)

	defer fp.Close()

	messageBuf := lib.FileNameToByte(fileName)

	conn.Write(messageBuf)
	fmt.Println("Sent the filename")
	conn.SetDeadline(time.Now().Add(50 * time.Second))
	for {

		messageBuf = make([]byte, lib.SocketByte)
		messageLen, err := fp.Read(messageBuf[:lib.SocketByte])
		messageBuf = messageBuf[:messageLen]

		if messageLen == 0 {
			break
		}
		if lib.CheckError(err) == true {
			return false
		}

		_, err = conn.Write(messageBuf)
		if lib.CheckError(err) == true {
			return false
		}

	}
	fmt.Println("sent the file data")

	return DownloadStatus(conn)
}

func DownloadStatus(conn net.Conn) bool {
	messageBuff := make([]byte, 2)
	messageLen, err := conn.Read(messageBuff)

	if lib.CheckError(err) == true {
		return false
	}

	if reflect.DeepEqual(messageBuff[:messageLen], []byte{0}) {
		fmt.Println("Consistency: Yes")
		return true
	} else {
		fmt.Println("Consistency: No")
		return false
	}

}
