// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package eddsa

import (
	"crypto/sha256"
	"math/rand"
	"testing"

	crand "crypto/rand"

	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bls24-317/fr"
	"github.com/consensys/gnark-crypto/hash"
)

func Example() {
	// instantiate hash function
	hFunc := hash.MIMC_BLS24_317.New()

	// create a eddsa key pair
	privateKey, _ := GenerateKey(crand.Reader)
	publicKey := privateKey.PublicKey

	// note that the message is on 4 bytes
	msg := []byte{0xde, 0xad, 0xf0, 0x0d}

	// sign the message
	signature, _ := privateKey.Sign(msg, hFunc)

	// verifies signature
	isValid, _ := publicKey.Verify(signature, msg, hFunc)
	if !isValid {
		fmt.Println("1. invalid signature")
	} else {
		fmt.Println("1. valid signature")
	}

	// Output: 1. valid signature
}

func TestSerialization(t *testing.T) {

	src := rand.NewSource(0)
	r := rand.New(src)

	privKey1, err := GenerateKey(r)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := privKey1.PublicKey

	privKey2, err := GenerateKey(r)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := privKey2.PublicKey

	pubKeyBin1 := pubKey1.Bytes()
	pubKey2.SetBytes(pubKeyBin1)
	pubKeyBin2 := pubKey2.Bytes()
	if len(pubKeyBin1) != len(pubKeyBin2) {
		t.Fatal("Inconistent size")
	}
	for i := 0; i < len(pubKeyBin1); i++ {
		if pubKeyBin1[i] != pubKeyBin2[i] {
			t.Fatal("Error serialize(deserialize(.))")
		}
	}

	privKeyBin1 := privKey1.Bytes()
	privKey2.SetBytes(privKeyBin1)
	privKeyBin2 := privKey2.Bytes()
	if len(privKeyBin1) != len(privKeyBin2) {
		t.Fatal("Inconistent size")
	}
	for i := 0; i < len(privKeyBin1); i++ {
		if privKeyBin1[i] != privKeyBin2[i] {
			t.Fatal("Error serialize(deserialize(.))")
		}
	}
}

func TestEddsaMIMC(t *testing.T) {

	src := rand.NewSource(0)
	r := rand.New(src)

	// create eddsa obj and sign a message
	privKey, err := GenerateKey(r)
	if err != nil {
		t.Fatal(nil)
	}
	pubKey := privKey.PublicKey
	hFunc := hash.MIMC_BLS24_317.New()

	var frMsg fr.Element
	frMsg.SetString("44717650746155748460101257525078853138837311576962212923649547644148297035978")
	msgBin := frMsg.Bytes()
	signature, err := privKey.Sign(msgBin[:], hFunc)
	if err != nil {
		t.Fatal(err)
	}

	// verifies correct msg
	res, err := pubKey.Verify(signature, msgBin[:], hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("Verifiy correct signature should return true")
	}

	// verifies wrong msg
	frMsg.SetString("44717650746155748460101257525078853138837311576962212923649547644148297035979")
	msgBin = frMsg.Bytes()
	res, err = pubKey.Verify(signature, msgBin[:], hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if res {
		t.Fatal("Verfiy wrong signature should be false")
	}

}

func TestEddsaSHA256(t *testing.T) {

	src := rand.NewSource(0)
	r := rand.New(src)

	hFunc := sha256.New()

	// create eddsa obj and sign a message
	// create eddsa obj and sign a message

	privKey, err := GenerateKey(r)
	pubKey := privKey.PublicKey
	if err != nil {
		t.Fatal(err)
	}

	signature, err := privKey.Sign([]byte("message"), hFunc)
	if err != nil {
		t.Fatal(err)
	}

	// verifies correct msg
	res, err := pubKey.Verify(signature, []byte("message"), hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("Verifiy correct signature should return true")
	}

	// verifies wrong msg
	res, err = pubKey.Verify(signature, []byte("wrong_message"), hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if res {
		t.Fatal("Verfiy wrong signature should be false")
	}

}

// benchmarks

func BenchmarkVerify(b *testing.B) {

	src := rand.NewSource(0)
	r := rand.New(src)

	hFunc := hash.MIMC_BLS24_317.New()

	// create eddsa obj and sign a message
	privKey, err := GenerateKey(r)
	pubKey := privKey.PublicKey
	if err != nil {
		b.Fatal(err)
	}
	var frMsg fr.Element
	frMsg.SetString("44717650746155748460101257525078853138837311576962212923649547644148297035978")
	msgBin := frMsg.Bytes()
	signature, _ := privKey.Sign(msgBin[:], hFunc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pubKey.Verify(signature, msgBin[:], hFunc)
	}
}
