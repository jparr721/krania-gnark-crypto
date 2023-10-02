// Copyright 2023 Tianyi Liu and Tiancheng Xie
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

package dkzg

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/jparr721/goMPI/mpi"
)

// testSRS re-used accross tests of the KZG scheme
var testSRS []*SRS

const srsSize = 230

func init() {
	testSRS = make([]*SRS, mpi.WorldSize)
	for i := 0; i < int(mpi.WorldSize); i++ {
		var err error
		testSRS[i], err = newSRS(ecc.NextPowerOfTwo(srsSize), []*big.Int{big.NewInt(42), big.NewInt(int64(27))}, uint64(i))
		if err != nil {
			panic(err)
		}
	}
}

func newSRS(size uint64, tau []*big.Int, num uint64) (*SRS, error) {

	_, _, gen1Aff, gen2Aff := bn254.Generators()
	t := num

	tau0 := new(fr.Element).SetBigInt(tau[0])
	// Lagrange Polynomial
	lagTau0 := lagrangeCalc(t, *tau0, nil)

	var srs SRS

	var alpha fr.Element
	alpha.SetBigInt(tau[1])

	srs.G2[0] = gen2Aff
	srs.G2[1].ScalarMultiplication(&gen2Aff, tau[1])

	lagBigInt := new(big.Int)
	lagTau0.ToBigIntRegular(lagBigInt)
	srs.G1 = make([]bn254.G1Affine, size)
	srs.G1[0].ScalarMultiplication(&gen1Aff, lagBigInt)

	alphas := make([]fr.Element, size-1)
	alphas[0] = alpha
	alphas[0].Mul(&alphas[0], &lagTau0)
	for i := 1; i < len(alphas); i++ {
		alphas[i].Mul(&alphas[i-1], &alpha)
	}
	for i := 0; i < len(alphas); i++ {
		alphas[i].FromMont()
	}
	g1s := bn254.BatchScalarMultiplicationG1(&gen1Aff, alphas)
	copy(srs.G1[1:], g1s)
	return &srs, nil
}

func TestDividePolyByXminusA(t *testing.T) {

	const pSize = 230

	// build random polynomial
	pol := make([]fr.Element, pSize)
	pol[0].SetRandom()
	for i := 1; i < pSize; i++ {
		pol[i] = pol[i-1]
	}

	// evaluate the polynomial at a random point
	var point fr.Element
	point.SetRandom()
	evaluation := eval(pol, point)

	// probabilistic test (using Schwartz Zippel lemma, evaluation at one point is enough)
	var randPoint, xminusa fr.Element
	randPoint.SetRandom()
	polRandpoint := eval(pol, randPoint)
	polRandpoint.Sub(&polRandpoint, &evaluation) // f(rand)-f(point)

	// compute f-f(a)/x-a
	h := dividePolyByXminusA(pol, evaluation, point)
	pol = nil // h reuses this memory

	if len(h) != 229 {
		t.Fatal("inconsistant size of quotient")
	}

	hRandPoint := eval(h, randPoint)
	xminusa.Sub(&randPoint, &point) // rand-point

	// f(rand)-f(point)	==? h(rand)*(rand-point)
	hRandPoint.Mul(&hRandPoint, &xminusa)

	if !hRandPoint.Equal(&polRandpoint) {
		t.Fatal("Error f-f(a)/x-a")
	}
}

func TestSerializationSRS(t *testing.T) {

	// create a SRS
	srs, err := NewSRS(64, []*big.Int{big.NewInt(42), big.NewInt(27)}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// serialize it...
	var buf bytes.Buffer
	_, err = srs.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}

	// reconstruct the SRS
	var _srs SRS
	_, err = _srs.ReadFrom(&buf)
	if err != nil {
		t.Fatal(err)
	}

	// compare

	if !reflect.DeepEqual(srs.G2, _srs.G2) {
		t.Fatal("scheme serialization failed G2")
	}
	if !reflect.DeepEqual(srs.G1, _srs.G1) {
		t.Fatal("scheme serialization failed G1")
	}

}

func TestCommit(t *testing.T) {

	// create a polynomial
	f := polynomial(60, mpi.SelfRank)

	fmt.Println("poly", mpi.SelfRank)
	for i := 0; i < len(f); i++ {
		fmt.Print(f[i], " ")
	}
	fmt.Println()

	// commit using the method from KZG
	kzgCommit, err := Commit(f, testSRS[mpi.SelfRank])
	if err != nil {
		t.Fatal(err)
	}
	if mpi.SelfRank != 0 {
		return
	}

	// check commitment using manual commit
	var x fr.Element
	x.SetString("27")

	ps := make([][]fr.Element, mpi.WorldSize)
	es := make([]fr.Element, mpi.WorldSize)
	bs := make([]bn254.G1Affine, mpi.WorldSize)
	for i := 0; i < int(mpi.WorldSize); i++ {
		ps[i] = polynomial(60, uint64(i))
		es[i] = eval(ps[i], x)
		bs[i] = testSRS[i].G1[0]
	}

	config := ecc.MultiExpConfig{ScalarsMont: true}
	var manualCommit bn254.G1Affine
	if _, err = manualCommit.MultiExp(bs, es, config); err != nil {
		t.Fatal(err)
	}

	// compare both results
	if !kzgCommit.Equal(&manualCommit) {
		t.Fatal("error KZG commitment")
	}
}

