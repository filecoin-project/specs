---
title: 'Glossary'
weight: 6
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Glossary

## Account Actor

The [_Account Actor_](sysactors) is responsible for user accounts.

## Actor

The Actor is the Filecoin equivalent of the smart contract in Ethereum.

An actor is an on-chain object with its own state and set of methods. An actor's state is persisted in the on-chain state tree, keyed by its address. All actors (miner actors, the storage market actor, account actors) have an address. Actor's methods are invoked by crafting messages and getting miners to include them in blocks.

There eleven (11) builtin [System Actors](sysactors) in total in the Filecoin System.

## Address

An address is an identifier that refers to an actor in the Filecoin state.

In the Filecoin network, an _address_ is a unique cryptographic value that serves to publicly identify a user. This value, which is a public key, is paired with a corresponding private key. The mathematical relationship between the two keys is such that access to the private key allows the creation of a signature that can be verified with the public key. Filecoin employs the Boneh–Lynn–Shacham (BLS) signature scheme for this purpose.

## Ask

An _ask_ contains the terms on which a miner is willing to provide services. Storage asks, for example, contain price and other terms under which a given storage miner is willing to lease its storage. The word comes from stock market usage, shortened from asking price.

## Block

In a blockchain, a [_block_](block) is the fundamental unit of record. Each block is cryptographically linked to one or more previous blocks. Blocks typically contain messages that apply changes to the previous state (for example, financial records) tracked by the blockchain. A block represents the state of the network at a given point in time.

## Block Height

The _height_ of a block corresponds to the number of epochs elapsed from genesis before the block was added to the blockchain. That said, `height` and `epoch` are synonymous. The height of the Filecoin blockchain is defined to be the maximum height of any block in the blockchain.

## Block Reward

