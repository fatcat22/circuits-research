pragma circom 2.1.2;

include "./lib/binary_merkle_tree.circom";
include "./lib/account.circom";

template Deposit(mktLevels) {
    signal input l1DepositRoot;
    signal input pubKey[2];
    signal input balance;

    signal input isNewAccount;
    signal input oldL1DepositRoot;
    signal input oldBalance;
    signal input oldNonce;
    signal input merkleProofHashs[mktLevels][1];
    signal input merkleProofHashsPos[mktLevels];

    signal input oldRoot;
    signal input newRoot;

    // check user exist if isNewAccount
    component oldAccountLeafHash = AccountHasher();
    oldAccountLeafHash.l1DepositRoot <== oldL1DepositRoot;
    oldAccountLeafHash.pubKeyX <== pubKey[0];
    oldAccountLeafHash.pubKeyY <== pubKey[1];
    oldAccountLeafHash.nonce <== oldNonce;
    oldAccountLeafHash.balance <== oldBalance;

    component accountExistChecker = CheckLeafExists(mktLevels);
    accountExistChecker.root <== oldRoot;
    accountExistChecker.leaf <== oldAccountLeafHash.hash;
    for (var i = 0; i < mktLevels; i++) {
        accountExistChecker.pathElements[i][0] <== merkleProofHashs[i][0];
        accountExistChecker.pathIndex[i] <== merkleProofHashsPos[i];
    }
    accountExistChecker.enabled <== isNewAccount;

    // check new merkle root
    component newAccountLeafHash = AccountHasher();
    newAccountLeafHash.l1DepositRoot <== l1DepositRoot;
    newAccountLeafHash.pubKeyX <== pubKey[0];
    newAccountLeafHash.pubKeyY <== pubKey[1];
    newAccountLeafHash.nonce <== (oldNonce + 1) * isNewAccount;
    newAccountLeafHash.balance <== oldBalance*isNewAccount + balance;

    component newRootChecker = CheckLeafExists(mktLevels);
    newRootChecker.root <== newRoot;
    newRootChecker.leaf <== newAccountLeafHash.hash;
    for (var i = 0; i < mktLevels; i++) {
        newRootChecker.pathElements[i][0] <== merkleProofHashs[i][0];
        newRootChecker.pathIndex[i] <== merkleProofHashsPos[i];
    }
    newRootChecker.enabled <== 1;

}

component main {public [oldRoot, newRoot]} = Deposit(32);