# Storage Market Actor (DEPRECATED)

* **Code Cid**: `<codec:raw><mhType:identity><"smarket">`

The storage market actor is the central point for the Filecoin storage market. It is responsible for registering new miners to the system, and maintaining the power table. The Filecoin storage market is a singleton that lives at a specific well-known address.

```sh
type StorageMarketActorState struct {
    miners {Address:Null}<Hamt>
    totalStorage BytesAmount
}
```

## Methods

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

## `Constructor`

**Parameters**

```sh
type StorageMarketConstructor struct {}
```

**Algorithm**

## `CreateStorageMiner`

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

## `SlashConsensusFault`

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

## `UpdateStorage`

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

## `GetTotalStorage`

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

## `PowerLookup`

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

## `IsMiner`

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

## `StorageCollateralForSize`


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
