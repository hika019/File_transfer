package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type SecretKey struct {
	d *big.Int
}

type PublicKey struct {
	n *big.Int
	e *big.Int
}

//k=ビット数(2048~)
func genKey(k int) (SecretKey, PublicKey) {
	fmt.Println("genKey")
	p, err := rand.Prime(rand.Reader, k/2)
	CheckErrorExit(err)
	q := p
	for p == q {
		q, err = rand.Prime(rand.Reader, k/2)
		CheckErrorExit(err)
	}

	n := new(big.Int)
	n.Mul(p, q)

	phi := new(big.Int)
	phi.Mul(p.Sub(p, big.NewInt(1)), q.Sub(q, big.NewInt(1)))

	gcdAns := big.NewInt(2)
	e := big.NewInt(65537)
	gcdAns = gcd(e, phi)
	for big.NewInt(1).Cmp(gcdAns) != 0 {
		tmp := new(big.Int)
		e, err = rand.Int(rand.Reader, tmp.Sub(phi, big.NewInt(2)))
		e.Add(e, big.NewInt(2))
		CheckErrorExit(err)
		gcdAns = gcd(e, phi)
	}

	d := new(big.Int)
	d.Exp(e, big.NewInt(-1), phi)

	secretKey := SecretKey{d: d}
	publickKey := PublicKey{n: n, e: e}

	return secretKey, publickKey
}

func gcd(m *big.Int, n *big.Int) *big.Int {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	z.GCD(x, y, m, n)
	return z
}

func enCrypt(key PublicKey, s *big.Int) *big.Int {
	c := new(big.Int)

	c.Exp(s, key.e, key.n)
	return c
}

func deCrypt(sKey SecretKey, pKey PublicKey, c *big.Int) *big.Int {
	s := new(big.Int)

	s.Exp(c, sKey.d, pKey.n)
	return s
}
