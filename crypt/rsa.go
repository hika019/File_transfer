package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/hika019/File_transfer/lib"
)

type SecretKey struct {
	d []byte
}

type PublicKey struct {
	n []byte
	e []byte
}

//k=ビット数(2048~)
func GenKeyRSA(k int) (SecretKey, PublicKey) {
	fmt.Println("genKey")
	p, err := rand.Prime(rand.Reader, k/2)
	lib.CheckErrorExit(err)
	q := p
	for p == q {
		q, err = rand.Prime(rand.Reader, k/2)
		lib.CheckErrorExit(err)
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
		lib.CheckErrorExit(err)
		gcdAns = gcd(e, phi)
	}

	d := new(big.Int)
	d.Exp(e, big.NewInt(-1), phi)
	fmt.Println()
	fmt.Printf("%x\n", d.Bytes())

	secretKey := SecretKey{d: d.Bytes()}
	publickKey := PublicKey{n: n.Bytes(), e: e.Bytes()}

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

	c.Exp(ByteToBigInt(s), ByteToBigInt(key.e), ByteToBigInt(key.n))
	return c.Bytes()
}

func ByteToBigInt(b []byte) *big.Int {
	a := new(big.Int)
	a.SetBytes(b)
	return a
}

func DeCryptRSA(sKey SecretKey, pKey PublicKey, c []byte) []byte {
	s := new(big.Int)

	s.Exp(ByteToBigInt(c), ByteToBigInt(sKey.d), ByteToBigInt(pKey.n))
	return s.Bytes()
}
