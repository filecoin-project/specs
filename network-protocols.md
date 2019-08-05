# Filecoin Specific Network Protocols

All filecoin network protocols are implemented as libp2p protocols. This document will assume that all data is communicated between peers on a libp2p stream as outlined in [networking](networking.md)

## Ipld dag-cbor RPC

Filecoin uses many pre-existing protocols from ipfs and libp2p, and also implements several new protocols of its own. For these Filecoin specific protocols, we will use the Ipld dag-cbor RPC protocol format, defined below.

This format consists of series of IPLD objects, serialized using `dag-cbor`. Whenever a filecoin protocol says "send X", it means "IPLD serialize the object X using `dag-cbor`, then write the serialized bytes".

{{% notice Note %}}
**Note:** Implementations should limit the maximum number of bytes when reading IPLD objects from the wire. We suggest `1MB` as a sane limit.
{{% /notice %}}

**Links**

- Ipld: https://github.com/ipld/specs/
- dag-cbor: https://github.com/ipld/specs/blob/master/block-layer/codecs/DAG-CBOR.md


## Hello Handshake

- **Name**: Hello
- **Protocol ID**: `/fil/hello/1.0.0`

> The Hello protocol is used when two filecoin nodes initially connect to each other in order to determine information about the other node.

Whenever a node gets a new connection, it opens a new stream on that connection and "says hello". This is done by crafting a `HelloMessage`, sending it to the other peer using CBOR RPC and finally, closing the stream.

```sh
type HelloMessage struct {
	heaviestTipSet [&Block]
	heaviestTipSetWeight UInt
	genesisHash &Block
}
```
Upon receiving a "hello" stream from another node, you should read off the CBOR RPC message, and then check that the genesis hash matches you genesis hash. If it does not, that node is not part of your network, and should probably be disconnected from. Next, the `HeaviestTipSet`, claimed `HeaviestTipSetWeight`, and peerID of the other node should be passed to the chain sync subsystem.

## Storage Deal Make

- **Name**: Storage Deal Make
- **Protocol ID**: `/fil/storage/mk/1.0.0`

> The storage deal protocol is used by any client to store data with a storage miner.

The protocol starts with storage client (which in this case may be a normal storage client, or a broker). It is assumed that the client has their data already prepared into a `piece` prior to executing this protocol. For more details on initial data processing, see [client data](client-data.md).

First the client sends a `SignedStorageDealProposal` to the storage miner:

```sh
type Commitment Bytes
```

```sh
type SerializationMode enum {
     | "UnixFs"
    ## no transformations applied
    | "Raw"
    ## Serialized as IPLD, encoding is specified in the CID stored in `pieceRef`
    | "IPLD"
}
```

```sh
type StorageDealProposal struct {
	## PieceRef is the hash of the data in native structure. This will be used for
	## certifying the data transfer.
    ## Reference for transit.
	pieceRef String

	## Specifies how the graph referenced by 'PieceRef' gets transformed
	## into the data that will be packed into a sector.
	serializationMode SerializationMode

	## The data hashed in a form that is compatible with the proofs system.
    ## Reference for actual storage in a sector.
	commP Commitment

	size BytesAmount

	totalPrice TokenAmount

	## Duration is how long the file should be stored for
	Duration NumBlocks

	## A reference to the mechanism that the proposer will use to pay the miner. It should be
    ## verifiable by the miner using on-chain information.
	payment PaymentInfo

	## MinerAddress is the address of the storage miner in the deal proposal
	minerAddress Address

	clientAddress Address
}
```

```sh
type SignedStorageDealProposal struct {
	proposal StorageDealProposal

	## Signature over the the encoded StorageDealProposal signed by the client.
	signature Signature
}
```

```sh
type PaymentInfo struct {
	## The address of the payment channel actor that will be used to facilitate payments.
	payChActor Address

	## Reference to the message used to create the payment channel. This allows the miner to wait until the
	## channel is accepted on chain. (optional)
	channelMessage &Message

  ## Set of payments from the client to the miner that can be cashed out contingent on the agreed
  ## upon data being provably within a live sector in the miners control on-chain.
	vouchers [SignedVoucher]
}
```

