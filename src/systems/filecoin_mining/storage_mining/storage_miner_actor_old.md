# Storage Miner Actor (DEPRECATED)

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

    ## Faulty sectors reported since last SubmitPost, up to the current proving period's challenge time.
    currentFaultSet BitField

    ## Faults submitted after the current proving period's challenge time, but before the PoSt for that period
    ## is submitted. These become the currentFaultSet when a PoSt is submitted.
    nextFaultSet BitField

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
```

```sh
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

## Methods

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
| `PaymentVerifyInclusion` | 17 |
| `PaymentVerifySector` | 18 |
| `AddFaults` | 19 |

## `Constructor`

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

## `CommitSector`

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

## `SubmitPoSt`

**Parameters**

```sh
type SubmitPost struct {
    proofs PoStProof
    doneSet Bitfield
} representation tuple
```

**Algorithm**

```go
func SubmitPost(proofs PoStProof, doneSet Bitfield) {
  if msg.From != self.Worker {
    Fatal("not authorized to submit post for miner")
  }

  feesRequired := 0
    nextProvingPeriodEnd := self.ProvingPeriodEnd + ProvingPeriodDuration(self.SectorSize)

    // TODO: rework fault handling, for now anything later than 2 proving periods is invalid
    if chain.now() > nextProvingPeriodEnd {
        Fatal("PoSt submited too late")
  } else if chain.Now() > self.ProvingPeriodEnd {
    feesRequired += ComputeLateFee(self.power, chain.Now() - self.provingPeriodEnd)
  }

  feesRequired += ComputeTemporarySectorFailureFee(self.sectorSize, self.currentFaultSet)

  if msg.Value < feesRequired {
    Fatal("not enough funds to pay post submission fees")
  }

  // we want to ensure that the miner can submit more fees than required, just in case
  if msg.Value > feesRequired {
    TransferFunds(msg.From, msg.Value-feesRequired)
  }

    var seed
    if chain.Now() < self.ProvingPeriodEnd {
      // good case, submitted in time
      seed = GetRandFromBlock(self.ProvingPeriodEnd - POST_CHALLENGE_TIME)
    } else {
      // bad case, submitted late, need to take new proving period end as reference
      seed = GetRandFromBlock(nextPovingPeriodEnd - POST_CHALLENGE_TIME)
    }

    faultSet := self.currentFaultSet

  if !VerifyPoSt(self.SectorSize, self.provingSet, seed, proof, faultSet) {
    Fatal("proof invalid")
  }

    // The next fault set becomes the current one
    self.currentFaultSet = self.nextFaultSet
    self.nextFaultSet = EmptySectorSet()

    // TODO: penalize for faults

  // Remove doneSet from the current sectors
  self.Sectors.Subtract(doneSet)

  // Update miner power to the amount of data actually proved during the last proving period.
  oldPower := self.Power

  self.Power = (self.ProvingSet.Size() - faultSet.Count()) * self.SectorSize
  StorageMarket.UpdateStorage(self.Power - oldPower)

  self.ProvingSet = self.Sectors

  // Updating proving period given a fixed schedule, independent of late submissions.
  self.ProvingPeriodEnd = nextProvingPeriodEnd

  // update next done set
  self.NextDoneSet = done
  self.ArbitratedDeals.Clear()
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

## `SlashStorageFault`

**Parameters**

```sh
type SlashStorageFault struct {
    miner Address
} representation tuple
```

**Algorithm**

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

## `GetCurrentProvingSet`

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

## `ArbitrateDeal`

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

## `DePledge`

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

## `GetOwner`

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

## `GetWorkerAddr`

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

## `GetPower`

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

## `GetPeerID`

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

## `GetSectorSize`

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

## `UpdatePeerID`

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

## `ChangeWorker`

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

## `IsLate`

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

## `IsSlashed`

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

## `PaymentVerifyInclusion`

Verifies a storage market payment channel voucher's 'Extra' data by validating piece inclusion proof.

**Parameters**

```sh
type PaymentVerify struct {
    Extra Bytes
    Proof Bytes
} representation tuple

type PieceInclusionVoucherData struct {
    CommP Bytes
    PieceSize BigInt
} representation tuple

