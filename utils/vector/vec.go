package vector

import "github.com/ncw/gmp"

// Vector struct
type Vector struct {
	v []*gmp.Int
}

// New returns a new vector
func New(len int) Vector {
	if len < 0 {
		panic("len < 0")
	}

	v := Vector{
		make([]*gmp.Int, len),
	}

	for i := range v.v {
		v.v[i] = gmp.NewInt(0)
	}

	return v
}

// FromInt64 returns a vector from int64 array
func FromInt64(elems ...int64) Vector {
	v := New(len(elems))

	for i, elem := range elems {
		v.v[i].SetInt64(int64(elem))
	}

	return v
}

// GetPtr returns a pointer to a vector
func (vec Vector) GetPtr() []*gmp.Int {
	return vec.v
}
