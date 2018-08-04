package piButterfly

import (
	"math"
	"strconv"
	"strings"

	"github.com/czebos/pi/generation"
	"github.com/czebos/pi/supportCode"
)

// Input: lambda - security parameter, parties - list of parties,
// and B - importat variable
// Output : the paths for the PiButterfly

// This function calculates the appropriate values for d, h d,
// then choses people based on the people in their subnet, randomly
func GenPaths(lambda int, parties []int, B int) []supportCode.Path {
	d := float64(supportCode.RoundB(math.Pow(math.Log2(float64(lambda)), 2))) // branching factor
	h := CalculateH(parties, int(d))                                          //change of base to d
	D := math.Pow(d, h)
	chosenParties := make([]int, 0) // not all participants can be servers
	for i := 0; i < int(D); i++ {
		chosenParties = append(chosenParties, i) // choose the first D parties
	}
	paths := make([]supportCode.Path, 0)
	for i := 0; i < int(B); i++ {
		nonce := generation.FUNCTIONFORCODE() // random nonce
		firstParticipant := chosenParties[supportCode.RandomNumber(int(D))]
		path := make([]int, 0)
		path = append(path, firstParticipant)
		nonces := make([]int, 0)
		nonces = append(nonces, nonce)
		path, nonces = pickPath(0, chosenParties, int(h-1), firstParticipant, int(d), int(d), nonce) //pick the path from the first participant
		paths = append(paths, supportCode.Path{path, nonces})
	}
	return paths
}

// This function takes in the interval, the base, the number, and the max
// number to calculated the different possible subnets. This depends on
// the base of the number and what base to change(works with butterfly)
func calculateNumbers(interval int, base int, number int, maxNumber int) []int {
	stringNumber := strconv.FormatInt(int64(number), base)
	stringNumberMax := strconv.FormatInt(int64(maxNumber), base)
	arrayString := strings.Split(stringNumber, "")
	arrayStringMax := strings.Split(stringNumberMax, "") // split the numbers into a new base into array
	arrayNumber := make([]int, 0)
	for i := 0; i < len(arrayStringMax)-len(arrayString); i++ { // add zeros to the front
		arrayNumber = append(arrayNumber, 0)
	}
	for i := 0; i < len(arrayString); i++ { // create array of numbers
		str, _ := strconv.Atoi(arrayString[i])
		arrayNumber = append(arrayNumber, str)
	}
	numbers := make([]int, 0) // possible canidates
	for i := 0; i < base; i++ {
		copyOf := make([]int, 0)
		for j := 0; j < len(arrayNumber); j++ { //copy the array
			copyOf = append(copyOf, arrayNumber[j])
		}
		copyOf[interval] = i
		base10Number := changeToBase10(copyOf, base) // change back to base 10, and add
		numbers = append(numbers, base10Number)
	}
	return numbers
}

// This function takes the array and changes the number to base10
func changeToBase10(number []int, base int) int {
	total := 0
	for index := 0; index < len(number); index++ {
		total += number[len(number)-(index+1)] * int(math.Pow(float64(base), float64(index)))
	}
	return total
}

// Input: interval - which base to change, parties - parties to mix from
// lastInterval - where it should stop, spot - location of the number,
// base - the base, d the divisor, firstNonce for shortening runtime
func pickPath(interval int, parties []int, lastInterval int, spot int, base int, d int, firstNonce int) ([]int, []int) {
	subNet := calculateNumbers(interval, base, spot, parties[len(parties)-1]) // make subnet
	people := make([]int, 0)
	nonces := make([]int, 0)
	if interval == 0 { // add first peerson to shorten runttime
		peopleFirst, noncesFirst := supportCode.CreatePath(subNet, d)
		people = append(people, peopleFirst...)
		nonces = append(nonces, noncesFirst...)
	} else {
		people, nonces = supportCode.CreatePath(subNet, d)
	}
	if interval == lastInterval { // base case
		return people, nonces
	} else {
		people2, nonces2 := pickPath(interval+1, parties, lastInterval, len(people), base, d, 0) //recurse with the next path, but one interval lower
		people = append(people, people2...)
		nonces = append(nonces, nonces2...)
		return people, nonces
	}

}

// Calculates H based on the function
func CalculateH(parties []int, d int) float64 {
	return math.Floor(math.Log2(float64(len(parties))) / math.Log2(float64(d)))
}
