package p521

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/nikamn/BC-SSE/utils/polyring"
)

// Curve uses elliptic curve P521
var Curve = elliptic.P521()

// ECPoint struct
type ECPoint struct {
	x *big.Int
	y *big.Int
}

// NewECPoint returns a new point (x,y)
func NewECPoint(x, y *big.Int) ECPoint {
	return ECPoint{
		x: x,
		y: y,
	}
}

// GobEncode function encodes ECPoint to bytes
func (ecp ECPoint) GobEncode() ([]byte, error) {
	return Marshal(Curve, ecp.x, ecp.y), nil
}

// GobDecode checks decoding of encoded bytes is equal to ECpoint or not
func (ecp *ECPoint) GobDecode(buf []byte) error {
	ecp.x, ecp.y = Unmarshal(Curve, buf)
	return nil
}

// Equals function
func (ecp ECPoint) Equals(other ECPoint) bool {
	return 0 == ecp.x.Cmp(other.x) && 0 == ecp.y.Cmp(other.y)
}

func (ecp ECPoint) String() string {
	return fmt.Sprintf("(%s, %s)", ecp.x.String(), ecp.y.String())
}

// PolyCommit struct
// a commitment to a polynomial {a0, a1, ..., at} is g^at
// ai are from the multiplicative group of integers modulo p
type PolyCommit struct {
	c []ECPoint
}

// GobEncode returns an encoded commitment to bytes
func (comm PolyCommit) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	lenC := int32(len(comm.c))

	if err := enc.Encode(lenC); err != nil {
		return nil, err
	}

	if err := enc.Encode(comm.c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GobDecode returns a decoded commitment from input bytes
func (comm *PolyCommit) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(r)

	var lenC int32 = 0

	if err := dec.Decode(&lenC); err != nil {
		return err
	}

	comm.c = make([]ECPoint, lenC)
	if err := dec.Decode(&comm.c); err != nil {
		return err
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
		s += fmt.Sprintf("%s, ", comm.c[i].String())
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
		make([]ECPoint, len(allCoeff)),
	}

	for i, coeff := range allCoeff {
		x, y := Curve.ScalarBaseMult(coeff.Bytes())
		comm.c[i] = ECPoint{
			x: x,
			y: y,
		}
	}

	return comm
}

// Verify verifies whether commitment to a polynomial is correct or not
func (comm PolyCommit) Verify(poly polyring.Polynomial) bool {
	allCoeff := poly.GetAllCoefficients()

	commCheck := PolyCommit{
		make([]ECPoint, len(allCoeff)),
	}

	for i, coeff := range allCoeff {
		x, y := Curve.ScalarBaseMult(coeff.Bytes())
		commCheck.c[i] = ECPoint{
			x: x,
			y: y,
		}
		if !commCheck.c[i].Equals(comm.c[i]) {
			return false
		}
	}

	return true
}

// VerifyEval verifies a commitment using (x,y)
func (comm PolyCommit) VerifyEval(x *big.Int, y *big.Int) bool {
	gYRef := NewECPoint(Curve.ScalarBaseMult(y.Bytes()))

	xx := big.NewInt(1)

	gPxx := big.NewInt(0)
	gPxy := big.NewInt(0)
	for i := range comm.c {
		px, py := Curve.ScalarMult(comm.c[i].x, comm.c[i].y, xx.Bytes())

		gPxx, gPxy = Curve.Add(gPxx, gPxy, px, py)

		xx.Mul(xx, x)
		xx.Mod(xx, Curve.Params().N)
	}

	gPx := NewECPoint(gPxx, gPxy)

	return gPx.Equals(gYRef)
}

// AdditiveHomomorphism return a commitment to Q+R
func AdditiveHomomorphism(commQ, commR PolyCommit) PolyCommit {
	if len(commQ.c) != len(commR.c) {
		panic("mismatch degree")
	}

	comm := PolyCommit{
		c: make([]ECPoint, len(commQ.c)),
	}

	for i := range comm.c {
		x, y := Curve.Add(commQ.c[i].x, commQ.c[i].y, commR.c[i].x, commR.c[i].y)
		comm.c[i] = NewECPoint(x, y)
	}

	return comm
}
