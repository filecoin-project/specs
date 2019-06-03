# The Filecoin Storage Market

### What is the Filecoin Storage Market

The Filecoin `storage market` is the underlying system used to discover, negotiate and form `storage contracts` between clients and storage providers called `storage miners` in a Filecoin network. The `storage market` itself is an `actor` that helps to mediate certain operations in the market, including adding new miners, and punishing faulty ones, it does not directly mediate any actual storage deals. The `storage contracts` between clients and miners specify that a given `piece` will be stored for a given time duration. It is assumed that the `client`, or some delegate of the client, remains online to monitor the `storage miner` and `slash` it in the case that the agreed upon data is removed from the miners proving set before the deal is finished.

The creation of such a storage market is motivated by the need to provide a fast, reliable and inexpensive solution to data generated worldwide. The cost and difficulty involved in starting datacenters around the world make a decentralized solution attractive here, enabling clients and miners to interact directly, forming agreements for storage ad-hoc around the world. Geography is only one such aspect in which a decentralized market can be made competitive. You can read more about the underlying motivations for building a storage market [here](https://www.youtube.com/watch?v=EClPAFPeXIQ).

In the current design of the `storage market`, `storage miners` post `asks` indicating the price they are willing to accepts, and `clients` select (either manually, or via some locally run algorithm) a set of storage miners to store their data with. They then contact the `storage miners` who programmatically either accept or deny their `deal proposals`. In the future, we may allow miners to search for clients and propose deals to them, but for now, for simplicity, we stick with the model described above.

### Visualization of the Filecoin Storage Market

TODO: This is a high level overview of how the storage market interacts with components

## The Market Interface

The market is defined and operated by a built-in network actor. See [storage market actor](actors.md#storage-market-actor).

## The Filecoin Storage Market Operation

The Filecoin storage market operates as follows. Miners providing storage submit ask orders, asking for a certain price for their available storage space, and clients with files to store look through the asks and select a miner they wish to use. Clients negotiate directly with the storage miner that owns that ask, off-chain. Storage is priced in terms of Filecoin per byte per block (note: we may change the units here).

### Market Datastructures

TODO: This storage-market.md doc should try to describe the high level interfaces of the storage market, details about storage and specific method behavior should go in a separate place where we talk about individual system actors.

The storage market contains the following data:

- StorageMiners - The storage market keeps track of the set of the addresses of all storage miners in the storage market. All miners referenced here were created by the storage market via the `CreateStorageMiner` method.
- TotalComittedStorage - This is a tally of all the committed storage in the network. This is both a nice metric to see how much data is being stored by the filecoin network, and a critical piece of information used by mining routine to compute each miners storage ratio.

## Market Flow

This section describes the flow required to store a single piece with a single storage miner. Most use-cases will involve performing this process for multiple pieces, with different miners.

#### Before Deal

1. **Data Preparation:** The client prepares their input data. See [client data](client-data.md) for more details.
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
     - The PieceInclusionProof proves that the piece in the deal is contained in a sector whose commitment is on chain. The 'Merkle Translation' hash from earlier is used here. See [piece inclusion proof for more details](proofs.md#piece-inclusion-proof)
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

The storage market actor maintains a total committed storage amount, aka "total power". The "power table" is a logical construct which is the fraction of this total power that each miner contributes. This is used by the consensus algorithm.

Each individual miner reports its power through its on-chain actor. A miner's power (and hence the total power) is updated only upon PoSt submission (when the set of sectors being proved may change) or when a miner is slashed.


### Power Updates

A miner's power is updated only when they submit a valid PoSt to the chain, or when slashed. See [storage market actor](actors.md#storage-market-actor).

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
