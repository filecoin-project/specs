# `Retrieval Market` V0

### What is the `Retrieval Market`

### What the `Retrieval Market` affects

### Dependencies

## Components

Version 0 of the `retrieval market` protocol is what we (tentatively) will launch the filecoin network with. It is version zero because it will only be good enough to fit the bill as a way to pay another node for a file.

The main components are as follows:

- A payment channel actor (See [payment channels](payment-channels.md) for details)
- 'retrieval-v0' `libp2p` services
- A chain-based content routing interface
- A set of commands to interact with the above

## Retrieval V0 `libp2p` Services

The v0 `retrieval market` will initially be implemented as two `libp2p` services. It will be request response based, where the client who is requesting a file sends a `retrieval deal proposal` to the miner. The miner chooses whether or not to accept it, sends their response which (if they accept the proposal) includes a `signed retrieval deal`, followed by the actual requested content, streamed as a series of bitswap block messages, using a pre-order traversal of the dag. Each block should use the [bitswap block message format](https://github.com/ipfs/go-ipfs/blob/master/exchange/bitswap/message/message.go#L216). This way, the client should be able to verify the data incrementally as it receives it. Once the client has received all the data, it should then send a payment channel SpendVoucher of the proposed amount to the miner. This protocol may be easily extended to include payments from the client to the miner every N blocks, but for now we omit that feature.

```sh
type RetDealProposal struct {
	## Reference to the data being retrieved.
	ref Link

	## The total amount that the client is willing to pay for the retrieval of the data.
	price TokenAmount

	## The info from the client to the retrieval miner for the data
	payment PaymentInfo
}

type RetDealResponse union {
    | AcceptedResponse 0
    | RejectedResponse 1
    | ErrorResponse 2
} representation keyed

type AcceptedResponse struct {}
type RejectedResponse struct {
    message optional String
}

type ErrorResponse RejectedResponse

type Block struct {
	## Cid prefix parameters for this block. It describes how to
	## hash the block to verify it matches the expected value.
	refix CidPrefix
	data Bytes
}

## Represents all the metadata of a Cid.
## It does not contains  any actual content information.
type CidPrefix struct {
	version  UInt
	codec    UInt
	mhType   UInt
	mhLength UInt
}
```

`Retrieval miners` should also support a query service that allows clients to request pricing information from a miner.

The query should include the CID of the piece that the client is interested in retrieving. The response contains whether or not the miner will serve that data, the price they will accept for it.

```sh
type RetQuery struct {
    ## TODO: what exactly does this link to?
	piece Link
}

type RetQueryResponse union {
    | AvailableResponse
    | UnavailableResponse
} representation keyed

type AvailableResponse struct {
	minPrice TokenAmount
}

type UnavailableResponse struct {}
```

## Chain Based Content Routing

For the version 0 protocol. We should implement a small helper service that looks up which miners have a given piece based on deals made in the blockchain. The service should first look the content up in the blockchain (or in some client index) to find the chain address of the miner, then use the lookup service to map that to a `libp2p` `peerID` and `multiaddr`.

The interface should match the exist libp2p content routing interface:

```go
type ChainContentRouting interface {
	FindProvidersAsync(ref Cid, count int) <-chan pstore.PeerInfo
}
```

## Retrieval Market Commands

We will need to add a few commands to allow the user to interact with the `retrieval market`, and for developers to be able to script higher level applications on top of it.

The command names here are not final, and are definitely subject to change later on once we are able to sit and think through proper UX.

```text
USAGE
  filecoin retr get <piece-cid> - Retrieve a piece from a miner.

SYNOPSIS
  filecoin retr get [--price=<amt>] [--miner=<peerID>] [--] <piece-cid>

ARGUMENTS

  <piece-cid> - Content ID of piece to retrieve.

OPTIONS

  --price                string - Amount of filecoin to offer for this data.
  --miner                string - Optional Peer ID of miner to connect to. (If unspecified, the chain routing service will be used)
```

```text
USAGE
  filecoin retr lookup <piece-cid> - Print a list of miners who have the given piece.

SYNOPSIS
  filecoin retr lookup [--sort=<sorttype>] [--] <piece-cid>

ARGUMENTS

  <piece-cid>... - Content ID of piece to find.

OPTIONS

  --sort                string - Output sorting scheme.
```

```text
USAGE
  filecoin retr query <minerID> [<piece-cid>] - Query the given retrieval miner.

SYNOPSIS
  filecoin retr query [--] <miner-id> [<piece-cid>]

ARGUMENTS

  <miner-id>  - ID of miner to query.
  [<piece-cid>] - Optional cid of piece to query for.
```
