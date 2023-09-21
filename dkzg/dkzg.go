// Package kzg provides constructor for curved-typed KZG SRS
//
// For more details, see ecc/XXX/fr/kzg package
package dkzg

import (
	"io"

	"github.com/jparr721/krania-gnark-crypto/ecc"

	dkzg_bn254 "github.com/jparr721/krania-gnark-crypto/ecc/bn254/fr/dkzg"
)

// SRS ...
type SRS interface {
	io.ReaderFrom
	io.WriterTo
}

// NewSRS returns an empty curved-typed SRS object
// that implements io.ReaderFrom and io.WriterTo interfaces
func NewSRS(curveID ecc.ID) SRS {
	switch curveID {
	case ecc.BN254:
		return &dkzg_bn254.SRS{}
	default:
		panic("not implemented")
	}
}
