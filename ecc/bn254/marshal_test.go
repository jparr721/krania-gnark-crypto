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

// Code generated by jparr721/krania-gnark-crypto DO NOT EDIT

package bn254

import (
	"bytes"
	"io"
	"math/big"
	"math/rand"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"

	"github.com/jparr721/krania-gnark-crypto/ecc/bn254/fp"
	"github.com/jparr721/krania-gnark-crypto/ecc/bn254/fr"
	"github.com/jparr721/krania-gnark-crypto/ecc/bn254/internal/fptower"
)

const (
	nbFuzzShort = 10
	nbFuzz      = 100
)

func TestEncoder(t *testing.T) {
	t.Parallel()
	// TODO need proper fuzz testing here

	var inA uint64
	var inB fr.Element
	var inC fp.Element
	var inD G1Affine
	var inE G1Affine
	var inF G2Affine
	var inG []G1Affine
	var inH []G2Affine
	var inI []fp.Element
	var inJ []fr.Element

	// set values of inputs
	inA = rand.Uint64()
	inB.SetRandom()
	inC.SetRandom()
	inD.ScalarMultiplication(&g1GenAff, new(big.Int).SetUint64(rand.Uint64()))
	// inE --> infinity
	inF.ScalarMultiplication(&g2GenAff, new(big.Int).SetUint64(rand.Uint64()))
	inG = make([]G1Affine, 2)
	inH = make([]G2Affine, 0)
	inG[1] = inD
	inI = make([]fp.Element, 3)
	inI[2] = inD.X
	inJ = make([]fr.Element, 0)

	// encode them, compressed and raw
	var buf, bufRaw bytes.Buffer
	enc := NewEncoder(&buf)
	encRaw := NewEncoder(&bufRaw, RawEncoding())
	toEncode := []interface{}{inA, &inB, &inC, &inD, &inE, &inF, inG, inH, inI, inJ}
	for _, v := range toEncode {
		if err := enc.Encode(v); err != nil {
			t.Fatal(err)
		}
		if err := encRaw.Encode(v); err != nil {
			t.Fatal(err)
		}
	}

	testDecode := func(t *testing.T, r io.Reader, n int64) {
		dec := NewDecoder(r)
		var outA uint64
		var outB fr.Element
		var outC fp.Element
		var outD G1Affine
		var outE G1Affine
		outE.X.SetOne()
		outE.Y.SetUint64(42)
		var outF G2Affine
		var outG []G1Affine
		var outH []G2Affine
		var outI []fp.Element
		var outJ []fr.Element

		toDecode := []interface{}{&outA, &outB, &outC, &outD, &outE, &outF, &outG, &outH, &outI, &outJ}
		for _, v := range toDecode {
			if err := dec.Decode(v); err != nil {
				t.Fatal(err)
			}
		}

		// compare values
		if inA != outA {
			t.Fatal("didn't encode/decode uint64 value properly")
		}

		if !inB.Equal(&outB) || !inC.Equal(&outC) {
			t.Fatal("decode(encode(Element) failed")
		}
		if !inD.Equal(&outD) || !inE.Equal(&outE) {
			t.Fatal("decode(encode(G1Affine) failed")
		}
		if !inF.Equal(&outF) {
			t.Fatal("decode(encode(G2Affine) failed")
		}
		if (len(inG) != len(outG)) || (len(inH) != len(outH)) {
			t.Fatal("decode(encode(slice(points))) failed")
		}
		for i := 0; i < len(inG); i++ {
			if !inG[i].Equal(&outG[i]) {
				t.Fatal("decode(encode(slice(points))) failed")
			}
		}
		if (len(inI) != len(outI)) || (len(inJ) != len(outJ)) {
			t.Fatal("decode(encode(slice(elements))) failed")
		}
		for i := 0; i < len(inI); i++ {
			if !inI[i].Equal(&outI[i]) {
				t.Fatal("decode(encode(slice(elements))) failed")
			}
		}
		if n != dec.BytesRead() {
			t.Fatal("bytes read don't match bytes written")
		}
	}

	// decode them
	testDecode(t, &buf, enc.BytesWritten())
	testDecode(t, &bufRaw, encRaw.BytesWritten())

}

