package main

import (
	"fmt"
	"math"
	"time"

	"github.com/czebos/pi/batching"
	"github.com/czebos/pi/crypto"
	"github.com/czebos/pi/generation"
	"github.com/czebos/pi/onionDS"
	"github.com/czebos/pi/onionForming"
	"github.com/czebos/pi/piButterfly"
	"github.com/czebos/pi/supportCode"
)

var _lambda = 8
var _corruptionRate = float64(0)
var _e = .00005
var _participantSize = 16
var _msg = []byte("Security is great!")

func main() {
	StartProtocol(_lambda, _corruptionRate, _e, _participantSize, _msg)
	fmt.Println()
	StartButterFly(_lambda, _corruptionRate, _e, _participantSize, _msg)
}

func StartProtocol(lambda int, corruptionRate float64, e float64, participantSize int, msg []byte) float64 {
	start := time.Now()

	clients := make([]crypto.Client, 0)
	pks := make([]supportCode.Point, 0)
	parties := make([]int, 0)

	fmt.Println("Setting up clients")
	for i := 0; i < participantSize; i++ {
		pk, sk := onionDS.GenKey()
		clients = append(clients, crypto.Client{i, sk})
		pks = append(pks, pk)
		parties = append(parties, i)
	}
	keys := onionDS.Init(parties, pks)

	fmt.Println("Creating Onions")
	formStart := time.Now()
	B := batching.CalculateB(participantSize, lambda, corruptionRate)
	h := generation.CalculateHB(B)
	nonces, expected := generation.GenChpts(lambda, 0, B, parties, h)
	onions := make([][]onionDS.Onion, 0)
	for i := 0; i < participantSize; i++ {
		recip := supportCode.RandomNumber(participantSize)
		onions = append(onions, make([]onionDS.Onion, 0))
		onion := onionForming.FormOnions(lambda, corruptionRate, 0, keys, parties, msg, recip, nonces[i])
		onions[i] = append(onions[i], onion...)
		fmt.Print("Person ")
		fmt.Print(i)
		fmt.Println("'s onions formed")
	}
	fmt.Print("OnionForming time: ")
	fmt.Println(time.Now().Sub(formStart).Seconds())
	R := batching.CalculateRoundTotal(lambda, B)
	execStart := time.Now()
	msgs := batching.Execution(lambda, corruptionRate, e, onions, clients, expected, B, R, false)
	fmt.Print("Execution time: ")
	fmt.Println(time.Now().Sub(execStart).Seconds())

	fmt.Println()
	fmt.Println()
	fmt.Print("Securely Transfered and Decrypted Message: \"")
	fmt.Print(batching.BytesToString(msgs[0]))
	fmt.Println("\"")
	fmt.Println()
	fmt.Print("Time taken: ")
	fmt.Print(time.Now().Sub(start))

	return (time.Now().Sub(start).Seconds())

}

func StartButterFly(lambda int, corruptionRate float64, e float64, participantSize int, msg []byte) float64 {
	start := time.Now()

	clients := make([]crypto.Client, 0)
	pks := make([]supportCode.Point, 0)
	parties := make([]int, 0)

	fmt.Println("Setting up clients")
	for i := 0; i < participantSize; i++ {
		pk, sk := onionDS.GenKey()
		clients = append(clients, crypto.Client{i, sk})
		pks = append(pks, pk)
		parties = append(parties, i)
	}
	keys := onionDS.Init(parties, pks)

	fmt.Println("Creating Onions")
	A := int(math.Pow(math.Log2(float64(lambda)), 2))
	A = supportCode.RoundB(float64(A))
	B := (int(piButterfly.CalculateH(parties, A)) + int(math.Log2(float64(A)))) * A
	RoundTotal := (int(piButterfly.CalculateH(parties, A)) + generation.CalculateHB(B) + 1) * A
	fmt.Print("Round Total: ")
	fmt.Println(RoundTotal - 1)
	formStart := time.Now()
	nonces, expected := generation.GenChpts(lambda, 0, B, parties, (int(piButterfly.CalculateH(parties, A)) + generation.CalculateHB(B)))
	onions := make([][]onionDS.Onion, 0)
	for i := 0; i < participantSize; i++ {
		recip := supportCode.RandomNumber(participantSize)
		onions = append(onions, make([]onionDS.Onion, 0))
		onion := piButterfly.FormOnions(lambda, corruptionRate, 0, keys, parties, msg, recip, nonces[i], B, RoundTotal)
		onions[i] = append(onions[i], onion...)
		fmt.Print("Person ")
		fmt.Print(i)
		fmt.Println("'s onions formed")
	}
	fmt.Print("OnionForming time: ")
	fmt.Println(time.Now().Sub(formStart).Seconds())
	execStart := time.Now()
	msgs := batching.Execution(lambda, corruptionRate, e, onions, clients, expected, B, RoundTotal, true)
	fmt.Print("Execution time: ")
	fmt.Println(time.Now().Sub(execStart).Seconds())

	fmt.Println()
	fmt.Println()
	fmt.Print("Securely Transfered and Decrypted Message: \"")
	fmt.Print(batching.BytesToString(msgs[0]))
	fmt.Println("\"")
	fmt.Println()
	fmt.Print("Time taken: ")
	fmt.Print(time.Now().Sub(start))

	return (time.Now().Sub(start).Seconds())

}
