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

	tmp := 0
	conn.SetReadDeadline(time.Now().Add(50 * time.Second))
	for {
		messageBuf := make([]byte, SocketSize)
		messageLen, err := conn.Read(messageBuf)
		//fmt.Println(messageBuf)
		CheckError(err)

		var dataSize uint16 = 0

		//fmt.Println(messageBuf)

		if messageLen != 0 {
			dataSize = ByteToInt(messageBuf)
		}

		if uint16(SocketDataSize) < dataSize {
			//そこそこ大きい画像を送るときに最初に謎データがくっつくでスキップ
			fmt.Println("data_size が無効")
			dataSize = 0
			continue
		}

		if dataSize == 0 {
			fmt.Println("Downloaded file data")
			fmt.Println(tmp)
			break
		}

		fmt.Println(tmp)

		fmt.Printf("%v byte", dataSize)
		fmt.Println("######")
		fmt.Fprintf(fp, "%v", string(messageBuf[:dataSize]))

		//ファイルに書き込み
		tmp++
	}
	fp.Close()

	//hash := file_hash(file_name)
	//conn.Write(hash)

}