func TestVerifySinglePoint(t *testing.T) {

	// create a polynomial
	f := polynomial(60, mpi.SelfRank)

	// commit the polynomial
	digest, err := Commit(f, testSRS[mpi.SelfRank])
	if err != nil {
		t.Fatal(err)
	}

	// compute opening proof at a random point
	var point fr.Element
	point.SetString("4321")
	proof, _, err := Open(f, point, testSRS[mpi.SelfRank])
	if err != nil {
		t.Fatal(err)
	}

	if mpi.SelfRank != 0 {
		return
	}

	ps := make([][]fr.Element, mpi.WorldSize)
	es := make([]fr.Element, mpi.WorldSize)
	bs := make([]bn254.G1Affine, mpi.WorldSize)
	for i := 0; i < int(mpi.WorldSize); i++ {
		ps[i] = polynomial(60, uint64(i))
		es[i] = eval(ps[i], point)
		bs[i] = testSRS[i].G1[0]
	}

	var expectedGroup bn254.G1Affine
	if _, err := expectedGroup.MultiExp(bs, es, ecc.MultiExpConfig{ScalarsMont: true}); err != nil {
		t.Fatal(err)
	}

	// verify the claimed valued
	if !proof.ClaimedDigest.Equal(&expectedGroup) {
		t.Fatal("inconsistant claimed value")
	}

	// verify correct proof
	err = Verify(&digest, &proof, point, testSRS[mpi.SelfRank])
	if err != nil {
		t.Fatal(err)
	}
	{
		// verify wrong proof
		var nexpectedGroup bn254.G1Affine
		nexpectedGroup.Add(&expectedGroup, &expectedGroup)
		proof.ClaimedDigest = nexpectedGroup

		err = Verify(&digest, &proof, point, testSRS[mpi.SelfRank])
		if err == nil {
			t.Fatal("verifying wrong proof should have failed")
		}
	}
	{
		// verify wrong proof with quotient set to zero
		// see https://cryptosubtlety.medium.com/00-8d4adcf4d255
		proof.H.X.SetZero()
		proof.H.Y.SetZero()
		err = Verify(&digest, &proof, point, testSRS[mpi.SelfRank])
		if err == nil {
			t.Fatal("verifying wrong proof should have failed")
		}
	}
}

func TestBatchOpenSinglePoint(t *testing.T) {
	// create a polynomial
	num := 10
	fs := make([][]fr.Element, num)
	digests := make([]bn254.G1Affine, num)

	var err error
	for i := 0; i < num; i++ {
		fs[i] = polynomial(60, mpi.SelfRank, i)
		digests[i], err = Commit(fs[i], testSRS[mpi.SelfRank])
		if err != nil {
			t.Fatal(err)
		}
	}

	// compute opening proof at a random point
	var point fr.Element
	point.SetString("4321")
	hfunc := sha256.New()
	proof, evals, err := BatchOpenSinglePoint(fs, digests, point, hfunc, testSRS[mpi.SelfRank])
	if err != nil {
		t.Fatal(err)
	}

	if mpi.SelfRank != 0 {
		return
	}

	ps := make([][][]fr.Element, num)
	es := make([][]fr.Element, num)
	bs := make([]bn254.G1Affine, mpi.WorldSize)
	for i := 0; i < num; i++ {
		ps[i] = make([][]fr.Element, mpi.WorldSize)
		es[i] = make([]fr.Element, mpi.WorldSize)
		for j := 0; j < int(mpi.WorldSize); j++ {
			ps[i][j] = polynomial(60, uint64(j), i)
			es[i][j] = eval(ps[i][j], point)
			if !es[i][j].Equal(&evals[i][j]) {
				t.Fatal("inconsistant evals")
			}
		}
	}

	for i := 0; i < int(mpi.WorldSize); i++ {
		bs[i] = testSRS[i].G1[0]
	}

	expectedGroups := make([]bn254.G1Affine, num)
	for i := 0; i < num; i++ {
		if _, err := expectedGroups[i].MultiExp(bs, es[i], ecc.MultiExpConfig{ScalarsMont: true}); err != nil {
			t.Fatal(err)
		}
		if !proof.ClaimedDigests[i].Equal(&expectedGroups[i]) {
			t.Fatal("inconsistant claimed digests for evaluation")
		}
	}

	// verify correct proof
	if err := BatchVerifySinglePoint(digests, &proof, point, hfunc, testSRS[mpi.SelfRank]); err != nil {
		t.Fatal(err)
	}
	{
		// verify wrong proof
		var nexpectedGroup bn254.G1Affine
		nexpectedGroup.Add(&proof.ClaimedDigests[0], &proof.ClaimedDigests[0])
		proof.ClaimedDigests[0] = nexpectedGroup

		if err := BatchVerifySinglePoint(digests, &proof, point, hfunc, testSRS[mpi.SelfRank]); err == nil {
			t.Fatal("verifying wrong claimed digests should have failed")
		}
	}
	{
		// verify wrong proof with quotient set to zero
		// see https://cryptosubtlety.medium.com/00-8d4adcf4d255
		proof.H.X.SetZero()
		proof.H.Y.SetZero()
		if err := BatchVerifySinglePoint(digests, &proof, point, hfunc, testSRS[mpi.SelfRank]); err == nil {
			t.Fatal("verifying wrong quotient digest should have failed")
		}
	}
}

