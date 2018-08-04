package piButterfly

import (
	"math"

	"github.com/czebos/pi/onionDS"
	"github.com/czebos/pi/supportCode"
)

// Input: lambda - security parameter, κ - corruption rate, sid - session ID,
// keys - mapping from peoplesID's to their public key, parties - all involved peoplesID,
// message - the bytes of the message, recipitant - person reciving the message)
// Output: the onions created from MBOS and FormDOs

// Caluclates B, an important variable for security, then uses the paramters to GenChpts,
// and use that to FormMBOs and FormDOs. Appends all onions to a list and returns it

// MODIFIED TO ACCOMIDATE PIBUTTERFLY
func FormOnions(lambda int, k float64, sid int, keys map[int]supportCode.Point, parties []int, msg []byte, recip int, nonces []supportCode.Nonce, B int, rounds int) []onionDS.Onion {
	onions := make([]onionDS.Onion, 0)
	paths := GenPaths(lambda, parties, B)
	paths2 := GenTree(lambda, B, parties, recip)
	for i := 0; i < len(paths); i++ {
		paths[i].Clients = append(paths[i].Clients, paths2[i].Clients...)
		paths[i].Nonces = append(paths[i].Nonces, paths2[i].Nonces...)
	}
	onionsToBeAdded, noncesLeft := FormMBOs(lambda, B, keys, parties, nonces, msg, recip, paths)

	onions = append(onions, onionsToBeAdded...)
	onionsToBeAdded = FormDOs(lambda, B, keys, parties, noncesLeft, rounds) //make sure nonces are used

	onions = append(onions, onionsToBeAdded...)
	return onions
}

// Input: lambda - security parameter, B - important constant based on security
// keys - mapping from peoplesID's to their public key, parties all involved peoplesID
// nonces - the struct created to designate a specific location and time for onions2
// msg - the message in byte, recip - the person reciving the onions
// Output: the message bearing onions designated for the recipitant

// Creates the paths through the binary tree, and uses all of the paths(Based on B)
// to change one specific location(based on the nonces), so that a person will expect
//to recieve an onion on a given round

// MODIFIED TO ACCOMIDATE PIBUTTERFLY
func FormMBOs(lambda int, B int, keys map[int]supportCode.Point, parties []int, nonces []supportCode.Nonce, msg []byte, recip int, paths []supportCode.Path) ([]onionDS.Onion, []supportCode.Nonce) {
	onions := make([]onionDS.Onion, 0)
	for i := range paths {
		path := paths[i]
		if len(nonces) == 0 { //abort if no nonces left
			panic("no nonces left")
		} else {
			randInt := supportCode.RandomNumber(len(nonces)) //grab a random spot
			nonce := nonces[randInt]
			nonces = supportCode.RemoveNonce(nonces, randInt)

			R := (nonce.L+int(CalculateH(parties, B)))*supportCode.RoundB(math.Pow(math.Log2(float64(lambda)), 2)) - 1 //calculate spot to change
			path.Clients[int(R)] = nonce.Recip
			path.Nonces[int(R)] = nonce.Nonce

			onion := onionDS.FormOnion(msg, keys, path.Clients, path.Nonces) //form onions
			onions = append(onions, onion)
		}
	}
	return onions, nonces
}

// Input:  lambda - security parameter, B - important constant based on security
// keys - mapping from peoplesID's to their public key, parties all involved peoplesID
// nonces - the struct created to designate a specific location and time for onions2
// this is the left over ones after MBOS
// Output: the dummy onions used for the protocol

// This runs similarly to FormMBOs. It used the rest of the Nonces to create
// an onion based on those, but it also makes it so that someone should expect
// to see am onion at a specific spot. It then returns those onions
// NOTE: There are 2 different kind of nonces going on here,
// The random number kind, and the set of instructions kind

// MODIFIED TO ACCOMIDATE PIBUTTERFLY
func FormDOs(lambda int, B int, keys map[int]supportCode.Point, parties []int, nonces []supportCode.Nonce, R int) []onionDS.Onion {
	onions := make([]onionDS.Onion, 0)
	for len(nonces) > 0 { //use the rest of the nonces
		randInt := supportCode.RandomNumber(len(nonces))
		nonce := nonces[randInt]
		nonces = supportCode.RemoveNonce(nonces, randInt) //grab a random nonce

		clientPath, noncePath := supportCode.CreatePath(parties, int(R)+1)
		noncePath = noncePath[:len(noncePath)-1] // the recipitant doesnt need a nonce

		changeSpot := int(float64(nonce.L)*math.Pow(math.Log2(float64(lambda)), 2)) - 1 //calculate spot to change
		clientPath[changeSpot] = nonce.Recip
		noncePath[changeSpot] = nonce.Nonce
		onion := onionDS.FormOnion(make([]byte, 0), keys, clientPath, noncePath) //form onions
		onions = append(onions, onion)
	}
	return onions
}
