package batching

import (
	"math"

	"github.com/czebos/pi/crypto"
	"github.com/czebos/pi/onionDS"
	"github.com/czebos/pi/supportCode"
)

const FACTOR = 100000

// Input: lambda - security parameter, corruptionRate - rate at which nodes are corrupted
// e - epsilon, round - round number to batch onions, onionsIn, onions to be batched for the
// individual, client - the client batching the onions, B -important security constant,
// expected - expected nonces(interval by clientID by nonces)
// Output - the onions to be batched out

// Batch onions takes in these parameters and batches each onion. If it is a diagonstic round,
// keeps track of the nonces it should expect. Then, it computes the threshold for Security,
// and if it doesnt get as many nonces as it expects, security was compromised
func BatchOnions(lambda int, corruptionRate float64, e float64, round int, onionsIn []onionDS.Onion, client crypto.Client, B int, expected [][][]int, isPiButterfly bool) []onionDS.Onion {
	onionsOut := make([]onionDS.Onion, 0)
	nonces := make([]int, 0)
	prevNonces := make(map[int]bool)
	var diagonosticRound bool
	if !isPiButterfly {
		diagonosticRound = checkDiagnosticRound(lambda, round)
	} else {
		diagonosticRound = (round+1)%B == 0
	}
	for len(onionsIn) != 0 { // for all the onions to be batched

		onions2, onion := pickOnion(onionsIn) //picks and removes an onion
		onionsIn = onions2
		prevNonce := onion.Nonce
		onion = onionDS.ProcOnion(onion, client.Sk)

		if !prevNonces[prevNonce] || round <= B { // if onionNonce isnt the same as before
			onionsOut = append(onionsOut, onion)
		}

		prevNonces[prevNonce] = true
		if diagonosticRound { // if it is a diagonstic round
			nonces = append(nonces, prevNonce) // append nonces
		}
	}
	if diagonosticRound {
		var interval int
		if !isPiButterfly {
			interval = calculateInterval(lambda, round)
		} else {
			interval = (round + 1) / B
		}
		threshold := calculateThreshold(lambda, round, e, corruptionRate, interval, B)
		if float64(len(setDifference(expected[interval-1][client.Id], nonces))) > threshold {
			panic("Abort Protocol, security compromised")
		}
	}
	return onionsOut
}

// This picks an onion from the list of onions, returns it,
// and removes it from the list
func pickOnion(onions []onionDS.Onion) ([]onionDS.Onion, onionDS.Onion) {
	idx := supportCode.RandomNumber(len(onions))
	onion := onions[idx]
	onions[len(onions)-1], onions[idx] = onions[idx], onions[len(onions)-1]
	return onions[:len(onions)-1], onion

}

// This calulates the magnitude of the set difference
// -----> |a\b|
func setDifference(a []int, b []int) []int {
	m := make(map[int]bool)
	diff := make([]int, 0)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// Checks diagonstic round based on calculations
func checkDiagnosticRound(lambda int, round int) bool {
	return (round+1)%int(math.Pow(math.Log2(float64(lambda)), 2)) == 0
}

// Calculates interval based on calculations
func calculateInterval(lambda int, round int) int {
	return (round + 1) / int(math.Pow(math.Log2(float64(lambda)), 2))
}

// Calculates threshold based on interval and calculations
func calculateThreshold(lambda int, round int, e float64, corruptionRate float64, interval int, B int) float64 {
	threshold := float64(0)
	if interval == 1 {
		threshold = math.Pow(math.Log2(float64(lambda)), 2)*corruptionRate*9*math.Pow(1-e, 5) + FACTOR
	} else {
		threshold = ((1 - e) * math.Sqrt(float64(B))) + FACTOR
	}
	return threshold
}

// Convertes bytes to String
func BytesToString(data []byte) string {
	return string(data[:])
}
