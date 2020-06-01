---
title: "Glossary"
---

# Updates to definitions

To make any updates to these definitions please submit a pull request with the changes, or open an issue and one of the maintainers will do it for you.

# Notes

- Want to split all repair stuff to separate doc
- Let's refer to Filecoin system rather than network. In a sense, the network is an instantiation of the system (this protocol). We can however refer to the Filecoin VM separately which means the system by which we apply changes to the state of the system at a point in time.
- Asterisks indicate that the definition requires updating by any affected party.

# Definitions

## Actor

An actor is an on-chain object with its own state and set of methods. An actors state is persisted on-chain in the state tree, keyed by its address. All actors (miner actors, the storage market actor, account actors) have an address. Actors methods are invoked by crafting messages and getting miners to include them in blocks.

Actors are very similar to smart contracts in Ethereum.

## Address

An address is an identifier that refers to an actor in the Filecoin state.

## Ask
An ask contains the terms on which a miner is willing to provide services. Storage asks, for example, contain price and other terms under which a given miner is willing to sell its storage. The word comes from stock market usage of ask, shortened from asking price.

## Block

A block in the Filecoin blockchain is a chunk of data appended to the shared history of the network including transactions, messages, etc. and representing the state of the storage network at a given point in time.

See [Data Structures]()

## Bootstrapping

## Chain weight

## Challenge sampling

## Challenge Ticket

Derived from a `Partial Ticket` as part of ElectionPoSt and used to run leader election.

## Cid

