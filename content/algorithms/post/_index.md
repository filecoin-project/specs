---
title: PoRep & PoSt
bookCollapseSection: true
weight: 3
dashboardWeight: 2
dashboardState: wip
dashboardAudit: wip
dashboardTests: 0
---

# Proof-of-Storage

## Preliminaries

Storage miners in the Filecoin network have to prove that they hold a copy of the data at any given point in time. This is realised through the [Storage Miner Actor](storage_mining#storage_miner_actor) who is the main player in the [Storage Mining](filecoin_mining#storage_mining) subsystem. The proof that a storage miner indeed keeps a copy of the data they have promised to store is achieved through "challenges", that is, by providing answers to specific questions posed by the system. In order for the system to be able to prove that a challenge indeed proves that the miner stores the data, the challenge has to: i) target a random part of the data and ii) be requested at a time interval such that it is not profitable for the miner to discard the copy of data and refetch it when challenged.

General _Proof-of-Storage_ schemes allow a user to check if a storage provider is storing the outsourced data at the time
of the challenge. How can we use **PoS** schemes to prove that some data was being stored throughout a period
of time? A natural answer to this question is to require the user to repeatedly (e.g. every minute) send
challenges to the storage provider. However, the communication complexity required in each interaction can
be the bottleneck in systems such as Filecoin, where storage providers are required to submit their proofs to
the blockchain network.

To address this question, we introduce a new proof, called _Proof-of-Spacetime_, where a verifier can check if a prover
has indeed stored the outsourced data they committed to over (storage) Space and over a period of Time (hence, the name Spacetime). 
The miner is required to: (1) generate
sequential Proofs-of-Storage (in our case Proof-of-Replication), as a way to determine time (2) recursively
compose the executions to generate a short proof.

Recall that in the Filecoin network, miners are storing data in fixed-size [sectors](filecoin_mining#sector). Sectors are filled with client data agreed through deals in the [Storage Market](), through [verified deals](algorithms#verified_clients), or with random client data in case of [Committed Capacity sectors](filecoin_mining#sector).

## Proof-of-Replication (PoRep)

In order to register a sector with the Filecoin network, the sector has to be sealed. Sealing is a computation-heavy process that produces a unique representation of the data in the form of a proof, called **_Proof-of-Replication_** or PoRep.

The PoRep proof ties together: i) the data itself, ii) the miner actor that performs the sealing and iii) the time when the specific data has been sealed by the specific miner. In other words, if the same miner attempts to seal the same data at a later time, then this will result in a different PoRep proof. Time is included as the blockchain height when sealing took place and the corresponding chain reference is called `SealRandomness`.

Once the proof has been generated, the miner runs a SNARK on the proof in order to compress it and submits the result to the blockchain. This constitutes a certification that the miner has indeed replicated a copy of the data they agreed to store.

The PoRep process includes the following two phases:

