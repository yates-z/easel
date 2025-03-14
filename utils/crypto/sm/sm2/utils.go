package sm2

import (
	"math/big"
)

func Add(x, y *big.Int) *big.Int {
	var z big.Int
	z.Add(x, y)
	return &z
}

func Sub(x, y *big.Int) *big.Int {
	var z big.Int
	z.Sub(x, y)
	return &z
}

func Mod(x, y *big.Int) *big.Int {
	var z big.Int
	z.Mod(x, y)
	return &z
}

func ModInverse(x, y *big.Int) *big.Int {
	var z big.Int
	z.ModInverse(x, y)
	return &z
}

func Mul(x, y *big.Int) *big.Int {
	var z big.Int
	z.Mul(x, y)
	return &z
}

func Lsh(x *big.Int, n uint) *big.Int {
	var z big.Int
	z.Lsh(x, n)
	return &z
}

func SetBit(x *big.Int, i int, b uint) *big.Int {
	var z big.Int
	z.SetBit(x, i, b)
	return &z
}

func And(x, y *big.Int) *big.Int {
	var z big.Int
	z.And(x, y)
	return &z
}

func IsEcPointInfinity(x, y *big.Int) bool {
	if x.Sign() == 0 && y.Sign() == 0 {
		return true
	}
	return false
}

func ZForAffine(x, y *big.Int) *big.Int {
	z := new(big.Int)
	if x.Sign() != 0 || y.Sign() != 0 {
		z.SetInt64(1)
	}
	return z
}
