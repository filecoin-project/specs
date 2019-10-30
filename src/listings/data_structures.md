---
title: "Data Structures"
---


# Address

An address is an identifier that refers to an actor in the Filecoin state. All [actors](actors.md) (miner actors, the storage market actor, account actors) have an address. An address encodes information about:

- Network this address belongs to
- Type of data the address contains
- The data itself
- Checksum (depending on the type of address)

For more detail, see the full {{<sref app_address "address spec">}}.


# Block

A block header contains information relevant to a particular point in time over which the network may achieve consensus. The block header contains:

- The address of the miner that mined the block
- A ticket associated to this block's creation to be used as randomness elsewhere in the protocol (see
{{<sref leader_election>}} for more details)
- An electionProof showing this miner was eligible to mine
- The set of parent blocks and aggregate {{<sref chain_selection "chain weight">}} of the parents
- This block's height
- Merkle root of the state tree (after applying the messages -- state transitions -- included in this block)
- Merkle root of the messages (state transitions) in this block
- Merkle root of the message receipts in this block
- Timestamp

{{% notice note %}}
**Note:** A block is functionally the same as a block header in the Filecoin protocol. While a block header contains Merkle links to the full system state, messages, and message receipts, a block can be thought of as the full set of this information (not just the Merkle roots, but rather the full data of the state tree, message tree, receipts tree, etc.). Because a full block is quite large, our chain consists of block headers rather than full blocks. We often use the terms `block` and `block header` interchangeably.
{{% /notice %}}

{{<goFile Block>}}


# Tipset

For more on Tipsets, see {{<sref expected_consensus "the Expected Consensus spec">}}. Implementations may choose not to create a Tipset data structure, instead representing its operations in terms of the underlying blocks.

{{<goFile TipSet>}}


# VRF Personalization

We define VRF personalizations as follow, to enable domain separation across operations that make use of the same VRF (e.g. `Ticket` and
`ElectionProof`).

| Type          | Prefix |
| ------------- | ------ |
| Ticket        | `0x01` |
| ElectionProof | `0x02` |


# Ticket

A ticket contains a shared random value referenced by a particular `Block` in the Filecoin blockchain.
Every miner must produce a new `Ticket` for each ticket used in a leader election attempt.

