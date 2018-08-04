package generation

import (
	"testing"
)

func TestGenPaths(*testing.T) {
	lambda := 16
	parties := make([]int, 0)
	for i := 0; i < 10; i++ {
		parties = append(parties, i)
	}
	B := 8
	paths := GenPaths(lambda, B, parties, 22)
	if len(paths) != (B) {
		panic("Incorrect number of paths")
	}
	if len(paths[2].Clients)-1 != len(paths[4].Nonces) {
		panic("Size incorrect")
	}
	if paths[0].Clients[17] != paths[1].Clients[17] {
		panic("Should match")
	}
	if paths[0].Clients[33] != paths[2].Clients[33] {
		panic("Should match")
	}

}

func TestGenChpts(*testing.T) {
	lambda := 4
	sid := 0
	B := 4
	parties := make([]int, 0)
	for i := 0; i < 4; i++ {
		parties = append(parties, i)
	}
	h := CalculateHB(B)
	nonces, expected := GenChpts(lambda, sid, B, parties, h)
	if len(nonces) < B {
		panic("Should be more nonces!!!")
	}
	hasNonce := false
	for i := 0; i < len(expected[nonces[0][0].L-1][nonces[0][0].Recip]); i++ {
		if expected[nonces[0][0].L-1][nonces[0][0].Recip][i] == nonces[0][0].Nonce {
			hasNonce = true
		}
	}
	if !hasNonce {
		panic("Should have nonce!")
	}
}
