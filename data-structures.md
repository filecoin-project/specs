# Data Structures

This document serves as an entry point for understanding all of the data structures in filecoin. These structures include

- `CID`
- `Address`
- `Block`
- `Message`
- `SignedMessage`
- State Tree

TODO: this should also include, or reference, how each data structure is serialized precisely.

## CID

A Content Identifier or `CID` is a per-object unique identifier used by Filecoin objects to [TODO: Identify how CIDs are used w.r.t. FIL and data object consumption]. Each `CID` is composed of a hash value,  a hash function (multihash), and a few extra labels to that inform object consumers how to deserialize the data associated with a `CID` and coupled FIlecoin object. For more specific information about CIDs or to learn more, take a look at the [CID Spec](https://github.com/ipld/cid). For most objects referenced by Filecoin, a Content Identifier (CID for short) is used. 

In Filecoin, `CID`s are serialized by applying a binary multi-base encoding, then encoding that as a CBOR byte array with a tag of `42`. For more information about CIDs, serialization, and an example in go-filecoin see [TODO: Add link]

## Address

An address is an identifier that refers to an actor in the Filecoin state.

```go
// Hashlength is the number of bytes in the output of Blake2-160
Hashlength = 20

// Length is the total length of the address
Length = 1 + 1 + Hashlength

// An actor's address in the Filecoin network
type address [Length]byte

type FilecoinAddress interface {
    Hash( PubKey string ) [Hashlength]byte
    New( [Hashlength]byte ) address
}
```

In Filecoin, `Hashlength` is defined as the 160-bit (20 Byte) output of Blake2b-160. For more information about addresses see the following examples,

- `Hashlength` -  [go-filecoin Hashlength example](https://github.com/filecoin-project/go-filecoin/blob/master/address/constants.go)
- `address` -  [go-filecoin address example](https://github.com/filecoin-project/go-filecoin/blob/master/address/address.go)

TODO: Update address spec from PR here: https://github.com/filecoin-project/specs/pull/19


## Block

A `Block` represents the canonical representation of the Filecoin network and state at a given  point in time. Moreover, a `Block` is the result of consensus. Each `Block` contains  the full state of the system at epoch `E` (via Merkle links [TODO: add reference to Merkle spec]), reference to the previous state, and a notion of a 'weight' [TODO: add reference to weight] which is used by miners to decide which blocks upon which to mine.

```go
// Block is a block in the blockchain.
type Block struct {
	// Miner is the address of the miner actor that mined this block. 
    // address.Address is a byte array of length = 2+Hashlength. For more 
    // information about Address see [TODO: add link to glossary]
	Miner address.Address

	// Ticket is the winning ticket that was submitted with this block.
    // For more information about Tickets see [TODO: add link to glossary]
	Ticket Signature

	// Parents is the set of parent blocks upon which this block was 
    // is based. The number of parents is expected to be one but can be 
    // greater than one when the result of Expected Consensus (EC) is 
    // also greater than one due to multiple winning ticket-holders for a 
    // round. For more informaiton about EC and tickets see: [TODO: add references]
	Parents []Cid

	// ParentWeightNum is the numerator of the aggregate chain weight
    // of the parent set.
	ParentWeightNum Uint64

	// ParentWeightDenom is the denominator of the aggregate chain 
    // weight of the parent set
	ParentWeightDenom Uint64 

	// Height is the chain height of this block.
	Height Uint64
    
    // StateRoot is a cid pointer to the state tree after application 
    // of the transactions state transitions.
	StateRoot Cid

	// Messages is the set of messages included in this block
	// TODO: should be a merkletree-ish thing
	Messages []SignedMessage

	// MessageReceipts is a set of receipts matching to the sending of
    // the `Messages`.
    // TODO: should be the same type of merkletree-list thing that the messages are
	MessageReceipts []MessageReceipt
}
```

### Serialization

`Block`s are currently serialized via CBOR marshaling using lower-camel-cased field names. Information about CBOR can be found here [TODO: add link]

## Message

```go
type Message struct {
	To   Address
	From Address
	
	// When receiving a message from a user account the nonce in
	// the message must match the expected nonce in the from actor.
	// This prevents replay attacks.
	Nonce Uint64

    // AttoFIL is a signed multi-precision integer quantity of
	// attofilecoin (atto is metric for 10**-18). The zero value for
	// AttoFIL represents the value 0.
	Value AttoFIL

	Method string
	
	Params []byte
}
```

[TODO: add definition of AttoFIL and link to go-filecoin example to Glossary]

### Parameter Encoding

TODO: discuss how method parameters get encoded

### Signed Message

A signed message is a wrapper type over the base message.

```go
type SignedMessage struct {
    Message Message
    Signature Signature
}
```

Where a `signature` is the serialized signature of the serialized representation of `Message`. For more details on how the signature is computed, see the [signatures spec](signatures.md).

### Serialization

`Message`s and `SignedMessage`s are currently serialized via CBOR marshaling using lower-camel-cased field names.

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

`Actor`s are currently serialized simply by CBOR marshaling them, using lower-camel-cased field names.

## State Tree

The state tree keeps track of all state in Filecoin. It is effectively a map of addresses to `actors` in the system. It is implemented using a HAMT.

## HAMT

TODO: link to spec for our CHAMP HAMT