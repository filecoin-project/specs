---
title: VM Interpreter
---

# VM Interpreter


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



# Pledge Collateral

Filecoin includes a concept of "Pledge Collateral", which is FIL collateral that storage miners must lock up when participating as miners.

Pledge collateral serves several functions in Filecoin. It:

- makes it possible to slash misbehaving or slow miners
- ensures that miners have skin in the game (for the Filecoin network as a whole)
- increases the cost of launching a 51% attack


## Computing Pledge Collateral

The total pledge collateral across all miners is a fixed proportion of available FIL.
Available FIL is computed as the total amount of FIL that has been mined, plus the total amount of FIL that's been vested, minued the amount of FIL which has been burned.

```go
availableFil := minedFil + vestedFil - burnedFil
```

Pledge collateral is subdivided into two kinds: power collateral and per-capita collateral.
Power collateral is split across miners according to their share of the total network power, and per-capita collateral is split across miners evenly.
Two parameters, `POWER_COLLATERAL_PROPORTION` and `PER_CAPITA_COLLATERAL_PROPORTION`, relate the total amount of collateral to the `availableFil`.

```go
totalPowerCollateral := availableFil * POWER_COLLATERAL_PROPORTION
totalPerCapitaCollateral := availableFil * PER_CAPITA_COLLATERAL_PROPORTION
totalPledgeCollateral := totalPowerCollateral + totalPerCapitaCollateral
```

Power-based collateral ensures that miners' collateral is proportional to their economic size and to their expected rewards.
The presence of per-capital collateral acts as a deterrent against Sibyl attacks.
We intend for the `POWER_COLLATERAL_PROPORTION` to be several times larger than the `PER_CAPITA_COLLATERAL_PROPORTION`.

To calculate any particular miner's collateral requirements, we need to know the miner's power, the total network power, and the total number of miners in the network.

```go
minerPowerCollateral := totalPowerCollateral * minerPower / totalNetworkPower
minerPerCapitaCollateral := totalPerCapitaCollateral / numMiners
```

Putting all these variables together, we have each miner's individual collateral requirement:
```go
minerPlegeCollateral := availableFil * ( POWER_COLLATERAL_PROPORTION * minerPower / totalNetworkPower PER_CAPITA_COLLATERAL_PROPORTION / numMiners)
```

## Dealing with Undercollateralization

In the course of normal events, miners may become undercollateralized.

They cannot directly undercollateralized themselves by adding more power, as commitSector will fail if they do not have sufficient collateral to cover their power requirements.
However, their collateral requirement could increase due to growth in availableFil, a reduction in the total network power, or a reduction in the total number of miners.
In such cases, the miner may continue to submit PoSts and mine blocks. When they win blocks, their block rewards will be garnished while they remain undercollateralized.

## Parameter Choices

We provisionally propose the following two parameters choices:

```go
POWER_COLLATERAL_PROPORTION := 0.2
PER_CAPITA_COLLATERAL_PROPORTION := 0.05
```

These are subject to change before launch.
