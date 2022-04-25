package lib

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"os"
)

const SocketByte int = 1200
const SocketDataByte int = SocketByte - 4
const DataSizeBytePos0 int = SocketDataByte + 0
const DataSizeBytePos1 int = SocketDataByte + 1
const DataSizeBytePos2 int = SocketDataByte + 2

const SHA256ByteLen int = 32

func IntToByte(bin []byte, i uint16) []byte {
	bin[DataSizeBytePos1] = uint8(i % 256)
	bin[DataSizeBytePos2] = uint8((i / 256) % 256)

	return bin[:]
}

func ByteToInt(byteData []byte) uint16 {
	dataLen := int(byteData[DataSizeBytePos1])
	dataLen += int(byteData[DataSizeBytePos2]) * 256
	return uint16(dataLen)
}

func CheckError(err error) bool {
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal error: ", err.Error())
		return true
	}
	return false
}

func CheckErrorExit(err error) {
	if CheckError(err) {
		os.Exit(1)
	}
}

func CreateSHA256(fileName string) []byte {
	r, err := os.Open(fileName)
	CheckError(err)

	hash := sha256.New()

	_, err = io.Copy(hash, r)
	CheckError(err)

	v := hash.Sum(nil)
	return v
}

func MyAddr() string {
	interfaces, err := net.Interfaces()
	CheckErrorExit(err)

	for _, inter := range interfaces {
		addr, err := inter.Addrs()
		CheckErrorExit(err)

		for _, a := range addr {
			if ipnet, ok := a.(*net.IPNet); ok {
				if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}

	}

	return ""
}

func FileNameToByte(f string) []byte {
	data := make([]byte, SocketByte)
	fByte := []byte(f)

	if SocketDataByte-SHA256ByteLen < len(fByte) {
		fmt.Printf("err: strStaticByte()/ strを%d byte以下にしてください\n", SocketDataByte)
		os.Exit(1)
	}

	for i, v := range fByte {
		data[i] = v
	}

	data = IntToByte(data, uint16(len(fByte)))
	hash := CreateSHA256(f)

	for i, v := range hash {
		data[SocketDataByte-SHA256ByteLen+i] = v
	}

	return data[:]
}

func ByteToFileName(data []byte) (string, []byte) {
	fileNameLen := ByteToInt(data)
	filename := data[:fileNameLen]
	hash := data[SocketDataByte-SHA256ByteLen : SocketDataByte]

	return string(filename), hash
}
