package supportCode

import (
	"math"
	"testing"
)

func TestTree(*testing.T) {
	lambda := 8
	parties := make([]int, 0)
	for i := 0; i < 6; i++ {
		parties = append(parties, i)
	}
	B := 8
	nodes := CreateTree(parties, B, lambda)
	if nodes[0].Parent != &nodes[8] || nodes[8].Parent != &nodes[12] || nodes[9].Parent != &nodes[12] || nodes[14].Parent != nil || nodes[13].Parent != nil {
		panic("Wrong Parent")
	}
	currentNode := nodes[3]
	path := make([]int, 0)
	for currentNode.Parent != nil {
		path = append(path, currentNode.Path...)
		currentNode = *currentNode.Parent
	}
	path = append(path, currentNode.Path...)
	if float64(len(path)) != 3*math.Pow(math.Log2(float64(lambda)), 2) {
		panic("Path length incorrect!")
	}

}

func TestRemoveNonce(*testing.T) {
	nonces := make([]Nonce, 0)
	for i := 0; i < 3; i++ {
		nonces = append(nonces, Nonce{0, 0, i})
	}
	nonces = RemoveNonce(nonces, 1)
	if nonces[0].Nonce != 0 || nonces[1].Nonce != 2 {
		panic("Swap incorrect")
	}
}

func TestRoundB(*testing.T) {
	if RoundB(53) != 64 {
		panic("Wrong number")
	}
	if RoundB(3) != 4 {
		panic("Wrong number")
	}
	if RoundB(127) != 128 {
		panic("Wrong number")
	}
}