CID is short for Content Identifier, a self describing content address used throughout the IPFS ecosystem. For more detailed information, see [the github documentation for it](https://github.com/ipld/cid).

## Client

A client is any user with an account who wishes to store data with a miner. A client's account is used to pay for the storage, and helps to prove the clients ability to pay.

## Collateral

Collateral is Filecoin tokens pledged by an actor as a commitment to a promise. If the promise is respected, the collateral is returned. If the promise is broken, the collateral is not returned in full. For instance:

- In becoming a Filecoin storage miner: the miner will put up collateral alongside their SEAL.
- In a Filecoin deal: both the miner and client put up collateral to ensure their respect of deal terms.

## Commitment

See [Filecoin Proofs](proofs.md)

## Confirmation

## Consensus

## Deal

*** *A deal in a Filecoin market is made when a bid and ask are matched, corresponding to an agreement on a service and price between a miner and client.

## Erasure coding

Erasure coding is a strategy through which messages can be lengthened so as to be made recoverable in spite of errors.

See [Wikipedia](https://en.wikipedia.org/wiki/Erasure_code)

## Fault

A fault occurs when a proof is not posted in the Filecoin system within the proving period, denoting another malfunction such as loss of network connectivity, storage malfunction, malicious miner, etc.

## Fair

## File

Files are what clients bring to the filecoin system to store. A file is split up into `pieces`, which are what is actually stored by the network.

## Finality

## Piece Inclusion Proof

See [Filecoin Proofs](proofs.md)

## Gas, Fees, Prices

## Generation Attack Threshold

Security parameter. Number of rounds within which a new Proof-of-Storage must be submitted in order for a miner to retain power in the network (and avoid getting slashed). This number must be be smaller than the minimum time it takes for an adversarial miner to generate a replica of the data (thereby not storing it undetectably for some period of time).

The Generation Attack Threshold is equal to the Polling Time + some Grace Period after which miners get slashed.

## GHOST

[GHOST](https://eprint.iacr.org/2013/881.pdf) is an acronym for `Greedy Heaviest Observable SubTree`, a class of blockchain structures in which multiple blocks can validly be included in the chain at any given height or round. GHOSTy protocols produce blockDAGs rather than blockchains and use a weighting function for fork selection, rather than simply picking the longest chain.

## Height

`Height` and `epoch` are synonymous and used interchangeably in this spec.

## Leader

A leader, in the context of Filecoin consensus, is a node that is chosen to propose the next block in the blockchain.

## Leader election

Leader election is the process by which the Filecoin network agrees who gets to create the next block.

## Message

A message is a call to an actor in the Filecoin VM.

## Miner

A miner is an actor in the Filecoin system performing a service in the network for a reward.

There are multiple types of miners in Filecoin:

- Storage miners - storage miners
- Retrieval miners:
- Repair miners (to be split out):

## Node

*** *A node is a communication endpoint that implements the Filecoin protocol. (also mention IPLD Node?)

## On-chain/off-chain

## Online/offline

## Partial Ticket

Ticket produced as part of the ElectionPoSt process used to both prove leader elections and prove storage of a given data replica through PoSt. At least one is included in every block header.

## Payment Channel

A payment channel is set up between actors in the Filecoin system to enable off-chain payments with on-chain guarantees, making settlement more efficient.

## Piece

A piece is a portion of a file that gets fitted into a sector.

## Pledge

****The initial commitment of a storage miner to provide a number of sectors to the system.

## Polling Time

Security Parameter. Polling time is the time between two online PoReps in a PoSt proof.

## Power

See `Power Fraction`.

## Power Fraction

A miner's `Power Fraction` or `Power` is the ratio of their committed storage as of their last PoSt submission over Filecoin's total committed storage as of the current block. It is used in leader election.

##  Power table

The {{<sref power_table>}} is an abstraction provided by the Filecoin storage market that lists the `power` of every miner in the system.

## Protocol

## Proving Period

The period of time during which storage miners must compute Proofs of Spacetime. At the end of the period they must submit their PoSt. Put another way, it is the duration of a PoSt.

## Proving Set

The elements used as input by a proof of Spacetime to enable a proof to be generated.

**** elements necessary to generate a SEAL, or elements necessary to generate a proof

## Proof of Replication

Proof that a unique encoding of data exists in physical storage.

Used in the Filecoin system to generate SEALed sectors through which storage miners prove they hold client data.

## Proof of Spacetime

Proof that a given encoding of data existed in physical storage continuously over a period of time.

Used in the Filecoin system by a storage miner to prove that client data was kept over the contract duration.

## Random(ness)

****Source of unpredictability used in the Filecoin system to ensure fairness and prevent malicious actors from gaining an advantage over the system.

TODO add a note to distinguish predictability from randomness

## Randomness Ticket

See Ticket.

## Election Randomness Lookback

Security parameter. A number of rounds to sample back from when choosing randomness for use in leader election. A higher number turns a more localized lottery into a more global one since a miner wins or loses on all descendants of a given randomness, but enables miners to look-ahead and know whether they will be elected in the future.

Also referred to as `K` in consensus settings.

## Repair

Repair refers to the processes and protocols by which the Filecoin network ensures that data that is partially lost (by, for example, a miner disappearing) can be re-constructed and re-added to the network.

## Round

See `Height` for definition. They are synonymous.

## SEAL/UNSEAL

See [Filecoin Proofs](proofs.md)

## Sector

A sector is a contiguous array of bytes that a miner puts together, seals, and performs Proofs of Spacetime on.

## Slashing

Filecoin implements two kinds of slashing: **Storage Fault Slashing** and **Consensus Fault Slashing**. Below are definitions of both types.

### Storage fault slashing

Storage Fault Slashing is a term that is used to encompass a broader set of penalties, including (but not limited to) Fault Fees, Sector Penalties, and Termination Fees. These penalties are to be paid by miners if they fail to provide sector reliability or decide to voluntarily exit the network. 

* **Fault Fee (FF):** A penalty that a miner incurs for each day a miner's sector is offline.
* **Sector Penalty (SP):**  A penalty that a miner incurs for a faulted sector that was not declared faulted before a WindowPoSt check occurs. 
    * The sector will pay FF after incurring an SP when the fault is detected.
* **Termination Fee (TP):** A penalty that a miner incurs when a sector is voluntarily or involuntarily terminated and is removed from the network.

### Consensus fault slashing

Consensus Fault Slashing is the penalty that a miner incurs for committing consensus faults. This penalty is applied to miners that have acted maliciously against the network's consensus functionality.




## Smart contracts

## Storage

Storage widely refers to a place in which to store data in a given system.

In the context of:

- The Filecoin miner: sotrage refers to disk sectors made available to the network.
- The Filecoin chain: storage refers to the way in which system state is tracked through time on-chain through blocks.
- Actor: the struct that defines an actor.

## State

****Refers to The shared history of the Filecoin system contains actors and their storage, deals, etc. State is deterministically generated from the initial state and the set of messages generated by the system.

## Ticket

A `ticket` is used as a source of randomness in EC leader election. Every block depends on one or more `ChallengeTicket` derived from a `PartialTicket` using input from a `RandomnessTicket` (naming PR welcome). One such new `RandomnessTicket` or `Ticket` is produced with every new block and included in its header

## Ticket Chain

Each chain in Filecoin can be associated with a given `ticket chain`. The `ticket chain` is assembled by taking the tickets (usually one) contained by the block with the smallest final ticket in each of the chain's `TipSet`s.

Ticket comparison is done by interpreting the tickets' Bytes as unsigned integers (little endian representation).

## TipSet

A `TipSet` is a set of blocks that have the same parent set and same number of `tickets`, which implies they will have been mined at the same `height`. A `TipSet` can contain multiple blocks if more than one miner successfully mines a block at the same `height` as another miner.


## Verifiable

Something that is verifiable can be checked for correctness by a third party.

## VDF

A verifiable function that guarantees a time delay given some hardware assumptions and a small set of requirements. These requirements are efficient proof verification, random output, and strong sequentiality. Verifiable delay functions are formally defined by [\[BBBF]](https://eprint.iacr.org/2018/601).

```text
{proof, value} <-- VDF(public parameters, seed)
```

## VM

Virtual Machine. The Filecoin VM refers to the system by which changes are applied to the Filecoin system's state. The VM takes messages as input, and outputs state.

## Voucher

Held by an actor as part of a payment channel to complete settlement when the counterparty defaults.

## VRF

A verifiable random function that receives {Secret Key (SK), seed} and outputs {proof of correctness, output value}. VRFs must yield a proof of correctness and a unique & efficiently verifiable output.

```text
{proof, value} <-- VRF(SK, seed)
```

## Weight

Every mined block has a computed `weight`. Together, the `weights` of all the blocks in a branch of the chain determines the cumulative `weight` of that branch. Filecoin's Expected Consensus is a GHOSTy or heaviest-chain protocol, where chain selection is done on the basis of an explicit weighting function. Filecoin’s `weight` function currently seeks to incentivize collaboration amongst miners as well as the addition of storage to the network. The specific weighting function is defined in [Chain Weighting](expected-consensus.md#chain-weighting).

## zkSNARK

Zero Knowledge Succinct ARguments of Knowledge. A way of producing a small 'proof' that convinces a 'verifier' that some computation was done correctly.
