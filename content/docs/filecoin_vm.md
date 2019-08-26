


## Filecoin VM

The majority of Filecoin's user facing functionality (payments, storage market, power table, etc) is managed through the Filecoin State Machine. The network generates a series of blocks, and agrees which 'chain' of blocks is the correct one. Each block contains a series of state transitions called `messages`, and a checkpoint of the current `global state` after the application of those `messages`.

The `global state` here consists of a set of `actors`, each with their own private `state`.

An `actor` is the Filecoin equivalent of Ethereum's smart contracts, it is essentially an 'object' in the filecoin network with state and a set of methods that can be used to interact with it. Every actor has a Filecoin balance attributed to it, a `state` pointer, a `code` CID which tells the system what type of actor it is, and a `nonce` which tracks the number of messages sent by this actor. (TODO: the nonce is really only needed for external user interface actors, AKA `account actors`. Maybe we should find a way to clean that up?)

### Method Invocation

There are two routes to calling a method on an `actor`.

First, to call a method as an external participant of the system (aka, a normal user with Filecoin) you must send a signed `message` to the network, and pay a fee to the miner that includes your `message`.  The signature on the message must match the key associated with an account with sufficient Filecoin to pay for the messages execution. The fee here is equivalent to transaction fees in Bitcoin and Ethereum, where it is proportional to the work that is done to process the message (Bitcoin prices messages per byte, Ethereum uses the concept of 'gas'. We also use 'gas').

Second, an `actor` may call a method on another actor during the invocation of one of its methods.  However, the only time this may happen is as a result of some actor being invoked by an external users message (note: an actor called by a user may call another actor that then calls another actor, as many layers deep as the execution can afford to run for).

### Sending Funds

As all messages carry a method ID, the method ID '0' is reserved for simple
transfers of funds. Funds specified by the value field are always transferred,
but specifying a method ID of '0' ensures that no other side effects occur.

### State Representation

