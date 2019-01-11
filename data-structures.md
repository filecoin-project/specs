# Data Structures

This document serves as an entry point for understanding all of the data structures in filecoin.

TODO: this should also include, or reference, how each data structure is serialized precisely.

## CID

For most objects referenced by Filecoin, a Content Identifier (CID for short) is used. This is effectively a hash value, prefixed with its hash function (multihash) prepended with a few extra labels to inform applications about how to deserialize the given data. To learn more, take a look at the [CID Spec](https://github.com/ipld/cid). 

CIDs are serialized by applying binary multibase encoding, then encoding that as a CBOR byte array with a tag of 42.

## Block

A block represents an individual point in time that the network may achieve consensus on. It contains (via merkle links) the full state of the system, references to the previous state, and some notion of a 'weight' for deciding which block is the 'best'.

```go
// Block is a block in the blockchain.
type Block struct {
	// Miner is the address of the miner actor that mined this block.
	Miner address.Address

	// Ticket is the winning ticket that was submitted with this block.
	Ticket Signature

	// Parents is the set of parents this block was based on. Typically one,
	// but can be several in the case where there were multiple winning ticket-
	// holders for a round.
	Parents SortedCidSet

	// ParentWeightNum is the numerator of the aggregate chain weight of the parent set.
	ParentWeightNum Uint64

	// ParentWeightDenom is the denominator of the aggregate chain weight of the parent set
	ParentWeightDenom Uint64 

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
	Nonce Uint64

	Value AttoFIL

	Method string
	
	Params []byte
}
```

### Parameter Encoding

TODO: discuss how method parameters get encoded

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

## Actor

```go
type Actor struct {
    // Code is a pointer to the code object for this actor
	Code    Cid
    
    // Head is a pointer to the root of this actors state
	Head    Cid
    
    // Nonce is a counter of the number of messages this actor has sent
	Nonce   Uint64
    
    // Balance is this actors current balance of filecoin
	Balance AttoFIL
}
```




### Serialization

Actors are currently serialized simply by CBOR marshaling them, using lower-camel-cased field names.

## State Tree

The state trie keeps track of all state in Filecoin. It is effectively a map of addresses to `actors` in the system. It is implemented using a HAMT.

## HAMT

TODO: link to spec for our CHAMP HAMT