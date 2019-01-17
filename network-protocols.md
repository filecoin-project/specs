# Filecoin Network Protocols

TODO: table of contents

---

All filecoin network protocols are implemented as libp2p protocols. This document will assume that all data is communicated between peers on a libp2p stream.

## CBOR RPC

Filecoin uses many pre-existing protocols from ipfs and libp2p, and also implements several new protocols of its own. For these Filecoin specific protocols, we will try to use a CBOR RPC protocol format. This format is effectively just a leb128 varint length delimeted series of cbor serialized objects. Whenever a filecoin protocol says "send X", it means "cbor serialize X, write its length encoded using leb128, then write the serialized bytes".

# Hello Handshake

The Hello protocol is used when two filecoin nodes initially connect to eachother in order to determine information about the other node. The libp2p protocol ID for this protocol is `/fil/hello/1.0.0`.

Whenever a node gets a new connection, it opens a new stream on that connection and 'says hello'. This is done by crafting a `HelloMessage`, json encoding it (TODO: BAD! switch to CBOR-RPC), writing it over the stream, and finally, closing the stream.
```go
type HelloMessage struct {
    HeaviestTipSet []Cid
    HeaviestTipSetHeight uint64
    GenesisHash Cid
}
```

```go
func SayHello(p PeerID) {
    s := OpenStream(p)
    mes := GetHelloMessage()
    serialized := json.Marshal(mes)
    
    s.Write(serialized)
    s.Close()
}
```

Upon receiving a 'hello' stream from another node, you should read off the json serialized hello message, and then check that the genesis hash matches you genesis hash. If it does not, that node is not part of your network, and should probably be disconnected from. Next, the `HeaviestTipSet`, claimed `HeaviestTipSetHeight`, and peerID of the other node should be passed to the chain sync subsystem.

# Storage Deal

The storage deal protocol is used by any client to store data with a storage miner. The libp2p protocol ID for this protocol is `/fil/storage/mk/1.0.0`.

The protocol starts with storage client (which in this case may be a normal storage client, or a broker). It is assumed that the client has their data already prepared into a `piece` prior to executing this protocol.

First the client sends a signed `StorageDealProposal` to the storage miner:

```go
type StorageDealProposal struct {
    // PieceRef is the hash of the data in native structure. This will be used for
    // certifying the data transfer
    PieceRef Cid
    
    // TranslatedRef is the data hashed in a form that is compatible with the proofs system
    // TODO: this *could* possibly be combined with the PieceRef 
    TranslatedRef Cid
    
    Size NumBytes
    
    TotalPrice TokenAmount
    
    // Duration is how long the file should be stored for
    Duration NumBlocks
    
    // PaymentRef is a reference to the mechanism that the proposer
    // will use to pay the miner. It should be verifiable by the
    // miner using on-chain information.
    Payment PaymentInfo
    
    Signature Signature
}

type PaymentInfo struct {
    // PayChActor is the address of the payment channel actor
    // that will be used to facilitate payments
    PayChActor Address
    
    // Channel is the ID of the specific channel the client will
    // use to pay the miner. It must already have sufficient funds locked up
    Channel ChannelID
    
    // Vouchers is a set of payments from the client to the miner that can be
    // cashed out contingent on the agreed upon data being provably within a
    // live sector in the miners control on-chain
    Vouchers []PaymentVouchers
}
```

### Deal State Values
Legal values for `DealState` are as follows:

```go
const (
	// Unset implies a programmer error. This value should never appear
    // in an actual message
    Unset = 0
    
	// Unknown signifies an unknown negotiation
	Unknown = 1

	// Rejected means the deal was rejected for some reason
	Rejected = 2

	// Accepted means the deal was accepted but hasnt yet started
	Accepted = 3

	// Started means the deal has started and the transfer is in progress
	Started = 4

	// Failed means the deal has failed for some reason
	Failed = 5

	// Complete means the deal is complete, and the sector that the deal is contained
    // in has been sealed and its commitment posted on chain.
	Complete = 6
    
    // Staged is used by the storage deal protocol to indicate the data has been
    // received and staged into a sector, but is not sealed yet
    Staged = 7
)
```



TODO: possibly also include a starting block height here, to indicate when this deal may be started (implying you could select a value in the future). After the first response, both parties will have signed agreeing that the deal started at that point. This could possibly be used to challenge either party in the event of a stall.

The miner then decides whether or not to accept the deal, and sends back a response:

