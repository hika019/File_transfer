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
	messageBuf := strStaticByte(fileName)
	fmt.Println(messageBuf)
	tmp := 0

	conn.Write(messageBuf)
	fmt.Println("Sent the file name")

	for {
		conn.SetDeadline(time.Now().Add(10 * time.Second))
		messageBuf = make([]byte, SocketByte)
		messageLen, err := fp.Read(messageBuf[:SocketDataByte])
		if messageLen == 0 {
			break
		}
		CheckError(err)
		//fmt.Println(messageLen)
		messageBuf = IntToByte(messageBuf, uint16(messageLen))
		messageBuf[DataSizeBytePos0] = uint8(1)

		if messageLen == 0 {
			break
		}

		if tmp%200 == 0 {
			time.Sleep(500 * time.Millisecond)
			fmt.Println(tmp)
		}

		tmp++
		//fmt.Println(messageBuf)
		conn.Write(messageBuf)

	}
	fmt.Println("sent the file data")
	fmt.Println(tmp)

	hash := CreateSHA256(fileName)
	fmt.Println(hash)

	//FIXME:ハッシュの確認の部分を実行するとファイルが送れない
	/*
		DownloadHash := make([]byte, SHA256ByteLen)

		conn.SetDeadline(ti me.Now().Add(5 * time.Second))
		DownloadHashLen, err := conn.Read(DownloadHash)
		CheckError(err)

		if reflect.DeepEqual(hash, DownloadHash[:DownloadHashLen]) {
			fmt.Println("Complete File Transefer")
			conn.Write([]byte{0})
		} else {
			fmt.Println("NOT Complete File Transefer!!")
			conn.Write([]byte{1})
		}
	*/
}

func DownloadHash(conn net.Conn) []byte {
	DownloadHash := make([]byte, SHA256ByteLen)

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err := conn.Read(DownloadHash)
	CheckError(err)
	return DownloadHash
}

//受信側のステータスの初期値が0のため正常を2とする
func SendFileStatus(conn net.Conn, fileName string) {
	hash := CreateSHA256(fileName)
	fmt.Println(hash)
	DownloadHash := DownloadHash(conn)

	if reflect.DeepEqual(hash, DownloadHash) {
		fmt.Println("Complete File Transefer")
		conn.Write([]byte{2})
	} else {
		fmt.Println("NOT Complete File Transefer!!")
		conn.Write([]byte{1})
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
