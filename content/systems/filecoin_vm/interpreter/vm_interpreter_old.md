---
title: VM Interpreter
bookhidden: true
---

# VM Interpreter


## Sending Funds

As all messages carry a method ID, the method ID '0' is reserved for simple
transfers of funds. Funds specified by the value field are always transferred,
but specifying a method ID of '0' ensures that no other side effects occur.

## State Representation

The `global state` is modeled as a map of actor `ID`s to actor structs. This map is implemented by an ipld HAMT (TODO: link to spec for our HAMT) with the 'key' being the serialized ID address (every actor has an ID address that can be looked up via the `InitActor`), and the value is an [`Actor`](actor) object with the actors information. Within each `Actor` object is a field called `state` that is an ipld pointer to a graph that can be entirely defined by the actor.

## Actor Creation

There are two mechanisms by which an actor can be created. By explicitly invoking `exec` on the `Init` actor, and by sending a message to a `Public Key` typed `Address`.

Calling `exec` to create an actor should generate an Actor address, and register it in the state tree (see [Init Actor](sysactors#initactor) for more details).

Sending a message to a non-existant account via a public key address causes the creation of an account actor for that address. The `To` address should be placed into the actor storage for later use in validating messages sent from this actor.

This second route for creating an actor is allowed to avoid the necessity of an explicit 'register account' step for creating new accounts.

## Execution (Calling a method on an Actor)

Message execution currently relies entirely on 'built-in' code, with a common external interface. The method and actor to call it on are specified in the `Method` and `To` fields of a message, respectively. Method parameters are encoded and put into the `Params` field of a message. The encoding is technically actor dependent, but for all built-in Filecoin actors it is the dag-cbor ipld encoding of the parameters struct for each method defined in actors.

### Storage

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


## Pledge Collateral

Filecoin includes a concept of "Pledge Collateral", which is FIL collateral that storage miners must lock up when participating as miners.

Pledge collateral serves several functions in Filecoin. It:

- makes it possible to slash misbehaving or slow miners
- ensures that miners have skin in the game (for the Filecoin network as a whole)
- increases the cost of launching a 51% attack


### Computing Pledge Collateral

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

### Dealing with Undercollateralization

In the course of normal events, miners may become undercollateralized.

They cannot directly undercollateralized themselves by adding more power, as commitSector will fail if they do not have sufficient collateral to cover their power requirements.
However, their collateral requirement could increase due to growth in availableFil, a reduction in the total network power, or a reduction in the total number of miners.
In such cases, the miner may continue to submit PoSts and mine blocks. When they win blocks, their block rewards will be garnished while they remain undercollateralized.

### Parameter Choices

We provisionally propose the following two parameters choices:

```go
POWER_COLLATERAL_PROPORTION := 0.2
PER_CAPITA_COLLATERAL_PROPORTION := 0.05
```

These are subject to change before launch.
