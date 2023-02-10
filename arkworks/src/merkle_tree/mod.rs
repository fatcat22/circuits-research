use crate::crh::TwoToOneCRHScheme;
use crate::{crh::CRHScheme, Error};
use ark_serialize::{CanonicalDeserialize, CanonicalSerialize};
use ark_std::borrow::Borrow;
use ark_std::hash::Hash;

/// Convert the hash digest in different layers by converting previous layer's output to
/// `TargetType`, which is a `Borrow` to next layer's input.
pub trait DigestConverter<From, To: ?Sized> {
    type TargetType: Borrow<To>;
    fn convert(item: From) -> Result<Self::TargetType, Error>;
}

/// A trivial converter where digest of previous layer's hash is the same as next layer's input.
pub struct IdentityDigestConverter<T> {
    _prev_layer_digest: T,
}

/// Merkle tree have three types of hashes.
/// * `LeafHash`: Convert leaf to leaf digest
/// * `TwoLeavesToOneHash`: Convert two leaf digests to one inner digest. This one can be a wrapped
/// version `TwoHashesToOneHash`, which first converts leaf digest to inner digest.
/// * `TwoHashesToOneHash`: Compress two inner digests to one inner digest
pub trait Config {
    type Leaf: ?Sized; // merkle tree does not store the leaf
                       // leaf layer
    type LeafDigest: Clone
        + Eq
        + core::fmt::Debug
        + Hash
        + Default
        + CanonicalSerialize
        + CanonicalDeserialize;
    // transition between leaf layer to inner layer
    type LeafInnerDigestConverter: DigestConverter<
        Self::LeafDigest,
        <Self::TwoToOneHash as TwoToOneCRHScheme>::Input,
    >;
    // inner layer
    type InnerDigest: Clone
        + Eq
        + core::fmt::Debug
        + Hash
        + Default
        + CanonicalSerialize
        + CanonicalDeserialize;

    // Tom's Note: in the future, if we want different hash function, we can simply add more
    // types of digest here and specify a digest converter. Same for constraints.

    /// leaf -> leaf digest
    /// If leaf hash digest and inner hash digest are different, we can create a new
    /// leaf hash which wraps the original leaf hash and convert its output to `Digest`.
    type LeafHash: CRHScheme<Input = Self::Leaf, Output = Self::LeafDigest>;
    /// 2 inner digest -> inner digest
    type TwoToOneHash: TwoToOneCRHScheme<Output = Self::InnerDigest>;
}