```sh
type DealState enum {
    ## Signifies an unknown negotiation.
    | Unknown 0
    ## The deal was rejected for some reason.
    | Rejected 1
    ## The deal was accepted but hasn't yet started.
    | Accepted 2
    ## The deal has started and the transfer is in progress.
    | Started 3
    ## The deal has failed for some reason.
    | Failed 4
    ## The data has been received and staged into a sector, but is not sealed yet.
    | Staged 5
    ## The data is being sealed and a `PieceInclusionProof` is available.
    | Sealing 6
    ## Deal is complete, and the sector that the deal is contained in has been sealed and its
    ## commitment posted on chain.
    | Complete 7
}
```

{{% notice todo %}}
**TODO**: possibly also include a starting block height here, to indicate when this deal may be started (implying you could select a value in the future). After the first response, both parties will have signed agreeing that the deal started at that point. This could possibly be used to challenge either party in the event of a stall. This starting block height also gives the miner time to seal and post the commitment on chain. Otherwise a weird condition exists where a client could immediately slash a miner for not having their data stored.
{{% /notice %}}


```sh
type StorageDealResponse union {
    | UnknownParams
    | RejectedParams
    | AcceptedParams
    | StartedParams
    | FailedParams
    | StagedParams
    | SealingParams
    | CompleteParams
} representation keyed

type UnknownParams struct {
	## Message is an optional message to add context to any given response
	message optional String
}

type RejectedParams struct {
    message optional String

	## A reference to the proposal this is the response to.
    proposal &SignedStorageDealProposal
}

type AcceptedParams RejectedParams
type FailedParams RejectedParams
type StagedParams RejectedParams

type SealingParams struct {
	## The proof needed to convince the client that the miner has sealed the data into a sector.
	## Note: the miner doesnt necessarily have to have committed the sector at this point
	## they just need to have staged it into a sector, and be committed to putting it at
	## that place in the sector.
	pieceInclusionProof PieceInclusionProof
}

type CompleteParams struct {
	## A reference to the message that was sent to submit the sector containing this data to the chain.
	sectorCommitMessage &Message
}
```

```sh
type SignedStorageDealResponse struct {
	response StorageDealResponse

	## Signature is a signature from the miner over the cbor encoded response
	signature Signature
}
```

### Process

1. [Client] send `SignedStorageDealProposal`.
2. [Miner]  send `SignedStorageDealResponse`, either accepting or rejecting the deal.
3. [Client] If `response.state` is `Accepted` then transfer the data in question.
4. [Miner] Once the miner receives all the data they validate it. On success they set the `DealState` to `Staged` (internally).
5. [Miner] When the sector gets sealed, the state gets set to `Sealing`.
6. [Miner] When the commitment is posted on chain, the state gets set to `Complete`.
6. [Client] Once the deal makes it to the `Sealing` state, they are able to query and get the `PieceInclusionProof` that they need to verify that the miner is indeed storing their data.

At any point in time the client can query (using the query protocol) the miner to get the current state of the deal.

{{% notice Note %}}
**Note:** The data transfer operation happens out of band from this protocol, and can be a simple bitswap transfer at first. Support for other more 'exotic' 'protocols' such as mailing hard drives is an explicit goal.
{{% /notice %}}


## Storage Deal Query

- **Name**: Storage Deal Query
- **Protocol ID**: `/fil/storage/qry/1.0.0`

This is the basic protocol for querying the current state of a given storage deal.
At any point, the client in this flow may query the miner for the state of a given proposal. To query, they send a `StorageDealQuery` that looks like this:

```sh
type StorageDealQuery struct {
	## ProposalCid is the cid of the proposal for the deal that we are querying
	## the state of
	proposal &SignedStorageDealProposal

	baseState DealState
}
```

If `baseState` is `Unset` or a terminal state (`Complete`, `Rejected`, or `Failed`) then the current state of the deal in question is returned. If the `baseState` is different than the current state of the deal, the current state of the deal is also returned immediately. In the case that the `baseState` matches the current state of the deal, then the stream is held open until the state changes, at which point the new state of the deal is returned.

