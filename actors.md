# Filecoin State Machine Actors

Any implementations of the Filecoin actors must be exactly byte for byte compatible with the go-filecoin actor implementations. The pseudocode below tries to capture the important logic, but capturing all the detail would require embedding exactly the code from go-filecoin, so for now, its simply informative pseudocode. The algorithms below are correct, and all implementations much match it (including go-filecoin), but details omitted from here should be looked for in the go-filecoin code.

This spec describes a set of actors that operate within the [Filecoin State Machine](state-machine.md). All types are defined in [the basic type encoding spec](data-structures.md#basic-type-encodings).

- [Init Actor](#init-actor)
- [Storage Market Actor](#storage-market-actor)
- [Storage Miner Actor](#storage-miner-actor)
- [Payment Channel Broker Actor](#payment-channel-broker-actor)

## Built In Actors

Some state machine actors are 'system' actors that get instantiated in the genesis block, and have their IDs allocated at that point.

| ID   | Actor              | Name                    |
| ---- | ------------------ | ----------------------- |
| 0    | InitActor          | Network Init            |
| 1    | AccountActor       | Network Treasury        |
| 2    | StorageMarketActor | Filecoin Storage Market |
|  99 | AccountActor | Burnt Funds |

## Init Actor

The init actor is responsible for creating new actors on the filecoin network. This is a built-in actor and cannot be replicated. In the future, this actor will be responsible for loading new code into the system (for user programmable actors). ID allocation for user instantiated actors starts at 100. This means that `NextID` will initially be set to 100.

```go
type InitActor struct {
	// Mapping from Address to ID, for lookups.
	AddressMap map[Address]BigInt

	NextID BigInt
}
```

### Code Cid
`<codec:raw><mhType:identity><"init">`

| Index     | Method Name       |
| -------- | ---------- |
| 1  | `Exec`     |
| 2 | `GetIdForAddress` |

### `Exec(code Cid, params []Param) Address`

>  This method is the core of the `Init Actor`. It handles instantiating new actors and assigning them their IDs.

#### Parameters

| Name     | Type       | Description                                                  |
| -------- | ---------- | ------------------------------------------------------------ |
| `code`   | `Cid`      | A pointer to the location at which the code of the actor to create is stored. |
| `params` | `[] Param` | The parameters passed to the constructor of the actor.       |



`Param` is the type representing any valid arugment that can be passed to a function.

TODO: Find a better place for this definition.


```go
func Exec(code Cid, params []byte) Address {
	// Get the actor ID for this actor.
	actorID = self.NextID
	self.NextID++

	// Make sure that only the actors defined in the spec can be launched.
	if !IsBuiltinActor(code) {
		Fatal("cannot launch actor instance that is not a builtin actor")
	}

	// Ensure that singeltons can be only launched once.
	// TODO: do we want to enforce this? If so how should actors be marked as such?
	if IsSingletonActor(code) {
		Fatal("cannot launch another actor of this type")
	}

	// This generates a unique address for this actor that is stable across message
	// reordering
	addr := VM.ComputeActorAddress()

	// Set up the actor itself
	actor := Actor{
		Code:    code,
		Balance: msg.Value,
	}

	// The call to the actors constructor will set up the initial state
	// from the given parameters
	actor.Constructor(params)

	VM.GlobalState.Set(actorID, actor)

	// Store the mapping of address to actor ID.
	self.AddressMap[addr] = actorID

	return addr
}

func IsSingletonActor(code Cid) bool {
	return code == StorageMarketActor || code == InitActor
}

// TODO: find a better home for this logic
func (VM VM) ComputeActorAddress(creator Address, nonce Integer) Address {
	return NewActorAddress(bytes.Concat(creator.Bytes(), nonce.BigEndianBytes()))
}
```

### `GetIdForAddress(addr Address) BigInt`

> This method allows for fetching the corresponding ID of a given Address

#### Parameters

| Name   | Type      | Description           |
| ------ | --------- | --------------------- |
| `addr` | `Address` | The address to lookup |



```go
func GetIdForAddress(addr Address) BigInt {
	id := self.AddressMap[addr]
	if id == nil {
		Fault("unknown address")
	}
	return id
}
```

## Account Actor

The Account actor is the actor used for normal keypair backed accounts on the filecoin network.

```go
type AccountActor struct {
	// Address contains the public key based address that this account was created with. If unset, this account may not send funds by normal means.
	Address Address
}
```

### Code Cid
`<codec:raw><mhType:identity><"account">`

| Index     | Method Name       |
| -------- | ---------- |
| 1  | `GetAddress` |

## Storage Market Actor

The storage market actor is the central point for the Filecoin storage market. It is responsible for registering new miners to the system, and maintaining the power table. The Filecoin storage market is a singleton that lives at a specific well-known address.

```go
type StorageMarketActor struct {
	Miners AddressSet

	// TODO: Determine correct unit of measure. Could be denominated in the
	// smallest sector size supported by the network.
	TotalStorage BytesAmount
}
```
### Code Cid
`<codec:raw><mhType:identity><"smarket">`

| Index     | Method Name       |
| -------- | ---------- |
| 1   | `CreateStorageMiner`     |
| 2 | `SlashConsensusFault` |
| 3 | `UpdateStorage` |
| 4 | `GetTotalStorage` |


#### CreateStorageMiner

Parameters:

- worker Address
- pledge BytesAmount
- sectorSize BytesAmount
- pid PeerID

Return: Address

```go
func CreateStorageMiner(worker Address, pledge, sectorSize BytesAmount, pid PeerID) Address {
	if !SupportedSectorSize(sectorSize) {
		Fatal("Unsupported sector size")
	}

	if pledge < MinimumPledge(sectorSize) {
		Fatal("Pledge too low")
	}

	if msg.Value < MinimumCollateral(pledge) {
		Fatal("not enough funds to cover required collateral")
	}

	newminer := InitActor.Exec(MinerActorCodeCid, EncodeParams(pubkey, pledge, sectorSize, pid))

	self.Miners.Add(newminer)

	return newminer
}
```

### SlashConsensusFault

Parameters:

- block1 BlockHeader
- block2 BlockHeader

Return: None

```go
func shouldSlash(block1, block2 BlockHeader) bool {
	// First slashing condition, blocks have the same ticket round
	if sameTicketRound(block1, block2) {
		return true
	}

	// Second slashing condition, miner ignored own block when mining
	// Case A: block2 could have been in block1's parent set but is not
	block1ParentTipSet := parentOf(block1)
	if !block1Parent.contains(block2) &&
		block1ParentTipSet.Height == block2.Height &&
		block1ParentTipSet.ParentCids == block2.ParentCids {
		return true
	}

	// Case B: block1 could have been in block2's parent set but is not
	block2ParentTipSet := parentOf(block2)
	if !block2Parent.contains(block1) &&
		block2ParentTipSet.Height == block1.Height &&
		block2ParentTipSet.ParentCids == block1.ParentCids {
		return true
	}

	return false
}

func SlashConsensusFault(block1, block2 BlockHeader) {
	if !ValidateSignature(block1.Signature) || !ValidSignature(block2.Signature) {
		Fatal("invalid blocks")
	}

	if AuthorOf(block1) != AuthorOf(block2) {
		Fatal("blocks must be from the same miner")
	}

	// see the "Consensus Faults" section of the faults spec (faults.md)
	// for details on these slashing conditions.
	if !shouldSlash(block1, block2) {
		Fatal("blocks do not prove a slashable offense")
	}

	miner := AuthorOf(block1)

	// TODO: Some of the slashed collateral should be paid to the slasher

	// Burn all of the miners collateral
	miner.BurnCollateral()

	// Remove the miner from the list of network miners
	self.Miners.Remove(miner)
	self.UpdateStorage(-1 * miner.Power)

	// Now delete the miner (maybe this is a bit harsh, but i'm okay with it for now)
	miner.SelfDestruct()
}
```

### UpdateStorage

UpdateStorage is used to update the global power table.

Parameters:

- delta BytesAmount

Return: None

```go
func UpdateStorage(delta BytesAmount) {
	if !self.Miners.Has(msg.From) {
		Fatal("update storage must only be called by a miner actor")
	}

	self.TotalStorage += delta
}
```

### GetTotalStorage

Parameters: None

Return: BytesAmount

```go
func GetTotalStorage() BytesAmount {
	return self.TotalStorage
}
```

## Storage Miner Actor

```go
type StorageMiner struct {
	// Owner is the address of the account that owns this miner. Income and returned
	// collateral are paid to this address. This address is also allowed to change the
	// worker address for the miner.
	Owner Address

	// Worker is the address of the worker account for this miner.
	// This will be the key that is used to sign blocks created by this miner, and
	// sign messages sent on behalf of this miner to commit sectors, submit PoSts, and
	// other day to day miner activities.
	Worker Address

	// PeerID is the libp2p peer identity that should be used to connect
	// to this miner
	PeerID peer.ID

	// PledgeBytes is the amount of space being offered by this miner to the network
	PledgeBytes BytesAmount

	// SectorSize is the amount of space in each sector committed to the network
	// by this miner.
	SectorSize BytesAmount

	// Collateral is locked up filecoin the miner has available to commit to storage.
	// When miners commit new sectors, tokens are moved from here to 'ActiveCollateral'
	// The sum of collateral here and in activecollateral should equal the required amount
	// for the size of the miners pledge.
	Collateral TokenAmount

	// ActiveCollateral is the amount of collateral currently committed to live storage
	ActiveCollateral TokenAmount

	// DePledgedCollateral is collateral that is waiting to be withdrawn
	DePledgedCollateral TokenAmount

	// DePledgeTime is the time at which the depledged collateral may be withdrawn
	DePledgeTime BlockHeight

	// Sectors is the set of all sectors this miner has committed
	Sectors SectorSet

	// ProvingSet is the set of sectors this miner is currently mining. It is only updated
	// when a PoSt is submitted (not as each new sector commitment is added)
	ProvingSet SectorSet

	// NextDoneSet is a set of sectors reported during the last PoSt submission as
	// being 'done'. The collateral for them is still being held until the next PoSt
	// submission in case early sector removal penalization is needed.
	NextDoneSet SectorSet

	// ArbitratedDeals is the set of deals this miner has been slashed for since the
	// last post submission
	ArbitratedDeals CidSet

	// TODO: maybe this number is redundant with power
	LockedStorage Integer

	// Power is the amount of power this miner has
	Power Integer

	// Asks are the set of active asks this miner has available
	Asks AskSet
}
```

### Code Cid
`<codec:raw><mhType:identity><"sminer">`

| Index     | Method Name       |
| -------- | ---------- |
| 1   | `AddAsk`     |
| 2 | `CommitSector` |
| 3 | `SubmitPoSt` |
| 4 | `IncreasePledge` |
| 5 | `SlashStorageFault` |
| 6 | `GetCurrentProvingSet` |
| 7 | `ArbitrateDeal` |
| 8 | `DePledge` |
| 9 | `GetOwner` |
| 10 | `GetWorkerAddr` |
| 11 | `GetPower` |
| 12 | `GetPeerID` |
| 13 | `GetSectorSize` |
| 14 | `UpdatePeerID` |
| 15 |  `ChangeWorker` |

### Constructor

Along with the call, the actor must be created with exactly enough filecoin for the collateral necessary for the pledge.
```go
func StorageMinerActor(worker Address, pledge BytesAmount, pid PeerID) {
	if msg.Value < CollateralForPledgeSize(pledge) {
		Fatal("not enough collateral given")
	}

	self.Owner = message.From
	self.Worker = worker
	self.PeerID = pid
	self.PledgeBytes = pledge
}
```


### AddAsk

Parameters:
- price TokenAmount
- ttl Integer

Return: AskID

```go
func AddAsk(price TokenAmount, ttl Integer) AskID {
	if msg.From != self.Worker {
		Fatal("Asks may only be added via the worker address")
	}

	// Filter out expired asks
	self.Asks.FilterExpired()

	askid := self.NextAskID
	self.NextAskID++

	self.Asks.Append(Ask{
		Price:  price,
		Expiry: CurrentBlockHeight + ttl,
		ID:     askid,
	})

	return askid
}
```

Note: this may be moved off chain soon, don't worry about testing it too heavily.



### CommitSector

Parameters:
- sectorID SectorID
- commD []byte
- commR []byte
- commRStar []byte
- proof SealProof

Return: None


```go
// NotYetSpeced: ValidatePoRep, EnsureSectorIsUnique, CollateralForSector, Commitment
func CommitSector(sectorID SectorID, commD, commR, commRStar []byte, proof SealProof) SectorID {
	if !miner.ValidatePoRep(miner.SectorSize, comm, miner.Worker, proof) {
		Fatal("bad proof!")
	}

	// make sure the miner isnt trying to submit a pre-existing sector
	if !miner.EnsureSectorIsUnique(comm) {
		Fatal("sector already committed!")
	}

	// make sure the miner has enough collateral to add more storage
	coll = CollateralForSector(miner.SectorSize)

	if coll < miner.Collateral {
		Fatal("not enough collateral")
	}

	miner.Collateral -= coll
	miner.ActiveCollateral += coll

	// Note: There must exist a unique index in the miner's sector set for each
	// sector ID. The `faults`, `recovered`, and `done` parameters of the
	// SubmitPoSt method express indices into this sector set.
	miner.Sectors.Add(sectorID, commR, commD)

	// if miner is not mining, start their proving period now
	// Note: As written here, every miners first PoSt will only be over one sector.
	// We could set up a 'grace period' for starting mining that would allow miners
	// to submit several sectors for their first proving period. Alternatively, we
	// could simply make the 'CommitSector' call take multiple sectors at a time.
	//
	// Note: Proving period is a function of sector size; small sectors take less
	// time to prove than large sectors do. Sector size is selected when pledging.
	if miner.ProvingSet.Size() == 0 {
		miner.ProvingSet = miner.Sectors
		miner.ProvingPeriodEnd = chain.Now() + ProvingPeriodDuration(miner.SectorSize)
	}
}
```

### SubmitPoSt

Parameters:
- proofs []PoStProof
- faults []FaultSet
- recovered Bitfield
- done Bitfield

Return: None

```go
// NotYetSpeced: ValidateFaultSets, GenerationAttackTime, ComputeLateFee
func SubmitPost(proofs []PoStProof, faults []FaultSet, recovered BitField, done BitField) {
	if msg.From != miner.Worker {
		Fatal("not authorized to submit post for miner")
	}

	// ensure recovered is a subset of the combined fault sets, and that done
	// does not intersect with either, and that all sets only reference sectors
	// that currently exist
	if !miner.ValidateFaultSets(faults, recovered, done) {
		Fatal("fault sets invalid")
	}

	var feesRequired TokenAmount

	if chain.Now() > miner.ProvingPeriodEnd+GenerationAttackTime(miner.SectorSize) {
		// TODO: determine what exactly happens here. Is the miner permanently banned?
		Fatal("Post submission too late")
	} else if chain.Now() > miner.ProvingPeriodEnd {
		feesRequired += ComputeLateFee(chain.Now() - miner.ProvingPeriodEnd)
	}

	feesRequired += ComputeTemporarySectorFailureFee(recovered)

	if msg.Value < feesRequired {
		Fatal("not enough funds to pay post submission fees")
	}

	// we want to ensure that the miner can submit more fees than required, just in case
	if msg.Value > feesRequired {
		Refund(msg.Value - feesRequired)
	}

	if !CheckPostProofs(miner.SectorSize, proofs, faults) {
		Fatal("proofs invalid")
	}

	permLostSet = AggregateBitfields(faults).Subtract(recovered)

	// adjust collateral for 'done' sectors
	miner.ActiveCollateral -= CollateralForSectors(miner.SectorSize, miner.NextDoneSet)
	miner.Collateral += CollateralForSectors(miner.SectorSize, miner.NextDoneSet)

	// penalize collateral for lost sectors
	miner.ActiveCollateral -= CollateralForSectors(miner.SectorSize, permLostSet)

	// burn funds for fees and collateral penalization
	BurnFunds(miner, CollateralForSectors(miner.SectorSize, permLostSet)+feesRequired)

	// update sector sets and proving set
	miner.Sectors.Subtract(done)
	miner.Sectors.Subtract(permLostSet)

	// update miner power to the amount of data actually proved during
	// the last proving period.
	oldPower := miner.Power
	miner.Power = SizeOf(Filter(miner.ProvingSet, faults)) * miner.SectorSize
	StorageMarket.UpdateStorage(miner.Power - oldPower)

	miner.ProvingSet = miner.Sectors

	// NEEDS REVIEW: early submission of PoSts may give the miner extra time for
	// their next PoSt, which could compound. Does the beacon reseeding for Posts
	// address this well enough?
	miner.ProvingPeriodEnd = miner.ProvingPeriodEnd + ProvingPeriodDuration(miner.SectorSize)

	// update next done set
	miner.NextDoneSet = done
	miner.ArbitratedDeals.Clear()
}

func ValidateFaultSets(faults []FaultSet, recovered, done BitField) bool {
	var aggregate BitField
	for _, fs := range faults {
		aggregate = aggregate.Union(fs.BitField)
	}

	// all sectors marked recovered must have actually failed
	if !recovered.IsSubsetOf(aggregate) {
		return false
	}

	// the done set cannot intersect with the aggregated faults
	// you can't mark a fault as 'done'
	if aggregate.Intersects(done) {
		return false
	}

	for _, bit := range aggregate.Bits() {
		if !miner.HasSectorByID(bit) {
			return false
		}
	}

	for _, bit := range done.Bits() {
		if !miner.HasSectorByID(bit) {
			return false
		}
	}

	return true
}

func ProvingPeriodDuration(sectorSize uint64) Integer {
	// TODO: eventually, this needs to be different for different sector sizes
	// The research team should give us concrete numbers
	return 24 * 60 * 60 * 2 // number of blocks in one day
}
```

### IncreasePledge

Parameters:

- addspace BytesAmount

Return: None


```go
func IncreasePledge(addspace BytesAmount) {
	// Note: msg.Value is implicitly transferred to the miner actor
	if miner.Collateral+msg.Value < CollateralForPledge(addspace+miner.PledgeBytes) {
		Fatal("not enough total collateral for the requested pledge")
	}

	miner.Collateral += msg.Value
	miner.PledgeBytes += addspace
}
```

### SlashStorageFault

Parameters:

- miner Address

Return: None

```go
func SlashStorageFault() {
	if self.SlashedAt > 0 {
		Fatal("miner already slashed")
	}

	if chain.Now() <= miner.ProvingPeriodEnd+GenerationAttackTime(miner.SectorSize) {
		Fatal("miner is not yet tardy")
	}

	if miner.ProvingSet.Size() == 0 {
		Fatal("miner is inactive")
	}

	// Strip miner of their power
	StorageMarketActor.UpdateStorage(-1 * self.Power)
	self.Power = 0

	// TODO: make this less hand wavey
	BurnCollateral(self.ConsensusCollateral)

	self.SlashedAt = CurrentBlockHeight
}
```

### GetCurrentProvingSet

Parameters: None

Return: `[][]byte`

```go
func GetCurrentProvingSet() [][]byte {
	return self.ProvingSet
}
```

Note: this is unlikely to ever be called on-chain, and will be a very large amount of data. We should reconsider the need for a list of all sector commitments (maybe fixing with accumulators?)

### ArbitrateDeal

This may be called by anyone to penalize a miner for dropping the data of a deal they committed to before the deal expires. Note: in order to call this, the caller must have the signed deal between the client and the miner in question, this would require out of band communication of this information from the client to acquire.

Parameters:
- deal Deal

Return: None

```go
func AbitrateDeal(d Deal) {
	if !ValidateSignature(d, self.Worker) {
		Fatal("invalid signature on deal")
	}

	if CurrentBlockHeight < d.StartTime {
		Fatal("Deal not yet started")
	}

	if d.Expiry < CurrentBlockHeight {
		Fatal("Deal is expired")
	}

	if !self.NextDoneSet.Has(d.PieceCommitment.Sector) {
		Fatal("Deal agreement not broken, or arbitration too late")
	}

	if self.ArbitratedDeals.Has(d.PieceRef) {
		Fatal("cannot slash miner twice for same deal")
	}

	pledge, storage := CollateralForDeal(d)

	// burn the pledge collateral
	self.BurnFunds(pledge)

	// pay the client the storage collateral
	TransferFunds(d.ClientAddr, storage)

	// make sure the miner can't be slashed twice for this deal
	self.ArbitratedDeals.Add(d.PieceRef)
}
```

TODO(scaling): This method, as currently designed, must be called once per sector. If a miner agrees to store 1TB (1000 sectors) for a particular client, and loses all that data, the client must then call this method 1000 times, which will be really expensive.

### DePledge

Parameters:
- amt TokenAmount

Return: None

```go
func DePledge(amt TokenAmount) {
	if msg.From != miner.Worker && msg.From != miner.Owner {
		Fatal("Not authorized to call DePledge")
	}

	if miner.DePledgeTime > 0 {
		if miner.DePledgeTime > CurrentBlockHeight {
			Fatal("too early to withdraw collateral")
		}

		TransferFunds(miner.Owner, miner.DePledgedCollateral)
		miner.DePledgeTime = 0
		miner.DePledgedCollateral = 0
		return
	}

	if amt > miner.Collateral {
		Fatal("Not enough free collateral to withdraw that much")
	}

	miner.Collateral -= amt
	miner.DePledgedCollateral = amt
	miner.DePledgeTime = CurrentBlockHeight + DePledgeCooldown
}
```

### GetOwner

Parameters: None

Return: Address

```go
func GetOwner() Address {
	return self.Owner
}
```

### GetWorkerAddr

Parameters: None

Return: Address

```go
func GetWorkerAddr() Address {
	return self.Worker
}
```

### GetPower

Parameters: None

Return: BytesAmount

```go
func GetPower() BytesAmount {
	return self.Power
}
```

### GetPeerID

Parameters: None

Return: PeerID

```go
func GetPeerID() PeerID {
	return self.PeerID
}
```

### GetSectorSize

Parameters: None

Return: BytesAmount

```go
func GetSectorSize() BytesAmount {
	return self.SectorSize
}
```

### UpdatePeerID

Parameters:
- pid PeerID

Return: None

```go
func UpdatePeerID(pid PeerID) {
	if msg.From != self.Worker {
		Fatal("only the mine worker may update the peer ID")
	}

	self.PeerID = pid
}
```

### ChangeWorker

Changes the worker address. Note that since Sector Commitments take the miners worker key as an input, any sectors sealed with the old key but not yet submitted to the chain will be invalid. All future sectors must be sealed with the new worker key.

Parameters:
- addr Address

Return: None

```go
func ChangeWorker(addr Address) {
	if msg.From != self.Owner {
		Fatal("only the owner can change the worker address")
	}

	self.Worker = addr
}
```

## Payment Channel Broker Actor

TODO

## Multisig Account Actor

A basic multisig account actor. Allows sending of messages like a normal account actor, but with the requirement of M of N parties agreeing to the operation. Completed and/or cancelled operations stick around in the actors state until explicitly cleared out. Proposers may cancel transactions they propose, or transactions by proposers who are no longer approved signers.

Self modification methods (add/remove signer, change requirement) are called by
doing a multisig transaction invoking the desired method on the contract itself. This means the 'signature
threshold' logic only needs to be implemented once, in one place.

The [init actor](#init-actor) is used to create new instances of the multisig.

#### State

```go
type Multisig struct {
	Signers  []Address
	Required uint

	NextTxID     uint64
	Transactions map[int]Transaction
}

type Transaction struct {
	Created   uint64
	TxID      uint64
	To        Address
	Value     TokenAmount
	Method    string
	Params    []byte
	Approved  []Address
	Completed bool
	Canceled  bool
}
```

### Code Cid
`<codec:raw><mhType:identity><"multisig">`

| Index     | Method Name       |
| -------- | ---------- |
| 1   | `Propose`     |
| 2 | `Approve` |
| 3 | `Cancel` |
| 4 | `ClearCompleted` |
| 5 | `AddSigner` |
| 6 | `RemoveSigner` |
| 7 | `SwapSigner` |
| 8 | `ChangeRequirement` |

#### Constructor

>  This method sets up the initial state for the multisig account

#### Parameters

| Name     | Type       | Description                                         |
| -------- | ---------- | ------------------------------------------------------------ |
| `signers`   | `[]Address`      | The addresses that will be the signatories of this wallet. |
| `required` | `uint` | The number of signatories required to perform a transaction.       |

```go
func Multisig(signers []Address, required uint) {
	self.Signers = signers
	self.Required = required
}
```



### Propose

> Propose is used to propose a new transaction to be sent by this multisig. The proposer must be a signer, and the proposal also serves as implicit approval from the proposer. If only a single signature is required, then the transaction is executed immediately.

| Name     | Type       | Description  |
| --- | --- | --- |
| `to` |  `Address` | The address of the target of the proposed transaction. |
|  `value` | `TokenAmount` | The amount of funds to send with the proposed transaction |
|  `method` | `string` | The method that will be invoked on the proposed transactions target. |
| `params` | `[]byte` | The parameters that will be passed to the method invocation on the proposed transactions target. |

```go
func Propose(to Address, value TokenAmount, method string, params []byte) uint64 {
	if !isSigner(msg.From) {
		Fatal("not authorized")
	}

	txid := self.NextTxID
	self.NextTxID++

	tx := Transaction{
		TxID:     txid,
		To:       to,
		Value:    value,
		Method:   method,
		Params:   params,
		Approved: []Address{msg.From},
	}

	self.Transactions.Append(tx)

	if self.Required == 1 {
		vm.Send(tx.To, tx.Value, tx.Method, tx.Params)
		tx.Complete = true
	}

	return txid
}
```

### Approve

> Approve is called by a signer to approve a given transaction. If their approval pushes the approvals for this transaction over the threshold, the transaction is executed.

| Name     | Type       | Description  |
| --- | --- | --- |
| `txid` |  `uint64` | The ID of the transaction to approve. |


```go
func Approve(txid uint64) {
	if !self.isSigner(msg.From) {
		Fatal("not authorized")
	}

	tx := self.getTransaction(txid)
	if tx.Complete {
		Fatal("transaction already completed")
	}
	if tx.Canceled {
		Fatal("transaction canceled")
	}

	for _, signer := range tx.Approved {
		if signer == msg.From {
			Fatal("already signed this message")
		}
	}

	tx.Approved.Append(msg.From)

	if len(tx.Approved) >= self.Required {
		Send(tx.To, tx.Value, tx.Method, tx.Params)
		tx.Complete = true
	}
}
```

### Cancel

```go
func Cancel(txid uint64) {
	if !self.isSigner(msg.From) {
		Fatal("not authorized")
	}

	tx := self.getTransaction(txid)
	if tx.Complete {
		Fatal("cannot cancel completed transaction")
	}
	if tx.Canceled {
		Fatal("transaction already canceled")
	}

	proposer := tx.Approved[0]
	if proposer != msg.From && isSigner(proposer) {
		Fatal("cannot cancel another signers transaction")
	}

	tx.Canceled = true
}
```

### ClearCompleted

```go
func ClearCompleted() {
	if !self.isSigner(msg.From) {
		Fatal("not authorized")
	}

	for tx := range self.Transactions {
		if tx.Completed || tx.Canceled {
			self.Transactions.Remove(tx)
		}
	}
}
```

### AddSigner

```go
func AddSigner(signer Address) {
	if msg.From != self.Address {
		Fatal("add signer must be called by wallet itself")
	}
	if self.isSigner(signer) {
		Fatal("new address is already a signer")
	}

	self.Signers.Append(signer)
}
```

### RemoveSigner

```go
func RemoveSigner(signer Address) {
	if msg.From != self.Address {
		Fatal("remove signer must be called by wallet itself")
	}
	if !self.isSigner(signer) {
		Fatal("given address was not a signer")
	}

	self.Signers.Remove(signer)
}
```

### SwapSigner

```go
func SwapSigner(old, new Address) {
	if msg.From != self.Address {
		Fatal("swap signer must be called by wallet itself")
	}
	if !self.isSigner(old) {
		Fatal("given old address was not a signer")
	}
	if self.isSigner(new) {
		Fatal("given new address was already a signer")
	}

	self.Signers.Remove(old)
	self.Signers.Append(new)
}
```

### ChangeRequirement

```go
func ChangeRequirement(req int) {
	if msg.From != self.Address {
		Fatal("change requirement must be called by wallet itself")
	}
	if req < 1 {
		Fatal("requirement must be at least 1")
	}

	self.Required = req
}
```

### Helper Methods

The various helper methods called above are defined here.

```go
func isSigner(a Address) bool {
	for signer := range self.Signers {
		if a == signer {
			return true
		}
	}
	return false
}

func getTransaction(txid int) Transaction {
	tx, ok := self.Transactions[txid]
	if !ok {
		Fatal("no such transaction")
	}

	return tx
}
```
