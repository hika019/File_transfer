package main

import (
	"crypto/aes"
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

	s, p := lib.GenKeyRSA()

	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}

		go handleClient(conn, p, s)

	}
}

func handleClient(conn net.Conn, p lib.PublicKey, s lib.SecretKey) {
	useCrypt := false
	block := lib.InitAESBlock()

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
	fmt.Println(messageLen)
	//EOFエラー回避
	if messageLen == 0 {
		return
	}
	if lib.CheckError(err) {
		return
	}

	if reflect.DeepEqual(messageBuf[:messageLen], []byte{255, 192, 0, 0, 255}) {

		fmt.Println("Use Crypt")
		conn.Write(lib.PublicKeyToByte(p))
		useCrypt = true
		//fmt.Println(p)
		messageBuf = make([]byte, lib.SocketByte+aes.BlockSize)
		messageLen, err = conn.Read(messageBuf)
		//fmt.Println(messageBuf)
		key := lib.DeCryptRSA(s, p, messageBuf[:messageLen])
		fmt.Println(key)
		block = lib.GenAESBlock(key)

		//EOFエラー回避
		if messageLen == 0 {
			return
		}
		if lib.CheckError(err) {
			return
		}

		messageBuf = make([]byte, lib.SocketByte+aes.BlockSize)
		messageLen, err = conn.Read(messageBuf)
		//EOFエラー回避
		if messageLen == 0 {
			return
		}
		if lib.CheckError(err) {
			return
		}
	}

	fmt.Println("Crypt: ", useCrypt)

	fileName, hash := lib.ByteToFileName(lib.DecryptAES(block, messageBuf[:messageLen], useCrypt))

	fileName = senderIP + "/" + fileName
	//fmt.Println("filename: ", fileName, hash)

	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if lib.CheckError(err) {
		return
	}
	defer fp.Close()

	buff := make([]byte, 0, (lib.SocketByte+aes.BlockSize)*2)
	buffLen := 0

	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		messageBuf = make([]byte, lib.SocketByte+aes.BlockSize)
		messageLen, err = conn.Read(messageBuf)

		buff = append(buff, messageBuf...)
		buffLen += messageLen

		//EOFエラー回避
		if messageLen == 0 {
			//meta := make([]byte, lib.SocketByte+aes.BlockSize-buffLen)
			//messageBuf = append(messageBuf, meta...)

			if useCrypt {
				_, err = fp.Write(lib.DecryptAES(block, buff, useCrypt)[:buffLen-aes.BlockSize])
			} else {
				_, err = fp.Write(lib.DecryptAES(block, buff, useCrypt)[:buffLen])
			}

			lib.CheckError(err)
			fmt.Println("download file")
			break
		}

		lib.CheckErrorExit(err)

		if buffLen >= lib.SocketByte+aes.BlockSize {
			_, err = fp.Write(lib.DecryptAES(block, messageBuf[:messageLen], useCrypt))
			lib.CheckError(err)

			buff = buff[lib.SocketByte+aes.BlockSize:]
			buffLen -= lib.SocketByte + aes.BlockSize
		}

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
