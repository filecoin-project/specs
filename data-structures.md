# Data Structures

This document serves as an entry point for understanding the data structures in Filecoin. In this document, we give an introduction to each of the protocol data structures (linking to more detailed data structure specs in some places) and then explain how to encode these data structures for use in other parts of Filecoin such as the network protocols and the blockchain.

## Address

An address is an identifier that refers to an actor in the Filecoin state. All [actors](actors.md) (miner actors, the storage market actor, account actors) have an address. An address encodes information about:
- Network this address belongs to
- Type of data the address contains
- The data itself
- Checksum (depending on the type of address)

For more detail about the different types of addresses and how they are structured and used, take a look at the [address spec](address.md).

## Content Identifier (CID)

Most objects in Filecoin are referenced by their Content Identifiers (CID for short). The CID is a hash value, prefixed with its hash function (multihash), prepended with a few extra labels to inform applications about how to deserialize the given data. 

For more information about CIDs, take a look at the [CID Spec](https://github.com/ipld/cid).

## Block Header

A block header contains information relevant to a particular point in time over which the network may achieve consensus. The block header contains:
- The address of the miner that mined the block
- An array of the tickets that led to this particular miner being selected as the leader for this round (see the [Secret Leader Election portion of the Expected Consensus spec](expected-consensus.md#secret-leader-election) for more) and a signature on the winning ticket
- The set of parent blocks and aggregate [chain weight](expected-consensus.md#chain-weighting) of the parents
- This block's height
- Merkle root of the state tree (after applying the messages -- state transitions -- included in this block)
- Merkle root of the messages (state transitions) in this block
- Merkle root of the message receipts in this block
- Timestamp

Note: A block is functionally the same as a block header in the Filecoin protocol. While a block header contains Merkle links to the full system state, messages, and message receipts, a block can be thought of as the full set of this information (not just the Merkle roots, but rather the full data of the state tree, message tree, receipts tree, etc.). Because a full block is quite large, our chain consists of block headers rather than full blocks. We often use the terms `block` and `block header` interchangeably.

Below is a sample interface (written in Go) for defining a block header.

```go
// BlockHeader is a block header in the blockchain.
type BlockHeader struct {
	// Miner is the address of the miner actor that mined this block.
	Miner Address

	// Tickets are the winning ticket that were submitted with this block. The tickets
	// array should include all the intervening tickets generated between (but not including) the 
	// parents' block height and the current block height.
	Tickets []Ticket

	// ElectionProof is a signature over the final ticket that proves this miner
	// is the leader at this round
	ElectionProof Signature

	// Parents is the set of parents this block was based on. Typically one,
	// but can be several in the case where there were multiple winning ticket-
	// holders for a round.
	Parents []Cid

	// ParentWeight is the aggregate chain weight of the parent set.
	ParentWeight Integer

	// Height is the chain height of this block.
	Height Uint64

	// StateRoot is a cid pointer to the state tree after application of the
	// transactions state transitions.
	StateRoot Cid

	// Messages is the set of messages included in this block. This field is the Cid
	// of the root of a sharray of Messages.
	Messages Cid

	// MessageReceipts is a set of receipts matching to the sending of the `Messages`.
	// This field is the Cid of the root of a sharray of MessageReceipts.
	MessageReceipts Cid

	// The block Timestamp is used to enforce a form of block delay by honest miners.
	// Unix time UTC timestamp stored as an unsigned integer
	Timestamp Timestamp
}
```

### Ordering of Elements in Block Header Arrays

As you see in the interface defined above, a typical block header contains many different arrays. The ordering of these arrays is important and should be consistent across implementations.

| Array | Ordering |
| --- | --- |
| `Tickets` | Sorted by the order in which the tickets were created (each ticket in the array is created at a specific block height). |
| `Parents` | Sorted lexicographically |

### Sharded Messages and Receipts

In a block header, the `Messages` and `MessageReceipts` fields are each CIDs that reference [sharray](sharray.md) datastructures. The `Messages` sharray contains the CIDs of the messages that are included in the block. The `MessageReceipts` sharray contains the message receipts directly.

Read more about sharrays in the [sharray spec](sharray.md).

## Message

```go
type Message struct {
	To   Address
	From Address

	// When receiving a message from a user account the nonce in
	// the message must match the expected nonce in the from actor.
	// This prevents replay attacks.
	Nonce Uint64

	Value BigInteger

	GasPrice Integer
	GasLimit Integer

	Method uint64
	Params []byte
}
```

### Parameter Encoding

Parameters to methods get encoded as described in the [basic types](#basic-type-encodings) section below, and then put into a CBOR encoded array.
(TODO: thinking about this, it might make more sense to just have `Params` be an array of things)

### Signing

A signed message is a wrapper type over the base message.

```go
type SignedMessage struct {
	Message   Message
	Signature Signature
}
```

The signature is a serialized signature over the serialized base message. For more details on how the signature itself is done, see the [signatures spec](signatures.md).

## Message Receipt

```go
type MessageReceipt struct {
	ExitCode uint8

	Return []byte

	GasUsed Integer
}
```

### Serialization

Message receipts are currently serialized simply by CBOR marshaling them, using lower-camel-cased field names.



## Actor

```go
type Actor struct {
	// Code is a pointer to the code object for this actor
	Code Cid

	// Head is a pointer to the root of this actors state
	Head Cid

	// Nonce is a counter of the number of messages this actor has sent
	Nonce Uint64

	// Balance is this actors current balance of filecoin
	Balance BigInteger
}
```

## State Tree

The state trie keeps track of all state in Filecoin. It is a map of addresses to `actors` in the system. It is implemented using a HAMT.

## HAMT

{{% notice todo %}}
**TODO**: link to spec for our CHAMP HAMT
{{% /notice %}}


## Signature

All signatures in Filecoin come with a type that signifies which key type was used to create the signature.

```go
type Signature struct {
	Type int
	Data []byte
}
```

### `Type` Values

| Key Type | Value |
|-------|----------|
|  Secp256k1 | 1 |
| BLS12-381 ECDSA | 2 |

### Serialization

`<uvarint(Type)><Data>`

Note: As signatures should always be within wrapper types, length prefixing is not needed here.



# Basic Type Encodings


## CID Serialization

CIDs are serialized by applying binary multibase encoding, then encoding that as a CBOR byte array with a tag of 42.

Types that appear in messages or in state must be encoded as described here.

#### `PublicKey`

The public key type is simply an array of bytes. (TODO: discuss specific encoding of key types, for now just calling it bytes is sufficient)

#### `BytesAmount`
BytesAmount is just a re-typed Integer.

#### `PeerID`
PeerID is just the serialized bytes of a libp2p peer ID.

Spec incomplete, take a look at this PR: https://github.com/libp2p/specs/pull/100

#### `Integer`

Integers are encoded as LEB128 signed integers.

#### `BitField`

Bitfields are a set of bits. Encoding still TBD, but it needs to be very compact. We can assume that most often, ranges of bits will be set, or not set, and use that to our advantage here. Some form of run length encoding may work well.

#### `SectorSet`

TODO

#### `FaultSet`

A fault set is a BitField and a block height, encoding TBD.

#### `BlockHeader`

BlockHeader is a serialized `Block`.

#### `SealProof`

SealProof is a 384-element array of bytes.

#### `PoStProof`

PoStProof is a 192-element array of bytes.

#### `TokenAmount`

TokenAmount is a re-typed Integer.



## LEB128 Encoding Reference


This is taken from the Dwarf Standard 4, Appendix C

#### Encode unsigned LEB128

```c
do
{
  byte = low order 7 bits of value;
  value >>= 7;
  if (value != 0) /* more bytes to come */
    set high order bit of byte;
  emit byte;
} while (value != 0);

```

#### Encode signed LEB128

```c
more = 1;
negative = (value < 0);
size = no. of bits in signed integer;
while(more)
{
  byte = low order 7 bits of value;
  value >>= 7;
  /* the following is unnecessary if the
   * implementation of >>= uses an arithmetic rather
   * than logical shift for a signed left operand
   */
  if (negative)
    /* sign extend */
    value |= - (1 << (size - 7));
    /* sign bit of byte is second high order bit (0x40) */
  if ((value ==  0 && sign bit of byte is clear) ||
      (value == -1 && sign bit of byte is set))
     more = 0;
  else
    set high order bit of byte;
  emit byte;
}
```

#### Decode unsigned LEB128

```c
result = 0;
shift = 0;
while(true)
{
  byte = next byte in input;
  result |= (low order 7 bits of byte << shift);
  if (high order bit of byte == 0)
    break;
  shift += 7;
}
```

#### Decode signed LEB128

```c
result = 0;
shift = 0;
size = number of bits in signed integer;
while(true)
{
  byte = next byte in input;
  result |= (low order 7 bits of byte << shift);
  shift += 7;
  /* sign bit of byte is second high order bit (0x40) */
  if (high order bit of byte == 0)
  break;
}
if ((shift <size) && (sign bit of byte is set))
  /* sign extend */
  result |= - (1 << shift);
```

# Filecoin Compact Serialization

Datastructures in Filecoin are encoded as compactly as is reasonable. At a high level, each object is converted into an ordered array of its fields (ordered by their appearance in the struct declaration), then CBOR marshaled, and prepended with an object type tag.

| FCS Type | tag  |
|---|---|
| block v1 | 43  |
| message v1 | 44 |
| signedMessage v1 | 45 |

For example, a message would be encoded as:

```
tag<44>[
  msg.To,
  msg.From,
  msg.Nonce,
  msg.Value,
  msg.GasPrice,
  msg.GasLimit,
  msg.Method,
  msg.Params
]
```

Each individual type should be encoded as specified:

| type | encoding |
| --- | ---- |
| Uint64 | CBOR major type 0 |
| BigInteger | [CBOR bignum](https://tools.ietf.org/html/rfc7049#section-2.4.2) |
| Address | CBOR major type 2 |
| Uint8 | CBOR Major type 0 |
| []byte | CBOR Major type 2 |
| string | CBOR Major type 3 |
| bool | [CBOR Major type 7, value 20/21](https://tools.ietf.org/html/rfc7049#section-2.3) |

## Encoding Considerations

Objects should be encoded using [canonical CBOR](https://tools.ietf.org/html/rfc7049#section-3.9), and decoders should operate in [strict mode](https://tools.ietf.org/html/rfc7049#section-3.10).  The maximum size of an FCS Object should be 1MB (2^20 bytes). Objects larger than this are invalid.

Additionally, CBOR Major type 5 is not used. If an FCS object contains it, that object is invalid.

## IPLD Considerations

Cids for FCS objects should use the FCS multicodec (`0x1f`), and should use a blake2b-256 multihash.

## Vectors

Below are some sample vectors for each data type.

### Message

Encoded: `d82c865501fd1d0f4dfcd7e99afcb99a8326b7dc459d32c6285501b882619d46558f3d9e316d11b48dcf211327026a1875c245037e11d600666d6574686f644d706172616d73617265676f6f64`

Decoded:
```
To:     Address("f17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy")
From:   Address("f1xcbgdhkgkwht3hrrnui3jdopeejsoatkzmoltqy")
Nonce:  uint64(117)
Value:  BigInt(15000000000)
Method: string("method")
Params: []byte("paramsaregood")
```

### Block

Encoded: `d82b895501fd1d0f4dfcd7e99afcb99a8326b7dc459d32c628814a69616d617469636b6574566920616d20616e20656c656374696f6e2070726f6f6681d82a5827000171a0e40220ce25e43084e66e5a92f8c3066c00c0eb540ac2f2a173326507908da06b96f678c242bb6a1a0012d687d82a5827000171a0e40220ce25e43084e66e5a92f8c3066c00c0eb540ac2f2a173326507908da06b96f6788080`

Decoded:
```
Miner:           Address("f17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy")
Tickets:         [][]byte{"iamaticket"}
ElectionProof:   []byte("i am an election proof")
Parents:         []Cid{"zDPWYqFD5abn4FyknPm1PibXdJ2kwRNVPDabKyzfdXVJGjnDuq4B"}
ParentWeight:    NewInt(47978)
Height:          uint64(1234567)
StateRoot:       Cid("zDPWYqFD5abn4FyknPm1PibXdJ2kwRNVPDabKyzfdXVJGjnDuq4B")
Messages:        []SignedMessage{}
MessageReceipts: []MessageReceipt{}
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
      <blocks> ::= <block_single> | <block_short> | <block_long>
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
