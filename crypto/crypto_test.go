package crypto

import (
	"testing"

	"github.com/dedis/protobuf"
)

type testingstruct struct {
	Test int
}

func TestEncryption(*testing.T) {
	string1 := 4
	struct1 := testingstruct{string1}
	encoded, err := protobuf.Encode(&struct1)
	if err != nil {
		panic(err)
	}
	pk, sk := GenKey()
	encrypted := Encrypt(pk, encoded)
	decrypted := Decrypt(sk, encrypted)
	var decoded testingstruct
	protobuf.Decode(decrypted, &decoded)
	if decoded.Test != string1 {
		panic("Wrong String!")
	}

}
