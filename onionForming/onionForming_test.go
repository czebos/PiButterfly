package onionForming

import (
	"fmt"
	"testing"

	"github.com/czebos/pi/batching"
	"github.com/czebos/pi/crypto"
	"github.com/czebos/pi/generation"
	"github.com/czebos/pi/onionDS"
	"github.com/czebos/pi/supportCode"
)

func TestMBOsDOs(*testing.T) {
	lambda := 8
	sks := make([]crypto.Scalar, 0)
	pks := make([]supportCode.Point, 0)
	parties := make([]int, 0)
	msg := make([]byte, 3)
	for i := 0; i < 16; i++ {
		pk, sk := onionDS.GenKey()
		parties = append(parties, i)
		pks = append(pks, pk)
		sks = append(sks, sk)
	}
	B := batching.CalculateB(len(parties), lambda, 0)
	keys := onionDS.Init(parties, pks)
	h := generation.CalculateHB(B)
	nonces, _ := generation.GenChpts(lambda, 0, B, parties, h)
	fmt.Println(len(nonces[0]))
	fmt.Println(B)
	recip := 1
	onions, remainingNonces := FormMBOs(lambda, B, keys, parties, nonces[0], msg, recip)
	onions2 := FormDOs(lambda, B, keys, parties, remainingNonces)
	if len(nonces[0]) != len(remainingNonces)+int(B) {
		panic("Incorrect Amount of Nonces Left")
	}
	if len(onions) != int(B) {
		panic("Incorrect Amount of Onions")
	}
	if len(onions2) != len(remainingNonces) {
		panic("Shouldve used all nonces!")
	}
}

func TestFormOnions(*testing.T) {
	lambda := 8
	sks := make([]crypto.Scalar, 0)
	pks := make([]supportCode.Point, 0)
	parties := make([]int, 0)
	msg := make([]byte, 3)
	for i := 0; i < 16; i++ {
		pk, sk := onionDS.GenKey()
		parties = append(parties, i)
		pks = append(pks, pk)
		sks = append(sks, sk)
	}
	recip := 1
	keys := onionDS.Init(parties, pks)
	B := batching.CalculateB(len(parties), lambda, 0)
	fmt.Println(B)
	h := generation.CalculateHB(B)
	nonceList, _ := generation.GenChpts(lambda, 0, B, parties, h)
	fmt.Println()
	for i := 0; i < 10; i++ {
		fmt.Print("Nonce ")
		fmt.Print(i)
		fmt.Print(": ")
		fmt.Println(len(nonceList[i]))
	}
	onions := FormOnions(lambda, 0, 0, keys, parties, msg, recip, nonceList[3])
	if len(onions) <= 0 {
		panic("Onions were not formed")
	}

}
