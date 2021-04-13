package polypoint

import (
	"github.com/Nik-U/pbc"
	"github.com/ncw/gmp"
)

// PolyPoint struct
type PolyPoint struct {
	X       int32
	Y       *gmp.Int
	PolyWit *pbc.Element
}

// NewZeroPoint returns a (0,0,nil) polypoint
func NewZeroPoint() *PolyPoint {
	return &PolyPoint{
		X:       0,
		Y:       gmp.NewInt(0),
		PolyWit: nil,
	}
}

// NewPoint returns a polypoint (x,y,w)
func NewPoint(x int32, y *gmp.Int, w *pbc.Element) *PolyPoint {
	return &PolyPoint{
		X:       x,
		Y:       y,
		PolyWit: w,
	}
}
