pragma circom 2.1.2;

include "../circomlib/circuits/poseidon.circom";

template AccountHasher() {
    signal input l1DepositRoot;
    signal input pubKeyX;
    signal input pubKeyY;
    signal input nonce;
    signal input balance;
    
    signal output hash;

    component hasher = Poseidon(5);
    hasher.inputs[0] <== l1DepositRoot;
    hasher.inputs[1] <== pubKeyX;
    hasher.inputs[2] <== pubKeyY;
    hasher.inputs[3] <== nonce;
    hasher.inputs[4] <== balance;

    hash <== hasher.out;
}
