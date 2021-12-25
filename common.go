package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

const SocketByte int = 300 //あまり大きいとバッファオーバーフローが起きる
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

func CheckError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal error: ", err.Error())
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
