package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

const SocketSize int = 1024
const SocketDataSize int = SocketSize - 3
const DataSizeBytePos1 int = SocketDataSize + 1
const DataSizeBytePos2 int = SocketDataSize + 2

func IntToByte(i uint16) []byte {
	byteData := make([]byte, 2)

	byteData[0] = uint8(i % 256)
	byteData[1] = uint8((i / 256) % 256)
	return byteData[:]
}

func ByteToInt(byteData []byte) uint16 {
	intData := int(byteData[0])
	intData += int(byteData[1]) * 256
	return uint16(intData)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: error: ", err.Error())
		os.Exit(1)
	}
}

func FileHash(filename string) []byte {
	r, err := os.Open(filename)
	CheckError(err)

	hash := sha256.New()

	if _, err := io.Copy(hash, r); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: error: ", err.Error())
		os.Exit(1)
	}

	v := hash.Sum(nil)
	return v
}
