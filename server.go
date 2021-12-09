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
	file_name := string(messageBuf)
	file_name = file_name[:(strings.Index(file_name, ":"))]
	fmt.Println("filename: ", file_name)

	fp, err := os.OpenFile(file_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err)

	tmp := 0
	conn.SetReadDeadline(time.Now().Add(50 * time.Second))
	for {
		messageBuf := make([]byte, SocketSize)
		messageLen, err := conn.Read(messageBuf)
		//fmt.Println(messageBuf)
		CheckError(err)

		data_size_byte := make([]byte, 2)
		var data_size uint16 = 0

		fmt.Println(messageBuf)

		if messageLen != 0 {
			data_size_byte[0] = messageBuf[DataSizeBytePos1]
			data_size_byte[1] = messageBuf[DataSizeBytePos2]
			data_size = ByteToInt(data_size_byte)
		}

		if uint16(SocketDataSize) < data_size {
			//そこそこ大きい画像を送るときに最初に謎データがくっつくでスキップ
			fmt.Println("data_size が無効")
			data_size = 0
			continue
		}

		if data_size == 0 {
			fmt.Println("Downloaded file data")
			fmt.Println(tmp)
			break
		}

		//fmt.Printf("%d byte\n", messageLen)

		//fmt.Println(data_size)

		//fmt.Println(string(messageBuf[:data_size]))
		fmt.Println(tmp)

		fmt.Println(data_size)
		fmt.Println("######")
		fmt.Fprintf(fp, "%v", string(messageBuf[:data_size]))

		//ファイルに書き込み
		tmp++
	}
	fp.Close()

	//hash := file_hash(file_name)
	//conn.Write(hash)

}