{{% notice Note %}}
**Note:** In the future we may want something more complex that is able to multiplex waiting for notifications about a large set of deals over a single stream simultaneously. Upgrading to that from this should be relatively simple, so for now, we do the simple thing.
{{% /notice %}}

## Retrieve Piece for Free

- **Name**: Retrieve Piece for Free
- **Protocol ID**: `/fil/retrieval/free/0.0.0`

> The Retrieve Piece for Free protocol is used to coordinate the transfer of a piece from miner to client at no cost to the client.

The client initiates the protocol by opening a libp2p stream to the miner. To find and connect to the miner, the [address lookup service](lookup-service.md) should be used. Once connected, the client must send the miner a `RetrievePieceRequest` message.

When the miner receives the request, it responds with a `RetrievePieceResponse` message indicating that it has accepted or rejected the request.

If the miner does not accept the request, it sends a `RetrievePieceResponseFailure`, with the `message` field set to indicate a reason for the request being rejected.

If the miner accepts the request, it sends a `RetrievePieceResponseSuccess`. Then the miner sends the client ordered `RetrievePieceChunk` messages until all of the piece's data has been transferred, at which point the miner closes the stream.

Note: The client must be able to reconstruct a piece by concatenating the `Data`-bytes in the order that they were received.

Note: The miner divides the piece in to chunks containing a maximum of `256 << 8` bytes due to a limitation in our software which caps the size of CBOR-encoded messages at `256 << 10` bytes.

```sh
type RetrievePieceRequest struct {
	## Identifier of user data, typically received from the client while consummating a storage deal.
	pieceRef String
}
```

```sh
type RetrievePieceResponse union {
    ## Success means that the piece can be retrieved from the miner.
    | RetrievePieceResponseSuccess 0
	## Failure indicates that the piece can not be retrieved from the miner.
    | RetrievePieceResponseFailure 1
} representation keyed

type RetrievePieceResponseSuccess struct {}

type RetrievePieceResponseFailure struct {
	## A string explaining the cause for the rejection.
	message string
}
```

```sh
##  A chunk of a piece. The length must be > 0.
type RetrievePieceChunk Bytes
```

## BlockSync

- **Name**: Block Sync
- **Protocol ID**: `/fil/sync/blk/0.0.1`

The blocksync protocol is a small protocol that allows Filecoin nodes to request ranges of blocks from each other. It is a simple request/response protocol.

The request requests a chain of a given length by the hash of its highest block. The `Options` allow the requester to specify whether or not blocks and messages to be included.

The response contains the requested chain in reverse iteration order. Each item in the `Chain` array contains the blocks for that tipset if the `Blocks` option bit in the request was set, and if the `Messages` bit was set, the messages across all blocks in that tipset. The `MsgIncludes` array contains one array of integers for each block in the `Blocks` array. Each of the arrays in `MsgIncludes` contains a list of indexes of messages from the `Messages` array that are in each `Block` in the blocks array.

```sh
type BlockSyncRequest struct {
    ## The TipSet being synced from
	start [&Block]
    ## How many tipsets to sync
	requestLength UInt
    ## Query options
    options Options
}
```

```sh
type Options enum {
    # Include only blocks
    | Blocks 0
    # Include only messages
    | Messages 1
    # Include messages and blocks
    | BlocksAndMessages 2
}

type BlockSyncResponse struct {
	chain [TipSetBundle]
	status Status
}

type TipSetBundle struct {
  blocks [Blocks]
  secpMsgs [SignedMessage]
  secpMsgIncludes [[UInt]]

  blsMsgs [Message]
  blsMsgIncludes [[Uint]]
}

type Status enum {
    ## All is well.
    | Success 0
    ## Sent back fewer blocks than requested.
    | PartialResponse 101
    ## Request.Start not found.
    | BlockNotFound 201
    ## Requester is making too many requests.
    | GoAway 202
    ## Internal error occured.
    | InternalError 203
}
```

### Example

The TipSetBundle

```
Blocks: [b0, b1]
secpMsgs: [mA, mB, mC, mD]
secpMsgIncludes: [[0, 1, 3], [1, 2, 0]]
```

corresponds to:

```
Block 'b0': [mA, mB, mD]
Block 'b1': [mB, mC, mA]
```
