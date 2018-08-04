package generation

import (
	"math"
	"math/big"
	"math/rand"

	"github.com/czebos/pi/supportCode"
	"github.com/dedis/kyber/util/random"
)

const ALPHA = 2

// Input: lambda - security parameter, sid - sessionID, B - important variable for security
// parties - all involved peoplesID, selfID - personsID
// Output: array of struct Nonce, which contains l - interval, k - the individual to recieive
// and c - the nonce numner

// Goes through every level of the binary tree, then up until their ID in the array(so that
// symmetrical onions are created, but not that multiple symmetical onions for the same peoplesID
// are created), and creates symmetrical onions for the person itself and the other person. If
// It is the same person, the probability is doubled and it only creats one for itself. This is
// done for every interval(i) and returned as a struct with all the info including a random nonce
// nonce integer
//// TODO: loop through set instead
func GenChpts(lambda int, sid int, B int, parties []int, h int) ([][]supportCode.Nonce, [][][]int) { // do this before onion forming and pass it in?
	noncesList := make([][]supportCode.Nonce, 0)
	for i := 0; i < len(parties); i++ { // creates the 2D array
		nonceList := make([]supportCode.Nonce, 0)
		noncesList = append(noncesList, nonceList)
	}
	expected := make([][][]int, 0)                      //create the expected nonces: intervals x parties x nonces
	for i := 0; i < int(math.Log2(float64(B)))+1; i++ { //creates the 3D array
		firstLayer := make([][]int, 0)
		expected = append(expected, firstLayer)
		for j := 0; j < len(parties); j++ {
			secondLayer := make([]int, 0)
			expected[i] = append(expected[i], secondLayer)
		}
	}
	for i := 0; i < h; i++ { //This goes through every height of the Binary Tree
		//for every interval number
		for j := 0; j < len(parties); j++ {
			//for every party before selfID(so only checks once)
			for k := 0; k <= j; k++ {
				if j == k {
					createDummy := CheckToDummy(B, ALPHA*2, float64(len(parties))) //double the probability
					if createDummy {                                               //same as below, but only creating one
						nonce := supportCode.Nonce{i + 1, k, FUNCTIONFORCODE()}
						noncesList[j] = append(noncesList[j], nonce)
						expected[i][k] = append(expected[i][k], nonce.Nonce) //also add to expected
					}
				} else {
					createDummy := CheckToDummy(B, ALPHA*2, float64(len(parties))) //check with probablity based on code
					if createDummy {                                               //create symmetric nonces
						nonceNumber := FUNCTIONFORCODE() //change for actual function or fake it
						nonce1 := supportCode.Nonce{i + 1, j, nonceNumber}
						nonce2 := supportCode.Nonce{i + 1, k, nonceNumber}

						noncesList[k] = append(noncesList[k], nonce1)
						expected[i][k] = append(expected[i][k], nonceNumber) //also add to expected
						noncesList[j] = append(noncesList[j], nonce2)
						expected[i][j] = append(expected[i][j], nonceNumber) //also add to expected
					}
				}
			}
		}

	}
	return noncesList, expected
}

// Checks to dummy onion with probablity in the paper
//// TODO: create the correct number
func CheckToDummy(B int, A int, N float64) bool { //account for same row later
	s1 := rand.NewSource((random.Int(big.NewInt(supportCode.NonceNumber), random.New())).Int64())
	r1 := rand.New(s1)
	randNumber := r1.Float64()
	if randNumber < float64(B*A)/(N*(math.Log2(float64(B))+1)) { //ask megumi about number
		return true
	} else {
		return false
	}
}

//temp function for megumis random nonce function
func FUNCTIONFORCODE() int {
	return supportCode.RandomNumber(supportCode.NonceNumber)
}

// Calculates HB based on calculations
func CalculateHB(B int) int {
	return int(math.Ceil(float64(int(math.Log2(float64(B)))+1) / 2))
}
