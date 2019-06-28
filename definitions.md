# Definitions / Glossary

## Updates to definitions

To make any updates to these definitions please submit a pull request with the changes, or open an issue and one of the maintainers will do it for you.

## Notes

- Want to split all repair stuff to separate doc
- Let's refer to Filecoin system rather than network. In a sense, the network is an instantiation of the system (this protocol). We can however refer to the Filecoin VM separately which means the system by which we apply changes to the state of the system at a point in time.
- Asterisks indicate that the definition requires updating by any affected party.

## Definitions

#### Actor

An actor is an on-chain object with its own state and set of methods. An actors state is persisted on-chain in the state tree, keyed by its address. All actors (miner actors, the storage market actor, account actors) have an address. Actors methods are invoked by crafting messages and getting miners to include them in blocks.

Actors are very similar to smart contracts in Ethereum.

#### Address

An address is an identifier that refers to an actor in the Filecoin state.

#### Ask

#### Bid

#### Block

A block in the Filecoin blockchain is a chunk of data appended to the shared history of the network including transactions, messages, etc. and representing the state of the storage network at a given point in time.

See [Data Structures]()

#### Bootstrapping

#### Chain weight

#### Challenge sampling

#### Cid

CID is short for Content Identifier, a self describing content address used throughout the ipfs ecosystem. For more detailed information, see [the github documentation for it](https://github.com/ipld/cid).

#### Client

A client is any user with an account who wishes to store data with a miner. A client's account is used to pay for the storage, and helps to prove the clients ability to pay.

#### Collateral

Collateral is Filecoin tokens pledged by an actor as a commitment to a promise. If the promise is respected, the collateral is returned. If the promise is broken, the collateral is not returned in full. For instance:

- In becoming a Filecoin storage miner: the miner will put up collateral alongside their SEAL to 
- In a Filecoin deal: both the miner and client put up collateral to ensure their respect of deal terms.

#### Commitment

See [Filecoin Proofs](proofs.md)

#### Confirmation

#### Consensus

#### Deal

*** *A deal in a Filecoin market is made when a bid and ask are matched, corresponding to an agreement on a service and price between a miner and client.

#### Erasure coding

Erasure coding is a strategy through which messages can be lengthened so as to be made recoverable in spite of errors.

See [Wikipedia](https://en.wikipedia.org/wiki/Erasure_code)

#### Epoch

An epoch refers to the period over which a given random seed is used in Expected Consensus. In the current Filecoin implementation, each round is an epoch.

#### Fault

A fault occurs when a proof is not posted in the Filecoin system within the proving period, denoting another malfunction such as loss of network connectivity, storage malfunction, malicious miner, etc.

#### Fair

#### File

Files are what clients bring to the filecoin system to store. A file is split up into `pieces`, which are what is actually stored by the network.

#### Finality

#### Piece Inclusion Proof

See [Filecoin Proofs](proofs.md)

#### Gas, Fees, Prices

#### Generation Attack Threshold

Security parameter. Number of rounds within which a new Proof-of-Storage must be submitted in order for a miner to retain power in the network (and avoid getting slashed). This number must be be smaller than the minimum time it takes for an adversarial miner to generate a replica of the data (thereby not storing it undetectably for some period of time).

The Generation Attack Threshold is equal to the Polling Time + some Grace Period after which miners get slashed.

#### GHOST

#### Leader

A leader, in the context of Filecoin consensus, is a node that is chosen to propose the next block in the blockchain.

#### Leader election

Leader election is the process by which the Filecoin network agrees who gets to create the next block.

#### Message

A message is a call to an actor in the Filecoin VM.

#### Miner

A miner is an actor in the Filecoin system performing a service in the network for a reward.

There are multiple types of miners in Filecoin:

- Storage miners - storage miners 
- Retrieval miners:
- Repair miners (to be split out):

#### Node

*** *A node is a communication endpoint that implements the Filecoin protocol. (also mention IPLD Node?)

#### Null Blocks

A null block refers to a block with no content mined by default during an epoch in which no miner is elected leader.

#### On-chain/off-chain

#### Online/offline

#### Payment Channel

A payment channel is set up between actors in the Filecoin system to enable off-chain payments with on-chain guarantees, making settlement more efficient.

#### Piece

A piece is a portion of a file that gets fitted into a sector.

#### Pledge

****The initial commitment of a storage miner to provide a number of sectors to the system.

#### Polling Time

Security Parameter. Polling time is the time between two online PoReps in a PoSt proof.

#### Power table

The power table is an abstraction provided by the Filecoin storage market that lists the `power` of every miner in the system.

#### Power table lookback

Security parameter. A number of rounds to sample back from when determining miner `power` for use in leader election. A higher number can help secure the system by making potential attacks more costly (as power must be maintained for longer to take effect), but also makes mining more costly for the same reason.

Also referred to as `L` in consensus settings. 

#### Protocol

#### Proving Period

The period of time during which storage miners must compute Proofs of Spacetime. At the end of the period they must submit their PoSt. Put another way, it is the duration of a PoSt.

#### Proving Set

The elements used as input by a proof of Spacetime to enable a proof to be generated.

**** elements necessary to generate a SEAL, or elements necessary to generate a proof

#### Proof of Replication

Proof that a unique encoding of data exists in physical storage.

Used in the Filecoin system to generate SEALed sectors through which storage miners prove they hold client data.

#### Proof of Spacetime

Proof that a given encoding of data existed in physical storage continuously over a period of time.

Used in the Filecoin system by a storage miner to prove that client data was kept over the contract duration.

#### Random(ness)

****Source of unpredictability used in the Filecoin system to ensure fairness and prevent malicious actors from gaining an advantage over the system.

TODO add a note to distinguish predictability from randomness

#### Election Randomness Lookback

Security parameter. A number of rounds to sample back from when choosing randomness for use in leader election. A higher number turns a more localized lottery into a more global one since a miner wins or loses on all descendants of a given randomness, but enables miners to look-ahead and know whether they will be elected in the future.

Also referred to as `K` in consensus settings. 

#### Repair

Repair refers to the processes and protocols by which the Filecoin network ensures that data that is partially lost (by, for example, a miner disappearing) can be re-constructed and re-added to the network.

#### Round

A round refers to the period over which a new leader election occurs. Typically, an election will find a single leader. If a single leader is found, a single block is generated. If multiple leaders are found, multiple blocks are generated. If no leader is found, no block is generated.

#### SEAL/UNSEAL

See [Filecoin Proofs](proofs.md)

#### Sector

A sector is a contiguous array of bytes that a miner puts together, seals, and performs Proofs of Spacetime on.

#### Slashing

#### Smart contracts

#### Storage

Storage widely refers to a place in which to store data in a given system.

In the context of:

- The Filecoin miner: sotrage refers to disk sectors made available to the network.
- The Filecoin chain: storage refers to the way in which system state is tracked through time on-chain through blocks.
- Actor: the struct that defines an actor.

#### State

****Refers to The shared history of the Filecoin system contains actors and their storage, deals, etc. State is deterministically generated from the initial state and the set of messages generated by the system.

#### Ticket

Some unpredictable element generated by the system for two uses:

- as a random challenge to PoSTs.
- to elect a leader in Expected Consensus.

See more on [Expected Consensus](expected-consensus.md) and [PoSTs](proofs.md).

#### TipSet

A collection of blocks mined by different miners, each an elected leader of a given epoch. All tip sets have the same parent-set and epoch number (height). TODO add picture.

#### Verifiable

Something that is verifiable can be checked for correctness by a third party.

#### VM

Virtual Machine. The Filecoin VM refers to the system by which changes are applied to the Filecoin system's state. The VM takes messages as input, and outputs state.

#### Voucher

Held by an actor as part of a payment channel to complete settlement when the counterparty defaults.

#### zkSNARK

Zero Knowledge Succinct ARguments of Knowledge. A way of producing a small 'proof' that convinces a 'verifier' that some computation was done correctly.
