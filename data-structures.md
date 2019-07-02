# Data Structures

This document serves as an entry point for the data structures in filecoin.

## Address

An address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. It encodes information about the network it belongs to, the type of data it contains, the data itself, and depending on the type, a checksum.

The details are described in [address spec](address.md).

## Block

A block represents the data at individual point in time that the network may achieve consensus on. It contains (via merkle links) the full state of the system, references to the previous state, and the notion of a 'weight' for deciding which block is the 'best'.

```sh
type Block struct {
	## Miner is the address of the miner actor that mined this block.
	miner Address

	## Tickets is a chain (possibly singleton) of tickets ending with a winning ticket submitted with this block.
	tickets [Ticket]

	## ElectionProof is a signature over the final ticket that proves this miner
	## is the leader at this round
	electionProof Signature

	## Parents is an array of distinct CIDs of parents on which this block was based.
	## Typically one, but can be several in the case where there were multiple winning ticket-
	## holders for a round.
	## The order of parent CIDs is not defined.
	parents [&Block]

	## ParentWeight is the aggregate chain weight of the parent set.
	parentWeight UInt

	## Height is the chain height of this block.
	height UInt

	## StateRoot is a cid pointer to the state tree after application of the
	## transactions state transitions.
	stateRoot &StateTree

	## Messages is the set of messages included in this block. This field is the Cid
	## of the root of a sharray of Messages.
	messages &[&Message]<Sharray>

	## BLSAggregate is an aggregated BLS signature for all the messages in this block that
	## were signed using BLS signatures
	blsAggregate Signature

	## MessageReceipts is a set of receipts matching to the sending of the `Messages`.
	## This field is the Cid of the root of a sharray of MessageReceipts.
	messageReceipts &[MessageReceipt]

	## The block Timestamp is used to enforce a form of block delay by honest miners.
	## Unix time UTC timestamp (in seconds) stored as an unsigned integer.
	timestamp Timestamp

	## BlockSig is a signature over the hash of the entire block with the miners
	## worker key to ensure that it is not tampered with after creation
	blockSig Signature
} representation tuple
```

## Message

```sh
type Message union {
    | UnsignedMessage 0
    | SignedMessage 1
} representation keyed
```

```sh
type UnsignedMessage struct {
	to   Address
	from Address

	## When receiving a message from a user account the nonce in the message must match the expected
    ## nonce in the from actor. This prevents replay attacks.
	nonce UInt

	value UInt

	gasPrice UInt
	gasLimit UInt

	method &ActorMethod
} representation tuple
```

```sh
type SignedMessage struct {
	message   UnsignedMessage
	signature Signature
} representation tuple
```

The signature is a serialized signature over the serialized base message. For more details on how the signature itself is done, see the [signatures spec](signatures.md).

## State Tree

The state tree keeps track of all state in Filecoin. It is a map of addresses to `actors` in the system.
The `ActorState` is defined in the [actors spec](actors.md).

```sh
type StateTree map {ID:Actor}<Hamt>
```

## Message Receipt

```sh
type MessageReceipt struct {
	exitCode UInt
	return Bytes
	gasUsed UInt
} representation tuple
```

## Actor

```sh
type Actor struct {
	## Cid of the code object for this actor.
	code Cid

	## Reference to the root of this actors state.
	head &ActorState

	## Counter of the number of messages this actor has sent.
	nonce UInt

	## Current balance of filecoin of this actor.
	balance UInt
}
```

## Signature

All signatures in Filecoin come with a type that signifies which key type was used to create the signature.

For more details on signature creation, see [the signatures spec](signatures.md).

```sh
type Signature union {
	| Secp256k1Signature 0
	| Bls12_381Signature 1
} representation byteprefix

type Secp256k1Signature Bytes
type Bls12_381Signature Bytes
```

## FaultSet

FaultSets are used to denote which sectors failed at which block height.

```sh
type FaultSet struct {
	index    UInt
	bitField BitField
}
```

The `index` field is a block height offset from the start of the miners proving period (in order to make it more compact).


## Basic Types

### CID

For most objects referenced by Filecoin, a Content Identifier (CID for short) is used. This is effectively a hash value, prefixed with its hash function (multihash) prepended with a extra labels to inform applications about how to deserialize the given data. [CID Spec](https://github.com/ipld/cid) contains the detailed spec.

### Timestamp

```sh
type Timestamp UInt
```

### PublicKey

The public key type is simply an array of bytes.

```sh
type PublicKey Bytes
```

### BytesAmount

BytesAmount is just a re-typed Integer.
```sh
type BytesAmount UInt
```

### PeerId

The serialized bytes of a libp2p peer ID.

Spec incomplete, take a look at this PR: https://github.com/libp2p/specs/pull/100

```sh
type PeerId Bytes
```

### Bitfield

Bitfields are a set encoded using a custom run length encoding: RLE+.

```sh
type Bitfield Bytes
```

### SectorSet

{{% notice todo %}}
Define me
{{% /notice %}}


### SealProof

SealProof is an opaque, dynamically-sized array of bytes.

### PoStProof

PoStProof is an opaque, dynamically-sized array of bytes.

### TokenAmount

A type to represent an amount of filecoin tokens.

```sh
type TokenAmount UInt
```

## RLE+ Bitset Encoding

RLE+ is a lossless compression format based on [RLE](https://en.wikipedia.org/wiki/Run-length_encoding).
It's primary goal is to reduce the size in the case of many individual bits, where RLE breaks down quickly,
while keeping the same level of compression for large sets of contiugous bits.

In tests it has shown to be more compact than RLE iteself, as well as [Concise](https://arxiv.org/pdf/1004.0403.pdf) and [Roaring](https://roaringbitmap.org/).

### Format

The format consists of a header, followed by a series of blocks, of which there are three different types.

The format can be expressed as the following [BNF](https://en.wikipedia.org/wiki/Backus%E2%80%93Naur_form) grammar.

```
    <encoding> ::= <header> <blocks>
      <header> ::= <version> <bit>
     <version> ::= "00"
      <blocks> ::= <block> <blocks> | ""
       <block> ::= <block_single> | <block_short> | <block_long>
<block_single> ::= "1"
 <block_short> ::= "01" <bit> <bit> <bit> <bit>
  <block_long> ::= "00" <unsigned_varint>
         <bit> ::= "0" | "1"
```

An `<unsigned_varint>` is defined as specified [here](https://github.com/multiformats/unsigned-varint).

#### Header

The header indiciates the very first bit of the bit vector to encode. This means the first bit is always
the same for the encoded and non encoded form.

#### Blocks

The blocks represent how many bits, of the current bit type there are. As `0` and `1` alternate in a bit vector
the inital bit, which is stored in the header, is enough to determine if a length is currently referencing
a set of `0`s, or `1`s.

##### Block Single

If the running length of the current bit is only `1`, it is encoded as a single set bit.

##### Block Short

If the running length is less than `16`, it can be encoded into up to four bits, which a short block
represents. The length is encoded into a 4 bits, and prefixed with `01`, to indicate a short block.

##### Block Long

If the running length is `16` or larger, it is encoded into a varint, and then prefixed with `00` to indicate
a long block.


> **Note:** The encoding is unique, so no matter which algorithm for encoding is used, it should produce
> the same encoding, given the same input.

## Other Considerations

- The maximum size of an Object should be 1MB (2^20 bytes). Objects larger than this are invalid.
- Hashes should use a blake2b-256 multihash.
