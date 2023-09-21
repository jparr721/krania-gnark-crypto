package main

import (
	"fmt"

	"github.com/jparr721/krania-gnark-crypto/internal/field"
	"github.com/jparr721/krania-gnark-crypto/internal/field/generator"
)

//go:generate go run main.go
func main() {
	const modulus = "0xFFFFFFFF00000001"
	goldilocks, err := field.NewFieldConfig("goldilocks", "Element", modulus, true)
	if err != nil {
		panic(err)
	}
	if err := generator.GenerateFF(goldilocks, "../"); err != nil {
		panic(err)
	}
	fmt.Println("successfully generated goldilocks field")
}