To produce the ticket values,
we use an [EC-VRF per Goldberg et al.](https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#page-10)
with Secp256k1 and SHA-256 to obtain a deterministic, pseudorandom output.

{{<goFile Ticket>}}

## Ticket Comparison

The ticket is represented concretely by the `Ticket` data structure.
Whenever the Filecoin protocol refers to ticket values
(notably in crafting {{<sref post "PoSTs">}} or running leader election),
what is meant is that the bytes of the `VRFResult` field in the `Ticket` struct are used.
Specifically, tickets are compared lexicographically,
interpreting the bytes of the `VRFResult.Output` as an unsigned integer value (little-endian).

# ElectionProof

An election proof is generated from a past ticket (chosen based on public network parameters)
and a given epoch number (or height) by a miner during the leader election process.
Its output value determines whether the miner is elected as one of the leaders,
and hence is eligible to produce a block for the current epoch.
The inclusion of the `ElectionProof` in the block allows other network participants
to verify that the block was mined by a valid leader. With every leader election attempt for a given ticket,
(in cases where no blocks are found in a round) a miner increments the epoch value thus increasing
block height with every leader election attempt.

{{<goFile ElectionProof>}}


# Message

`Message` data structures in Filecoin describe operations that can be performed on the Filecoin VM state
(e.g., FIL transactions between accounts).
To facilitate the process of producing secure protocol implementations,
we explicitly distinguish between
{{<sref crypto_signatures "signed and unsigned">}} `Message` structures.

{{<goFile Message>}}
{{<goFile UnsignedMessage>}}
{{<goFile SignedMessage>}}
{{<goFile MessageReceipt>}}


# State Tree

The state tree keeps track of the entire state of the {{<sref vm>}} at any given point.
It is a map from `Address` structures to `Actor` structures, where each `Actor`
may also contain some additional `ActorState` that is specific to a given actor
type.

{{<goFile StateTree>}}


# Actor

{{<goFile Actor>}}


# Signature

{{<sref crypto_signatures "Cryptographic signatures">}} in Filecoin are represented
as byte arrays, and come with a tag that signifies what key type was used to create
the signature.

{{<goFile Signature>}}


# FaultSet

`FaultSet` data structures are used to denote which sectors failed at which block heights.

{{<goFile FaultSet>}}

In order to make the serialization more compact,
the `index` field denotes a block height offset from the start of the corresponding
miner's proving period.


# Basic Types

## CID
For most objects referenced by Filecoin, a Content Identifier (CID for short) is used.
This is effectively a hash value, prefixed with its hash function (multihash)
as well as extra labels to inform applications about how to deserialize the given data.
For a more detailed specification, we refer the reader to the
[IPLD repository](https://github.com/ipld/cid).


## Timestamp
{{<goFile Timestamp>}}


## PublicKey
The public key type is simply an array of bytes.
{{<goFile PublicKey>}}


## BytesAmount
BytesAmount is just a re-typed Integer.
{{<goFile BytesAmount>}}


## PeerId
The serialized bytes of a libp2p peer ID.
{{% todo %}} Spec incomplete; take a look at [this PR](https://github.com/libp2p/specs/pull/100).{{% /todo %}}
{{<goFile PeerId>}}


## Bitfield
Bitfields are a set encoded using a custom run length encoding: RLE+.
{{<goFile Bitfield>}}


## SectorSet
A sector set stores a mapping of sector IDs to the respective `commR`s.
{{<goFile SectorSet>}}

{{% todo %}}
Improve on this; see https://github.com/filecoin-project/specs/issues/116.
{{% /todo %}}


## SealProof
SealProof is an opaque, dynamically-sized array of bytes.
{{<goFile SealProof>}}


## PoSTProof
PoSTProof is an opaque, dynamically-sized array of bytes.
{{<goFile PoSTProof>}}


## TokenAmount
A type to represent an amount of Filecoin tokens.
{{<goFile TokenAmount>}}


## SectorID
Uniquely identifies a miner's sector.
{{<goFile SectorID>}}


# RLE+ Bitset Encoding

RLE+ is a lossless compression format based on [RLE](https://en.wikipedia.org/wiki/Run-length_encoding).
Its primary goal is to reduce the size in the case of many individual bits, where RLE breaks down quickly,
while keeping the same level of compression for large sets of contiugous bits.

In tests it has shown to be more compact than RLE itself, as well as [Concise](https://arxiv.org/pdf/1004.0403.pdf) and [Roaring](https://roaringbitmap.org/).

## Format

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

### Header

The header indicates the very first bit of the bit vector to encode. This means the first bit is always
the same for the encoded and non-encoded form.

### Blocks

The blocks represent how many bits, of the current bit type there are. As `0` and `1` alternate in a bit vector
the inital bit, which is stored in the header, is enough to determine if a length is currently referencing
a set of `0`s, or `1`s.

#### Block Single

If the running length of the current bit is only `1`, it is encoded as a single set bit.

#### Block Short

If the running length is less than `16`, it can be encoded into up to four bits, which a short block
represents. The length is encoded into a 4 bits, and prefixed with `01`, to indicate a short block.

#### Block Long

If the running length is `16` or larger, it is encoded into a varint, and then prefixed with `00` to indicate
a long block.

> **Note:** The encoding is unique, so no matter which algorithm for encoding is used, it should produce
> the same encoding, given the same input.

#### Bit Numbering

For Filecoin, byte arrays representing RLE+ bitstreams are encoded using [LSB 0](https://en.wikipedia.org/wiki/Bit_numbering#LSB_0_bit_numbering) bit numbering.


# Other Considerations

- The maximum size of an Object should be 1MB (2^20 bytes). Objects larger than this are invalid.
- Hashes should use a blake2b-256 multihash.
