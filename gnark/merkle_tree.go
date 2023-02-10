package src

import (
	"github.com/consensys/gnark/frontend"
)

type MerkleTree struct {
	Leaf         frontend.Variable
	PathIndex    []frontend.Variable
	PathIndexPos []frontend.Variable
	Root         frontend.Variable
}
