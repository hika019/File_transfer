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

	messageBuf := strStaticByte(fileName)
	fmt.Println(messageBuf)
	send(conn, fileName)

}

func send(conn net.Conn, fileName string) bool {
	defer conn.Close()
	fp, err := os.Open(fileName)
	CheckError(err)

	defer fp.Close()

	messageBuf := strStaticByte(fileName)
	//fmt.Println(messageBuf)
	tmp := 0

	conn.Write(messageBuf)
	fmt.Println("Sent the file name")

	for {

		messageBuf = make([]byte, SocketByte)
		messageLen, err := fp.Read(messageBuf[:SocketByte])
		if messageLen == 0 {
			conn.Write([]byte{})
			break
		}
		if CheckError(err) == true {
			return false
		}

		messageBuf = messageBuf[:messageLen]

		//messageBuf = IntToByte(messageBuf, uint16(messageLen))
		//messageBuf[DataSizeBytePos0] = uint8(1)

		tmp++

		conn.SetDeadline(time.Now().Add(1 * time.Second))
		dataLen, err := conn.Write(messageBuf)
		if CheckError(err) == true {
			return false
		}

		//接続が切断されたらbreak
		if dataLen == 0 {
			break
		}
	}
	fmt.Println("sent the file data")
	fmt.Println(tmp)
	conn.SetDeadline(time.Now().Add(1 * time.Second))
	return DownloadHashAndSendStatus(conn, fileName)
}

func DownloadHash(conn net.Conn) []byte {
	DownloadHash := make([]byte, SHA256ByteLen)

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

func strStaticByte(str string) []byte {
	data := make([]byte, SocketByte)
	strByte := []byte(str)

	if SocketDataByte < len(strByte) {
		fmt.Printf("err: strStaticByte()/ strを%d byte以下にしてください\n", SocketDataByte)
		os.Exit(1)
	}

	for i, v := range strByte {
		data[i] = v
	}

	data = IntToByte(data, uint16(len(strByte)))
	data[DataSizeBytePos0] = uint8(1)
	return data[:]
}