func TestIsCompressed(t *testing.T) {
	t.Parallel()
	var g1Inf, g1 G1Affine
	var g2Inf, g2 G2Affine

	g1 = g1GenAff
	g2 = g2GenAff

	{
		b := g1Inf.Bytes()
		if !isCompressed(b[0]) {
			t.Fatal("g1Inf.Bytes() should be compressed")
		}
	}

	{
		b := g1Inf.RawBytes()
		if isCompressed(b[0]) {
			t.Fatal("g1Inf.RawBytes() should be uncompressed")
		}
	}

	{
		b := g1.Bytes()
		if !isCompressed(b[0]) {
			t.Fatal("g1.Bytes() should be compressed")
		}
	}

	{
		b := g1.RawBytes()
		if isCompressed(b[0]) {
			t.Fatal("g1.RawBytes() should be uncompressed")
		}
	}

	{
		b := g2Inf.Bytes()
		if !isCompressed(b[0]) {
			t.Fatal("g2Inf.Bytes() should be compressed")
		}
	}

	{
		b := g2Inf.RawBytes()
		if isCompressed(b[0]) {
			t.Fatal("g2Inf.RawBytes() should be uncompressed")
		}
	}

	{
		b := g2.Bytes()
		if !isCompressed(b[0]) {
			t.Fatal("g2.Bytes() should be compressed")
		}
	}

	{
		b := g2.RawBytes()
		if isCompressed(b[0]) {
			t.Fatal("g2.RawBytes() should be uncompressed")
		}
	}

}

