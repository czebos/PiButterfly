package onionDS

import (
	"testing"

	"github.com/czebos/pi/crypto"
	"github.com/czebos/pi/supportCode"
)

func TestInit(*testing.T) {
	m := make([]int, 0)
	m = append(m, int(7))
	l := make([]supportCode.Point, 0)
	pk, _ := GenKey()
	l = append(l, pk)
	diction := Init(m, l)
	if diction[7] != pk {
		panic("error: wrong number!!")
	}

}

func TestWrapandProc(*testing.T) {
	pk := make([]supportCode.Point, 0)
	sk := make([]crypto.Scalar, 0)
	pi := make([]int, 0)
	ni := make([]int, 0)
	for i := 0; i < 3; i++ {
		pki, ski := GenKey()
		pk = append(pk, pki)
		sk = append(sk, ski)
		pi = append(pi, i)
		ni = append(ni, i)
	}
	cipher := make([]byte, 2)
	onion := Onion{cipher, 1, 2, 4}
	o1 := Wrap(onion, pk[0], pi[0], ni[0])
	o1 = Wrap(o1, pk[1], pi[1], ni[1])
	o1 = Wrap(o1, pk[2], pi[2], ni[2])
	o1 = ProcOnion(o1, sk[2])
	o1 = ProcOnion(o1, sk[1])
	o1 = ProcOnion(o1, sk[0])
	dict := Init(pi, pk)
	pi = append(pi, 5)
	o2 := FormOnion(make([]byte, 2), dict, pi, ni)
	o2 = ProcOnion(o2, sk[0])
	if o2.NextParticipant != 1 {
		panic("error wrong number")
	}
	if o1.NextParticipant != 2 {
		panic("error wrong number")
	}

}
