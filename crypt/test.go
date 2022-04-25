package main

import (
	"fmt"
	"math/big"
)

func main() {

	sKey, pKey := GenKey(1024 * 4)

	/*
		fmt.Println(sKey.d)
		fmt.Println(pKey.e)
		fmt.Println(pKey.n)
	*/

	hoge := big.NewInt(123456789)
	fmt.Println(hoge)
	c := EnCrypt(pKey, hoge)
	fmt.Println(c)
	s := DeCrypt(sKey, pKey, c)
	fmt.Println(s)

}
