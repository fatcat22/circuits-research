pragma circom 2.1.2;

include "./circomlib/circuits/poseidon.circom";
include "./circomlib/circuits/comparators.circom";
include "./lib/binary_merkle_tree.circom";
include "./lib/account.circom";

template Withdraw(mktLevels) {
    signal input l1DepositRoot;
    signal input pubKey[2];
    signal input nonce;
    signal input balance;

    signal input orgMkProofHashs[mktLevels][1];
    signal input orgMkProofHashsPos[mktLevels];

    signal input withdrawAmount;

    signal input withdrawMkProofHashs[mktLevels][1];
    signal input withdrawMkProofHashsPos[mktLevels];

    signal input oldRoot;
    signal input newRoot;

    assert(withdrawAmount <= balance);

    // TODO: verify sign

    // check original account exist
    component orgLeaf = AccountHasher();
    orgLeaf.l1DepositRoot <== l1DepositRoot;
    orgLeaf.pubKeyX <== pubKey[0];
    orgLeaf.pubKeyY <== pubKey[1];
    orgLeaf.nonce <== nonce;
    orgLeaf.balance <== balance;

    component orgLeafExistChecker = CheckLeafExists(mktLevels);
    orgLeafExistChecker.root <== oldRoot;
    orgLeafExistChecker.leaf <== orgLeaf.hash;
    for (var i = 0; i < mktLevels; i++) {
        orgLeafExistChecker.pathElements[i][0] <== orgMkProofHashs[i][0];
        orgLeafExistChecker.pathIndex[i] <== orgMkProofHashsPos[i];
    }
    orgLeafExistChecker.enabled <== 1;

    // doing withdraw
    var newBalance = balance - withdrawAmount;
    var newNonce = nonce + 1;

    // check orginal account node
    component orgNewLeaf = AccountHasher();
    orgNewLeaf.l1DepositRoot <== l1DepositRoot;
    orgNewLeaf.pubKeyX <== pubKey[0];
    orgNewLeaf.pubKeyY <== pubKey[1];
    orgNewLeaf.nonce <== newNonce + 1;
    orgNewLeaf.balance <== newBalance;

    component orgNewLeafExistChecker = CheckLeafExists(mktLevels);
    orgNewLeafExistChecker.root <== newRoot;
    orgNewLeafExistChecker.leaf <== orgNewLeaf.hash;
    for (var i = 0; i < mktLevels; i++) {
        orgNewLeafExistChecker.pathElements[i][0] <== orgMkProofHashs[i][0];
        orgNewLeafExistChecker.pathIndex[i] <== orgMkProofHashsPos[i];
    }
    orgNewLeafExistChecker.enabled <== 1;

    // check withdraw node
    component withdrawLeaf = AccountHasher();
    withdrawLeaf.l1DepositRoot <== 0;
    withdrawLeaf.pubKeyX <== pubKey[0];
    withdrawLeaf.pubKeyY <== pubKey[1];
    withdrawLeaf.nonce <== newNonce;
    withdrawLeaf.balance <== withdrawAmount;

    component withdrawLeafExistChecker = CheckLeafExists(mktLevels);
    withdrawLeafExistChecker.root <== newRoot;
    withdrawLeafExistChecker.leaf <== withdrawLeaf.hash;
    for (var i = 0; i < mktLevels; i++) {
        withdrawLeafExistChecker.pathElements[i][0] <== withdrawMkProofHashs[i][0];
        withdrawLeafExistChecker.pathIndex[i] <== withdrawMkProofHashsPos[i];
    }
    withdrawLeafExistChecker.enabled <== 1;
}

component main {public [oldRoot, newRoot]} = Withdraw(32);