package commitpbc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/Nik-U/pbc"
	"github.com/nikamn/BC-SSE/utils/conv"
	"github.com/nikamn/BC-SSE/utils/ecparam"
	"github.com/nikamn/BC-SSE/utils/polyring"
)

// Curve exported from ecparam
var Curve = ecparam.PBC256

// PolyCommit struct
// a commitment to a polynomial {a0, a1, ..., at} is g^at
// ai are from the multiplicative group of integers modulo p
type PolyCommit struct {
	c []*pbc.Element
}

// GobEncode returns an encoded commitment to bytes
func (comm PolyCommit) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	binary := make([][]byte, len(comm.c))

	for i := range binary {
		binary[i] = comm.c[i].CompressedBytes()
		if comm.c[i].Is0() {
			binary[i][0] = 0xff
		}
	}

	if err := enc.Encode(binary); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GobDecode returns a decoded commitment from input bytes
func (comm *PolyCommit) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(r)

	var binary [][]byte

	if err := dec.Decode(&binary); err != nil {
		return err
	}

	comm.c = make([]*pbc.Element, len(binary))

	for i := range binary {
		comm.c[i] = Curve.Pairing.NewG1()
		// handling infinity point specially
		if binary[i][0] == 0xff {
			comm.c[i].Set0()
		} else {
			comm.c[i].SetCompressedBytes(binary[i])
		}
	}

	return nil
}

// Equals checks if polyCommits are equal or not
func (comm PolyCommit) Equals(other PolyCommit) bool {
	if len(comm.c) != len(other.c) {
		return false
	}

	for i := range other.c {
		if !comm.c[i].Equals(other.c[i]) {
			return false
		}
	}

	return true
}

// Bytes encodes a commitment to binary bytes using GobEncode()
func (comm PolyCommit) Bytes() []byte {
	binary, err := comm.GobEncode()
	if err != nil {
		panic(err.Error())
	}

	return binary
}

// String converts a commitment to a string
func (comm PolyCommit) String() string {
	s := ""
	for i := range comm.c {
		s += fmt.Sprintf("%s, ", comm.c[i])
	}

	return s
}

// Print prints a commitment
func (comm PolyCommit) Print() {
	fmt.Println("comm =", comm.String())
}

// NewPolyCommit returns a commitment to an input polynomial
func NewPolyCommit(polynomial polyring.Polynomial) PolyCommit {
	allCoeff := polynomial.GetAllCoefficients()

	comm := PolyCommit{
		c: make([]*pbc.Element, len(allCoeff)),
	}

	for i, coeff := range allCoeff {
		comm.c[i] = Curve.Pairing.NewG1()
		pow := conv.GmpInt2BigInt(coeff)
		Curve.G.PowBig(comm.c[i], pow)
	}

	return comm
}

// Verify verifies whether commitment to a polynomial is correct or not
func (comm PolyCommit) Verify(poly polyring.Polynomial) bool {
	coeffs := poly.GetAllCoefficients()

	commCheck := PolyCommit{
		c: make([]*pbc.Element, len(coeffs)),
	}

	for i, coeff := range coeffs {
		commCheck.c[i] = Curve.Pairing.NewG1()
		Curve.G.PowBig(commCheck.c[i], conv.GmpInt2BigInt(coeff))
		if !commCheck.c[i].Equals(comm.c[i]) {
			return false
		}
	}

	return true
}

// VerifyEval verifies a commitment using (x,y)
func (comm PolyCommit) VerifyEval(x *big.Int, y *big.Int) bool {
	gYRef := Curve.Pairing.NewG1()
	Curve.G.PowBig(gYRef, y)

	xx := big.NewInt(1)

	gPx := Curve.Pairing.NewG1()
	gPx.Set1()

	tmp := Curve.Pairing.NewG1()
	for i := range comm.c {
		// tmp = g^ai^{x^i}
		tmp.PowBig(comm.c[i], xx)

		gPx.Mul(tmp, gPx)

		xx.Mul(xx, x)
		xx.Mod(xx, Curve.Nbig)
	}

	return gPx.Equals(gYRef)
}

// AdditiveHomomorphism return a commitment to Q+R
func AdditiveHomomorphism(commQ, commR PolyCommit) PolyCommit {
	if len(commQ.c) != len(commR.c) {
		panic("mismatch degree")
	}

	comm := PolyCommit{
		c: make([]*pbc.Element, len(commQ.c)),
	}

	for i := range comm.c {
		comm.c[i] = Curve.Pairing.NewG1()
		comm.c[i].Mul(commQ.c[i], commR.c[i])
	}

	return comm
}
