package conv

import (
	"math/big"

	"github.com/ncw/gmp"
)

// BigInt2GmpInt converts a Big integer to a Gmp integer
func BigInt2GmpInt(a *big.Int) *gmp.Int {
	b := gmp.NewInt(0)
	b.SetBytes(a.Bytes())

	return b
}

// GmpInt2BigInt converts a Gmp integer to a Big integer
func GmpInt2BigInt(a *gmp.Int) *big.Int {
	b := new(big.Int)
	b.SetBytes(a.Bytes())

	return b
}
