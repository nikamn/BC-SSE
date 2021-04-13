package interpolation

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/ncw/gmp"
	"github.com/nikamn/BC-SSE/utils/polyring"
	"github.com/stretchr/testify/assert"
)

const PolyOrder = 500
const RandSeed = 2

var largeStr string

func genPrime(p *gmp.Int, bitnum int) {
	var buffer bytes.Buffer
	for i := 0; i < bitnum; i++ {
		buffer.WriteString("0")
	}

	largeStr = "1"
	largeStr += buffer.String()

	p.SetString(largeStr, 10)
	// No next_prime method in go yet. Placeholder for now
	p.Set(gmp.NewInt(15486511))
	// p.Set(gmp.NewInt(7))
}

func TestLagrangeInterpolate(t *testing.T) {
	p := gmp.NewInt(0)
	genPrime(p, 256)
	r := rand.New(rand.NewSource(RandSeed))

	fmt.Printf("Prime p = %s\n", p.String())

	originalPoly, err := polyring.NewRand(PolyOrder, r, p)
	assert.Nil(t, err, "New")

	// Test EvalArray
	x := make([]*gmp.Int, PolyOrder+1)
	y := make([]*gmp.Int, PolyOrder+1)
	polyring.VecInit(x)
	polyring.VecInit(y)
	polyring.VecRand(x, p, r)

	originalPoly.EvalModArray(x, p, y)

	fmt.Println("Finished eval")
	fmt.Println("Starting interpolation")

	reconstructedPoly, err := LagrangeInterpolate(PolyOrder, x, y, p)
	assert.Nil(t, err, "New")

	//fmt.Printf("Original Poly ")
	//originalPoly.Print()

	//fmt.Printf("Reconstructed Poly ")
	//reconstructedPoly.Print()
	assert.True(t, reconstructedPoly.IsSame(originalPoly))
}
