package main

import (
	"fmt"
	"net"
	"os"
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

	file_name := os.Args[1]

	fp, err := os.Open(file_name)
	checkError(err)

	tcpAddr, err := net.ResolveTCPAddr(protocol, serverIP+":"+serverPort)
	checkError(err)

	myAddr := new(net.TCPAddr)
	myAddr.IP = net.ParseIP(myIP)
	myAddr.Port = myPort
	conn, err := net.DialTCP(protocol, myAddr, tcpAddr)
	checkError(err)

	defer conn.Close()

	defer fp.Close()
	sent_binary := make([]byte, soket_size)
	tmp := 0

	conn.SetDeadline(time.Now().Add(50 * time.Second))
	fmt.Println(file_name)
	conn.Write([]byte(file_name + ":"))
	fmt.Println("Sent the file name")

	for {
		bytes, err := fp.Read(sent_binary[:data_size])
		fmt.Println(bytes)
		bytes_size := int_to_byte(uint16(bytes))
		tmp++

		sent_binary[data_size_byte_pos1] = bytes_size[0]
		sent_binary[data_size_byte_pos2] = bytes_size[1]
		if bytes == 0 {
			break
		}
		checkError(err)

		//fmt.Println(sent_binary)
		fmt.Println(tmp)
		fmt.Println(sent_binary)
		conn.Write(sent_binary)
		//fmt.Println(bytes)
		//fmt.Println(string(sent_binary))
		//fmt.Println(buf)
	}
	fmt.Println("sent the file data")
	fmt.Println(tmp)
	conn.Write([]byte("end sent data"))

}
