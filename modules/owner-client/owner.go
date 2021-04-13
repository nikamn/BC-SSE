package main

import (
	"fmt"
	"math/rand"
	
	"github.com/ncw/gmp"
	"github.com/nikamn/BC-SSE/utils/commitment"
	"github.com/nikamn/BC-SSE/utils/polyring"
	"github.com/nikamn/BC-SSE/utils/polypoint"
	"github.com/nikamn/BC-SSE/utils/intrinsic"
	"github.com/nikamn/BC-SSE/utils/basic"
)

// MaxNodes is maximum number of nodes
const MaxNodes = 10

func main() {

	intrinsic.CreateDirIfNotExist("./output/params")
	intrinsic.CreateDirIfNotExist("./output/secretShares")

	/* User input theta */
	var Theta int
	fmt.Printf("Give number of nodes(theta). %d is the maximum number of nodes\n", MaxNodes)
	fmt.Println(MaxNodes/2, "< (theta) <", MaxNodes)
	fmt.Scanf("%d", &Theta)
	
	basic.CreateFile("./output/params/Theta", fmt.Sprintf("%d", Theta))
	/* User input theta taken */

	// degree of polynomial = theta
	polyOrder := Theta
	
	// hardcoded large prime p for Polyring
	p := new(gmp.Int)
	p.SetString("57896044618658097711785492504343953926634992332820282019728792006155588075521", 10)
	
	basic.CreateFile("./output/params/primeP", p.String())

	// random source seed
	rnd := rand.New(rand.NewSource(99))

	c := commitment.DLPolyCommit{}
	c.SetupFix2(polyOrder, "218882428714186575617")

	// Sample a Poly
	poly, _ := polyring.NewRand(polyOrder, rnd, p)
	basic.CreateFile("./output/params/poly", poly.String())
	
	C := c.NewG1()
	// PolyCommit
	c.Commit(C, poly)
	basic.CreateFile("./output/params/commitment", poly.String())

	// secret sharing with parties
	fmt.Printf("\nSharing secret with %d parties\n\n", MaxNodes)
	
	noOfParties := MaxNodes
	
	secretShares := make([]*polypoint.PolyPoint, noOfParties)
	
	xs := make([]int32, noOfParties)
	ys := make([]*gmp.Int, noOfParties)
	
	for i := 0; i < noOfParties; i++ {
		xs[i] = int32(i)
		ys[i] = gmp.NewInt(0)
		w := c.NewG1()
		poly.EvalMod(gmp.NewInt(int64(xs[i])), p, ys[i])
		c.CreateWitness(w, poly, gmp.NewInt(int64(xs[i])))
		secretShares[i] = polypoint.NewPoint(xs[i], ys[i], w)
		//fmt.Println("Party witness", i+1, secretShares[i].PolyWit)
		intrinsic.Save(fmt.Sprintf("./output/secretShares/party%d", i+1), ys[i])
	}

	fmt.Println("\n\nx value array\t", xs)
	fmt.Println("\ny value array\t", ys)

}