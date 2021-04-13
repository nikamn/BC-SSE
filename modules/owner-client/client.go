package main

import (
	"fmt"
	"os"
    "io/ioutil"	
	
	"github.com/ncw/gmp"
	"github.com/nikamn/BC-SSE/utils/polyring"
	"github.com/nikamn/BC-SSE/utils/polypoint"
	"github.com/nikamn/BC-SSE/utils/intrinsic"
	"github.com/nikamn/BC-SSE/utils/interpolation"

)

// MaxNodes is maximum number of nodes
const MaxNodes = 10

func main() {

	file, err := os.Open("./output/params/Theta")

	if err != nil {
		fmt.Println(err)
	}
	
	var Theta int
	_, err = fmt.Fscanf(file, "%d\n", &Theta)

	// degree of polynomial = theta
	polyOrder := Theta

	// hardcoded large prime p for Polyring
	p := new(gmp.Int)

	b, _ := ioutil.ReadFile("./output/params/primeP") // just pass the file name
    str := string(b) // convert content to a 'string'
    p.SetString(str, 10)

    b, _ = ioutil.ReadFile("./output/params/poly") // just pass the file name
    str = string(b) // convert content to a 'string'
	poly := polyring.FromString(str)

	fmt.Println("\noriginalPoly: ", poly)

	noOfParties := MaxNodes
	
	secretShares := make([]*polypoint.PolyPoint, noOfParties)
	
	X := make([]int32, noOfParties)
	Y := make([]*gmp.Int, noOfParties)

	for i := 0; i < noOfParties; i++ {
		X[i] = int32(i)
		Y[i] = gmp.NewInt(0)
		intrinsic.Load(fmt.Sprintf("./output/secretShares/party%d", i+1), &Y[i])
		secretShares[i] = polypoint.NewPoint(X[i], Y[i], nil)
	}

	var Xs, Ys []*gmp.Int
	for j := 0; j < polyOrder+1; j++ {
		Xs = append(Xs, gmp.NewInt(int64(X[j])))
		Ys = append(Ys, Y[j])
	}

	fmt.Println("\nchosen portion of x array", Xs)
	fmt.Println("\ncorresponding y array", Ys)
	
	// reconstruct the share
	reconstructedPoly, err := interpolation.LagrangeInterpolate(polyOrder, Xs, Ys, p)
	if err != nil {
		panic("can't recover the secret")
	}

	res3 := reconstructedPoly.IsSame(poly)
	fmt.Println("\nreconstructedPoly: ", reconstructedPoly)
	fmt.Println("\nreconstructedPoly is same as original poly : ", res3, "\n")


}