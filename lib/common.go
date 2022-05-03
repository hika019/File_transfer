package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const SocketByte int = 1200
const SocketDataByte int = SocketByte - 4
const DataSizeBytePos0 int = SocketDataByte + 0
const DataSizeBytePos1 int = SocketDataByte + 1
const DataSizeBytePos2 int = SocketDataByte + 2

const SHA256ByteLen int = 32

type config struct {
	SentIP    string `json:'sentIP'`
	ReceiveIP string `json:'receivIP'`
}

func IntToByte(i uint16) []byte {
	bin := make([]byte, 2)
	bin[0] = uint8(i % 256)
	bin[1] = uint8((i / 256) % 256)

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

	//config := *loadConfig()

	for _, inter := range interfaces {
		addr, err := inter.Addrs()
		CheckErrorExit(err)

		for _, a := range addr {
			if ipnet, ok := a.(*net.IPNet); ok {
				if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					fmt.Println(ipnet.Mask.Size())
					return ipnet.IP.String()
				}
			}
		}

	}

	return ""
}

func FileNameToByte(f string) []byte {
	data := make([]byte, SocketByte)
	hash := CreateSHA256(f)

	if strings.Contains(f, "/") {
		i := strings.LastIndex(f, "/")
		f = f[i+1:]
	}

	if strings.Contains(f, `\`) {
		i := strings.LastIndex(f, `\`)
		f = f[i+1:]
	}
	fByte := []byte(f)

	if SocketDataByte-SHA256ByteLen < len(fByte) {
		fmt.Printf("err: strStaticByte()/ strを%d byte以下にしてください\n", SocketDataByte)
		os.Exit(1)
	}

	for i, v := range fByte {
		data[i] = v
	}

	tmp := IntToByte(uint16(len(fByte)))

	data[DataSizeBytePos1] = tmp[0]
	data[DataSizeBytePos2] = tmp[1]

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

func LoadConfig() *config {
	f, err := os.Open("setting.json")
	CheckErrorExit(err)

	defer f.Close()

	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	CheckErrorExit(err)

	return &cfg
}

func MaskStr(IP string) string {
	mask, err := strconv.Atoi(IP[len(IP)-2:])
	CheckErrorExit(err)

	maskInt := uint32(4294967295)
	maskInt = maskInt >> uint32(32-mask) << uint32(32-mask)
	maskbin := make([]byte, 4)

	for i := 0; i < 2; i++ {
		tmp := uint16(maskInt >> (16 * i) & 65535)

		bin := IntToByte(tmp)
		maskbin[4-i*2-1] = bin[0]
		maskbin[4-i*2-2] = bin[1]
		fmt.Println(bin)

	}
	return hex.EncodeToString(maskbin)
}
