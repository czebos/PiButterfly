package onionDS

import (
	"github.com/czebos/pi/crypto"
	"github.com/czebos/pi/supportCode"
	"github.com/dedis/protobuf"
)

//Basic struct for an Onion
type Onion struct {
	Ciphertext      []byte //This can be message or encryptedOnion
	Id              int
	NextParticipant int
	Nonce           int
}

func EncodeOnion(o *Onion) ([]byte, error) {
	return protobuf.Encode(o)
}

func DecodeOnion(bytes []byte, o *Onion) error {
	return protobuf.Decode(bytes, o)
}

var id int

// Creates mapping of all participants to their public ID
// Arrays must be in order
func Init(participants []int, pkis []supportCode.Point) map[int]supportCode.Point {
	m := make(map[int]supportCode.Point)
	for i := 0; i < len(participants); i++ {
		m[participants[i]] = pkis[i]
	}
	return m
}

// Simplifies importing by calling GenKey from Crypto
func GenKey() (supportCode.Point, crypto.Scalar) {
	pk, sk := crypto.GenKey()
	return pk, sk
}

// Wraps the onion with the given public key pk, recipitant p, and nonce n
// Outputes the wrapped onion with the ciphertext of the old onion
func Wrap(O Onion, pk supportCode.Point, p int, n int) Onion {
	encodedOnion, err := EncodeOnion(&O) //must turn to byte form
	if err != nil {
		panic(err)
	}
	encryptedOnion := crypto.Encrypt(pk, encodedOnion)

	wrapped := Onion{encryptedOnion, id, p, n}
	id += int(1) //increment the onionID, shows what layer it is
	return wrapped
}

// Loops through wrap to create multiple layers of the onion. Starts
// From the end so that the onion encrypts correctly
// Complex wrap function, outputs the onion
func FormOnion(msg []byte, pki map[int]supportCode.Point, pi []int, ni []int) Onion {
	if len(pi)-1 != len(ni) {
		panic("Wrong lengths")
	}
	currOnion := Onion{msg, 0, pi[len(pi)-1], 0}
	for i := len(pi) - 2; i >= 0; i-- {
		currOnion = Wrap(currOnion, pki[pi[i]], pi[i], ni[i])
	}
	return currOnion
}

// Procs or "Peels" a layer of the onion. Makes the onion
// one layer less, or reveals the message
func ProcOnion(O Onion, sk crypto.Scalar) Onion {
	decryptedOnion := crypto.Decrypt(sk, O.Ciphertext) //then decrypts it
	var proccedOnion Onion
	DecodeOnion(decryptedOnion, &proccedOnion) //finally decodes the encoded bytes to an Onion
	return proccedOnion
}