The reward in [FIL](glossary#fil) given to [storage miners](glossary#storage-miner-actor) for contributing to the network with storage and proving that they have stored the files they have committed to store. The _Block Reward_ is allocated to the storage miners that mine blocks and extend the blockchain.

## Blockchain

A [_blockchain_](filecoin_blockchain) is a system of record in which new records, or blocks are cryptographically linked to preceding records. This construction is a foundational component of secure, verifiable, and distributed transaction ledgers.

## Bootstrapping

Bootstrapping traditionally refers to the process of starting a network. In the context of the Filecoin network _bootstrapping_ refers to the process of onboarding a new Filecoin node in the Filecoin network and relates to connecting the new node to other peers, synchronizing the blockchain and "catching up" with the current state.

## Capacity commitment

If a storage miner doesn't find any available deal proposals appealing, they can alternatively make a _capacity commitment_, filling a sector with arbitrary data, rather than with client data. Maintaining this sector allows the storage miner to provably demonstrate that they are reserving space on behalf of the network. Also referred to as Committed Capacity (CC).

## Challenge Sampling

An algorithm for challenge derivation used in [Proof of Replication](glossary#proof-of-replication-porep) or [Proof of SpaceTime](glossary#proof-of-spacetime-post).

## CID

[CID](multiformats#cids) is short for Content Identifier, a self describing content address used throughout the IPFS ecosystem. CIDs are used in Filecoin to identify files submitted to the decentralized storage network. For more detailed information, see [the github documentation for it](https://github.com/ipld/cid).

## Client

There are two types of [clients](filecoin_nodes#node_types) in Filecoin, the storage client and the retrieval client, both of which can be implemented as part of the same physical host. All clients have an account, which is used to pay for the storage or retrieval of data.

## Collateral

[Collateral](filecoin_mining#miner_collaterals) is Filecoin tokens pledged by an actor as a commitment to a promise. If the promise is respected, the collateral is returned. If the promise is broken, the collateral is not returned in full.

In order to enter into a [storage deal](glossary#deal), a [storage miner](glossary#storage-miner-actor) is required to provide [FIL](glossary#fil) as _collateral_, to be paid out as compensation to a client in the event that the miner fails to uphold their storage commitment.

## Consensus

The algorithm(s) and logic needed so that the state of the blockchain is agreed across all nodes in the network.

## Consensus Fault Slashing

Consensus Fault Slashing is the penalty that a miner incurs for committing consensus faults. This penalty is applied to miners that have acted maliciously against the network's consensus functionality.

## Cron Actor

The [_Cron Actor_](sysactors) is a scheduler actor that runs critical functions at every epoch.

## Deal

Two participants in the Filecoin network can enter into a [_deal_](storage_market#deal-flow) in which one party contracts the services of the other for a given price agreed between the two. The Filecoin specification currently details _storage deals_ (in which one party agrees to store data for the other for a specified length of time) and _retrieval deals_ (in which one party agrees to transmit specified data to the other).

## Deal Quality Multiplier

This factor is assigned to different deal types (committed capacity, regular deals, and verified client deals) to reward different content.

## Deal Weight

This weight converts spacetime occupied by deals into consensus power. Deal weight of verified client deals in a sector is called Verified Deal Weight and will be greater than the regular deal weight.

## DRAND

[DRAND](drand), short for Distributed Randomness, is a publicly verifiable random beacon protocol that Filecoin uses as a source of unbiasable entropy for [leader election](glossary#leader-election). See the [DRAND website](https://drand.love/) for more details.

## Election

On every [epoch](glossary#epoch), a small subset of Filecoin [storage miners](glossary#storage-miner-actor) are _elected_ to mine a new, or a few new [block(s)](glossary#block) for the Filecoin blockchain. A miner's probability of being elected is roughly proportional to the share of the Filecoin network's total storage capacity they contribute. Election in Filecoin is realized through [Expected Consensus](expected_consensus).

## Election Proof

Election Proof is used as a source of randomness in EC leader election. The election proof is created by calling [VRF](glossary#vrf) and giving the secret key of the miner's worker and the DRAND value of the current epoch as input.

## Epoch

Time in the Filecoin blockchain is discretized into _epochs_ that are currently set to thirty (30) seconds in duration. On every epoch, a subset of storage miners are elected to each add a new block to the Filecoin blockchain via [Winning Proof-of-Spacetime](glossary#winning-proof-of-spacetime-winningpost). Also referred to as [Round](glossary#round).

## Fault

A fault occurs when a proof is not posted in the Filecoin system within the proving period, denoting another malfunction such as loss of network connectivity, storage malfunction, or malicious behaviour.

When a [storage miner](glossary#storage-miner-actor) fails to complete [Window Proof-of-Spacetime](glossary#window-proof-of-spacetime-windowpost) for a given sector, the Filecoin network registers a _fault_ for that sector, and the miner is [_slashed_](glossary#slashing). If a storage miner does not resolve the fault quickly, the network assumes they have abandoned their commitment.

## FIL

_FIL_ is the name of the Filecoin unit of currency; it is alternatively denoted by the Unicode symbol for an integral with a double stroke (⨎).

## File

[Files](file) are what clients bring to the filecoin system to store. A file is converted to a UnixFS DAG and is placed in a [piece](glossary#piece). A piece is the basic unit of account in the storage network and is what is actually stored by the Filecoin network.

## Filecoin

The term _Filecoin_ is used generically to refer to the Filecoin project, protocol, and network.

## Finality

[Finality](expected_consensus#finality-in-ec) is a well known concept in blockchain environments and refers to the amount of time needed until having a reasonable guarantee that a message cannot be reversed or cancelled. It is measured in terms of delay, normally in epochs or rounds from the point when a message has been included in a block published on-chain.

## `fr32`

The term `fr32` is derived from the name of a struct that Filecoin uses to represent the elements of the arithmetic field of a pairing-friendly curve, specifically Bls12-381—which justifies use of 32 bytes. `F` stands for "Field", while `r` is simply a mathematic letter-as-variable substitution used to denote the modulus of this particular field.

## Gas, Gas Fees

_Gas_ is a property of a [message](glossary#message), corresponding to the resources involved in including that message in a given [block](glossary#block). For each message included in a block, the block's creator (i.e., miner) charges a fee to the message's sender.

## Genesis Block

The _genesis block_ is the first block of the Filecoin blockchain. As is the case with every blockchain, the genesis block is the foundation of the blockchain. The tree of any block mined in the future should link back to the genesis block.

## GHOST

[GHOST](https://eprint.iacr.org/2013/881.pdf) is an acronym for `Greedy Heaviest Observable SubTree`, a class of blockchain structures in which multiple blocks can validly be included in the chain at any given height or round. GHOSTy protocols produce blockDAGs rather than blockchains and use a weighting function for fork selection, rather than simply picking the longest chain.

## Height

Same as Block Height.

## Init Actor

The [_Init Actor_](sysactors) initializes new actors and records the network name.

## Lane

[_Lanes_](payment_channels#lanes) are used to split Filecoin Payment Channels as a way to update the channel state (e.g., for different services exchanged between the same end-points/users). The channel state is updated using [vouchers](glossary#voucher) with every lane marked with an associated `nonce` and amount of tokens it can be redeemed for.

## Leader

A leader, in the context of Filecoin consensus, is a node that is chosen to propose the next block in the blockchain during [Leader Electioni](glossary#leader-election).

## Leader election

Same as [Election](glossary#election).

## Message

A message is a call to an actor in the Filecoin VM. The term _message_ is used to refer to data stored as part of a [block](glossary#block). A block can contain several messages.

## Miner

A miner is an actor in the Filecoin system performing a service in the network for a reward.

There are three types of miners in Filecoin:

- [Storage miners](glossary#storage-miner-actor), who store files on behalf of clients.
- [Retrieval miners](glossary#retrieval-miner), who deliver stored files to clients.
- Repair miners, who replicate files to keep them available in the network, when a storage miner presents a fault.

## Multisig Actor

The [_Multisig Actor_](sysactors) (or Multi-Signature Wallet Actor) is responsible for dealing with operations involving the Filecoin wallet.

## Node

A [node](filecoin_nodes) is a communication endpoint that implements the Filecoin protocol.

## On-chain/off-chain

On-chain actions are those that change the state of the tree and the blockchain and interact with the Filecoin VM. Off-chain actions are those that do not interact with the Filecoin VM.

## Payment Channel

A [_payment channel_](filecoin_token#payment_channels) is set up between actors in the Filecoin system to enable off-chain payments with on-chain guarantees, making settlement more efficient. Payment channels are managed by the Payment Channel Actor, who is responsible for setting up and settling funds related to payment channels.

## Piece

The [Piece](piece) is the main unit of account and negotiation for the data that a user wants to store on Filecoin. In the Lotus implementation, a piece is a [CAR file](https://github.com/ipld/specs/blob/master/block-layer/content-addressable-archives.md#summary) produced by an IPLD DAG with its own _payload CID_ and _piece CID_. However, a piece can be produced in different ways as long as the outcome matches the _piece CID_. A piece is not a unit of storage and therefore, it can be of any size up to the size of a [sector](glossary#sector). If a piece is larger than a sector (currently set to 32GB or 64GB and chosen by the miner), then it has to be split in two (or more) pieces. For more details on the exact sizing of a Filecoin Piece as well as how it can be produced, see the [Piece section](piece).

## Pledged Storage

Storage capacity (in terms of [sectors](glossary#sector))that a miner has promised to reserve for the Filecoin network via [Proof-of-Replication](glossary#proof-of-replication-porep) is termed _pledged storage_.

## Power

See [Power Fraction](glossary#power-fraction).

## Power Fraction

A storage miner's `Power Fraction` or `Power` is the ratio of their committed storage, as of their last PoSt submission, over Filecoin's total committed storage as of the current block. It is used in [Leader Election](glossary#leader-election). It is the proportion of power that a storage miner has in the system as a fraction of the overall power of all miners.

## Power Table

The Power Table is an abstraction provided by the Filecoin storage market that lists the `power` of every [storage miner](glossary#storage-miner-actor) in the system.

## Protocol

Commonly refers to the "Filecoin Protocol".

## Proving Period

Commonly referred to as the "duration of a PoSt", the _proving period_ is the period of time during which storage miners must compute Proofs of Spacetime. By the end of the period they must have submitted their PoSt.

## Proving Set

The elements used as input by a proof of Spacetime to enable a proof to be generated.

## Proof of Replication (PoRep)

[_Proof-of-Replication_](pos#porep) is a procedure by which a [storage miner](glossary#storage-miner-actor) can prove to the Filecoin network that they have created a unique copy of some piece of data on the network's behalf. PoRep is used in the Filecoin system to generate sealed sectors through which storage miners prove they hold client data.

## Proof of Spacetime (PoSt)

[_Proof-of-Spacetime_](pos#post) is a procedure by which a [storage-miner](glossary#storage-miner-actor) can prove to the Filecoin network they have stored and continue to store a unique copy of some data on behalf of the network for a period of time. Proof-of-Spacetime manifests in two distinct varieties in the present Filecoin specification: [Window Proof-of-Spacetime](glossary#window-proof-of-spacetime-windowpost) and [Winning Proof-of-Spacetime](glossary#winning-proof-of-spacetime-winningpost).

## Quality-Adjusted Power

This parameter measures the consensus power of stored data on the network, and is equal to [Raw Byte Power](glossary#raw-byte-power) multiplied by [Sector Quality Multiplier](glossary#sector-quality-multiplier).

## Randomness

Randomness is used in Filecoin in order to generate random values for electing the next leader and prevent malicious actors from predicting future and gaining advantage over the system. Random values are drawn from a [DRAND](glossary#drand) beacon and appropriately formatted for usage.

## Randomness Ticket

See Ticket.

## Raw Byte Power

This measurement is the size of a sector in bytes.

## Retrieval miner

A [_retrieval miner_](retrieval_market#retrieval_provider) is a Filecoin participant that enters in retrieval [deals](glossary#deal) with clients, agreeing to supply a client with a particular file in exchange for [FIL](glossary#fil). Note that unlike [storage miners](glossary#storage-miner-actor), retrieval miners are not additionally rewarded with the ability to add blocks to (i.e., extend) the Filecoin blockchain; their only reward is the fee they extract from the client.

## Repair

Repair refers to the processes and protocols by which the Filecoin network ensures that data that is partially lost (by, for example, a miner disappearing) can be re-constructed and re-added to the network. Repairing is done by [Repair Miners](glossary#miner).

## Reward Actor

The [_Reward Actor_](sysactors) is responsible for distributing block rewards to [storage miners](glossary#storage-miner-actor) and token vesting.

## Round

A Round is synonymous to the [epoch](glossary#epoch) and is the time period during which new blocks are mined to extend the blockchain. The duration is of a round is set to 30 sec.

## Seal

Sealing is a cryptographic operation that transforms a sector packed with deals into a certified replica associated with: i) a particular miner’s cryptographic identity, ii) the sector's own identity.

_Sealing_ is one of the fundamental building blocks of the Filecoin protocol. It is a computation-intensive process performed over a [sector](glossary#sector) that results in a unique representation of the sector as it is produced by a specific miner. The properties of this new representation are essential to the [Proof-of-Replication](glossary#proof-of-replication-porep) and the [Proof-of-Spacetime](glossary#proof-of-spacetime-post) procedures.

## Sector

The [sector](sector) is the default unit of storage that miners put in the network (currently 32GBs or 64GBs). A sector is a contiguous array of bytes that a [storage miner](glossary#storage-miner-actor) puts together, seals, and performs Proofs of Spacetime on. Storage miners store data on behalf of the Filecoin network in fixed-size sectors.

Sectors can contain data from multiple deals and multiple clients. Sectors are also split in “Regular Sectors”, i.e., those that contain deals and “Committed Capacity” (CC), i.e., the sectors/storage that have been made available to the system, but for which a deal has not been agreed yet.

## Sector Quality Multiplier

Sector quality is assigned on Activation (the epoch when the miner starts proving theyʼre storing the file). The sector quality multiplier is computed as an average of deal quality multipliers (committed capacity, regular deals, and verified client deals), weighted by the amount of spacetime each type of deal occupies in the sector.

## Sector Spacetime

This measurement is the sector size multiplied by its promised duration in byte-epochs.

## Slashing

Filecoin implements two kinds of slashing: [**Storage Fault Slashing**](glossary#storage-fault-slashing) and [**Consensus Fault Slashing**](glossary#consensus-fault-slashing).

## Smart contracts

In the Filecoin blockchain smart contracts are referred to as [actors](glossary#actor).

## State

The _State_ or [_State Tree_](state_tree) refers to the shared history of the Filecoin system which contains actors and their storage power. The _State_ is deterministically generated from the initial state and the set of messages generated by the system.

## Storage Market Actor

The [_Storage Market Actor_](sysactors) is responsible for managing storage and retrieval deals.

## Storage Miner Actor

The [_Storage Miner Actor_](sysactors) commits storage to the network, stores data on behalf of the network and is rewarded in [FIL](glossary#fil) for the storage service. The storage miner actor is responsible for collecting proofs and reaching consensus on the latest state of the storage network. When they create a block, storage miners are rewarded with newly minted FIL, as well as the message fees they can levy on other participants seeking to include messages in the block.

## Storage Power Actor

The [_Storage Power Actor_](sysactors) is responsible for keeping track of the storage power allocated at each storage miner.

## Storage Fault Slashing

Storage Fault Slashing is a term that is used to encompass a broader set of penalties, including (but not limited to) Fault Fees, Sector Penalties, and Termination Fees. These penalties are to be paid by miners if they fail to provide sector reliability or decide to voluntarily exit the network.

- **Fault Fee (FF):** A penalty that a miner incurs for each day a miner's sector is offline.
- **Sector Penalty (SP):** A penalty that a miner incurs for a faulted sector that was not declared faulted before a WindowPoSt check occurs.
  - The sector will pay FF after incurring an SP when the fault is detected.
- **Termination Penalty (TP):** A penalty that a miner incurs when a sector is voluntarily or involuntarily terminated and is removed from the network.

## Ticket or VRF Chain

Tickets are generated as in [Election Proof](glossary#election-proof), but the input of every ticket includes the concatenation of the previous ticket, hence the term chain. This means that the new ticket is generated by running the VRF on the old ticket concatenated with the new DRAND value (and the key as with the Election Proof).

## Tipset

A [tipset](https://filecoin.io/blog/tipsets-family-based-approach-to-consensus/) is a set of [blocks](glossary#block) that each have the same [height](glossary#block-height) and parent tipset; the Filecoin [blockchain](glossary#blockchain) is a chain of tipsets, rather than a chain of blocks.

Each tipset is assigned a weight corresponding to the amount of storage the network is provided per the commitments encoded in the tipset's blocks. The consensus protocol of the network directs nodes to build on top of the heaviest chain.

By basing its blockchain on tipsets, Filecoin can allow multiple [storage miners](glossary#storage-miner-actor) to create blocks in the same [epoch](glossary#epoch), increasing network throughput. By construction, this also provides network security: a node that attempts to intentionally prevent the valid blocks of a second node from making it onto the canonical chain runs up against the consensus preference for heavier chains.

## Verified client

To further incentivize the storage of "useful" data over simple [capacity commitments](glossary#capacity-commitment), [storage miners](glossary#storage-miner-actor) have the additional opportunity to compete for special [deals](glossary#deal) offered by verified clients. Such clients are certified with respect to their intent to offer deals involving the storage of meaningful data, and the power a storage miner earns for these deals is augmented by a multiplier.

## Verified Registry Actor

The [_Verified Registry Actor_](sysactors) is responsible for managing [verified clients](glossary#verified-client).

## VDF

A Verifiable Delay Function that guarantees a random delay given some hardware assumptions and a small set of requirements. These requirements are efficient proof verification, random output, and strong sequentiality. Verifiable delay functions are formally defined by [BBBF](https://eprint.iacr.org/2018/601).

```text
{proof, value} <-- VDF(public parameters, seed)
```

## (Filecoin) Virtual Machine (VM)

The [Filecoin VM](actor) refers to the system by which changes are applied to the Filecoin system's state. The VM takes messages as input, and outputs updated state. The four main Actors interact with the Filecoin VM to update the state. These are: the `InitActor`, the `CronActor`, the `AccountActor` and the `RewardActor`.

## Voucher

[Vouchers](payment_channels#vouchers) are used as part of the Payment Channel Actor. Vouchers are signed messages exchanged between the channel creator and the channel recipient to acknowledge that a part of the service has been completed. Vouchers are the realisation of micropayments or checkpoints ini a payment channel. Vouchers are submitted to the blockchain and when `Collected`, funds are moved from the channel creator's account to the channel recipient's account.

## VRF

A Verifiable Random Function (VRF) that receives {Secret Key (SK), seed} and outputs {proof of correctness, output value}. VRFs must yield a proof of correctness and a unique & efficiently verifiable output.

```text
{proof, value} <-- VRF(SK, seed)
```

## Weight

Every mined block has a computed `weight`, also called its `WinCount`. Together, the `weights` of all the blocks in a branch of the chain determines the cumulative `weight` of that branch. Filecoin's Expected Consensus is a GHOSTy or heaviest-chain protocol, where chain selection is done on the basis of an explicit weighting function. Filecoin’s `weight` function currently seeks to incentivize collaboration amongst miners as well as the addition of storage to the network. The specific weighting function is defined in [Chain Weighting](expected_consensus#chain-weighting).

## Window Proof-of-Spacetime (WindowPoSt)

[_Window Proof-of-Spacetime_ (WindowPoSt)](post#windowpost) is the mechanism by which the commitments made by [storage miners](glossary#storage-miner-actor) are audited. It sees each 24-hour period broken down into a series of windows. Correspondingly, each storage miner's set of pledged [sectors](glossary#sector) is partitioned into subsets, one subset for each window. Within a given window, each storage miner must submit a [Proof-of-Spacetime](glossary#proof-of-spacetime-post) for each sector in their respective subset. This requires ready access to each of the challenged sectors, and will result in a [zk-SNARK-compressed](glossary#zksnark) proof published to the Filecoin [blockchain](glossary#blockchain) as a [message](glossary#message) in a [block](glossary#block). In this way, every sector of [pledged storage](glossary#pledged-storage) is audited at least once in any 24-hour period, and a permanent, verifiable, and public record attesting to each storage miner's continued commitment is kept.

The Filecoin network expects constant availability of stored data. Failing to submit WindowPoSt for a sector will result in a [fault](glossary#fault), and the storage miner supplying the sector will be [slashed](glossary#slashing).

## Winning Proof-of-Spacetime (WinningPoSt)

[_Winning Proof-of-Spacetime_ (WinningPoSt)](post#winningpost) is the mechanism by which [storage miners](glossary#storage-miner-actor) are rewarded for their contributions to the Filecoin network. At the beginning of each [epoch](glossary#epoch), a small number of storage miners are [elected](glossary#election) to each mine a new [block](glossary#block). As a requirement for doing so, each miner is tasked with submitting a compressed [Proof-of-Spacetime](glossary#proof-of-spacetime-post) for a specified [sector](glossary#sector). Each elected miner who successfully creates a block is granted [FIL](glossary#fil), as well as the opportunity to charge other Filecoin participants fees to include [messages](glossary#message) in the block.

Storage miners who fail to do this in the necessary window will forfeit their opportunity to mine a block, but will not otherwise incur penalties for their failure to do so.

## zk-SNARK

zk-SNARK stands for Zero-Knowledge Succinct Non-Interactive Argument of Knowledge.

An _argument of knowledge_ is a construction by which one party, called the _prover_, can convince another, the _verifier_, that the prover has access to some piece of information. There are several possible constraints on such constructions:

- A _non-interactive_ argument of knowledge has the requirement that just a single message, sent from the prover to the verifier, should serve as a sufficient argument.

A _zero-knowledge_ argument of knowledge has the requirement that the verifier should not need access to the knowledge the prover has access to in order to verify the prover's claim.

A _succinct_ argument of knowledge is one that can be "quickly" verified, and which is "small", for appropriate definitions of both of those terms.

A Zero-Knowledge Succinct Non-Interactive Argument of Knowledge (zk-SNARK) embodies all of these properties. Filecoin utilizes these constructions to enable its distributed network to efficiently verify that [storage miners](glossary#storage-miner-actor) are storing files they pledged to store, without requiring the verifiers to maintain copies of these files themselves.

In summary, Filecoin uses zk-SNARKs to produce a small 'proof' that convinces a 'verifier' that some computation on a stored file was done correctly, without the verifier needing to have access to the stored file itself.
