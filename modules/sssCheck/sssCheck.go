package main

import (
	"fmt"
	//"log"
	"math/rand"
	//"encoding/json"
	//"os"

	"github.com/ncw/gmp"
	"github.com/nikamn/BC-SSE/utils/commitment"
	"github.com/nikamn/BC-SSE/utils/interpolation"
	"github.com/nikamn/BC-SSE/utils/polypoint"
	"github.com/nikamn/BC-SSE/utils/polyring"
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
	//fmt.Println("\nc: ",c)

	// Sample a Poly and an x
	poly, _ := polyring.NewRand(polyOrder, rnd, p)
	basic.CreateFile("./output/params/poly", poly.String())
	fmt.Println("\npoly :", poly)
	
	fmt.Println("\n\n//---------Verifying polynomial commitment and evaluation at random point on polynomial---------")
	C := c.NewG1()
	// Test PolyCommit
	c.Commit(C, poly)
	basic.CreateFile("./output/params/commitment", poly.String())
	fmt.Println("\nCommit : ", C)

	// verify poly
	res := c.VerifyPoly(C, poly)
	fmt.Println("\nVerifyPoly : ", res)

	// Test EvalCommit
	// x is a random point
	x := new(gmp.Int)
	x.Rand(rnd, p)
	polyOfX := new(gmp.Int)
	w := c.NewG1()
	c.PolyEval(polyOfX, poly, x)
	c.CreateWitness(w, poly, x)
	fmt.Println("\npolyOfX: ", polyOfX)
	fmt.Println("witness: ", w)

	// verify polyOfX
	res2 := c.VerifyEval(C, x, polyOfX, w)
	fmt.Println("\nVerifyEval : ", res2)
	fmt.Println("\n------------Verification of polynomial complete---------------//\n")
		
	// secret sharing with parties
	fmt.Printf("\n\nSharing secret with %d parties\n\n", MaxNodes)
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
		fmt.Println("Party Secret", i+1, secretShares[i].PolyWit)
		intrinsic.Save(fmt.Sprintf("./output/secretShares/party%d", i+1), secretShares[i])
		/*json, err := json.Marshal(secretShares[i])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(json))*/
	}
	fmt.Println("\n\nx value array\t", xs)
	fmt.Println("\ny value array\t", ys)
	//sharesReceived := make(map[int32]*gmp.Int)
	var Xs, Ys []*gmp.Int
	for j := 0; j < polyOrder+1; j++ {
		Xs = append(Xs, gmp.NewInt(int64(xs[j])))
		Ys = append(Ys, ys[j])
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

// reconstruction
	//originalPoly, _ := polyring.NewRand(polyOrder, rnd, p)

	// Test EvalArray
	/*Xs := make([]*gmp.Int, polyOrder+1)
	Ys := make([]*gmp.Int, polyOrder+1)
	polyring.VecInit(Xs)
	polyring.VecInit(Ys)
	polyring.VecRand(Xs, MaxNodes, rnd)

	originalPoly.EvalModArray(Xs, p, Ys)

	fmt.Println("\nFinished eval")
	fmt.Println("Starting interpolation")

	reconstructedPoly, _ := interpolation.LagrangeInterpolate(polyOrder, Xs, Ys, p)
	res3 := reconstructedPoly.IsSame(originalPoly)
	fmt.Println("reconstructedPoly: ", reconstructedPoly)
	fmt.Println("\nreconstructedPoly is same as originalPoly : ", res3)*/