- **Sealing preCommit phase 1.** In this phase, PoRep SDR [encoding](sdr#encoding) and [replication](sdr#replication) takes place.
- **Sealing preCommit phase 2.** In this phase, [Merkle proof and tree generation](sdr#merkle-proofs) is performed using the Poseidon hashing algorithm.

## Proof-of-Spacetime (PoSt)

From this point onwards, miners have to prove that they continuously store the data they pledged to store. PoSt is a procedure during which miners are given cryptographic challenges that can only be correctly answered if the miner is actually storing a copy of the sealed data. The challenge must be answered within a limited amount of time, which is smaller than the time that a miner would need to seal the data on the spot (had they discarded the sealed copy) in order to answer the challenge.

There are two challenges (and their corresponding mechanisms) that are realised as part of the PoSt process: _WindowPoSt_ and _WinningPoSt_.

### WindowPoSt

WindowPoSt is the mechanism by which the commitments made by storage miners are audited. In _WindowPoSt_ every 24-hour period is broken down into a series of 30min, non-overlapping windows, making a total of 48 windows, otherwise called "proving periods". Each storage miner’s set of pledged sectors is partitioned into subsets, so that there is one subset for each 30min window. Within a given window, each storage miner must submit a PoSt for each sector in their respective subset. It follows that the more the sectors a miner has pledged to store, the bigger the subset of sectors that the miner will need to prove per window. This requires ready access to each of the challenged sectors. For each sector, the miner will have to produce a SNARK-compressed proof and publish it to the blockchain as a message in a block. This proves that the miner has indeed stored the pledged sector. In this way, _every sector of pledged storage is audited at least once in any 24-hour period_, and a permanent, verifiable, and public record attesting to each storage miner’s continued commitment is kept.

The Filecoin network expects constant availability of stored files. Failing to submit WindowPoSt for a sector will result in a fault, and the storage miner supplying the sector will be slashed – that is, a portion of their [pledge collateral](filecoin_mining#miner_collaterasl) will be forfeited, and their storage power will be reduced (see [Storage Power Consensus](filecoin_blockchain#storage_power_consensus).


PoSt includes the following two phases, which are a continuation from the "Sealing preCommit phases 1 & 2" in PoRep:

- **Sealing commit phase 1.** This is an intermediate phase that performs preparation necessary to generate a proof from [PoSt Challenges](sdr#post-challenges). 
- **Sealing commit phase 2.** This final sealing phase involves the creation of a [PoSt Circuit](sdr#post-circuit) or SNARK, which is used to compress the requisite proof before it is broadcast to the blockchain.


### WinningPoSt

WinningPoSt is Filecoin's mechanism by which storage miners are chosen to be rewarded for their storage contribution to the network. At the beginning of each epoch, a small number of storage miners are elected to mine new blocks. Recall that the Filecoin blockchain operates on the basis of tipsets, which are groups of blocks. This means that in the Filecoin blockchain multiple blocks can be mined at the same height. Each elected miner who successfully creates a block is granted the Filecoin Block Reward, as well as the opportunity to charge other nodes fees in order to include their messages in the block.

The probability of a storage miner being elected to mine a block is governed by Filecoin's [Expected Consensus](algorithms#expected_consensus) algorithm and guarantees that miners will be chosen (on expectation) proportionally to their _Quality Adjusted Power_ in the network, as reported n the power table `WinningPoStSectorLookback` epochs before the election. When a miner is chosen to mine a new block and extend the blockchain, they have to submit a proof similar to the one submitted for WindowPoSt before the end of the current epoch. If they miss the epoch end deadline, then the miner misses the opportunity to mine a block and get a Block Reward. No penalty is incurred in this case.

### Constants & Terminology

- **partition:** a group of 2350 sectors proven simultaneously.
- **proving period:** average period for proving all sectors (currently set to ~24 hours).
- **deadline:** one of multiple points during a proving period when the proofs for some partitions are due.
- **challenge window:** the period immediately before a deadline during which a challenge can be generated by the chain and the requisite proofs computed.
- **miner size:** the amount of proven storage maintained by a single miner actor.


Every miner must demonstrate availability of all claimed sectors on a 24hr basis. Constraints on individual proof computations limit a single proof to 2350 sectors (a partition), with 10 challenges each.

As an illustrative example of the current sector size, with a 32 GiB sector, 1 EiB of storage power requires:
- ~33 million sectors
- ~14 thousand partitions
- ~14 thousand proofs per period (day)

### Design

Each miner actor is allocated a 24-hr proving period at random upon creation. This proving period is divided into 48 non-overlapping half-hour deadlines. Each sector is assigned to one of these deadlines at the beginning of the first proving period when the miner registers the sector and never changes deadline. The sets of sectors due at each deadline is recorded in a collection of 48 bitfields.

At the beginning of each period, a cron invocation removes expired sectors and allocates any newly-proven sectors sequentially to one of the 48 deadlines. Sectors are first allocated to fill any deadline up to the next whole-partition multiple (2350) of sectors; next a new partition is started on the deadline with the fewest partitions. If all deadlines have the same number of sectors, a new partition is opened on a random deadline.

The per-deadline sector sets are frozen at the beginning of each proving period as proving set bitfields. The sector IDs are then (logically) divided sequentially into partitions, and the partitions across all deadlines for the miner logically numbered sequentially. Thus, a sector may move between partitions at the same deadline as other sectors fault or expire, but never changes deadline.

If a miner adds 48 partitions worth of sectors (~3.8 PiB) they will have one proof due for each deadline. When a miner has more than 48 partitions, some deadlines will have multiple proofs due at the same deadline. These simultaneous proofs are expected to be computed and submitted together in a single message, at least up to 10-20 partitions per message, but can be split arbitrarily between messages (which, however, will cost more gas).

A PoSt proof submission must indicate which deadline it targets and which partition indices for that deadline the proofs represent. The actor code receiving a submission maps the partition numbers through the deadline’s proving set bitfields to obtain the sector numbers. Faulty sectors are masked from the proving set by substituting a non-faulty sector number. The actor records successful proof verification for each of the partitions in a bitfield of partition indices (or records nothing if verification fails).

There are currently three types of Faults, the _Declared Fault_, the _Detected Fault_ and the _Skipped Fault_. They are discussed in more detail as part of the [Storage Mining subystem](storage_mining#faults-penalties-and-fee-debt).

Summarising:

- A miner maintains its sectors *active* by generating Proofs-of-Spacetime (PoSt) and submit `miner.SubmitWindowedPoSt` for their sectors in a timely manner.
- A PoSt proves that sectors are persistently stored through time.
- Each miner proves all of its sectors once per *proving period*; each sector must be proven by a particular time called deadline.
- Sectors are also assigned to a partition. A partition is a set of sectors that is not larger than the Seal Proof allowed number of sectors `sp.WindowPoStPartitionSectors`.
- Sectors are assigned to a partition at `miner.ProveCommitSector` and they can be re-arranged via `CompactPartitions`.
- Partitions are a by-product of our current proof mechanism. There is a limit in the number of sectors (`sp.WindowPoStPartitionSectors`) that can be proven in a single SNARK proof. If more than this amount is required to be proven, more than one SNARK proof is required, given that each SNARK proof represents a partition.
- A *proving period* is a period of `WPoStProvingPeriod` epochs in which a `Miner` actor is scheduled to prove its storage.
- A *proving period* is evenly divided in `WPoStPeriodDeadlines` *deadlines*.
- Each miner has a different start of proving period `ProvingPeriodStart` that is assigned at `Power.CreateMiner`.
- A *deadline* is a period of `WPoStChallengeWindow` epochs that divides a proving period.
- Sectors are assigned to a deadline on `miner.ProveCommitSector` and will remain assigned to it throughout their lifetime.
- In order to prove that they continuously store a sector, a miner must submit a `miner.SubmitWindowedPoSt` for each deadline.

There are four relevant epochs associated to a deadline, shown in the table below:

| Name          | Distance from `Open`      | Description                                                                                                                   |
|---------------|---------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `Open`        | `0`                       | Epoch from which a PoSt Proof for this deadline can be submitted.                                                             |
| `Close`       | `WPoStChallengeWindow`    | Epoch after which a PoSt Proof for this deadline will be rejected.                                                            |
| `FaultCutoff` | `-FaultDeclarationCutoff` | Epoch after which a `miner.DeclareFault` and `miner.DeclareFaultRecovered` for sectors in the upcoming deadline are rejected. |
| `Challenge`   | `-WPoStChallengeLookback` | Epoch at which the randomness for the challenges is available.                                                                |


