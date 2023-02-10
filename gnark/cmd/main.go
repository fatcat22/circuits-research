package main

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	research "github.com/fatcat22/circuits-research/gnark"
)

func main() {
	deposit := research.Deposit{
		L1DepositRoot: 0,
		PubKey:        [2]frontend.Variable{0, 0},
		Balance:       0,

		IsNewAccount:     0,
		OldL1DepositRoot: 0,
		OldNonce:         0,
		OldBalance:       0,

		MerklePath: [32]frontend.Variable{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		MerkleHelper: [32]frontend.Variable{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},

		OldRoot: 0,
		NewRoot: 0,
	}
	ccsDeposit, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &deposit)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("%#v\n", ccsDeposit)
	fmt.Printf("NbConstraints: %d\n", ccsDeposit.GetNbConstraints())

	withdraw := research.Withdraw{
		L1DepositRoot: 0,
		PubKey:        [2]frontend.Variable{0, 0},
		Nonce:         0,
		Balance:       0,
		OrgMkProofHashs: [32]frontend.Variable{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		OrgMkProofHashsPos: [32]frontend.Variable{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		WithdrawAmount: 0,
		WithdrawMkProofHashs: [32]frontend.Variable{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		WithdrawMkProofHashsPos: [32]frontend.Variable{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		OldRoot: 0,
		NewRoot: 0,
	}

	ccsWithdraw, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &withdraw)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("%#v\n", ccsWithdraw)
	fmt.Printf("NbConstraints: %d\n", ccsWithdraw.GetNbConstraints())
}
