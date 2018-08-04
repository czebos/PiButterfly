package piButterfly

import (
	"math"

	"github.com/czebos/pi/supportCode"
)

// Input: lambda - security parameter, B - important variable for security
// parties - all involved peoplesID, reci - the person expexting the message(Bob)
// Output: an array of paths struct that contains array on IDs and Nonces

// Creates a path of the intermediatery hops with a nonce number in parallel.
// Makes a path for every leaf of the binary tree. The first B nodes are the leaves,
//and creats a path up the tree to the root node.

// MODIFIED TO ACCOMIDATE PIBUTTERFLY
func GenTree(lambda int, B int, parties []int, recip int) []supportCode.Path {
	nodes := CreateTree(parties, B, lambda) //creates the tree with the right parents shuffle
	paths := make([]supportCode.Path, 0)
	for i := 0; i < int(B); i++ {
		path := make([]int, 0)
		nonces := make([]int, 0)
		currentNode := nodes[i]
		for currentNode.Parent != nil { //iterate through parents and add their nonces/paths
			path = append(path, currentNode.Path...)
			nonces = append(nonces, currentNode.Nonces...)
			currentNode = *currentNode.Parent
		}
		path = append(path, currentNode.Path...) //finalize for the last node
		path = append(path, recip)
		nonces = append(nonces, currentNode.Nonces...)
		paths = append(paths, supportCode.Path{path, nonces}) //add the path
	}
	return paths
}

// Creates a tree of Nodes based on the B passed, with length of path depending on lambda
// Then assigns the corresponding parent.
// the first b are the leaves

// MODIFIED TO ACCOMIDATE PIBUTTERFLY
func CreateTree(parties []int, B int, lambda int) []supportCode.Node {
	nodes := make([]supportCode.Node, (2*B)-1) //create length based on B
	for i := 0; i < len(nodes); i++ {
		path, nonces := supportCode.CreatePath(parties, supportCode.RoundB(math.Pow(math.Log2(float64(lambda)), 2))) //use path for each node
		nodes[i].Path = path
		nodes[i].Nonces = nonces
		//assigns nodes
	}
	factor := 0 //this makes the right parent index
	pivotNode := supportCode.CalculatePivotNode(B)
	for i := 0; i < pivotNode; i++ {
		nodes[i].Parent = &nodes[B+factor] //so the parent is the right node, and last is null
		if i%2 == 1 {
			factor += 1
		}
	}
	return nodes
}
