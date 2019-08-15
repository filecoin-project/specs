# Filecoin State Machine Actors

Any implementations of the Filecoin actors must be exactly byte for byte compatible with the go-filecoin actor implementations. The pseudocode below tries to capture the important logic, but capturing all the detail would require embedding exactly the code from go-filecoin, so for now, its simply informative pseudocode. The algorithms below are correct, and all implementations much match it (including go-filecoin), but details omitted from here should be looked for in the go-filecoin code.

This spec describes a set of actors that operate within the [Filecoin State Machine](state-machine.md). All types are defined in [the basic type encoding spec](data-structures.md#basic-type-encodings).

## Actor State

Each actor type defines their own structure for storing their state. We
represent each with an IPLD schema at the beginning of each actor section in
this document.

## System Actors

Some state machine actors are 'system' actors that get instantiated in the genesis block, and have their IDs allocated at that point.

| ID   | Actor              | Name                    |
| ---- | ------------------ | ----------------------- |
| 0    | InitActor          | Network Init            |
| 1    | AccountActor       | Network Treasury        |
| 2    | StorageMarketActor | Filecoin Storage Market |
| 99   | AccountActor       | Burnt Funds             |

## Built In Actors

### Init Actor

- **Code Cid**: `<codec:raw><mhType:identity><"init">`

The init actor is responsible for creating new actors on the filecoin network. This is a built-in actor and cannot be replicated. In the future, this actor will be responsible for loading new code into the system (for user programmable actors). ID allocation for user instantiated actors starts at 100. This means that `NextID` will initially be set to 100.

```sh
type InitActorState struct {
    addressMap {Address:ID}<Hamt>
    nextId UInt
}
```

#### Methods

| Name | Method ID |
|--------|-------------|
| `Constructor` | 1 |
| `Exec` | 2 |
| `GetIdForAddress` | 3 |

#### `Constructor`

**Parameters**

```sh
type InitConstructor struct {
}

```

**Algorithm**

#### `Exec`

This method is the core of the `Init Actor`. It handles instantiating new actors and assigning them their IDs.

**Parameters**

```sh
type Exec struct {
    ## Reference to the location at which the code of the actor to create is stored.
    code &Code
    ## Parameters passed to the constructor of the actor.
    params ActorMethod
} representation tuple
```

**Algorithm**

```go
func Exec(code Cid, params ActorMethod) Address {
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
	// TODO: where do `creator` and `nonce` come from?
	addr := VM.ComputeActorAddress(creator, nonce)

	// Set up the actor itself
	actor := Actor{
		Code:    code,
		Balance: msg.Value,
		Head:    nil,
		Nonce:   0,
	}

	// The call to the actors constructor will set up the initial state
	// from the given parameters, setting `actor.Head` to a new value when successfull.
	// TODO: can constructors fail?
	actor.Constructor(params)

	VM.GlobalState.Set(actorID, actor)

	// Store the mapping of address to actor ID.
	self.AddressMap[addr] = actorID

	return addr
}

func IsSingletonActor(code Cid) bool {
	return code == StorageMarketActor || code == InitActor
}
```

```go
// TODO: find a better home for this logic
func (VM VM) ComputeActorAddress(creator Address, nonce Integer) Address {
	return NewActorAddress(bytes.Concat(creator.Bytes(), nonce.BigEndianBytes()))
}
```

#### `GetIdForAddress`

This method allows for fetching the corresponding ID of a given Address

**Parameters**

```sh
type GetIdForAddress struct {
    addr Address
} representation tuple
```

**Algorithm**

```go
func GetIdForAddress(addr Address) UInt {
	id := self.AddressMap[addr]
	if id == nil {
		Fault("unknown address")
	}
	return id
}
```

### Account Actor

- **Code Cid**: `<codec:raw><mhType:identity><"account">`

The Account actor is the actor used for normal keypair backed accounts on the filecoin network.

```sh
type AccountActorState struct {
    address Address
}
```

#### Methods

| Name | Method ID |
|--------|-------------|
| `AccountConstructor` | 1 |
| `GetAddress` | 2 |

```
type AccountConstructor struct {
}
```

#### `GetAddress`

**Parameters**

```sh
type GetAddress struct {
} representation tuple
```

**Algorithm**

```go
func GetAddress() Address {
	return self.address
}
```

### Storage Market Actor

* **Code Cid**: `<codec:raw><mhType:identity><"smarket">`

The storage market actor is the central point for the Filecoin storage market. It is responsible for registering new miners to the system, and maintaining the power table. The Filecoin storage market is a singleton that lives at a specific well-known address.

```sh
type StorageMarketActorState struct {
    miners {Address:Null}<Hamt>
    totalStorage BytesAmount
}
```

#### Methods

| Name | Method ID |
|--------|-------------|
| `StorageMarketConstructor` | 1 |
| `CreateStorageMiner` | 2 |
| `SlashConsensusFault` | 3 |
| `UpdateStorage` | 4 |
| `GetTotalStorage` | 5 |
| `PowerLookup` | 6 |
| `IsMiner` | 7 |
| `StorageCollateralForSize` | 8 |

#### `Constructor`

**Parameters**

```sh
type StorageMarketConstructor struct {}
```

**Algorithm**

#### `CreateStorageMiner`

**Parameters**

```sh
type CreateStorageMiner struct {
    worker Address
    sectorSize BytesAmount
    peerId PeerId
} representation tuple
```

**Algorithm**

```go
func CreateStorageMiner(worker Address, owner Address, sectorSize BytesAmount, pid PeerID) Address {
	if !SupportedSectorSize(sectorSize) {
		Fatal("Unsupported sector size")
	}

	newminer := InitActor.Exec(MinerActorCodeCid, EncodeParams(worker, owner, pledge, sectorSize, pid))

	self.Miners.Add(newminer)

	return newminer
}
```

#### `SlashConsensusFault`

**Parameters**

```sh
type SlashConsensusFault struct {
    block1 &Block
    block2 &Block
} representation tuple
```

**Algorithm**

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

#### `UpdateStorage`

UpdateStorage is used to update the global power table.

**Parameters**

```sh
type UpdateStorage struct {
    delta BytesAmount
} representation tuple

```

**Algorithm**

```go
func UpdateStorage(delta BytesAmount) {
	if !self.Miners.Has(msg.From) {
		Fatal("update storage must only be called by a miner actor")
	}

	self.TotalStorage += delta
}
```

#### `GetTotalStorage`

**Parameters**

```sh
type GetTotalStorage struct {

} representation tuple
```

**Algorithm**

```go
func GetTotalStorage() BytesAmount {
	return self.TotalStorage
}
```

#### `PowerLookup`

**Parameters**

```sh
type PowerLookup struct {
    miner Address
} representation tuple
```

**Algorithm**

```go
func PowerLookup(miner Address) BytesAmount {
	if !self.Miners.Has(miner) {
		Fatal("miner not registered with storage market")
	}

	mact := LoadMinerActor(miner)

	return mact.GetPower()
}
```

#### `IsMiner`

**Parameters**

```sh
type IsMiner struct {
    addr Address
} representation tuple
```

**Algorithm**

```go
func IsMiner(addr Address) bool {
	return self.Miners.Has(miner)
}
```

#### `StorageCollateralForSize`


**Parameters**

```sh
type StorageCollateralForSize struct {
    size UInt
} representation tuple
```

**Algorithm**

```go
func StorageCollateralforSize(size UInt) TokenAmount {
	// TODO:
}
```

## Storage Miner Actor

* **Code Cid**: `<codec:raw><mhType:identity><"sminer">`

```sh
type StorageMinerActorState struct {
  ## contains mostly static info about this miner
  info &MinerInfo


    ## Collateral that is waiting to be withdrawn.
    dePledgedCollateral TokenAmount

	## Time at which the depledged collateral may be withdrawn.
    dePledgeTime BlockHeight

	## All sectors this miner has committed.
    sectors &SectorSet

	## Sectors this miner is currently mining. It is only updated
	## when a PoSt is submitted (not as each new sector commitment is added).
    provingSet &SectorSet

	## Sectors reported during the last PoSt submission as being 'done'. The collateral
    ## for them is still being held until the next PoSt submission in case early sector
    ## removal penalization is needed.
    nextDoneSet BitField

	## Deals this miner has been slashed for since the last post submission.
    arbitratedDeals {Cid:Null}

	## Amount of power this miner has.
    power UInt

    ## List of sectors that this miner was slashed for.
    slashedSet optional &SectorSet

    ## The height at which this miner was slashed at.
    slashedAt optional BlockHeight

    ## The amount of storage collateral that is owed to clients, and cannot be used for collateral anymore.
    owedStorageCollateral TokenAmount

    provingPeriodEnd BlockHeight
}

type MinerInfo struct {
	## Account that owns this miner.
    ## - Income and returned collateral are paid to this address.
    ## - This address is also allowed to change the worker address for the miner.
    owner Address

	## Worker account for this miner.
	## This will be the key that is used to sign blocks created by this miner, and
	## sign messages sent on behalf of this miner to commit sectors, submit PoSts, and
	## other day to day miner activities.
    worker Address

    ## Libp2p identity that should be used when connecting to this miner.
    peerId PeerId

    ## Amount of space in each sector committed to the network by this miner.
    sectorSize BytesAmount

}
```

#### Methods

| Name | Method ID |
|--------|-------------|
| `StorageMinerConstructor` | 1 |
| `CommitSector` | 2 |
| `SubmitPost` | 3 |
| `SlashStorageFault` | 4 |
| `GetCurrentProvingSet` | 5 |
| `ArbitrateDeal` | 6 |
| `DePledge` | 7 |
| `GetOwner` | 8 |
| `GetWorkerAddr` | 9 |
| `GetPower` | 10 |
| `GetPeerID` | 11 |
| `GetSectorSize` | 12 |
| `UpdatePeerID` | 13 |
| `ChangeWorker` | 14 |
| `IsSlashed` |  15 |
| `IsLate` | 16 |
| `PaymentVerify` | 17 |

#### `Constructor`

Along with the call, the actor must be created with exactly enough filecoin for the collateral necessary for the pledge.

**Parameters**

```sh
type StorageMinerConstructor struct {
    worker Address
    owner Address
    sectorSize BytesAmount
    peerId PeerId
} representation tuple
```

**Algorithm**

```go
func StorageMinerActor(worker Address, owner Address, sectorSize BytesAmount, pid PeerID) {
	self.info.owner = message.From
	self.info.worker = worker
	self.info.peerID = pid
	self.info.sectorSize = sectorSize

	self.sectors = EmptySectorSet()
	self.provingSet = EmptySectorSet()
}
```

#### `CommitSector`

**Parameters**

```sh
type CommitSector struct {
    sectorId SectorID
    commD Bytes
    commR Bytes
    commRStar Bytes
    proof SealProof
} representation tuple
```

**Algorithm**

{{% notice todo %}}
TODO: ValidatePoRep, EnsureSectorIsUnique, CollateralForSector, Commitment
{{% /notice %}}

```go
func CommitSector(sectorID SectorID, commD, commR, commRStar []byte, proof SealProof) SectorID {
	if !self.ValidatePoRep(self.info.sectorSize, comm, self.info.worker, proof) {
		Fatal("bad proof!")
	}

	// make sure the miner isnt trying to submit a pre-existing sector
	if !self.EnsureSectorIsUnique(comm) {
		Fatal("sector already committed!")
	}

	// Power of the miner after adding this sector
	futurePower = self.power + self.info.sectorSize
	collateralRequired = CollateralForPower(futurePower)

	if collateralRequired > vm.MyBalance() {
		Fatal("not enough collateral")
	}

	// ensure that the miner cannot commit more sectors than can be proved with a single PoSt
	if miner.Sectors.Size() >= POST_SECTORS_COUNT {
		Fatal("too many sectors")
	}

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

func CollateralForPower(power BytesAmount) TokenAmount {
	availableFil = FakeGlobalMethods.GetAvailableFil()
	totalNetworkPower = StorageMinerActor.GetTotalStorage()
	numMiners = StorageMarket.GetMinerCount()
	powerCollateral = availableFil * NetworkConstants.POWER_COLLATERAL_PROPORTION * power / totalNetworkPower
	perCapitaCollateral = availableFil * NetworkConstants.PER_CAPITA_COLLATERAL_PROPORTION / numMiners
	collateralRequired = math.Ceil(minerPowerCollateral + minerPerCapitaCollateral)
	return collateralRequired
}
```

#### `SubmitPoSt`

**Parameters**

```sh
type SubmitPost struct {
    proofs PoStProof
    faults [FaultSet]
    recovered Bitfield
    done Bitfield
} representation tuple
```

**Algorithm**

{{% notice todo %}}
TODO: GenerationAttackTime
{{% /notice %}}

```go
func SubmitPost(proofs PoStProof, faults []FaultSet, recovered Bitfield, done Bitfield) {
	if msg.From != miner.Worker {
		Fatal("not authorized to submit post for miner")
	}

	// ensure recovered is a subset of the combined fault sets, and that done
	// does not intersect with either, and that all sets only reference sectors
	// that currently exist
	allFaults = AggregateBitfields(faults)
	if !miner.ValidateFaultSets(faults, recovered, done) {
		Fatal("fault sets invalid")
	}

	feesRequired := 0

	if chain.Now() > self.ProvingPeriodEnd+GenerationAttackTime(self.SectorSize) {
		// slashing ourselves
		SlashStorageFault(self)
		return
	} else if chain.Now() > self.ProvingPeriodEnd {
		feesRequired += ComputeLateFee(self.power, chain.Now()-self.provingPeriodEnd)
	}

	feesRequired += ComputeTemporarySectorFailureFee(self.sectorSize, recovered)

	if msg.Value < feesRequired {
		Fatal("not enough funds to pay post submission fees")
	}

	// we want to ensure that the miner can submit more fees than required, just in case
	if msg.Value > feesRequired {
		TransferFunds(msg.From, msg.Value-feesRequired)
	}

	if !CheckPostProof(miner.SectorSize, proof, faults) {
		Fatal("proof invalid")
	}

	// combine all the fault set bitfields, and subtract out the recovered
	// ones to get the set of sectors permanently lost
	permLostSet = allFaults.Subtract(recovered)

	// burn funds for fees and collateral penalization
	self.BurnFunds(CollateralForSize(self.SectorSize*permLostSet.Size()) + feesRequired)

	// update sector sets and proving set
	miner.Sectors.Subtract(done)
	miner.Sectors.Subtract(permLostSet)

	// update miner power to the amount of data actually proved during
	// the last proving period.
	oldPower := miner.Power

	miner.Power = (miner.ProvingSet.Size() - allFaults.Count()) * miner.SectorSize
	StorageMarket.UpdateStorage(miner.Power - oldPower)

	miner.ProvingSet = miner.Sectors

	// Updating proving period given a fixed schedule, independent of late submissions.
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
	return 24 * 60 * 60 * 2 // number of blocks in one day
}

func ComputeLateFee(power Integer, blocksLate Integer) TokenAmount {
	return 4 // TODO: real collateral calculation, obviously
}

func ComputeTemporarySectorFailureFee(sectorSize BytesAmount, numSectors Integer) TokenAmount {
	return 4 // TODO: something tells me that 4 might not work in all situations. probably should find a better way to compute this
}
```

#### `SlashStorageFault`

**Parameters**

```sh
type SlashStorageFault struct {
    miner Address
} representation tuple
```

**Algorithm**


## late submission and resubmission

- If After provingPeriodEnd
  - post submission now requires paying pay late submission fee
    - (fixed cost, dependent on the total storage size)
  - [implicit] loose all power, can be explicitly slashed for
  - post submission returns your power immediately
- If After `GenerationAttackTimeout` (<< 1 proving period)
  - .... nothing changes atm
- If After `PoStTimeout` (< 1 proving period)
  - [explicit - slashStorageFault] loose all storage collateral
  - clients can arbitrate deals with the miner now
  - post submission now requires paying both late fee + lost storage collateral
  - the valid post submission returns your power immediately
- If After `SectorFailureTimeout` (> 1 proving period)
  - [explicit - slashStorageFault] loose all sectors
     - [implicit] resets proving period, as they need to resubmit all sectors after this
  - there is now no way to reintroduce the sectors, unless they are resubmitted, etc
  - power can not be recovered anymore

- If [miner] does not call postsubmit


Notes:
- Should post submission only return your power, after the following proving period, when after the generation attack timeout
- Is one proving period enough time for abitration of faulty deals for a client?

```go
func SlashStorageFault() {
	// You can only be slashed once for missing your PoSt.
	if self.SlashedAt > 0 {
		Fatal("miner already slashed")
	}

	// Only if the miner is actually late, they can be slashed.
	if chain.Now() <= self.ProvingPeriodEnd+GenerationAttackTime(self.SectorSize) {
		Fatal("miner is not yet tardy")
	}

	// Only a miner who is expected to prove, can be slashed.
	if self.ProvingSet.Size() == 0 {
		Fatal("miner is inactive")
	}

	// Strip the miner of their power.
	StorageMarketActor.UpdateStorage(-1 * self.Power)
	self.Power = 0

	self.slashedSet = self.ProvingSet
	// remove proving set from our sectors
	self.sectors.Substract(self.slashedSet)

	// clear proving set
	self.ProvingSet = nil

	self.owedStorageCollateral = StorageMarketActor.StorageCollateralForSize(
		self.slashedSet.Size() * self.SectorSize,
	)

	self.SlashedAt = CurrentBlockHeight
}
```

#### `GetCurrentProvingSet`

**Parameters**

```sh
type GetCurrentProvingSet struct {
} representation tuple
```

**Algorithm**

```go
func GetCurrentProvingSet() [][]byte {
	return self.ProvingSet
}
```

{{%notice note %}}
**Note**: this is unlikely to ever be called on-chain, and will be a very large amount of data. We should reconsider the need for a list of all sector commitments (maybe fixing with accumulators?)
{{% /notice %}}

#### `ArbitrateDeal`

This may be called by anyone to penalize a miner for dropping the data of a deal they committed to before the deal expires. Note: in order to call this, the caller must have the signed deal between the client and the miner in question, this would require out of band communication of this information from the client to acquire.

**Parameters**

```sh
type ArbitrateDeal struct {
    deal Deal
} representation tuple
```

**Algorithm**

```go
func AbitrateDeal(deal Deal) {
	if !VM.ValidateSignature(deal, self.Worker) {
		Fatal("invalid signature on deal")
	}

	if VM.CurrentBlockHeight() < deal.StartTime {
		Fatal("Deal not yet started")
	}

	if deal.Expiry < VM.CurrentBlockHeight() {
		Fatal("Deal is expired")
	}

	if !self.NextDoneSet.Has(deal.pieceInclusionProof.sectorID) {
		Fatal("Deal agreement not broken, or arbitration too late")
	}

	if self.ArbitratedDeals.Has(deal.commP) {
		Fatal("cannot slash miner twice for same deal")
	}

	if !deal.pieceInclusionProof.Verify(deal.commP, deal.size) {
		Fatal("invalid piece inclusion proof or size")
	}

	storageCollateral := StorageMarketActor.StorageCollateralForSize(deal.size)

	if self.owedStorageCollateral < storageCollateral {
		Fatal("math is hard, and we didnt do it right")
	}

	// pay the client the storage collateral
	VM.TransferFunds(storageCollateral, deal.client)

	// keep track of how much we have payed out
	self.owedStorageCollateral -= storageCollateral

	// make sure the miner can't be slashed twice for this deal
	self.ArbitratedDeals.Add(deal.commP)
}
```

{{% notice todo %}}
**TODO(scaling)**: This method, as currently designed, must be called once per sector. If a miner agrees to store 1TB (1000 sectors) for a particular client, and loses all that data, the client must then call this method 1000 times, which will be really expensive.
{{% /notice %}}

#### `DePledge`

**Parameters**

```sh
type DePledge struct {
    amount TokenAmount
} representation tuple
```

**Algorithm**

```go
func DePledge(amt TokenAmount) {
	if msg.From != self.info.Worker && msg.From != self.info.owner {
		Fatal("Not authorized to call DePledge")
	}

	if self.DePledgeTime > 0 {
		if self.DePledgeTime > VM.CurrentBlockHeight() {
			Fatal("too early to withdraw collateral")
		}

		TransferFunds(self.info.owner, self.DePledgedCollateral)
		self.DePledgeTime = 0
		self.DePledgedCollateral = 0
		return
	}

	collateralRequired = CollateralForPower(self.power)

	if amt+collateralRequired > vm.MyBalance() {
		Fatal("Not enough free collateral to withdraw that much")
	}

	self.DePledgedCollateral = amt
	self.DePledgeTime = CurrentBlockHeight + DePledgeCooldown
}
```

#### `GetOwner`

**Parameters**
```sh
type GetOwner struct {
} representation tuple
```

**Algorithm**

```go
func GetOwner() Address {
	return self.info.owner
}
```

#### `GetWorkerAddr`

**Parameters**

```sh
type GetWorkerAddr struct {
} representation tuple
```

**Algorithm**

```go
func GetWorkerAddr() Address {
	return self.info.worker
}
```

#### `GetPower`

**Parameters**

```sh
type GetPower struct {
} representation tuple
```

**Algorithm**

```go
func GetPower() BytesAmount {
	return self.power
}
```

#### `GetPeerID`

**Parameters**

```sh
type GetPeerID struct {
} representation tuple
```

**Algorithm**

```go
func GetPeerID() PeerID {
	return self.info.peerID
}
```

#### `GetSectorSize`

**Parameters**

```sh
type GetSectorSize struct {
} representation tuple
```

**Algorithm**

```go
func GetSectorSize() BytesAmount {
	return self.info.sectorSize
}
```

#### `UpdatePeerID`

**Parameters**

```sh
type UpdatePeerID struct {
    peerId PeerId
} representation tuple
```

**Algorithm**

```go
func UpdatePeerID(pid PeerID) {
	if msg.From != self.info.worker {
		Fatal("only the mine worker may update the peer ID")
	}

	self.info.peerID = pid
}
```

#### `ChangeWorker`

Changes the worker address. Note that since Sector Commitments take the miners worker key as an input, any sectors sealed with the old key but not yet submitted to the chain will be invalid. All future sectors must be sealed with the new worker key.

**Parameters**

```sh
type ChangeWorker struct {
    addr Address
} representation tuple
```

**Algorithm**

```go
func ChangeWorker(addr Address) {
	if msg.From != self.info.owner {
		Fatal("only the owner can change the worker address")
	}

	self.info.worker = addr
}
```

#### `IsLate`

IsLate checks whether the miner has submitted their PoSt on time (i.e. not after ProvingPeriodEnd).

**Parameters**

```sh
type IsLate struct {
} representation tuple
```

**Algorithm**

```go
func IsLate() (bool) {
    return self.provingPeriodEnd < VM.CurrentBlockHeight()
}
```

#### `IsSlashed`

Checks whether the miner has been slashed and not recovered. Note that if the miner is slashed and recovers, this will return False: it checks current state rather than historical occurence.

**Parameters**

```sh
type IsSlashed struct {
} representation tuple
```

**Algorithm**

```go
func IsSlashed() (bool) {
    # SlashedAt is reset on recovery
    return self.SlashedAt > 0
}
```

#### `PaymentVerify`

Verifies a storage market payment channel voucher's 'Extra' data.

**Parameters**

```sh
type PaymentVerify struct {
    Extra StorageVoucherData
    Proof Bytes
} representation tuple

type StorageVoucherData struct {
    DataCommitment Bytes
    RequiredSector Bytes
} representation tuple
```

**Algorithm**

```go
func PaymentVerify(extra StorageVoucherData, proof PieceInclusionProof) {
	if extra.DataCommitment != nil {
		if !self.ValidateInclusion(proof, extra.DataCommitment) {
			Fatal("piece inclusion proof was invalid")
		}
		return
	}

	if extra.RequiredSector != nil {
		if !self.HasSector(extra.RequiredSector) {
			Fatal("miner does not have required sector")
		}
		return
	}

	Fatal("voucher data contained neither data commitment or required sector")
}

func ValidateInclusion(pip PieceInclusionProof, dataCommitment Bytes) bool {
  if !self.HasSector(pip.SectorID) {
    Fatal("miner does not have sector referenced by piece inclusion proof")
  }

  return ValidatePIP(pip, dataCommitment)
}
```

### Payment Channel Actor

- **Code Cid:** `<codec:raw><mhType:identity><"paych">`

The payment channel actor manages the on-chain state of a point to point payment channel.

```sh
type PaymentChannel struct {
	from Address
	to   Address

	toSend       TokenAmount

	closingAt      UInt
	minCloseHeight UInt

	laneStates {UInt:LaneState}
} representation tuple

type SignedVoucher struct {
  TimeLock BlockHeight
  SecretPreimage Bytes
  Extra ModVerifyParams
  Lane Uint
  Nonce Uint
  Merges []Merge
  Amount TokenAmount
  MinCloseHeight Uint

  Signature Signature
}

type ModVerifyParams struct {
  Actor Address
  Method Uint
  Data Bytes
}

type Merge struct {
  Lane Uint
  Nonce Uint
}

type LaneState struct {
  Closed bool
  Redeemed TokenAmount
  Nonce Uint
}

type PaymentChannelMethod union {
  | PaymentChannelConstructor 0
  | UpdateChannelState 1
  | Close 2
  | Collect 3
} representation keyed
```

#### Methods

| Name | Method ID |
|--------|-------------|
| `Constructor` | 1 |
| `UpdateChannelState` | 2 |
| `Close` | 3 |
| `Collect` | 4 |

#### `Constructor`

**Parameters**

```sh
type PaymentChannelConstructor struct {
  to Address
}
```

**Algorithm**

{{% notice todo %}}

TODO: Define me

{{% /notice %}}

#### `UpdateChannelState`

**Parameters**

```sh
type UpdateChannelState struct {
  sv SignedVoucher
  secret Bytes
  proof Bytes
} representation tuple
```

**Algorithm**

```go
func UpdateChannelState(sv SignedVoucher, secret []byte, proof []byte) {
	if !self.validateSignature(sv) {
		Fatal("Signature Invalid")
	}

	if chain.Now() < sv.TimeLock {
		Fatal("cannot use this voucher yet!")
	}

	if sv.SecretPreimage != nil {
		if Hash(secret) != sv.SecretPreimage {
			Fatal("Incorrect secret!")
		}
	}

	if sv.Extra != nil {
		ret := vmctx.Send(sv.Extra.Actor, sv.Extra.Method, sv.Extra.Data, proof)
		if ret != 0 {
			Fatal("spend voucher verification failed")
		}
	}

	ls := self.LaneStates[sv.Lane]
	if ls.Closed {
		Fatal("cannot redeem a voucher on a closed lane")
	}

	if ls.Nonce > sv.Nonce {
		Fatal("voucher has an outdated nonce, cannot redeem")
	}

	var mergeValue TokenAmount
	for _, merge := range sv.Merges {
		if merge.Lane == sv.Lane {
			Fatal("voucher cannot merge its own lane")
		}

		ols := self.LaneStates[merge.Lane]
		if ols.Nonce >= merge.Nonce {
			Fatal("merge in voucher has outdated nonce, cannot redeem")
		}

		mergeValue += ols.Redeemed
		ols.Nonce = merge.Nonce
	}

	ls.Nonce = sv.Nonce
	balanceDelta = sv.Amount - (mergeValue + ls.Redeemed)
	ls.Redeemed = sv.Amount

	newSendBalance = self.ToSend + balanceDelta
	if newSendBalance < 0 {
		// TODO: is this impossible?
		Fatal("voucher would leave channel balance negative")
	}

	if newSendBalance > self.Balance {
		Fatal("not enough funds in channel to cover voucher")
	}

	self.ToSend = newSendBalance

	if sv.MinCloseHeight != 0 {
		if self.ClosingAt != 0 && self.ClosingAt < sv.MinCloseHeight {
			self.ClosingAt = sv.MinCloseHeight
		}
		if self.MinCloseHeight < sv.MinCloseHeight {
			self.MinCloseHeight = sv.MinCloseHeight
		}
	}
}

func Hash(b []byte) []byte {
	return blake2b.Sum(b)
}
```

#### `Close`

**Parameters**

```sh
type Close struct {
} representation tuple
```

**Algorithm**

```go
const ChannelClosingDelay = 6 * 60 * 2 // six hours

func Close() {
	if msg.From != self.From && msg.From != self.To {
		Fatal("not authorized to close channel")
	}
	if self.ClosingAt != 0 {
		Fatal("Channel already closing")
	}

	self.ClosingAt = chain.Now() + ChannelClosingDelay
	if self.ClosingAt < self.MinCloseHeight {
		self.ClosingAt = self.MinCloseHeight
	}
}
```

#### `Collect`

**Parameters**

```sh
type Collect struct {
} representation tuple
```

**Algorithm**

```go
func Collect() {
	if self.ClosingAt == 0 {
		Fatal("payment channel not closing or closed")
	}

	if chain.Now() < self.ClosingAt {
		Fatal("Payment channel not yet closed")
	}

	TransferFunds(self.From, self.Balance-self.ToSend)
	TransferFunds(self.To, self.ToSend)
  self.ToSend = 0
}
```

### Multisig Account Actor

- **Code Cid**: `<codec:raw><mhType:identity><"multisig">`

A basic multisig account actor. Allows sending of messages like a normal account actor, but with the requirement of M of N parties agreeing to the operation. Completed and/or cancelled operations stick around in the actors state until explicitly cleared out. Proposers may cancel transactions they propose, or transactions by proposers who are no longer approved signers.

Self modification methods (add/remove signer, change requirement) are called by
doing a multisig transaction invoking the desired method on the contract itself. This means the 'signature
threshold' logic only needs to be implemented once, in one place.

The [init actor](#init-actor) is used to create new instances of the multisig.

```sh
type MultisigActorState struct {
    signers [Address]
    required UInt
    nextTxId UInt
    transactions {UInt:Transaction}
}

type Transaction struct {
    txID UInt
    to Address
    value TokenAmount
    method &ActorMethod
    approved [Address]
    completed Bool
    canceled Bool
    retcode UInt
}
```

#### Methods

| Name | Method ID |
|--------|-------------|
| `MultisigConstructor` | 1 |
| `Propose` | 2 |
| `Approve` | 3 |
| `Cancel` | 4 |
| `ClearCompleted` | 5 |
| `AddSigner` | 6 |
| `RemoveSigner` | 7 |
| `SwapSigner` | 8 |
| `ChangeRequirement` | 9 |


#### `Constructor`

This method sets up the initial state for the multisig account

**Parameters**

```sh
type MultisigConstructor struct {
    ## The addresses that will be the signatories of this wallet.
    signers [Address]
    ## The number of signatories required to perform a transaction.
    required UInt
} representation tuple
```

**Algorithm**

```go
func Multisig(signers []Address, required UInt) {
	self.Signers = signers
	self.Required = required
}
```

#### `Propose`

Propose is used to propose a new transaction to be sent by this multisig. The proposer must be a signer, and the proposal also serves as implicit approval from the proposer. If only a single signature is required, then the transaction is executed immediately.

**Parameters**


```sh
type Propose struct {
    ## The address of the target of the proposed transaction.
    to Address
    ## The amount of funds to send with the proposed transaction.
    value TokenAmount
    ## The method and parameters that will be invoked on the proposed transactions target.
    method &ActorMethod
} representation tuple
```

**Algorithm**

```go
func Propose(to Address, value TokenAmount, method String, params Bytes) UInt {
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
		tx.RetCode = vm.Send(tx.To, tx.Value, tx.Method, tx.Params)
		tx.Complete = true
	}

	return txid
}
```

#### `Approve`

Approve is called by a signer to approve a given transaction. If their approval pushes the approvals for this transaction over the threshold, the transaction is executed.

**Parameters**

```sh
type Approve struct {
    ## The ID of the transaction to approve.
    txid UInt
} representation tuple
```

**Algorithm**

```go
func Approve(txid UInt) {
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
		tx.RetCode = Send(tx.To, tx.Value, tx.Method, tx.Params)
		tx.Complete = true
	}
}
```

#### `Cancel`

**Parameters**

```sh
type Cancel struct {
    txid UInt
} representation tuple
```

**Algorithm**

```go
func Cancel(txid UInt) {
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

#### `ClearCompleted`

**Parameters**

```sh
type ClearCompleted struct {
} representation tuple
```

**Algorithm**

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

#### `AddSigner`

**Parameters**

```sh
type AddSigner struct {
    signer Address
    increaseReq bool
} representation tuple
```

**Algorithm**

```go
func AddSigner(signer Address, increaseReq bool) {
	if msg.From != self.Address {
		Fatal("add signer must be called by wallet itself")
	}
	if self.isSigner(signer) {
		Fatal("new address is already a signer")
	}
	if increaseReq {
		self.Required = self.Required + 1
	}

	self.Signers.Append(signer)
}
```

#### `RemoveSigner`

**Parameters**

```sh
type RemoveSigner struct {
    signer Address
    decreaseReq bool
} representation tuple
```

**Algorithm**

```go
func RemoveSigner(signer Address, decreaseReq bool) {
	if msg.From != self.Address {
		Fatal("remove signer must be called by wallet itself")
	}
	if !self.isSigner(signer) {
		Fatal("given address was not a signer")
	}
	if decreaseReq || len(self.Signers)-1 < self.Required {
		// Reduce Required outherwise the wallet is locked out
		self.Required = self.Required - 1
	}

	self.Signers.Remove(signer)
}
```

#### `SwapSigner`

**Parameters**

```sh
type SwapSigner struct {
    old Address
    new Address
} representation tuple
```

**Algorithm**

```go
func SwapSigner(old Address, new Address) {
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

#### `ChangeRequirement`

**Parameters**

```sh
type ChangeRequirement struct {
    requirement UInt
} representation tuple
```

**Algorithm**

```go
func ChangeRequirement(req UInt) {
	if msg.From != self.Address {
		Fatal("change requirement must be called by wallet itself")
	}
	if req < 1 {
		Fatal("requirement must be at least 1")
	}
	if req > len(self.Signers) {
		Fatal("requirement must be less than number of signers")
	}

	self.Required = req
}
```

## Helper Methods

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
```

```go
func getTransaction(txid UInt) Transaction {
	tx, ok := self.Transactions[txid]
	if !ok {
		Fatal("no such transaction")
	}

	return tx
}
```

```go
func AggregateBitfields(faults []FaultSet) Bitfield {
	var out Bitfield
	for _, f := range faults {
		out = out.Union(f.bitField)
	}
	return out
}
```

```go
func BurnFunds(amt TokenAmount) {
	TransferFunds(BurntFundsAddress, amt)
}
```
