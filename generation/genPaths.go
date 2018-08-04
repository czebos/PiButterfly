package generation

import (
	"github.com/czebos/pi/supportCode"
)

// Input: lambda - security parameter, B - important variable for security
// parties - all involved peoplesID, reci - the person expexting the message(Bob)
// Output: an array of paths struct that contains array on IDs and Nonces

// Creates a path of the intermediatery hops with a nonce number in parallel.
// Makes a path for every leaf of the binary tree. The first B nodes are the leaves,
//and creats a path up the tree to the root node.
func GenPaths(lambda int, B int, parties []int, recip int) []supportCode.Path {
	nodes := supportCode.CreateTree(parties, B, lambda) //creates the tree with the right parents shuffle
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
