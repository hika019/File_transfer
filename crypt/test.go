package main

import (
	"fmt"
	"math/big"
)

func main() {

	sKey, pKey := GenKeyRSA(1024 * 4)

	/*
		fmt.Println(sKey.d)
		fmt.Println(pKey.e)
		fmt.Println(pKey.n)
	*/

	hoge := big.NewInt(123456789)
	fmt.Println(hoge)
	c := EnCryptRSA(pKey, hoge)
	fmt.Println(c)
	s := DeCryptRSA(sKey, pKey, c)
	fmt.Println(s)

	plainText := []byte("Bob loves Alice. But Alice hate Bob...")
	key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")
	// Create new AES cipher block
	block := GenAESBlock(key)

	cipherText := EnCryptAES(block, plainText)
	fmt.Printf("Cipher text: %x \n", cipherText)

	block2 := GenAESBlock(key)
	fmt.Printf("Decrypted text: %s\n", string(DecrptAES(block2, cipherText)))
}
