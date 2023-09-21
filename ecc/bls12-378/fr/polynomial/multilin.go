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

package polynomial

import (
	"github.com/jparr721/krania-gnark-crypto/ecc/bls12-378/fr"
)

// MultiLin tracks the values of a (dense i.e. not sparse) multilinear polynomial
// The variables are X₁ through Xₙ where n = log(len(.))
// .[∑ᵢ 2ⁱ⁻¹ bₙ₋ᵢ] = the polynomial evaluated at (b₁, b₂, ..., bₙ)
// It is understood that any hypercube evaluation can be extrapolated to a multilinear polynomial
type MultiLin []fr.Element

// Fold is partial evaluation function k[X₁, X₂, ..., Xₙ] → k[X₂, ..., Xₙ] by setting X₁=r
func (m *MultiLin) Fold(r fr.Element) {
	mid := len(*m) / 2

	bottom, top := (*m)[:mid], (*m)[mid:]

	// updating bookkeeping table
	// knowing that the polynomial f ∈ (k[X₂, ..., Xₙ])[X₁] is linear, we would get f(r) = f(0) + r(f(1) - f(0))
	// the following loop computes the evaluations of f(r) accordingly:
	//		f(r, b₂, ..., bₙ) = f(0, b₂, ..., bₙ) + r(f(1, b₂, ..., bₙ) - f(0, b₂, ..., bₙ))
	for i := 0; i < mid; i++ {
		// table[i] ← table[i] + r (table[i + mid] - table[i])
		top[i].Sub(&top[i], &bottom[i])
		top[i].Mul(&top[i], &r)
		bottom[i].Add(&bottom[i], &top[i])
	}

	*m = (*m)[:mid]
}

// Evaluate extrapolate the value of the multilinear polynomial corresponding to m
// on the given coordinates
func (m MultiLin) Evaluate(coordinates []fr.Element) fr.Element {
	// Folding is a mutating operation
	bkCopy := m.Clone()

	// Evaluate step by step through repeated folding (i.e. evaluation at the first remaining variable)
	for _, r := range coordinates {
		bkCopy.Fold(r)
	}

	return bkCopy[0]
}

// Clone creates a deep copy of a book-keeping table.
// Both multilinear interpolation and sumcheck require folding an underlying
// array, but folding changes the array. To do both one requires a deep copy
// of the book-keeping table.
func (m MultiLin) Clone() MultiLin {
	tableDeepCopy := Make(len(m))
	copy(tableDeepCopy, m)
	return tableDeepCopy
}

// Add two bookKeepingTables
func (m *MultiLin) Add(left, right MultiLin) {
	size := len(left)
	// Check that left and right have the same size
	if len(right) != size {
		panic("Left and right do not have the right size")
	}
	// Reallocate the table if necessary
	if cap(*m) < size {
		*m = make([]fr.Element, size)
	}

	// Resize the destination table
	*m = (*m)[:size]

	// Add elementwise
	for i := 0; i < size; i++ {
		(*m)[i].Add(&left[i], &right[i])
	}
}

// EvalEq computes Eq(q₁, ... , qₙ, h₁, ... , hₙ) = Π₁ⁿ Eq(qᵢ, hᵢ)
// where Eq(x,y) = xy + (1-x)(1-y) = 1 - x - y + xy + xy interpolates
//      _________________
//      |       |       |
//      |   0   |   1   |
//      |_______|_______|
//  y   |       |       |
//      |   1   |   0   |
//      |_______|_______|
//
//              x
// In other words the polynomial evaluated here is the multilinear extrapolation of
// one that evaluates to q' == h' for vectors q', h' of binary values
func EvalEq(q, h []fr.Element) fr.Element {
	var res, nxt, one, sum fr.Element
	one.SetOne()
	for i := 0; i < len(q); i++ {
		nxt.Mul(&q[i], &h[i]) // nxt <- qᵢ * hᵢ
		nxt.Double(&nxt)      // nxt <- 2 * qᵢ * hᵢ
		nxt.Add(&nxt, &one)   // nxt <- 1 + 2 * qᵢ * hᵢ
		sum.Add(&q[i], &h[i]) // sum <- qᵢ + hᵢ	TODO: Why not subtract one by one from nxt? More parallel?

		if i == 0 {
			res.Sub(&nxt, &sum) // nxt <- 1 + 2 * qᵢ * hᵢ - qᵢ - hᵢ
		} else {
			nxt.Sub(&nxt, &sum) // nxt <- 1 + 2 * qᵢ * hᵢ - qᵢ - hᵢ
			res.Mul(&res, &nxt) // res <- res * nxt
		}
	}
	return res
}

