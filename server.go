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
	checkError(err)

	listner, err := net.ListenTCP(protocol, tcpAddr)
	checkError(err)

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
	messageBuf := make([]byte, soket_size)
	//fmt.Println(messageBuf)

	_, err := conn.Read(messageBuf)
	file_name := string(messageBuf)
	file_name = file_name[:(strings.Index(file_name, ":"))]
	fmt.Println("filename: ", file_name)

	fp, err := os.OpenFile(file_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError(err)

	defer fp.Close()

	tmp := 0
	conn.SetReadDeadline(time.Now().Add(50 * time.Second))
	for {
		messageBuf := make([]byte, soket_size)
		messageLen, err := conn.Read(messageBuf)
		//fmt.Println(messageBuf)
		checkError(err)

		data_size_byte := make([]byte, 2)
		var data_size uint16 = 0

		if messageLen != 0 {
			data_size_byte[0] = messageBuf[data_size_byte_pos1]
			data_size_byte[1] = messageBuf[data_size_byte_pos2]
			data_size = byte_to_int(data_size_byte)
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
		fmt.Println(messageBuf)
		fmt.Fprintf(fp, "%v", string(messageBuf[:data_size]))

		//ファイルに書き込み
		tmp++

	}

}