type InclusionProof struct {
    Sector BigInt // for CommD, also verifies the sector is in sector set
    Proof  Bytes
} representation tuple
```

**Algorithm**

```go
func PaymentVerifyInclusion(extra PieceInclusionVoucherData, proof InclusionProof) {
  has, commD := self.GetSector(proof.Sector)
  if !has {
    Fatal("miner does not have required sector")
  }

  return ValidatePIP(self.SectorSize, extra.PieceSize, extra.CommP, commD, proof.Proof)
}
```


## `PaymentVerifySector`

Verifies a storage market payment channel voucher's 'Extra' data by checking for presence of a specified sector in miner's sector set.

Miners should prefer payment vouchers with this method used for validation over `PaymentVerifyInclusion`, because posting them to the chain will be much cheaper.

Clients should only create such vouchers after verifying that miners have related sectors in their sector set, and after checking piece inclusion proof.

Miners can incentivize clients to produce such vouchers by applying small 'discount' to amount of token clients have to pay.

**Parameters**

```sh
type PaymentVerify struct {
    Extra Bytes
    Proof Bytes
} representation tuple
```

**Algorithm**

```go
func PaymentVerifyInclusion(extra BigInt, proof Bytes) {
  if len(proof) > 0 {
    Fatal("unexpected proof bytes")
  }

  return self.HasSector(extra)
}
```

## `AddFaults`

**Parameters**

```sh
type AddFaults struct {
    faults FaultSet
} representation tuple
```

**Algorithm**

```go
func AddFaults(faults FaultSet) {
    challengeBlockHeight := self.ProvingPeriodEnd - POST_CHALLENGE_TIME

    if VM.CurrentBlockHeight() < challengeBlockHeight {
        // Up to the challenge time new faults can be added.
        self.currentFaultSet = Merge(self.currentFaultSet, faults)
    } else {
        // After that they are only accounted for in the next proving period
        self.nextFaultSet = Merge(self.nextFaultSet, faults)
    }
}
```

/*
type CollateralVault struct {
    Pledged(Collateral) TokenAmount
    Pledge(Collateral, TokenAmount)
    DePledge(Collateral, TokenAmount)

    pledgedStorageCollateral UVarint
    pledgedConsensusCollateral UVarint
}

type Collateral union {
    | StorageDealCollateral
    | ConsensusCollateral
}

type TokenAmount UVarint # What are the units? Attofil?

type MinedSector {
    SectorID           UInt
    CommR              sector.Commitment
    FaultStatus(Epoch) Fault
}

type Fault union {
    | None
    | GracePeriod
    | Fault
}

<!-- type StorageMinerActor struct {
  // Amount of power this miner has.
  power UInt

  provingPeriodEnd Epoch


  // Collateral that is waiting to be withdrawn.
  dePledgeCollateral TokenAmount

  // Time at which the depledged collateral may be withdrawn.
  dePledgeTime Epoch

  // All sectors this miner has committed.
  sectors &SectorSet

  // Sectors this miner is currently mining. It is only updated
  // when a PoSt is submitted (not as each new sector commitment is added).
  provingSet &SectorSet

  // Faulty sectors reported since last SubmitPost, up to the current proving period's challenge time.
  currentFaultSet FaultSet

  // Faults submitted after the current proving period's challenge time, but before the PoSt for that period
  // is submitted. These become the currentFaultSet when a PoSt is submitted.
  nextFaultSet FaultSet

  // Sectors reported during the last PoSt submission as being 'done'. The collateral
  // for them is still being held until the next PoSt submission in case early sector
  // removal penalization is needed.
  nextDoneSet DoneSet

  // List of sectors that this miner was slashed for.
  slashedSet optional &SectorSet

  // Deals this miner has been slashed for since the last post submission.
  arbitratedDeals {Cid:Null} // TODO

  // The height at which this miner was slashed at.
  slashedAt optional Epoch

  // The amount of storage collateral that is owed to clients, and cannot be used for collateral anymore.
  owedStorageCollateral TokenAmount

  // Internal methods
  verifySeal(sectorID SectorID, comm SealCommitment, proof SealProof)
  verifyPoSt(proofs base.PoStProof, doneSet Bitfield)

  // Getters
  GetOwner() address.Address
  GetWorkerAddr() address.Address
  GetPower() BytesAmount
  GetPeerID() PeerID
  GetSectorSize() BytesAmount
  GetCurrentProvingSet() BitField

  // SubmitPost verifies the PoSt
  SubmitPost(proofs base.PoStProof, doneSet DoneSet) bool // TODO: rename to ProvePower?
  DePledge(amt TokenAmount)

  AddCollateral()

  AbitrateDeal (deal Deal)
  SlashStorageFault() // TODO maybe add CheckStorageFault?
  UpdateFaults(faults FaultSet) // TODO rename into ReportFaults
  IsLate() (bool)
  IsSlashed() (bool)

  UpdatePeerID(pid PeerID)
  ChangeWorker(addr address.Address)

  PaymentVerifyInclusion(extra PieceInclusionVoucherData, proof InclusionProof) (bool)
  PaymentVerifyInclusion(extra BigInt, proof Bytes) (bool)
} -->
