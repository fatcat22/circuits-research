package src

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash/mimc"
)

type Deposit struct {
	L1DepositRoot frontend.Variable
	PubKey        [2]frontend.Variable
	Balance       frontend.Variable

	IsNewAccount     frontend.Variable
	OldL1DepositRoot frontend.Variable
	OldNonce         frontend.Variable
	OldBalance       frontend.Variable

	MerklePath, MerkleHelper [32]frontend.Variable

	OldRoot, NewRoot frontend.Variable
}

func (dp *Deposit) Define(api frontend.API) error {
	mimc, _ := mimc.NewMiMC(api)

	// check user exist if isNewAccount
	mimc.Write(dp.OldL1DepositRoot, dp.PubKey[0], dp.PubKey[1], dp.OldNonce, dp.OldBalance)
	oldLeafHash := mimc.Sum()

	sum := oldLeafHash
	for i := 1; i < len(dp.MerklePath); i++ {
		api.AssertIsBoolean(dp.MerkleHelper[i-1])
		d1 := api.Select(dp.MerkleHelper[i-1], sum, dp.MerklePath[i])
		d2 := api.Select(dp.MerkleHelper[i-1], dp.MerklePath[i], sum)
		mimc.Write(d1, d2)
		sum = mimc.Sum()
	}
	api.AssertIsEqual(oldLeafHash, api.Select(dp.IsNewAccount, oldLeafHash, sum))

	// check new merkle root
	mimc.Write()
	leaf := []frontend.Variable{dp.L1DepositRoot, dp.PubKey[0], dp.PubKey[1], 0, dp.Balance}
	dp.MerklePath[0] = leaf
	//merklePath := append([]frontend.Variable{leaf}, dp.MerklePath...)
	merkle.VerifyProof(api, mimc, dp.NewRoot, dp.MerklePath[:], dp.MerkleHelper[:])

	return nil
}