const benchSize = 1 << 16

func BenchmarkKZGCommit(b *testing.B) {
	benchSRS, err := NewSRS(ecc.NextPowerOfTwo(benchSize), []*big.Int{big.NewInt(42), big.NewInt(27)}, nil)
	if err != nil {
		b.Fatal(err)
	}
	// random polynomial
	p := polynomial(benchSize/2, mpi.SelfRank)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Commit(p, benchSRS)
	}
}

func BenchmarkDivideByXMinusA(b *testing.B) {
	const pSize = 1 << 22

	// build random polynomial
	pol := make([]fr.Element, pSize)
	pol[0].SetRandom()
	for i := 1; i < pSize; i++ {
		pol[i] = pol[i-1]
	}
	var a, fa fr.Element
	a.SetRandom()
	fa.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dividePolyByXminusA(pol, fa, a)
		pol = pol[:pSize]
		pol[pSize-1] = pol[0]
	}
}

func BenchmarkKZGOpen(b *testing.B) {
	benchSRS, err := NewSRS(ecc.NextPowerOfTwo(benchSize), []*big.Int{big.NewInt(42), big.NewInt(27)}, nil)
	if err != nil {
		b.Fatal(err)
	}

	// random polynomial
	p := polynomial(benchSize/2, mpi.SelfRank)
	var r fr.Element
	r.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = Open(p, r, benchSRS)
	}
}

func BenchmarkKZGVerify(b *testing.B) {
	benchSRS, err := NewSRS(ecc.NextPowerOfTwo(benchSize), []*big.Int{big.NewInt(42), big.NewInt(27)}, nil)
	if err != nil {
		b.Fatal(err)
	}

	// random polynomial
	p := polynomial(benchSize/2, mpi.SelfRank)
	var r fr.Element
	r.SetRandom()

	// commit
	comm, err := Commit(p, benchSRS)
	if err != nil {
		b.Fatal(err)
	}

	// open
	openingProof, _, err := Open(p, r, benchSRS)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := Verify(&comm, &openingProof, r, benchSRS)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKZGBatchOpen10(b *testing.B) {
	benchSRS, err := NewSRS(ecc.NextPowerOfTwo(benchSize), []*big.Int{big.NewInt(42), big.NewInt(27)}, nil)
	if err != nil {
		b.Fatal(err)
	}

	// 10 random polynomials
	var ps [10][]fr.Element
	for i := 0; i < 10; i++ {
		ps[i] = polynomial(benchSize/2, mpi.SelfRank)
	}

	// commitments
	var commitments [10]Digest
	for i := 0; i < 10; i++ {
		commitments[i], _ = Commit(ps[i], benchSRS)
	}

	// pick a hash function
	hf := sha256.New()

	var r fr.Element
	r.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := BatchOpenSinglePoint(ps[:], commitments[:], r, hf, benchSRS)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKZGBatchVerify10(b *testing.B) {
	benchSRS, err := NewSRS(ecc.NextPowerOfTwo(benchSize), []*big.Int{big.NewInt(42), big.NewInt(27)}, nil)
	if err != nil {
		b.Fatal(err)
	}

	// 10 random polynomials
	var ps [10][]fr.Element
	for i := 0; i < 10; i++ {
		ps[i] = polynomial(benchSize/2, mpi.SelfRank)
	}

	// commitments
	var commitments [10]Digest
	for i := 0; i < 10; i++ {
		commitments[i], _ = Commit(ps[i], benchSRS)
	}

	// pick a hash function
	hf := sha256.New()

	var r fr.Element
	r.SetRandom()

	proof, _, err := BatchOpenSinglePoint(ps[:], commitments[:], r, hf, benchSRS)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := BatchVerifySinglePoint(commitments[:], &proof, r, hf, benchSRS)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func polynomial(size int, num uint64, idx ...int) []fr.Element {
	f := make([]fr.Element, size)
	for i := 0; i < size; i++ {
		tmp := fr.NewElement(num)
		tmp2 := fr.NewElement(uint64(i))
		f[i].Add(&tmp2, &tmp)
		if len(idx) > 0 {
			tmp3 := fr.NewElement(uint64(idx[0]*3 + 1))
			f[i].Add(&f[i], &tmp3)
		}
	}
	return f
}