func TestG1AffineSerialization(t *testing.T) {
	t.Parallel()
	// test round trip serialization of infinity
	{
		// compressed
		{
			var p1, p2 G1Affine
			p2.X.SetRandom()
			p2.Y.SetRandom()
			buf := p1.Bytes()
			n, err := p2.SetBytes(buf[:])
			if err != nil {
				t.Fatal(err)
			}
			if n != SizeOfG1AffineCompressed {
				t.Fatal("invalid number of bytes consumed in buffer")
			}
			if !(p2.X.IsZero() && p2.Y.IsZero()) {
				t.Fatal("deserialization of uncompressed infinity point is not infinity")
			}
		}

		// uncompressed
		{
			var p1, p2 G1Affine
			p2.X.SetRandom()
			p2.Y.SetRandom()
			buf := p1.RawBytes()
			n, err := p2.SetBytes(buf[:])
			if err != nil {
				t.Fatal(err)
			}
			if n != SizeOfG1AffineUncompressed {
				t.Fatal("invalid number of bytes consumed in buffer")
			}
			if !(p2.X.IsZero() && p2.Y.IsZero()) {
				t.Fatal("deserialization of uncompressed infinity point is not infinity")
			}
		}
	}

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[G1] Affine SetBytes(RawBytes) should stay the same", prop.ForAll(
		func(a fp.Element) bool {
			var start, end G1Affine
			var ab big.Int
			a.ToBigIntRegular(&ab)
			start.ScalarMultiplication(&g1GenAff, &ab)

			buf := start.RawBytes()
			n, err := end.SetBytes(buf[:])
			if err != nil {
				return false
			}
			if n != SizeOfG1AffineUncompressed {
				return false
			}
			return start.X.Equal(&end.X) && start.Y.Equal(&end.Y)
		},
		GenFp(),
	))

	properties.Property("[G1] Affine SetBytes(Bytes()) should stay the same", prop.ForAll(
		func(a fp.Element) bool {
			var start, end G1Affine
			var ab big.Int
			a.ToBigIntRegular(&ab)
			start.ScalarMultiplication(&g1GenAff, &ab)

			buf := start.Bytes()
			n, err := end.SetBytes(buf[:])
			if err != nil {
				return false
			}
			if n != SizeOfG1AffineCompressed {
				return false
			}
			return start.X.Equal(&end.X) && start.Y.Equal(&end.Y)
		},
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG2AffineSerialization(t *testing.T) {
	t.Parallel()
	// test round trip serialization of infinity
	{
		// compressed
		{
			var p1, p2 G2Affine
			p2.X.SetRandom()
			p2.Y.SetRandom()
			buf := p1.Bytes()
			n, err := p2.SetBytes(buf[:])
			if err != nil {
				t.Fatal(err)
			}
			if n != SizeOfG2AffineCompressed {
				t.Fatal("invalid number of bytes consumed in buffer")
			}
			if !(p2.X.IsZero() && p2.Y.IsZero()) {
				t.Fatal("deserialization of uncompressed infinity point is not infinity")
			}
		}

		// uncompressed
		{
			var p1, p2 G2Affine
			p2.X.SetRandom()
			p2.Y.SetRandom()
			buf := p1.RawBytes()
			n, err := p2.SetBytes(buf[:])
			if err != nil {
				t.Fatal(err)
			}
			if n != SizeOfG2AffineUncompressed {
				t.Fatal("invalid number of bytes consumed in buffer")
			}
			if !(p2.X.IsZero() && p2.Y.IsZero()) {
				t.Fatal("deserialization of uncompressed infinity point is not infinity")
			}
		}
	}

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[G2] Affine SetBytes(RawBytes) should stay the same", prop.ForAll(
		func(a fp.Element) bool {
			var start, end G2Affine
			var ab big.Int
			a.ToBigIntRegular(&ab)
			start.ScalarMultiplication(&g2GenAff, &ab)

			buf := start.RawBytes()
			n, err := end.SetBytes(buf[:])
			if err != nil {
				return false
			}
			if n != SizeOfG2AffineUncompressed {
				return false
			}
			return start.X.Equal(&end.X) && start.Y.Equal(&end.Y)
		},
		GenFp(),
	))

	properties.Property("[G2] Affine SetBytes(Bytes()) should stay the same", prop.ForAll(
		func(a fp.Element) bool {
			var start, end G2Affine
			var ab big.Int
			a.ToBigIntRegular(&ab)
			start.ScalarMultiplication(&g2GenAff, &ab)

			buf := start.Bytes()
			n, err := end.SetBytes(buf[:])
			if err != nil {
				return false
			}
			if n != SizeOfG2AffineCompressed {
				return false
			}
			return start.X.Equal(&end.X) && start.Y.Equal(&end.Y)
		},
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// define Gopters generators

// GenFr generates an Fr element
func GenFr() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var elmt fr.Element

		if _, err := elmt.SetRandom(); err != nil {
			panic(err)
		}

		return gopter.NewGenResult(elmt, gopter.NoShrinker)
	}
}

// GenFp generates an Fp element
func GenFp() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var elmt fp.Element

		if _, err := elmt.SetRandom(); err != nil {
			panic(err)
		}

		return gopter.NewGenResult(elmt, gopter.NoShrinker)
	}
}

// GenE2 generates an fptower.E2 elmt
func GenE2() gopter.Gen {
	return gopter.CombineGens(
		GenFp(),
		GenFp(),
	).Map(func(values []interface{}) fptower.E2 {
		return fptower.E2{A0: values[0].(fp.Element), A1: values[1].(fp.Element)}
	})
}

// GenE6 generates an fptower.E6 elmt
func GenE6() gopter.Gen {
	return gopter.CombineGens(
		GenE2(),
		GenE2(),
		GenE2(),
	).Map(func(values []interface{}) fptower.E6 {
		return fptower.E6{B0: values[0].(fptower.E2), B1: values[1].(fptower.E2), B2: values[2].(fptower.E2)}
	})
}

// GenE12 generates an fptower.E6 elmt
func GenE12() gopter.Gen {
	return gopter.CombineGens(
		GenE6(),
		GenE6(),
	).Map(func(values []interface{}) fptower.E12 {
		return fptower.E12{C0: values[0].(fptower.E6), C1: values[1].(fptower.E6)}
	})
}

// GenBigInt generates a big.Int
func GenBigInt() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var s big.Int
		var b [fp.Bytes]byte
		_, err := rand.Read(b[:])
		if err != nil {
			panic(err)
		}
		s.SetBytes(b[:])
		genResult := gopter.NewGenResult(s, gopter.NoShrinker)
		return genResult
	}
}
