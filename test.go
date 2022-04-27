package main

import (
	"fmt"

	"github.com/hika019/File_transfer/lib"
)

func main() {
	data := lib.LoadConfig()
	sIP := data.SentIP
	rIP := data.ReceiveIP

	sMask := lib.MaskStr(sIP)
	rMask := lib.MaskStr(rIP)

	fmt.Println(sMask)
	fmt.Println(rMask)

	fmt.Println(sIP)

}
