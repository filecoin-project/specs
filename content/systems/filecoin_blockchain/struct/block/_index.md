---
title: Block
weight: 1
dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# Block

The Block is the main unit of the Filecoin blockchain.

The Block structure in the Filecoin blockchain is composed of: i) the Block Header, ii) the list of messages inside the block, and iii) the signed messages. This is represented inside the `FullBlock` abstraction. The messages indicate the required set of changes to apply in order to arrive at a deterministic state of the chain.

The Lotus implementation of the block can be found [here](https://github.com/filecoin-project/lotus/blob/master/chain/types/fullblock.go). It has the following `struct`:

```go
type FullBlock struct {
	Header        *BlockHeader
	BlsMessages   []*Message
	SecpkMessages []*SignedMessage
}
```

> **Note:** A block is functionally the same as a block header in the Filecoin protocol. While a block header contains Merkle links to the full system state, messages, and message receipts, a block can be thought of as the full set of this information (not just the Merkle roots, but rather the full data of the state tree, message tree, receipts tree, etc.). Because a full block is large in size, the Filecoin blockchain consists of block headers rather than full blocks. We often use the terms `block` and `block header` interchangeably.

A `BlockHeader` is a canonical representation of a block. BlockHeaders are propagated between miner nodes. From the blockcheader message, a miner has all the required information to apply the associated `FullBlock`'s state and update the chain. In order to be able to do this, the minimum set of information items that need to be included in the `BlockHeader` are shown below and include among others: the miner's address, the Ticket, the [Proof of SpaceTime](post), the CID of the parents where this block evolved from in the IPLD DAG, as well as the messages' own CID.

The Lotus implementation of the block header can be found [here](https://github.com/filecoin-project/lotus/blob/master/chain/types/blockheader.go). It has the following `struct`s:

```go
type BlockHeader struct {
	Miner address.Address // 0
	Ticket *Ticket // 1
	ElectionProof *ElectionProof // 2
	BeaconEntries []BeaconEntry // 3
	WinPoStProof []proof.PoStProof // 4
	Parents []cid.Cid // 5
	ParentWeight BigInt // 6
	Height abi.ChainEpoch // 7
	ParentStateRoot cid.Cid // 8
	ParentMessageReceipts cid.Cid // 9
	Messages cid.Cid // 10
	BLSAggregate *crypto.Signature // 11
	Timestamp uint64 // 12
	BlockSig *crypto.Signature // 13
	ForkSignaling uint64 // 14
	// ParentBaseFee is the base fee after executing parent tipset
	ParentBaseFee abi.TokenAmount // 15
	// internal
	validated bool // true if the signature has been validated
}

type Ticket struct {
	VRFProof []byte
}

type ElectionProof struct {
	VRFProof []byte
}

type BeaconEntry struct {
	Round uint64
	Data  []byte
}
```

The `BlockHeader` structure has to refer to the TicketWinner of the current round ++ which eensures the correct winner is passed to [ChainSync](chainsync).

```go
func IsTicketWinner(vrfTicket []byte, mypow BigInt, totpow BigInt) bool
```

The `Message` structure has to include the source and destination addresses, a `Nonce` and the `GasPrice`.

The Lotus implementation of the message can be found [here](https://github.com/filecoin-project/lotus/blob/master/chain/types/message.go). It has the following structure:

```go
type Message struct {
	Version uint64

	To   address.Address
	From address.Address

	Nonce uint64

	Value abi.TokenAmount

	GasLimit   int64
	GasFeeCap  abi.TokenAmount
	GasPremium abi.TokenAmount

	Method abi.MethodNum
	Params []byte
}
```

The message is also validated before it is passed to the [chain synchronization logic](chainsync):

```go
func (m *Message) ValidForBlockInclusion(minGas int64) error
```

## Block syntax validation

Syntax validation refers to validation that may be performed on a block and its messages without reference to outside information such as the parent state tree.

An invalid block must not be transmitted or referenced as a parent.

A syntactically valid block header must decode into fields matching the type definition below. 

A syntactically valid header must have:

- between 1 and `5*ec.ExpectedLeaders` `Parents` CIDs if `Epoch` is greater than zero (else empty `Parents`),
- a non-negative `ParentWeight`,
- a `Miner` address which is an ID-address,
- a non-negative `Epoch`,
- a positive `Timestamp`,
- a `Ticket` with non-empty `VRFResult`,
- `ElectionPoStOutput` containing:
   - a `Candidates` array with between 1 and `EC.ExpectedLeaders` values (inclusive),
   - a non-empty `PoStRandomness` field,
   - a non-empty `Proof` field,
- a non-empty `ForkSignal` field.
   
A syntactically valid full block must have:

- all referenced messages syntactically valid,
- all referenced parent receipts syntactically valid,
- the sum of the serialized sizes of the block header and included messages is no greater than `block.BlockMaxSize`,
- the sum of the gas limit of all explicit messages is no greater than `block.BlockGasLimit`.


Note that validation of the block signature requires access to the miner worker address and public key from the parent tipset state, so signature validation forms part of semantic validation. Similarly, message signature validation requires lookup of the public key associated with each message's `From` account actor in the block's parent state.

## Block semantic validation

Semantic validation refers to validation that requires reference to information outside the block header and messages themselves, in particular related to the parent tipset and state on which the block is built.

The semantic validation of a block is carried out by the [Chain Synchronization](chainsync) module.

```go
func (syncer *Syncer) ValidateBlock(ctx context.Context, b *types.FullBlock) error
```

The semantic validation process for Filecoin blocks includes the following:

- Checks for nils on ElectionProof, Ticket, BlockSig and BLSAggregate
- Loads parent Tipset via `TipSetKey`. This is needed in order to link to the history of the block being validated.
- Gets "lookback" tipset and loads the lookback `TipSetState`
- Gets latest beacon entry
- Check that a block isn't seen too far into the future (there appears to be an acceptable window into the future for which a block might be accepted)

In parallel, the [ChainSync logic](chainsync) checks that:
    - the block message is valid
    - it comes from a valid and not slashed miner
    - the state root of the parent tipset matches `BlockHeader.ParentStateRoot`
    - the parent message receipts match `BlockHeader.ParentMessageReceipts`

The ChainSync logic also validates `BlockHeader.ElectionProof`. In order to achieve that it has to:
        - Get block randomness
        - Verify election proof
        - Check that the miner has not been slashed
        - Check that the miner has power
        - Check if block was a winner
 	  	- Check block signature validation
        - Signature validation in lib/sigs/sigs.go:
        ```go
        func CheckBlockSignature(ctx context.Context, blk *types.BlockHeader, worker address.Address) error```
	    - Validate drand beacon
	    - Verify `BlockHeader.WinPoStProof`, which includes:
    	    - Validate block ticket proofs
        	- Verify Winning PoSt proof

If all of the above tests are successful, the block is marked as validated. Ultimately, an invalid block must not be propagated further or validated as a parent node.

