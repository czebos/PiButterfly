package supportCode

import (
	"github.com/dedis/kyber"
	"github.com/dedis/kyber/group/edwards25519"
	"github.com/dedis/kyber/xof/blake2xb"
)

var SUITE = edwards25519.NewBlakeSHA256Ed25519WithRand(blake2xb.New(nil))
var refPt = SUITE.Point()

type Nonce struct {
	L     int
	Recip int
	Nonce int
}

type Node struct {
	Parent *Node
	Path   []int
	Nonces []int
}

type Path struct {
	Clients []int
	Nonces  []int
}

type Point struct {
	P kyber.Point
}