// Eq sets m to the representation of the polynomial Eq(q₁, ..., qₙ, *, ..., *) × m[0]
func (m *MultiLin) Eq(q []fr.Element) {
	n := len(q)

	if len(*m) != 1<<n {
		n := Make(1 << n)
		n[0].Set(&(*m)[0])
		//TODO: Dump m?
		*m = n
	}

	//At the end of each iteration, m(h₁, ..., hₙ) = Eq(q₁, ..., qᵢ₊₁, h₁, ..., hᵢ₊₁)
	for i, qI := range q { // In the comments we use a 1-based index so qI = qᵢ₊₁
		// go through all assignments of (b₁, ..., bᵢ) ∈ {0,1}ⁱ
		for j := 0; j < (1 << i); j++ {
			j0 := j << (n - i)                 // bᵢ₊₁ = 0
			j1 := j0 + 1<<(n-1-i)              // bᵢ₊₁ = 1
			(*m)[j1].Mul(&qI, &(*m)[j0])       // Eq(q₁, ..., qᵢ₊₁, b₁, ..., bᵢ, 1) = Eq(q₁, ..., qᵢ, b₁, ..., bᵢ) Eq(qᵢ₊₁, 1) = Eq(q₁, ..., qᵢ, b₁, ..., bᵢ) qᵢ₊₁
			(*m)[j0].Sub(&(*m)[j0], &(*m)[j1]) // Eq(q₁, ..., qᵢ₊₁, b₁, ..., bᵢ, 0) = Eq(q₁, ..., qᵢ, b₁, ..., bᵢ) Eq(qᵢ₊₁, 0) = Eq(q₁, ..., qᵢ, b₁, ..., bᵢ) (1-qᵢ₊₁)
		}
	}
}

func init() {
	//TODO: Check for whether already computed in the Getter or this?
	lagrangeBasis = make([][]Polynomial, maxLagrangeDomainSize+1)

	//size = 0: Cannot extrapolate with no data points

	//size = 1: Constant polynomial
	lagrangeBasis[1] = []Polynomial{make(Polynomial, 1)}
	lagrangeBasis[1][0][0].SetOne()

	//for size ≥ 2, the function works
	for size := uint8(2); size <= maxLagrangeDomainSize; size++ {
		lagrangeBasis[size] = computeLagrangeBasis(size)
	}
}

func getLagrangeBasis(domainSize int) []Polynomial {
	//TODO: Precompute everything at init or this?
	/*if lagrangeBasis[domainSize] == nil {
		lagrangeBasis[domainSize] = computeLagrangeBasis(domainSize)
	}*/
	return lagrangeBasis[domainSize]
}

const maxLagrangeDomainSize uint8 = 12

var lagrangeBasis [][]Polynomial

// computeLagrangeBasis precomputes in explicit coefficient form for each 0 ≤ l < domainSize the polynomial
// pₗ := X (X-1) ... (X-l-1) (X-l+1) ... (X - domainSize + 1) / ( l (l-1) ... 2 (-1) ... (l - domainSize +1) )
// Note that pₗ(l) = 1 and pₗ(n) = 0 if 0 ≤ l < domainSize, n ≠ l
func computeLagrangeBasis(domainSize uint8) []Polynomial {

	constTerms := make([]fr.Element, domainSize)
	for i := uint8(0); i < domainSize; i++ {
		constTerms[i].SetInt64(-int64(i))
	}

	res := make([]Polynomial, domainSize)
	multScratch := make(Polynomial, domainSize-1)

	// compute pₗ
	for l := uint8(0); l < domainSize; l++ {

		// TODO: Optimize this with some trees? O(log(domainSize)) polynomial mults instead of O(domainSize)? Then again it would be fewer big poly mults vs many small poly mults
		d := uint8(0) //n is the current degree of res
		for i := uint8(0); i < domainSize; i++ {
			if i == l {
				continue
			}
			if d == 0 {
				res[l] = make(Polynomial, domainSize)
				res[l][domainSize-2] = constTerms[i]
				res[l][domainSize-1].SetOne()
			} else {
				current := res[l][domainSize-d-2:]
				timesConst := multScratch[domainSize-d-2:]

				timesConst.Scale(&constTerms[i], current[1:]) //TODO: Directly double and add since constTerms are tiny? (even less than 4 bits)
				nonLeading := current[0 : d+1]

				nonLeading.Add(nonLeading, timesConst)

			}
			d++
		}

	}

	// We have pₗ(i≠l)=0. Now scale so that pₗ(l)=1
	// Replace the constTerms with norms
	for l := uint8(0); l < domainSize; l++ {
		constTerms[l].Neg(&constTerms[l])
		constTerms[l] = res[l].Eval(&constTerms[l])
	}
	constTerms = fr.BatchInvert(constTerms)
	for l := uint8(0); l < domainSize; l++ {
		res[l].ScaleInPlace(&constTerms[l])
	}

	return res
}

// InterpolateOnRange performs the interpolation of the given list of elements
// On the range [0, 1,..., len(values) - 1]
// TODO: Am I crazy or is this EXTRApolation and not INTERpolation
func InterpolateOnRange(values []fr.Element) Polynomial {
	nEvals := len(values)
	lagrange := getLagrangeBasis(nEvals)

	var res Polynomial
	res.Scale(&values[0], lagrange[0])

	temp := make(Polynomial, nEvals)

	for i := 1; i < nEvals; i++ {
		temp.Scale(&values[i], lagrange[i])
		res.Add(res, temp)
	}

	return res
}
