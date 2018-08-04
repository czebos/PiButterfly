package batching

import (
	"fmt"
	"math"
	"testing"

	"github.com/czebos/pi/crypto"
	"github.com/czebos/pi/generation"
	"github.com/czebos/pi/onionDS"
	"github.com/czebos/pi/onionForming"
	"github.com/czebos/pi/supportCode"
)

func TestSetDifference(*testing.T) {
	set1 := make([]int, 0)
	set2 := make([]int, 0)
	set1 = append(set1, 2, 3, 4, 5, 6)
	set2 = append(set2, 2, 4, 5, 8, 9)
}

func TestExecution(*testing.T) {
	lambda := 8
	corruptionRate := float64(0)
	e := .00005
	clients := make([]crypto.Client, 0)
	pks := make([]supportCode.Point, 0)
	parties := make([]int, 0)
	for i := 0; i < 16; i++ {
		pk, sk := onionDS.GenKey()
		clients = append(clients, crypto.Client{i, sk})
		pks = append(pks, pk)
		parties = append(parties, i)
	}
	keys := onionDS.Init(parties, pks)
	fmt.Println("Keys and people initiated")
	A := math.Sqrt(float64(len(clients)) * math.Pow(math.Log2(float64(lambda)), 3))
	temp := (A / float64((1 - corruptionRate))) * (math.Log2(A+1) / float64(2))
	B := supportCode.RoundB(temp)
	h := generation.CalculateHB(B)
	nonces, expected := generation.GenChpts(lambda, 0, B, parties, h)
	onions := make([][]onionDS.Onion, 0)
	for i := 0; i < 16; i++ {
		onions = append(onions, make([]onionDS.Onion, 0))
		msg := []byte("Security is great!")
		j := 0
		onion := onionForming.FormOnions(lambda, corruptionRate, 0, keys, parties, msg, j, nonces[i])
		onions[i] = append(onions[i], onion...)
		fmt.Print("Person ")
		fmt.Print(i)
		fmt.Println("'s onions formed")
	}
	R := CalculateRoundTotal(lambda, B)
	msgs := Execution(lambda, corruptionRate, e, onions, clients, expected, B, R, false)
	if len(msgs) <= 0 {
		panic("there should be msgs!")
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("Securely Transfered and Decrypted Message:", BytesToString(msgs[0]))
}
