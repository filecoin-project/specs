# Data Structures

This document serves as an entry point for understanding all of the data structures in filecoin.

TODO: this should also include, or reference, how each data structure is serialized precisely.

## Address

An address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. An address encodes information about the network it belongs to, the type of data it contains, the data itself, and depending on the type, a checksum.

```go
type Address struct {
    // 0: main-net
    // 1: test-net
    network byte

    // 0: SECP256K1 Public Key
    // 1: ActorID
    // 2: Builtin Actor
    // 3: BLS Public Key
    typ byte

    // raw bytes associated with typ
    data []byte
}
```

## CID

For most objects referenced by Filecoin, a Content Identifier (CID for short) is used. This is effectively a hash value, prefixed with its hash function (multihash) prepended with a few extra labels to inform applications about how to deserialize the given data. To learn more, take a look at the [CID Spec](https://github.com/ipld/cid). 

CIDs are serialized by applying binary multibase encoding, then encoding that as a CBOR byte array with a tag of 42.

## Block

A block represents an individual point in time that the network may achieve consensus on. It contains (via merkle links) the full state of the system, references to the previous state, and some notion of a 'weight' for deciding which block is the 'best'.

```go
// Block is a block in the blockchain.
type Block struct {
	// Miner is the address of the miner actor that mined this block.
	Miner Address

	// Ticket is the winning ticket that was submitted with this block.
	Ticket Signature

	// Parents is the set of parents this block was based on. Typically one,
	// but can be several in the case where there were multiple winning ticket-
	// holders for a round.
	Parents []Cid

	// ParentWeightNum is the numerator of the aggregate chain weight of the parent set.
	ParentWeightNum Integer

	// ParentWeightDenom is the denominator of the aggregate chain weight of the parent set
	ParentWeightDenom Integer 

	// Height is the chain height of this block.
	Height Uint64
    
    // StateRoot is a cid pointer to the state tree after application of the
	// transactions state transitions.
	StateRoot Cid

	// Messages is the set of messages included in this block
	// TODO: should be a merkletree-ish thing
	Messages []SignedMessage

	// MessageReceipts is a set of receipts matching to the sending of the `Messages`.
    // TODO: should be the same type of merkletree-list thing that the messages are
	MessageReceipts []MessageReceipt
}
```

### Serialization

Blocks are currently serialized simply by CBOR marshaling them, using lower-camel-cased field names.

## Message

```go
type Message struct {
	To   Address
	From Address
	
	// When receiving a message from a user account the nonce in
	// the message must match the expected nonce in the from actor.
	// This prevents replay attacks.
	Nonce Integer

	Value Integer

	Method string
	
	Params []byte
}
```

### Parameter Encoding

Parameters to methods get encoded as described in the [basic types](#basic-type-encodings) section below, and then put into a CBOR encoded array.

### Signing

A signed message is a wrapper type over the base message.

```go
type SignedMessage struct {
    Message Message
    Signature Signature
}
```

The signature is a serialized signature over the serialized base message. For more details on how the signature itself is done, see the [signatures spec](signatures.md).

### Serialization

Messages and SignedMessages are currently serialized simply by CBOR marshaling them, using lower-camel-cased field names.

## Message Receipt

```go
type MessageReceipt struct {
    ExitCode uint8

    Return []byte
}
```

### Serialization

Message receipts are currently serialized simply by CBOR marshaling them, using lower-camel-cased field names.



## Actor

```go
type Actor struct {
    // Code is a pointer to the code object for this actor
	Code    Cid
    
    // Head is a pointer to the root of this actors state
    Head    Cid
    
    // Nonce is a counter of the number of messages this actor has sent
    Nonce   Integer
    
    // Balance is this actors current balance of filecoin
    Balance AttoFIL
}
```




### Serialization

Actors are currently serialized simply by CBOR marshaling them, using lower-camel-cased field names.

## State Tree

The state trie keeps track of all state in Filecoin. It is a map of addresses to `actors` in the system. It is implemented using a HAMT.

## HAMT

TODO: link to spec for our CHAMP HAMT

# Basic Type Encodings

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

SealProof is an array of bytes.

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