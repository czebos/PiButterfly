package crypto

import (
	"github.com/czebos/pi/supportCode"
	"github.com/dedis/kyber"
	"github.com/dedis/kyber/encrypt/ecies"
	"github.com/dedis/kyber/util/random"
)

type Scalar struct {
	s kyber.Scalar
}

type Client struct {
	Id int
	Sk Scalar
}

// Encrypts a message or ciphertext of type []byte by using the public key
// Outputs the ciphertext. Imported from Kyber
func Encrypt(pk supportCode.Point, msg []byte) []byte { //maybe add special type to make it look nicer
	ciphertext, err := ecies.Encrypt(supportCode.SUITE, pk.P, msg, supportCode.SUITE.Hash)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

// Decrypts a ciphertext of type []byte by using the secret key
// Outputs the ciphertext/message. Imported from Kyber
func Decrypt(sk Scalar, ciphertext []byte) []byte {
	msg, err := ecies.Decrypt(supportCode.SUITE, sk.s, ciphertext, supportCode.SUITE.Hash)
	if err != nil {
		panic(err)
	}
	return msg
}

//Generates a keyPair that may be used, imported from Kyber
func GenKey() (supportCode.Point, Scalar) {
	sk := supportCode.SUITE.Scalar().Pick(random.New())
	pk := supportCode.SUITE.Point().Mul(sk, nil)
	return supportCode.Point{pk}, Scalar{sk}

}
