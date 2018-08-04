package batching

import (
	"fmt"
	"math"
	"time"

	"github.com/czebos/pi/crypto"
	"github.com/czebos/pi/generation"
	"github.com/czebos/pi/onionDS"
	"github.com/czebos/pi/supportCode"
)

// Input: lambda - security parameter, corruptionRate - rate at which nodes are corrupted
// e - epsilon, onionsOut - onions to be organized for batching, clients - list of clients
// expected - expected nonces(interval by clientID by nonces)
// Output - the messages, fully decrypted

// Execution is the entire process of the protocol. First it organizes the onionsMade for
// their first recipitant. Then it goes through every round, batches the onions and orgznizes them
// then recreates and returns the message
func Execution(lambda int, corruptionRate float64, e float64, onionsOut [][]onionDS.Onion, clients []crypto.Client, expected [][][]int, B int, R int, isPiButterfly bool) [][]byte {
	onionsIn := sendOnions(onionsOut) // organize onions for recipitant
	fmt.Print("Round total: ")
	fmt.Println(R - 1)
	totalTime := float64(0)
	pivotRound := ((math.Floor(math.Log2(float64(len(clients)))/math.Log2(float64(B))) + 1) * math.Pow(math.Log2(float64(lambda)), 2)) - 1
	counter := 0
	for i := 0; i < int(R)-1; i++ {
		start := time.Now()
		if i == int(pivotRound)+1 {
			totalTime = float64(0)
		}
		fmt.Print("On round: ")
		fmt.Println(i)
		for j := 0; j < len(clients); j++ {
			onionsOut[j] = BatchOnions(lambda, corruptionRate, e, i, onionsIn[j], clients[j], B, expected, isPiButterfly)
			counter++
		}
		onionsIn = sendOnions(onionsOut) //organizes onions for next recipitant
		totalTime += (time.Now().Sub(start).Seconds())
		if i == int(pivotRound) && isPiButterfly {
			fmt.Print("Time of Butterfly: ")
			fmt.Println(float64(totalTime) / float64((int(pivotRound) + 1)))
		}
		if i == R-2 && isPiButterfly {
			fmt.Print("Time of Tree: ")
			fmt.Println(float64(totalTime) / float64((i - int(pivotRound))))
		}
	}

	messages := make([][]byte, 0)
	for j := 0; j < len(clients); j++ {
		for k := 0; k < len(onionsIn[j]); k++ {
			message := onionDS.ProcOnion(onionsIn[j][k], clients[j].Sk).Ciphertext //procs the last layer an returns the message
			if len(message) > 0 {                                                  //if it is not a dummy message
				messages = append(messages, message)
			}
		}
	}
	fmt.Print("Counter: ")
	fmt.Println(counter)
	return messages
}

// This reorganizes the onions. It recieves any order of the onions, and puts
// them in the correct spot of the array for the next recipitant and returns
// the array with all the organized onions
func sendOnions(onionsOut [][]onionDS.Onion) [][]onionDS.Onion {
	onionsIn := make([][]onionDS.Onion, 0)
	for i := 0; i < len(onionsOut); i++ {
		onionsIn = append(onionsIn, make([]onionDS.Onion, 0))
	}
	for i := 0; i < len(onionsOut); i++ {
		for j := 0; j < len(onionsOut[i]); j++ {
			onionsIn[onionsOut[i][j].NextParticipant] = append(onionsIn[onionsOut[i][j].NextParticipant], onionsOut[i][j])
		}
	}
	return onionsIn
}

// Calculates B based on caluclations
func CalculateB(clientLength int, lambda int, corruptionRate float64) int {
	A := math.Sqrt(float64(clientLength) * math.Pow(math.Log2(float64(lambda)), 3))
	temp := (A / float64((1 - corruptionRate))) * ((math.Log2(A) + 1) / float64(2))
	B := supportCode.RoundB(temp)
	return B
}

// Calculates the number of rounds based on calculations
func CalculateRoundTotal(lambda int, B int) int {
	return int(math.Pow(math.Log2(float64(lambda)), 2) * (float64(generation.CalculateHB(B)) + 1))
}