```go
type StorageDealResponse struct {
    State DealState
    
    // Message is an optional message to add context to any given response
    Message string
    
    // ProposalCid is the cid of the StorageDealProposal object this response is for
    ProposalCid Cid
    
    // PieceConfirmation is a collection of information needed to convince the client that
    // the miner has sealed the data into a sector. 
    PieceConfirmation PieceConfirmation
    
    // Signature is a signature from the miner over the response
    Signature Signature
}
```

If `response.State` is `Accepted` then the client should proceed to transfer the data in question to the storage miner. This operation happens out of band from this protocol, and can be a simple bitswap transfer at first. Support for other more 'exotic' 'protocols' such as mailing hard drives is an explicit goal.

Next, when the miner receives all the data and validates it, they set the deals state to `Staged`. When the sector gets sealed, and the commitment is posted on chain, the state gets set to `Complete` and the deals `PieceConfirmation` field should be set to the appropriate values.

Once the deal makes it to this state, the client should be able to query and get the `PieceConfirmation` that they need to complete their proofs of repair for the data.

## Query

Here we describe a basic protocol for querying the current state of a given storage deal. In the future we may want something more complex that is able to multiplex waiting for notifications about a large set of deals over a single stream simultaneously. Upgrading to that from this should be relatively simple, so for now, we do the simple thing.

The libp2p protocol ID for this protocol is `/fil/storage/qry/1.0.0`

At any point, the client in this flow may query the miner for the state of a given proposal. To query, they send a 'StorageDealQuery' that looks like this:

```go
type StorageDealQuery struct {
    // ProposalCid is the cid of the proposal for the deal that we are querying
    // the state of
    ProposalCid *cid.Cid
    
    BaseState DealState
}
```

If `BaseState` is `Unset` or a terminal state (`Complete`, `Rejected`, or `Failed`) then the current state of the deal in question is returned. If the `BaseState` is different than the current state of the deal, the current state of the deal is also returned immediately. In the case that the `BaseState` matches the current state of the deal, then the stream is held open until the state changes, at which point the new state of the deal is returned.


# Retrieve Piece for Free

The Retrieve Piece for Free protocol is used to coordinate the transfer of a piece from miner to client at no cost to the client.

The client initiates the protocol by opening a libp2p stream to the miner using the `/fil/retrieval/free/0.0.0` protocol id. To find and connect to the miner, the [address lookup service](lookup-service.md) should be used. Once connected, the client must send the miner a `RetrievePieceRequest` message using the [CBOR RPC](#CBOR-RPC) protocol format.

The `RetrievePieceRequest` is specified as follows:

```
type RetrievePieceRequest struct {
    // `PieceRef` identifies a piece of user data, typically received from the
    // client while consummating a storage deal.
    PieceRef *cid.Cid
}
```

When the miner receives the request, it responds with a `RetrievePieceResponse` message indicating that it has accepted or rejected the request.

The `RetrievePieceResponse` message is specified as follows:

```
type RetrievePieceResponse struct {
    // `Status` communicates the miner's willingness to send a piece back to a
    // client. The value of the `Status` field must be one of: `Failure` or
    // `Success`.
    Status RetrievePieceStatus

    // If `Status` is `Failure`, `ErrorMessage` should contain a string
    // explaining the cause for the rejection.
    ErrorMessage string
}
```

Legal values for `RetrievePieceStatus` are as follows:

```
const (
	// Unset implies a programmer error. This value should never appear in an
	// actual message.
	Unset = RetrievePieceStatus(iota)

	// Failure indicates that the piece can not be retrieved from the miner.
	Failure

	// Success means that the piece can be retrieved from the miner.
	Success
)
```

If the miner does not accept the request, it sends a `RetrievePieceResponse` with the `Status` field set to `Failure`. The miner should set the `ErrorMessage` field to indicate a reason for the request being rejected.

If the miner accepts the request, it sends a `RetrievePieceResponse` with the `Status` field set to `Success`. The miner then sends the client ordered `RetrievePieceChunk` messages until all of the piece's data has been transferred, at which point the miner closes the stream.

Note: The client must be able to reconstruct a piece by concatenating the `Data`-bytes in the order that they were received.

Note: The miner divides the piece in to chunks containing a maximum of `256 << 8` bytes due to a limitation in our software which caps the size of CBOR-encoded messages at `256 << 10` bytes.

The `RetrievePieceChunk` message is specified as follows:

```
type RetrievePieceChunk struct {
    // The `Data` field contains a chunk of a piece. The length of `Data` must
    // be > 0.
    Data []byte
}
```

TODO: document the query deal interaction


