package supportCode

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"

	"github.com/dedis/kyber/util/random"
)

const NonceNumber = 10000000 // Changes Number of Constants to Gen from

// Swaps the index I with the last one and removes it from the array(Nonce specific)
func RemoveNonce(arr []Nonce, i int) []Nonce {
	arr[len(arr)-1], arr[i] = arr[i], arr[len(arr)-1]
	return arr[:len(arr)-1]
}

// Creates random number based on the number passed
func RandomNumber(n int) int {
	s1 := rand.NewSource((random.Int(big.NewInt(NonceNumber), random.New())).Int64())
	r1 := rand.New(s1)
	return r1.Intn(n)
}

// Creates a random path iid based on the length and options to passed
// Outputs the random path and random corresponding nonces
func CreatePath(options []int, length int) ([]int, []int) {
	if options == nil || length == 0 {
		panic("error: null")
	}
	nonces := make([]int, 0)
	path := make([]int, 0)
	for i := 0; i < length; i++ { //for every length of the path
		nonces = append(nonces, RandomNumber(NonceNumber))       //get a random number
		path = append(path, options[RandomNumber(len(options))]) // use a random index and get that person
	}
	return path, nonces
}

// Creates a tree of Nodes based on the B passed, with length of path depending on lambda
// Then assigns the corresponding parent.
// the first b are the leaves
func CreateTree(parties []int, B int, lambda int) []Node {
	nodes := make([]Node, (2*B)-1) //create length based on B
	for i := 0; i < len(nodes); i++ {
		path, nonces := CreatePath(parties, int(math.Pow(math.Log2(float64(lambda)), 2))) //use path for each node
		nodes[i].Path = path
		nodes[i].Nonces = nonces
		//assigns nodes
	}
	factor := 0 //this makes the right parent index
	pivotNode := CalculatePivotNode(B)
	for i := 0; i < pivotNode; i++ {
		nodes[i].Parent = &nodes[B+factor] //so the parent is the right node, and last is null
		if i%2 == 1 {
			factor += 1
		}
	}
	return nodes

}

// Calculates the node at which to pivot in the tree.
// Increases efficiency
func CalculatePivotNode(B int) int {
	hb := int(math.Ceil(float64(int(math.Log2(float64(B)))+1) / 2))
	total := 0
	for i := 0; i < hb; i++ {
		total += B
		B = B / 2
	}
	return total
}

// testing method
func returnIndex(n *Node, nodes []Node) int {
	idx := -1
	for i := 0; i < len(nodes); i++ {
		if n == &nodes[i] {
			idx = i
		}
	}
	return idx
}

// testing method
func PrintNodeParents(nodes []Node) {
	for i := 0; i < len(nodes); i++ {
		fmt.Print("Parent: ")
		fmt.Print(returnIndex(nodes[i].Parent, nodes))
	}
}

// Takes B and rounds it to the next power of 2. This is needed
// So the binary tree is correctly made
func RoundB(B float64) int {
	lastPowerOf2 := math.Floor(math.Log2(B))
	newB := int(math.Pow(2, lastPowerOf2))
	for math.Log2(float64(newB)) != lastPowerOf2+1 {
		newB += 1
	}
	return newB
}
