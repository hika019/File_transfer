package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/hika019/File_transfer/lib"
)

func GenAESBlock(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	lib.CheckErrorExit(err)
	return block
}

func EnCryptAES(block cipher.Block, s []byte) []byte {
	// Create IV
	cipherText := make([]byte, aes.BlockSize+len(s))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Printf("err: %s\n", err)
	}

	// Encrypt
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], s)
	return cipherText
}

func DecrptAES(block cipher.Block, c []byte) []byte {
	// Decrpt
	decryptedText := make([]byte, len(c[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, c[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, c[aes.BlockSize:])
	return decryptedText
}
