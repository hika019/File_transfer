package lib

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type SecretKey struct {
	D []byte
}

type PublicKey struct {
	N []byte
	E []byte
}

//GenKeyRSAのビット数
const KBitLen int = 2048
const KByteLen int = KBitLen / 8

func PublicKeyToByte(p PublicKey) []byte {
	d := make([]byte, KByteLen*2)

	j := KByteLen - 1
	for i := len(p.N) - 1; 0 <= i; i-- {
		d[j] = p.N[i]
		j--
	}

	j = KByteLen*2 - 1
	for i := len(p.E) - 1; 0 <= i; i-- {
		d[j] = p.E[i]
		j--
	}
	return d
}

func ByteToPublickKey(d []byte) PublicKey {
	p := new(PublicKey)

	p.N = d[0:KByteLen]
	p.E = d[KByteLen : KByteLen*2]

	return *p
}

//k=ビット数(2048)
func GenKeyRSA() (SecretKey, PublicKey) {
	fmt.Println("genKey")
	p, err := rand.Prime(rand.Reader, KBitLen/2)
	CheckErrorExit(err)
	q := p
	for p == q {
		q, err = rand.Prime(rand.Reader, KBitLen/2)
		CheckErrorExit(err)
	}

	n := new(big.Int)
	n.Mul(p, q)

	phi := new(big.Int)
	phi.Mul(p.Sub(p, big.NewInt(1)), q.Sub(q, big.NewInt(1)))

	e := big.NewInt(65537)
	gcdAns := gcd(e, phi)
	for big.NewInt(1).Cmp(gcdAns) != 0 {
		tmp := new(big.Int)
		e, err = rand.Int(rand.Reader, tmp.Sub(phi, big.NewInt(2)))
		e.Add(e, big.NewInt(2))
		CheckErrorExit(err)
		gcdAns = gcd(e, phi)
	}

	d := new(big.Int)
	d.Exp(e, big.NewInt(-1), phi)
	//fmt.Println()
	//fmt.Printf("%x\n", d.Bytes())

	secretKey := SecretKey{D: d.Bytes()}
	publickKey := PublicKey{N: n.Bytes(), E: e.Bytes()}

	fmt.Println("Gen RSA key")
	return secretKey, publickKey
}

func gcd(m *big.Int, n *big.Int) *big.Int {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	z.GCD(x, y, m, n)
	return z
}

func EnCryptRSA(key PublicKey, s []byte) []byte {
	c := new(big.Int)

	c.Exp(ByteToBigInt(s), ByteToBigInt(key.E), ByteToBigInt(key.N))
	return c.Bytes()
}

func ByteToBigInt(b []byte) *big.Int {
	a := new(big.Int)
	a.SetBytes(b)
	return a
}

func DeCryptRSA(sKey SecretKey, pKey PublicKey, c []byte) []byte {
	s := new(big.Int)

	s.Exp(ByteToBigInt(c), ByteToBigInt(sKey.D), ByteToBigInt(pKey.N))
	return s.Bytes()
}
