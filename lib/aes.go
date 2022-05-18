package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

const AESKeyLen int = 32

func GenAESKey(keyLen int) []byte {
	key := make([]byte, keyLen)
	for i := 0; i < keyLen; i++ {
		tmp, err := rand.Int(rand.Reader, big.NewInt(256))
		CheckErrorExit(err)
		key[i] = uint8(tmp.Int64() % 256)
	}
	return key
}

func GenAESBlock(key []byte) cipher.Block {
	fmt.Println("call -> GenAESBlock")
	block, err := aes.NewCipher(key)
	CheckErrorExit(err)
	fmt.Println("end -> GenAESBlock")
	return block
}

func InitAESBlock() cipher.Block {
	k, _ := hex.DecodeString("645E739A7F9F162725C1533DC2C5E827")
	return GenAESBlock(k)
}

func EnCryptAES(block cipher.Block, s []byte, useCrypt bool) []byte {
	//fmt.Println("call -> EnCryptAES")
	if !useCrypt {
		return s
	}

	// Create IV
	cipherText := make([]byte, aes.BlockSize+len(s))
	iv := cipherText[:aes.BlockSize]
	_, err := io.ReadFull(rand.Reader, iv)
	CheckErrorExit(err)

	// Encrypt
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], s)
	return cipherText
}

func DecryptAES(block cipher.Block, c []byte, useCrypt bool) []byte {
	//fmt.Println("call -> DecryptAES")
	if !useCrypt {
		return c
	}

	// Decrpt
	decryptedText := make([]byte, len(c[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, c[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, c[aes.BlockSize:])
	//fmt.Println("end -> DecryptAES")
	return decryptedText
}
