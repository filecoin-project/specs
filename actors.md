# Filecoin State Machine Actors

Any implementations of the Filecoin actors must be exactly byte for byte compatible with the go-filecoin actor implementations. The pseudocode below tries to capture the important logic, but capturing all the detail would require embedding exactly the code from go-filecoin, so for now, its simply informative pseudocode. The algorithms below are correct, and all implementations much match it (including go-filecoin), but details omitted from here should be looked for in the go-filecoin code.

This spec decsribes a set of actors that operate within the [Filecoin State Machine](state-machine.md). All types are defined in [the basic type encoding spec](data-structures.md#basic-type-encodings).


- [Storage Market Actor](#storage-market-actor)
- [Storage Miner Actor](#storage-miner-actor)

## Init Actor
The init actor is responsible for creating new actors on the filecoin network. This is a built-in actor and cannot be replicated. In the future, this actor will be responsible for loading new code into the system (for user programmable actors).

```go
type InitActor struct {
    // Mapping from Address to ID, for lookups.
		AddressMap map[Address]BigInt
}
```

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
    actorID = len(self.AddressMap)

    // Make sure that only the actors defined in the spec can be launched.
    if !IsBuiltinActor(code) {
        Fatal("cannot launch actor instance that is not a builtin actor")
    }

    // Ensure that singeltons can be only launched once.
    // TODO: do we want to enforce this? If so how should actors be marked as such?
    if IsSingletonActor(code) && self.hasInstance(code) {
      Fatal("cannot launch singelton, which is already launched")
    }

  	// This call will insert the actor into the global state tree, under the provided ID.
    addr := VM.CreateNewActor(actorID, code, params)

    // Store the mapping of address to actor ID.
    self.AddressMap[addr] = actorID

    return addr
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





## Storage Market Actor

The storage market actor is the central point for the Filecoin storage market. It is responsible for registering new miners to the system, and maintaining the power table. The Filecoin storage market is a singleton that lives at a specific well-known address.

```go
type StorageMarketActor struct {
    Miners AddressSet

    TotalStorage Integer
}
```



#### CreateStorageMiner

Parameters:

- pubkey PublicKey

- pledge BytesAmount

- pid PeerID

Return: Address

```go
func CreateStorageMiner(pubkey PublicKey, pledge BytesAmount, pid PeerID) Address {
    if pledge < MinimumPledge {
        Fatal("Pledge too low")
    }

    if msg.Value < MinimumCollateral(pledge) {
        Fatal("not enough funds to cover required collateral")
    }

    newminer := InitActor.Exec(MinerActorCodeCid, EncodeParams(pubkey, pledge, pid))

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
func SlashConsensusFault(block1, block2 BlockHeader) {
	if block1.Height != block2.Height {
        Fatal("cannot slash miner for blocks of differing heights")
    }

    if !ValidateSignature(block1.Signature) || !ValidateSignature(block2.Signature) {
        Fatal("Invalid blocks")
    }

    if AuthorOf(block1) != AuthorOf(block2) {
        Fatal("blocks must be from the same miner")
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

- delta Integer

Return: None

```go
func UpdateStorage(delta Integer) {
    if !self.Miners.Has(msg.From) {
        Fatal("update storage must only be called by a miner actor")
    }

    self.TotalStorage += delta
}
```

### GetTotalStorage

Parameters: None

Return: Integer

```go
func GetTotalStorage() Integer {
    return self.TotalStorage
}
```

## Storage Miner Actor

```go
type StorageMiner struct {
    // Owner is the address of the account that owns this miner
    Owner Address

    // Worker is the address of the worker account for this miner
    Worker Address

    // PeerID is the libp2p peer identity that should be used to connect
    // to this miner
    PeerID peer.ID

    // PublicKey is the public portion of the key that the miner will use to sign blocks
    PublicKey PublicKey

    // PledgeBytes is the amount of space being offered by this miner to the network
    PledgeBytes BytesAmount

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
    LockedStorage BytesAmount

    // Power is the amount of power this miner has
    Power BytesAmount
}
```

### Constructor

Along with the call, the actor must be created with exactly enough filecoin for the collateral necessary for the pledge.
```go
func StorageMinerActor(pubkey PublicKey, pledge BytesAmount, pid PeerID) {
    if msg.Value < CollateralForPledgeSize(pledge) {
        Fatal("not enough collateral given")
    }

    self.Owner = message.From
    self.PublicKey = pubkey
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
        Price: price,
        Expiry: CurrentBlockHeight + ttl,
        ID: askid,
    })

    return askid
}
```

Note: this may be moved off chain soon, don't worry about testing it too heavily.



### CommitSector

Parameters:
- commD []byte
- commR []byte
- commRStar []byte
- proof SealProof

Return: SectorID


```go
// NotYetSpeced: ValidatePoRep, EnsureSectorIsUnique, CollateralForSector, Commitment
func CommitSector(comm Commitment, proof *SealProof) SectorID {
    if !miner.ValidatePoRep(comm, miner.PublicKey, proof) {
        Fatal("bad proof!")
    }

    // make sure the miner isnt trying to submit a pre-existing sector
    if !miner.EnsureSectorIsUnique(comm) {
        Fatal("sector already committed!")
    }

    // make sure the miner has enough collateral to add more storage
    // currently, all sectors are the same size, and require the same collateral
    // in the future, we may have differently sized sectors and need special handling
    coll = CollateralForSector()

    if coll < miner.Collateral {
        Fatal("not enough collateral")
    }

    miner.Collateral -= coll
    miner.ActiveCollateral += coll

    sectorId = miner.Sectors.Add(commR)
    // TODO: sectors IDs might not be that useful. For now, this should just be the number of
    // the sector within the set of sectors, but that can change as the miner experiences
    // failures.

    // if miner is not mining, start their proving period now
    // Note: As written here, every miners first PoSt will only be over one sector.
    // We could set up a 'grace period' for starting mining that would allow miners
    // to submit several sectors for their first proving period. Alternatively, we
    // could simply make the 'CommitSector' call take multiple sectors at a time.
    if miner.ProvingSet.Size() == 0 {
       miner.ProvingSet = miner.Sectors
       miner.ProvingPeriodEnd = chain.Now() + ProvingPeriodDuration
    }

    return sectorId
}
```

### SubmitPoSt

Parameters:
- proofs []PoStProof
- faults []FailureSet
- recovered SectorSet
- done SectorSet

Return: None

```go
// NotYetSpeced: ValidateFaultSets, GenerationAttackTime, ComputeLateFee
func SubmitPost(proofs []PoStProof, faults []FaultSet, recovered BitField, done BitField) {
    if msg.From != miner.Worker {
        Fatal("not authorized to submit post for miner")
    }

    // ensure the fault sets properly stack, recovered is a subset of the combined
    // fault sets, and that done does not intersect with either, and that all sets
    // only reference sectors that currently exist
    if !miner.ValidateFaultSets(faults, recovered, done) {
        Fatal("fault sets invalid")
    }

    var feesRequired TokenAmount

    if chain.Now() > miner.ProvingPeriodEnd + GenerationAttackTime {
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


    if !CheckPostProofs(proofs, faults) {
        Fatal("proofs invalid")
    }

    permLostSet = AggregateBitfields(faults).Subtract(recovered)

    // adjust collateral for 'done' sectors
    miner.ActiveCollateral -= CollateralForSectors(miner.NextDoneSet)
    miner.Collateral += CollateralForSectors(miner.NextDoneSet)

    // penalize collateral for lost sectors
    miner.ActiveCollateral -= CollateralForSectors(permLostSet)

    // burn funds for fees and collateral penalization
    BurnFunds(miner, CollateralForSectors(permLostSet) + feesRequired)

    // update sector sets and proving set
    miner.Sectors.Subtract(done)
    miner.Sectors.Subtract(permLostSet)

    // update miner power to the amount of data actually proved during
    // the last proving period.
    oldPower := miner.Power
    miner.Power = SizeOf(Filter(miner.ProvingSet, faults))
    StorageMarket.UpdateStorage(miner.Power - oldPower)

    miner.ProvingSet = miner.Sectors

    // NEEDS REVIEW: early submission of PoSts may give the miner extra time for
    // their next PoSt, which could compound. Does the beacon reseeding for Posts
    // address this well enough?
    miner.ProvingPeriodEnd = miner.ProvingPeriodEnd + ProvingPeriodDuration

    // update next done set
    miner.NextDoneSet = done
    miner.ArbitratedDeals.Clear()
}
```

### IncreasePledge

Parameters:

- addspace BytesAmount

Return: None


```go
func IncreasePledge(addspace BytesAmount) {
    // Note: msg.Value is implicitly transferred to the miner actor
    if miner.Collateral + msg.Value < CollateralForPledge(addspace + miner.PledgeBytes) {
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

    if chain.Now() <= miner.ProvingPeriodEnd + GenerationAttackTime {
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
    // TODO: Do both the worker and the owner have the right to call this?
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

Return: Integer

```go
func GetPower() Integer {
    return self.Power
}
```

### GetKey

Parameters: None

Return: PublicKey

```go
func GetKey() PublicKey {
    return self.PublicKey
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

## Payment Channel Broker Actor

TODO
