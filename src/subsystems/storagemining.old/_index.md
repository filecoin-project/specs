---
title: Storage Market
entries:
  - protocol
  - actors
  - components
  - libp2p
# suppressMenu: true
---


## The Filecoin Storage Market

The Filecoin `storage market` is the underlying system used to discover, negotiate and form `storage contracts` between clients and storage providers called `storage miners` in a Filecoin network. The `storage market` itself is an `actor` that helps to mediate certain operations in the market, including adding new miners, and punishing faulty ones, it does not directly mediate any actual storage deals. The `storage contracts` between clients and miners specify that a given `piece` will be stored for a given time duration. It is assumed that the `client`, or some delegate of the client, remains online to monitor the `storage miner` and `slash` it in the case that the agreed upon data is removed from the miners proving set before the deal is finished.

The creation of such a storage market is motivated by the need to provide a fast, reliable and inexpensive solution to data generated worldwide. The cost and difficulty involved in starting datacenters around the world make a decentralized solution attractive here, enabling clients and miners to interact directly, forming agreements for storage ad-hoc around the world. Geography is only one such aspect in which a decentralized market can be made competitive. You can read more about the underlying motivations for building a storage market [here](https://www.youtube.com/watch?v=EClPAFPeXIQ).

In the current design of the `storage market`, `storage miners` post `asks` indicating the price they are willing to accepts, and `clients` select (either manually, or via some locally run algorithm) a set of storage miners to store their data with. They then contact the `storage miners` who programmatically either accept or deny their `deal proposals`. In the future, we may allow miners to search for clients and propose deals to them, but for now, for simplicity, we stick with the model described above.

### Visualization of the Filecoin Storage Market

TODO: This is a high level overview of how the storage market interacts with components

## The Market Interface

This interface, written using Go type notation, defines the set of methods that are callable on the storage market actor. The storage market actor is a built-in network actor. For more information about Actors, see [actors.md](actors.md).

```go
type StorageMarket interface {
	// CreateStorageMiner registers a new storage miner with the given public key and a
	// pledge of the given size. The miners collateral is set by the value in the message.
	// The public key must match the private key used to sign off on blocks created
	// by this miner. This key is the 'worker' key for the miner.
	// The libp2p peer ID specified should reference the libp2p identity that the
	// miner is operating. This is the ID that clients will connect to to propose deals
	// TODO: maybe rename to 'RegisterStorageMiner'?
	CreateStorageMiner(pubk PublicKey, pledge BytesAmount, pid libp2p.PeerID) Address

	// SlashConsensusFault is used to slash a misbehaving miner who submitted two different
	// blocks at the same block height. The signatures on each block are validated
	// and the offending miner has their entire collateral slashed, including the
	// invalidation of any any all storage they are providing. The caller is rewarded
	// a small amount to compensate for gas fees (TODO: maybe it should be more?)
	SlashConsensusFault(blk1, blk2 BlockHeader)

	// SlashStorageFault slashes a storage miner for not submitting their PoSTs within
	// the correct [time window](#TODO-link-to-faulty-submission). This may be called by anyone who detects the faulty behavior.
	// The slashed miner then loses all of their staked collateral, and also loses all
	// of their power, and as a result, is no longer a candidate leader for extending the chain.
	SlashStorageFault(miner Address)

	// UpdateStorage is called by a miner to adjust the storage market actors
	// accounting of the total storage in the storage market.
	UpdateStorage(delta BytesAmount)

	// GetTotalStorage returns the total committed storage in the system. This number is
	// also used as the 'total power' in the system for the purposes of the power table
	GetTotalStorage() BytesAmount
}
```


## The Filecoin Storage Market Operation

The Filecoin storage market operates as follows. Miners providing storage submit ask orders, asking for a certain price for their available storage space, and clients with files to store look through the asks and select a miner they wish to use. Clients negotiate directly with the storage miner that owns that ask, off-chain. Storage is priced in terms of Filecoin per byte per block (note: we may change the units here).

### Market Datastructures

The storage market contains the following data:

- StorageMiners - The storage market keeps track of the set of the addresses of all storage miners in the storage market. All miners referenced here were created by the storage market via the `CreateStorageMiner` method.
- TotalComittedStorage - This is a tally of all the committed storage in the network. This is both a nice metric to see how much data is being stored by the filecoin network, and a critical piece of information used by mining routine to compute each miners storage ratio.

## Market Flow

This section describes the flow required to store a single piece with a single storage miner. Most use-cases will involve performing this process for multiple pieces, with different miners.

#### Before Deal

1. **Data Preparation:** The client prepares their input data. See [client data](client-data.md) for more details.
2. **Miner Selection:** The client looks at asks on the network, and then selects a storage miner to store their data with.
   - Note: this is currently a manual process.
3. **Payment Channel Setup:** The client calls [`Payment.Setup`](#payments) with the piece and the funds they are going to pay the miner with. All payments between clients and storage providers use payment channels.


#### Deal

Note: The details of this protocol including formats, datastructures, and algorithms, can be found [here](network-protocols.md#storage-deal).


1. **Storage Deal Staging:** The client now runs the ['make storage deal'](network-protocols.md#storage-deal) protocol, as follows:
  - The client sends a `StorageDealProposal` for the piece in question
    - This contains updates for the payment channel that the client may close at any time, unless the piece gets confirmed (see next section), in which case the miner is able to extend the channel.
  - The miner decides whether or not to accept the deal and sends back a `StorageDealResponse`
    - Note: Different implementations may come up with different ways of making a decision on a given deal.
  - If the miner accepts, the client now sends the data to the miner
  - Once the miner receives the data:
    - They validate that the data matches the storage market hash claimed by the client
    - They stage it into a sector and set the deal state to `Staged`


2. **Storage Deal Start**: Clients makes sure data is in a [sector](definitions.md#sector)

    - **PieceInclusionProof:** Once the miner seals the sector, they update the PieceInclusionProof in the deal state, which the client then gets the next time they query that state.
     - The PieceInclusionProof proves that the piece in the deal is contained in a sector whose commitment is on chain. The `commP` hash from earlier is used here. See [piece inclusion proof for more details](proofs.md#piece-inclusion-proof)
    -  Note: a client that is not interested in staying online to wait for PieceInclusionProof can leave immediately, however, they run the risk that their files don't actually get stored (but if their data is not stored, the miner will not be able to claim payment for it).
     - Note: In order to provide the piece inclusion proof, the miner needs to fill the sector. This may take some time. So there is a wait between the time the data is transferred to the miner, and when the piece inclusion proof becomes available.
   - **Mining**: Miner posts `seal commitment` and associated proof on chain by calling `CommitSector` and starts running `proofs of spacetime`. See [storage mining cycle](mining.md#storage-mining-cycle) for more details.

3. **Storage Deal Abort:** If the miner doesn't provide the PieceInclusionProof, the client can invalidate the payment channel.
   - This is done by invoking the 'close' method on the channel on-chain. This process starts a timer that, on finishing, will release the funds to the client. 
   - If a client attempts to abort a deal that they have actually made with a miner, the miner can submit a payment channel update to force the channel to stay open for the length of the agreement.

4. **Storage Deal Complete:** The client periodically queries the miner for the deals status until the deal is 'complete', at which point the client knows that the data is properly replicated.
   - The client should store the returned 'PieceInclusionProof' for later validation.

TODO: 'complete' isnt really the right word here, as it implies that the deal is over.


5. **Income Withdrawal**: When the miner wishes to withdraw funds, they call [`Payment.RedeemVoucher`](#payments).


## The Power Table

The `power table` is exported by the storage market for use by consensus. There isn't actually a concrete object that is the power table (though the concept is conceptually helpful), instead, the [storage market actor](actors.md#storage-market-actor) exports the `GetTotalStorage` and `PowerLookup`  methods which can be used to lookup total network power and a miner's power, respectively. 
Each individual miner reports its power through their actor.

To check the power of a given miner, use the following:

```go
func GetMinersPowerAt(ts TipSet, m Address) Integer {
  curState := GetStateTree(ts)
  miner := curState.GetMiner(m)
  if miner.IsSlashed() || miner.IsLate() {
    return 0
  }
  
  # lookback to the last valid PoSt put up by the miner
  lookbackTipset := WalkBack(ts, miner.provingPeriodEnd - provingPeriodDuration(miner.SectorSize))
  lbState := GetStateTree(lookbackTipset)
  
  sm := lbState.GetStorageMarket()
  
  return sm.PowerLookup(m)
}
```

### Power Updates

Whenever a new [PoSt](proofs.md) or [Fault](faults.md) makes it on chain, the storage market updates the underlying power values appropriately.

Specifically, a miner's power is initialized/maintained when they [submit a valid PoSt](actors.md#submitPoSt) to the chain, and decreases if they are slashed (for a [storage fault](actors.md#slashStorageFault) or a [consensus fault](actors.md#slashConsensusFault)).

Power is deducted when miners remove sectors by reporting the sector 'missing' or 'done' in a PoSt.


## Payments

The storage market expects a payments system to allow clients to pay miners for storage. Any payments system that has the following capabilities may be used:

```go
type Payments interface {
	// Setup sets up a payment from the caller to the target address. The payment
	// MUST be contingent on the miner being able to prove that they have the data
	// referenced by 'piece'. The total amount of Filecoin that may be transfered by
	// this payment is specified by 'value'
	Setup(target Address, piece Cid, value TokenAmount) ID

	// MakeVouchers creates a set of vouchers redeemable by the target of the
	// previously created payment. It creates 'count' vouchers, each of which is
	// redeemable only after an certain block height, evenly spaced out between
	// start and end. Each voucher should be redeemable for proportionally more
	// Filecoin, up to the total amount specified during the payment setup.
	MakeVouchers(id ID, start, end BlockHeight, count int) []Voucher

	// Redeem voucher is called by the target of a given payment to claim the
	// funds represented by it. The voucher can only be redeemed after the block
	// height that is attributed to the voucher, and also only if the proof given
	// proves that the target is correctly storing the piece referenced in the
	// payment setup.
	RedeemVoucher(v Voucher, proof Proof)
}
```

For details on the implementation of the payments system, see [the payments doc](payments.md).



## Future Protocol Improvements

- Slashable Commitments
  - When miners initially receive the data for a deal with a client, that signed response statement can be used to slash the miner in the event that they never include that data in a sector.

# Open questions

- Storage time should likely be designated in terms of proving period. Where a proving period is the number of blocks in which every miner must submit a proof for their sectors. Not doing this makes accounting hard: "when exactly did this sector fail?"






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







## Piece Commitment

### commP

A piece commitment (`commP`) is the root hash of a piece that a client wants to store in Filecoin. It is generated using `RepHash` (as described in [Proof-of-Replication](zigzag-porep.md)) on some raw data which has been zero-padded to a multiple of 127 bytes, then preprocessed yielding `Fr32 padded` data which is a multiple of 128 bytes. 









## Transfer formats

The transfer format is the format to transfer a file over the network. This format SHALL be used for the initial transfer (from clients to storage miners) and for later retrievals (from storage miners to the clients).

The default transfer format is `unixfsv1`. Cliens MAY agree to use other formats of their preference.

### `unixfsv1`

The default transfer format is Unixfsv1 with the following parameters:

- Chunking: Fixed, 1MB
- Leaf Format: Raw
- Max Branch Width: 1024

For details on how UnixfsV1 works, see its spec [here](https://github.com/ipfs/specs/tree/master/unixfs).




## Storage Formats

The Storage Format MUST be use for generating Filecoin proofs and hashing sectors data. 

The current required storage format is `paddedfr32v1`.

### `paddedfr32v1`

A correctly formatted `paddedfr32v1` data must have:

- **Fr32 Padding**: Every 32 bytes blocks MUST contain two zeroes in the most significant bits (every 254 bits must be followed by 2 zeroes if interpreted as little-endian number). That is, for each block, `0x11000000 & block[31] == 0`.
- **Piece Padding**: In order to generate minimal `PieceInclusionProofs`, blocks of 32 zero bytes MUST be added so that the total number of blocks (including *piece padding*) is a power of two. **Piece Padding** can be omitted if the prover wishes to generate unaligned proofs. [NOTE: not yet fully specified.]

**Why do we need a special Storage Encoding Format?** In the Filecoin proofs we do operations in an arithmetic field of size `p`, where `p` is a prime of size `2^255`, hence the size of the data blocks must be smaller than `p`. We cautiously decide to have data blocks of size 254 to avoid possible overflows (data blocks numerical representation is bigger than `p`). 



## Miners Claiming Earnings

Storage Miners claim their Storage Market earnings via payment channels.

The client proposes the cadence of the earnings for a deal by creating `SignedVoucher`-s. Each vouchers specify how often Storage Miners can claim earnings and how much each earning should be, more precisely, each voucher has some tokens assigned and can be redeemed only at a particular block height. The vouchers are part of the `PaymentInfo` included in the `StorageDealProposal`. When receiving a proposal, a Storage Miner can review and accept these terms by completing the deal protocol.

After the block defined in each `SignedVoucher` is passed, the Storage Miner could claim the earning by updating the payment channel calling `UpdateChannelState` on the `PaymentChannel` actor for a particular `SignedVoucher`. This call passes if the Storage Miner is still storing the piece in sector and if the Storage Miner is not late in their PoSt submission and if the time specified in the `SignedVoucher` has passed.



## Storage Miner Payments

TODO: these bits were pulled out of a different doc, and describe strategies by which client payments to a miner might happen. We need to organize 'clients paying miners' better, unclear if it should be the same doc that talks about payment channel constructions.

1. **Updates Contingent on Inclusion Proof**
   - In this case, the miner must provide an inclusion proof that shows the client data is contained in one of the miners sectors on chain, and submit that along with the payment channel update.
   - This can be pretty expensive for smaller files, and ideally, we make it to one of the latter two options
   - This option does however allow clients to upload their files and leave.
2. **Update Contingent on CommD Existence**
   - For this, the client needs to wait around until the miner finishes packing a sector, and computing its commD. The client then signs a set of payment channel updates that are contingent on the given commD existing on chain.
   - This route makes it difficult for miners to re-seal smaller files (really, small files just suck)
3. **Reconciled Payment**
   - In either of the above cases, the miner may go back to the client and say "Look, these payment channel updates you gave me are able to be cashed in right now, could you take them all and give me back a single update for a slightly smaller amount?".
   - The slightly smaller amount could be the difference in transaction fees, meaning the client saves money, and the miner gets the same amount.
