package src

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type Withdraw struct {
	L1DepositRoot frontend.Variable
	PubKey        [2]frontend.Variable
	Nonce         frontend.Variable
	Balance       frontend.Variable

	OrgMkProofHashs    [32]frontend.Variable
	OrgMkProofHashsPos [32]frontend.Variable

	WithdrawAmount frontend.Variable

	WithdrawMkProofHashs    [32]frontend.Variable
	WithdrawMkProofHashsPos [32]frontend.Variable

	OldRoot frontend.Variable
	NewRoot frontend.Variable
}

func (w *Withdraw) Define(api frontend.API) error {
	mimc, _ := mimc.NewMiMC(api)

	// TODO: verify sign

	// check original account exist
	mimc.Write(w.L1DepositRoot, w.PubKey[0], w.PubKey[1], w.Nonce, w.Balance)
	sum := mimc.Sum()
	for i := 1; i < len(w.OrgMkProofHashs); i++ {
		api.AssertIsBoolean(w.OrgMkProofHashsPos[i-1])
		d1 := api.Select(w.OrgMkProofHashsPos[i-1], sum, w.OrgMkProofHashs[i])
		d2 := api.Select(w.OrgMkProofHashsPos[i-1], w.OrgMkProofHashs[i], sum)
		mimc.Write(d1, d2)
		sum = mimc.Sum()
	}
	api.AssertIsEqual(sum, w.OldRoot)

	// doing withdraw
	newBalance := api.Sub(w.Balance, w.WithdrawAmount)
	newNonce := api.Add(w.Nonce, 1)

	// check orginal account node
	mimc.Write(w.L1DepositRoot, w.PubKey[0], w.PubKey[1], newNonce, newBalance)
	sum = mimc.Sum()
	for i := 1; i < len(w.OrgMkProofHashs); i++ {
		api.AssertIsBoolean(w.OrgMkProofHashsPos[i-1])
		d1 := api.Select(w.OrgMkProofHashsPos[i-1], sum, w.OrgMkProofHashs[i])
		d2 := api.Select(w.OrgMkProofHashsPos[i-1], w.OrgMkProofHashs[i], sum)
		mimc.Write(d1, d2)
		sum = mimc.Sum()
	}
	api.AssertIsEqual(sum, w.NewRoot)

	// check withdraw node
	mimc.Write(0, w.PubKey[0], w.PubKey[1], newNonce, w.WithdrawAmount)
	sum = mimc.Sum()
	for i := 1; i < len(w.WithdrawMkProofHashs); i++ {
		api.AssertIsBoolean(w.WithdrawMkProofHashsPos[i-1])
		d1 := api.Select(w.WithdrawMkProofHashsPos[i-1], sum, w.WithdrawMkProofHashs[i])
		d2 := api.Select(w.WithdrawMkProofHashsPos[i-1], w.WithdrawMkProofHashs[i], sum)
		mimc.Write(d1, d2)
		sum = mimc.Sum()
	}
	api.AssertIsEqual(sum, w.NewRoot)

	return nil
}
