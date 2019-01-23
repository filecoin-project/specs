# Filecoin State Machine Actors

Any implementations of the Filecoin actors must be exactly byte for byte compatible with the go-filecoin actor implementations. The pseudocode below tries to capture most of the important logic, but capturing all the detail would require embedding exactly the code from go-filecoin, so for now, its simply informative pseudocode.

## Storage Market Actor

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
    
    return VM.CreateActor(MinerActor, msg.Value, pubkey, pledge, pid)
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
    
    // Burn all of the miners collateral
    miner.BurnCollateral()
    
    // TODO: Some of the slashed collateral should be paid to the slasher
    
    // Remove the miner from the list of network miners
    self.Miners.Remove(miner)
    self.UpdateStorage(-1 * miner.Power)

    // Now delete the miner (maybe this is a bit harsh, but i'm okay with it for now)
    miner.SelfDestruct()
}
```


### UpdateStorage

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

### Constructor

Along with the call, the actor must be created with exactly enough filecoin for the collateral necessary for the pledge.
```go
func StorageMinerActor(pubkey PublicKey, pledge BytesAmount, pid PeerID) {
    if msg.Value < CollateralForPledgeSize(pledge) {
        Fatal("not enough collateral given")
    }
    
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
func CommitSector(commD, commR []byte, proof *SealProof) SectorID {
    if !miner.ValidatePoRep(commD, commR, miner.PublicKey, proof) {
        Fatal("bad proof!")
    }
    
    // make sure the miner isnt trying to submit a pre-existing sector
    if !miner.EnsureSectorIsUnique(commR) {
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
- proof PoStProof
- faults []FailureSet
- recovered SectorSet
- done SectorSet

Return: None

```go
func SubmitPost(proof PoSt, faults []FaultSet, recovered BitField, done BitField) {
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
    
    if !CheckPostProof(proof, faults) {
        Fatal("proof invalid")
    }
    
    permLostSet = AggregateBitfields(faults).Subtract(recovered)
    
    // adjust collateral for 'done' sectors
    miner.ActiveCollateral -= CollateralForSectors(miner.NextDoneSet)
    miner.Collateral += CollateralForSectors(miner.NextDoneSet)
    
    // penalize collateral for lost sectors
    miner.Collateral -= CollateralForSectors(permLostSet)
    miner.ActiveCollateral -= CollateralForSectors(permLostSet)
    
    // burn funds for fees and collateral penalization
    BurnFunds(miner, CollateralForSectors(permLostSet) + feesRequired)
    
    // update sector sets and proving set
    miner.Sectors.Subtract(miner.NextDoneSet)
    miner.Sectors.Subtract(permLostSet)
    
    // update miner power to the amount of data actually proved during
    // the last proving period.
    miner.Power = SizeOf(Filter(miner.ProvingSet, faults))
    
    miner.ProvingSet = miner.Sectors
    
    // NEEDS REVIEW: early submission of PoSts may give the miner extra time for
    // their next PoSt, which could compound. Does the beacon reseeding for Posts
    // address this well enough?
    miner.ProvingPeriodEnd = miner.ProvingPeriodEnd + ProvingPeriodDuration
    
    // update next done set
    miner.NextDoneSet = done
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
func SlashStorageFault(miner Address) {
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
    
    coll := CollateralForDeal(d)
    
    self.BurnFunds(coll)
    
    self.NextDoneSet.Remove(d.PieceCommitment.Sector)
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

## ABI Types and Encodings

Method parameters are encoded and put into the 'Params' field of a message. The encoding is simply a cbor array of each of the types individually encoded. The individual encodings for each type are as follows.

#### `PublicKey`

The public key type is simply an array of bytes. (TODO: discuss specific encoding of key types, for now just calling it bytes is sufficient)

#### `BytesAmount`
BytesAmount is just a re-typed Integer.

#### `PeerID`
PeerID is just the serialized bytes of a libp2p peer ID.

Spec incomplete, take a look at this PR: https://github.com/libp2p/specs/pull/100

#### `Integer`

Integers are encoded as LEB128 signed integers.

#### `BitField`

Bitfields are a set of bits. Encoding still TBD, but it needs to be very compact. We can assume that most often, ranges of bits will be set, or not set, and use that to our advantage here. Some form of run length encoding may work well.

#### `SectorSet`

TODO

#### `FaultSet`

A fault set is a BitField and a block height, encoding TBD.

#### `BlockHeader`

BlockHeader is a serialized `Block`, as described in [the data structures document](data-structures.md#block).

#### `SealProof`

SealProof is just an array of bytes.

#### `TokenAmount`

TokenAmount is just a re-typed Integer.