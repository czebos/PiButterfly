package piButterfly

import (
	"testing"
)

func TestNumbers(*testing.T) {
	numberArray := make([]int, 0)
	numberArray = append(numberArray, 1, 10, 19)
	numberArray2 := calculateNumbers(0, 3, 1, 26)
	for i := 0; i < len(numberArray2); i++ {
		if numberArray[i] != numberArray2[i] {
			panic("Wrong Number")
		}
	}
}

func TestGenPaths(*testing.T) {
	parties := make([]int, 0)
	for i := 0; i < 16; i++ {
		parties = append(parties, i)

	}
	paths := GenPaths(4, parties, 9)
	if paths[0].Clients[0]%4 != paths[0].Clients[1]%4 && paths[0].Clients[1]%4 != paths[0].Clients[2]%4 && paths[0].Clients[3]%4 != paths[0].Clients[2]%4 {
		panic("remainder should be the same!")
	}
}