The `global state` is modeled as a map of actor `ID`s to actor structs. This map is implemented by an ipld HAMT (TODO: link to spec for our HAMT) with the 'key' being the serialized ID address (every actor has an ID address that can be looked up via the `InitActor`), and the value is an [`Actor`](data-structures.md#actor) object with the actors information. Within each `Actor` object is a field called `state` that is an ipld pointer to a graph that can be entirely defined by the actor.

### Actor Creation

There are two mechanisms by which an actor can be created. By explicitly invoking `exec` on the `Init` actor, and by sending a message to a `Public Key` typed `Address`.

Calling `exec` to create an actor should generate an Actor address, and register it in the state tree (see [Init Actor](actors.md#init-actor) for more details).

Sending a message to a non-existant account via a public key address causes the creation of an account actor for that address. The `To` address should be placed into the actor storage for later use in validating messages sent from this actor.

This second route for creating an actor is allowed to avoid the necessity of an explicit 'register account' step for creating new accounts.

### Execution (Calling a method on an Actor)

Message execution currently relies entirely on 'built-in' code, with a common external interface. The method and actor to call it on are specified in the `Method` and `To` fields of a message, respectively. Method parameters are encoded and put into the `Params` field of a message. The encoding is technically actor dependent, but for all built-in Filecoin actors it is the dag-cbor ipld encoding of the parameters struct for each method defined in [the actors doc](actors.md).

These functions are given, as input, an `ExecutionContext` containing useful information for their execution.

```go
type VMContext interface {
	// Message is the message that kicked off the current invocation
	Message() Message

	// Storage provides access to the VM storage layer
	Storage() Storage

	// Origin is the address of the account that initiated the top level invocation
	Origin() Address

	// Send allows the current execution context to invoke methods on other actors in the system
	Send(to Address, method string, value AttoFIL, params []interface{}) ([][]byte, uint8, error)

	// BlockHeight returns the height of the block this message was added to the chain in
	BlockHeight() BlockHeight
}
```

If the execution completes successfully, changes to the state tree are saved. Otherwise, the message is marked as failed, and any state changes are reverted.

```go
func ApplyMessage(st StateTree, msg Message) MessageReceipt {
	st.Snapshot()
	fromActor, found := st.GetActor(msg.From)
	if !found {
		Fatal("no such from actor")
	}

	totalCost := msg.Value + (msg.GasLimit * msg.GasPrice)
	if fromActor.Balance < totalCost {
		Fatal("not enough funds")
	}

	if msg.Nonce() != fromActor.Nonce+1 {
		Fatal("invalid nonce")
	}

	toActor, found := st.GetActor(msg.To)
	if !found {
		toActor = TryCreateAccountActor(st, msg.To)
	}

	st.DeductFunds(msg.From, totalCost)
	st.DepositFunds(msg.To, msg.Value)

	vmctx := makeVMContext(st, msg)

	if msg.Method != 0 {
		ret, errcode := toActor.Invoke(vmctx, msg.Method, msg.Params)
		if errcode != 0 {
			// revert all state changes since snapshot
			st.Revert()
			st.DeductFunds(msg.From, vmctx.GasUsed()*msg.GasPrice)
		} else {
			// refund unused gas
			st.DepositFunds(msg.From, (msg.GasLimit-vmctx.GasUsed())*msg.GasPrice)
		}
	}

	// reward miner gas fees
	st.DepositFunds(BlockMiner, msg.GasPrice*vmctx.GasUsed())

	return MessageReceipt{
		ExitCode: errcode,
		Return:   ret,
		GasUsed:  vmctx.GasUsed(),
	}
}

func TryCreateAccountActor(st StateTree, addr Address) Actor {
	switch addr.Type() {
	case BLS:
		return NewBLSAccountActor(addr)
	case Secp256k1:
		return NewSecp256k1AccountActor(addr)
	case ID:
		Fatal("no actor with given ID")
	case Actor:
		Fatal("no such actor")
	}
}
```

#### Receipts

Every message execution generates a [receipt](data-structures.md#message-receipt). These receipts contain the encoded return value of the method invocation, and an exit code.

#### Storage

Actors are given acess to a `Storage` interface to fulfil their need for persistent storage. The `Storage` interface describes a content addressed block storage system (`Put` and `Get`) and a pointer into it (`Head` and `Commit`) that points to the actor's current state.

```go
type Storage interface {
	// Put writes the given object to the storage staging area and returns its CID
	Put(interface{}) (Cid, error)

	// Get fetches the given object from storage (either staging, or local) and returns
	// the serialized data.
	Get(Cid) ([]byte, error)

	// Commit updates the actual stored state for the actor. This is a compare and swap
	// operation, and will fail if 'old' is not equal to the current return value of `Head`.
	// This functionality is used to prevent issues with re-entrancy
	Commit(old Cid, new Cid) error

	// Head returns the CID of the current actor state
	Head() Cid
}
```

Actors can store state as a single block or implement any persistent
data structure that can be built upon a content addressed block store.
Implementations may provide data structure implementations to simplify
development. The current interface only supports CBOR-IPLD, but this
should soon expand to allow other types of IPLD data structures (as long
as the system has resolvers for them).

The current state of a given actor can be accessed first by calling `Head` to retrieve the CID of the root of the actors state, then by using `Get` to retrieve the actual object being referenced.

To store data, `Put` is used. Any number of objects may be `Put`, but only the object whose CID is committed, or objects that are linked to in some way by the committed object will be kept. All other objects are dropped after the method invocation returns. Objects stored via `Put` are first marshaled to CBOR-IPLD, and then stored, the returned CID is a 32 byte sha2-256 CBOR-IPLD content identifier.

### Burning Funds

In the case that an actor needs to provably burn funds, the funds should be transferred to the 'Burnt Funds Actor' (ID 99).  





## Filecoin State Machine Actors

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
| `PaymentVerifyInclusion` | 17 |
| `PaymentVerifySector` | 18 |
| `AddFaults` | 19 |

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

#### `SlashStorageFault`

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

#### `PaymentVerifyInclusion`

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


#### `PaymentVerifySector`

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

#### `AddFaults`

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
    initialBalance UInt
    startingBlock UInt
    unlockDuration UInt
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
    ## Unlock time (in blocks) of initial filecoin balance of this wallet. Unlocking is linear.
    unlockDuration UInt
} representation tuple
```

**Algorithm**

```go
func Multisig(signers []Address, required UInt, unlockDuration UInt) {
	self.Signers = signers
	self.Required = required
	self.initialBalance = msg.Value
	self.unlockDuration = unlockDuration
	self.startingBlock = VM.CurrentBlockHeight()
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
		if !self.canSpend(tx.value) {
			Fatal("transaction amount exceeds available")
		}
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
		if !self.canSpend(tx.Value) {
			Fatal("transaction amount exceeds available")
		}
		tx.RetCode = vm.Send(tx.To, tx.Value, tx.Method, tx.Params)
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

```go
func canSpend(amt TokenAmount) bool {
	if self.unlockDuration == 0 {
		return true
	}
	var MinAllowableBalance = (self.initialBalance / self.unlockDuration) * (VM.CurrentBlockHeight() - self.startingBlock)
	return MinAllowableBalance >= (vm.MyBalance() - amt)
}
```



# Faults

A fault is what happens when partcipants in the protocol are behaving incorrectly and that behavior needs to be punished. There are a number of possible faults in the Filecoin protocol, their details are all recorded below.

## Fault List

### Consensus Faults

- **Duplicate Block Submission Slashing:**
  - **Condition:** If any miner posts two blocks satisfying the slashing conditions defined in [Expected Consensus](expected-consensus.md).
  - **Reporting:** Anyone may call `SlashConsensusFault` and pass in the two offending block headers.
  - **Check:** The chain checks that both blocks are valid, correctly signed by the same miner, and satisfy the consensus slashing conditions.
  - **Penalization:** All of the miner's pledge collateral and all of their power is irrevocably slashed. This miner can never again produce blocks, even if they attempt to repost their collateral.

### Market Faults

- **Late submission penalty:**
  - **Condition**: If the miner posts their PoSt after the proving period ends, but before the generation attack threshold.
  - **Reporting:** The miner submits their PoSt as usual, but includes the late submission fee.
  - **Check:** The chain checks first that the submission is within the `generation attack threshold`, and then checks that the fee provided matches the required fee for how many blocks late the submission is.
  - **Penalization:** The miner is penalized proportionally to the delay. Penalizations are enforced by a standard PoSt submission.
    - *Economic penalization*: To determine the penalty amount, `ComputeLateFee(minerPower, numLate)` is called.
    - *Power penalization*: The miners' power is not reduced. Note that the current view of the power table is computed with the lookback parameter.
      - *Why are we accounting the power table with a lookback parameter ?* If we do not use the lookback parameter then, we need to penalize late miners for the duration that they are late. This is tricky to do efficiently. For xample, if miners A, B and C each have 1/3 of the networks power, and C is late in submitting their proofs, then for that duration, A and B should each have effectively half of the networks power (and a 50% chance each of winning the block).
  - TODO: write on the spec exact parameters for PoSt Deadline and Gen Attack threshold
- **Unreported storage fault slashing:**
  - **Condition:** If the miner does not submit their PoSt by the `generation attack threshold`.
  - **Reporting:** The miner can be slashed by anyone else in the network who calls `SlashStorageFaults`. We expect miners to report these faults.
    - Future design note: moving forward, we should either compensate the caller, or require this
    - Note: we could *require* the method be called, as part of the consensus rules (this gets complicated though). In this case, there is a DoS attack where if I make a large number of miners each with a single sector, and fail them all at the same time, the next block miner will be forced to do a very large amount of work. This would either need an extended 'gas limit', or some other method to avoid too long validation times.
  - **Check:** The chain checks that the miners last PoSt submission was before the start of their current proving period, and that the current block is after the generation attack threshold for their current proving period.
  - **Penalization:** Penalizations are enforced by `SlashStorageFault` on the `storage market` actor.
    - *Economic Penalization*: Miner loses all collateral.
    - *Power Penalization*: Miner loses all power.
    - Note: If a miner is in this state, where they have failed to submit a PoST, any block they attempt to mine will be invalid, even if the election function selects them. (the election function should probably be made to never select them)
    - Future design note: There is a way to tolerate Internet connection faults. A miner runs an Emergency PoSt which does not take challenges from the chain, if the miner gets reconnected before the VDF attack time (based on Amax), then, they can submit the Emergency PoSt and get pay a late penalization fee.
- **Reported storage fault penalty:**
  - **Condition:** The miner submits their PoSt with a non-empty set of 'missing sectors'.
  - **Reporting:** The miner can specify some sectors that they failed to prove during the proving period.
    - Note: These faults are output by the `ProveStorage` routine, and are posted on-chain when posting the proof. This occurs when the miner (for example) has a disk failure, or other local data corruption.
  - **Check:** The chain checks that the proof verifies with the missing sectors.
  - **Penalization:** The miner is penalized for collateral and power proportional to the number of missing sectors. The sectors are also removed from the miners proving set.
    - TODO: should the collateral lost here be proportional to the remaining time?
    - TODO(nicola): check if the time between posting two proofs allows for a generation attack if it does not then we might reconsider the sector not being lost
  - Note: if a sector is missed here, and they are recovered after the fact, the miner could simple 're-commit' the sector. They still have to pay the collateral, but the data can be quickly re-introduced into the system to avoid clients calling them out for breach of contract (this would only work because the sector commD/commR is the same)
  - Note: In the case where a miner is temporarily unable to prove some of their data, they can simply wait for the temporary unavailability to recover, and then continue proving, submitting the proofs a bit late if necessary (paying appropriate fees, as described above).
- **Breach of contract dispute:**
  - **Condition:** A client who has stored data with a miner, and the miner removes the sector containing that data before the end of the agreed upon time period.
  - **Reporting:** The client invokes `ArbitrateDeal` on the offending miner actor with a signed deal from that miner for the storage in question. Note: the reporting must happen within one proving period of the miner removing the storage erroneously.
  - **Check:** The chain checks that the deal was correctly signed by the miner in question, that the deal has not yet expired, and that the sector referenced by the deal is no longer in the miners proving set.
  - **Penalization:** The miner is penalized an amount proportional to the incorrectly removed sector. This penalty is taken from their pledged collateral .
  - Note: This implies that miners cannot re-seal data into different sectors. We could come up with a protocol where the client gives the miner explicit consent to re-seal, but that is more complicated and can be done later.
