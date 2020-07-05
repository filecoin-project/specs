---
title: Election PoSt
---

# Election PoSt
---


This document describes Election-PoSt, the Proof-of-Spacetime used in Filecoin.

At a high-level it marries `ElectionPoSt` with a `SurprisePoSt` fallback:

- By coupling leader election and PoSt, `ElectionPoSt` ensures that miners must do the work to prove their sectors at every round in order to earn block rewards.
- Small miners may not win on a regular basis however, `SurprisePoSt` thus comes in as a lower-bound to how often miners must PoSt and helps ensure the Power Table does not grow stale for its long-tail of smaller miners.

{{< hint danger >}}
Issue with label
{{< /hint >}}

{{</* label pledge_collateral */>}}
# Overview

Election PoSt couples the PoSt process with block production, meaning that in order to produce a block, the miner must produce a valid PoSt proof (snark output). Specifically, a subset of non-faulted sectors the miner is storing (i.e. eligible sectors) allows them to attempt a leader election using a PartialTicket any of which could yield a valid ChallengeTicket for leader election. The probability of scratching a winning ChallengeTicket is dependent on sector size and total storage. Miners are rewarded in proportion to the quantity of winning ChallengeTickets they generate in a given epoch, thereby incentivizing a miner to check as much of their storage as allowed in order to express their full power in a leader election. The number of election proofs a miner can generate in a given epoch will determine the block reward it earns.

An election proof is generated from a given PartialTicket by Hashing it, and using that hash to generate a value in [0,1]. Specifically `ChallengeTicket = H(PartialTicket)/2^len(H)`. 

This does mean that, in a given round, a lucky miner may succeed in generating a block by proving only a single sector, but again on expectation, a miner will have to prove their sectors at every single round in order for their full power to contribute to block generation. In the event one of the miner’s sectors cannot be proven (i.e. the miner does not have access to the nodes from that sector), no tickets will be returned. In order to prevent this, a miner can declare faults on their faulty sectors to avoid having to include them in the eligible sector set. Their power will be reduced accordingly.

The Election-PoSt construction includes a Surprise PoSt cleanup which will randomly challenge leaders who have not been elected in some time so they can prove their storage. Accordingly, miners will be challenged once per proving period on expectation. Any miner who fails to respond to the Challenge with a SurprisePoSt will see their power curbed across all their sectors (this is called a `Detected Fault`). If they are unable to prove their power within another two SurprisePoSts, their power will be terminated entirely.

## ElectionPoSt

To enable short PoSt response time, miners are required submit a PoSt when they are elected to mine a block, hence PoSt on election or ElectionPoSt. When miners win a block, they need to immediately generate a PoStProof and submit that in the block header along with the winning PartialTickets. Both the PartialTickets and PoStProof are checked at Block Validation by `StoragePowerConsenusSubsystem`. When a block is included, a special message is added that calls `SubmitElectionPoSt` which will process sector updates in the same way as successful `SubmitSurprisePoSt` do.

## SurprisePoSt Cleanup

In the absence of a posted `ElectionPoSt`, the chain randomly challenges a miner to submit a `SurprisePoSt` once per `ProvingPeriod` on expectation to ensure their storage is still accounted for. The process is largely the same as for `ElectionPoSt` but a successful `SurprisePoSt` does not entitle a miner to generate a block and gains them no reward. It allows them to maintain power in the power table.

For every miner challenged, a `NotifyOfPoStSurpriseChallenge` is issued and sent as an on-chain message to the chosen `StorageMinerActor`.
`PoStSurprise` will frequently be used by small miners who do not win blocks, and by miners as they are onboarding power to the network (since they will not be able to win ElectionPoSts to start).

# In Detail

## ElectionPoSt Generation

Filecoin's ElectionPoSt process makes use of two calls to the system library:

- `GenerateCandidates` which takes in the miner's sectors and a wanted number of Challenge Tickets and generates a number of inclusion proofs for a number of challenged sectors chosen randomly in proportion to the requested number of challenged tickets.
- `GeneratePoSt` takes a set of ChallengeTickets and generates a __*Proof of Spacetime*__ for them, proving the miner storage as a whole.

As stated above, a miner is incentivized to repeat this process at every block time in order to check whether they were elected leaders (see [Expected Consensus](\missing-link)). The rationality assumption made by ElectionPoSt is thus that storing files continuously and earning block rewards accordingly will be more profitable to miners than regenerating data at various epochs to sporadically participate in leader election.

At every epoch, each miner will challenge a portion of sectors `numSectorsSampled` from their `Proving Set` at random, according to some `ElectionPoStSampleRate` with each sector being issued K PoSt challenges (coverage may not be perfect).

By proving access to the challenged range of nodes (i.e. merkle tree leaf from the committed sector) in the sector, the miner can generate a set of valid ChallengeTickets in order to check them as part of leader election in EC (in order to find the winning tickets). The winning tickets will be stored on the block and used to generate a PoSt (using a SNARK). A block header will thus contain a number of “winning” PoStCandidates (each containing a partialTicket, SectorID and other elements, used to verify the leader election) and a PostProof generated from the ChallengeTickets.

This is all included in a field called `ElectionPoStOutput`

In order to simplify implementation and bound block header size, we can set a maximum number of possible election proofs for any block. For instance, for EC.E=5, we can cap challengeTicket submissions at 16 per block, which would cover more than 99.99% of cases (using Chernoff bounds) encountered by a 50% miner (i.e. much more in practice).

The epoch time is divided into three non-overlapping portions: (1) checking sectors for small challenge tickets, (2) computing the SNARK to produce a PoSt certifying the winning challenge tickets, and (3) block production/propagation. The target duration for each is still to be tuned. Proof parameter adjustments might allow a longer block propagation window while maintaining PoSt security (but trades off against sealing time and other things).

If no one has found a winning ticket in this epoch, increment epoch value as part of the post_randomness sampling and try again.

At every round:
1. **(sample randomness)**
The miner draws a randomness ticket from the randomness chain from a given epoch SPC_LOOKBACK_POST back and concats it with the minerID, and epoch for use as post_randomness (this is as it used to happen previously, except the SPC_LOOKBACK_POST is now also the EC randomness_lookback, which is small here):

    - `post_randomness = VRF_miner(ChainRandomness(currentBlockHeight - SPC_LOOKBACK_POST))`

2. **(select eligible sectors)**
The miner calls `GenerateCandidates` from proofs with their non-faulted (declared or detected) sectors, meaning those in their `proving set` (from the `StorageMinerActor`) along with `numSectorsSampled` `PartialTickets`. Note that the PoSt sectors are sampled over the `ProvingSet` and not just active sectors since if the PoSt is successful all PoSt in the `ProvingSet` will become active. That is, by including these sectors in their `ProvingSet`, miners are claiming that these are active sectors, which will be proven in the PoSt itself. Therefore we can use these sectors to generate the PoSt.

```text
    numSectorsMiner = len(miner.provingSet)
    numSectorsSampled = ceil(EPoStSampleRate * numSectorsMiner)
```

Note also that even with `numSectorsSampled == len(ProvingSet)`, this process may not sample all of a miner’s sectors in any given epoch, given how the data to be proven in challenged sectors is selected (there could be collisions, e.g. the same sectors selected multiple times).

3. **(generate Partial Ticket(s))** for each selected sector

    - first, generate K `PoStChallenges` (C_i), sampled at random from the chosen sector.
    There will be a chosen challenge-range size (a power of 2), called `ChallengeRangeSize`. Challenge ranges must be subtree-aligned and thus divide the sector into fixed-size data blocks. So challenge ranges can be indexed because allowable challenge ranges do not overlap. C_i is selected from the valid indexes into the challenge ranges.
    `C_i = HashToBlockNumber(post_randomness || S_j || i)`
    C_i corresponds to a given index of a known-sized data block (the sector)
    - for each C_i, miner reads the specified range of data from disk, which is of size `ChallengeRangeSize : C_i_output`
    - Miner generates a PartialTicket using all the PoSt witnesses
    - `PartialTicket = H(post_randomness || minerID ||S_ j || C_1_Output || … || C_K_Output)`

4. **(check Challenge Ticket(s) for winners)**
Given returned PartialTickets, miner checks them for winning tickets using the target set by expected consensus in [Expected Consensus](\missing-link) (per the `TicketIsWinner()` method below).

```text
winningTickets = []
def checkTicketForWinners(partialTickets):
    for partialTicket in partialTickets:
        challengeTicket = finalizeTicket(PartialTicket) 
        if TicketIsWinner(challengeTicket):
            winningTickets += partialTicket

def finalizeTicket(partialTicket):
    return H(ChallengeTicket) / 2^len(H)

```

A single winning ticket can be used to submit a block but a miner would want to check as many as possible to increase their potential rewards. The target ensures that on expectation, a miner's total power is expressed if they check all of their tickets, taking the `ElectionPoStSampleRate` into account.

5. **(generate a `PoStProof`)** for inclusion in the block header

    - Using the winning tickets from step 4 (there may be multiple tickets from the same `SectorNumber`), call `GeneratePoSt` from proofs to generate a single PoSt for the block

If no one has found a winning ticket in this epoch, increment epoch value as part of the post_randomness sampling and try again.

**Parameters:**

- `Randomness_ticket` --  a ticket drawn from the randomness ticket chain at a prior tipset 
- `Randomness_lookback` -- how far back to draw randomness from the randomness ticket chain - it will be as large as allowed by PoSt security, likely 1 or 2 epochs
- `K (e.g. 20-100s)` - number of  challenges per sector -- must be large enough such that the PoSpace is secure.
- `ChallengeRangeSize` - challenge read size (between 32B and 256KB)  -- based on security analysis.
- `EPoStsampleRate` - sector sampling fraction (e.g. 1, .10, .04) -- 1 to start-- It should be large enough to make it irrational to fully regenerate sectors. We may choose some subset if cost of verifying all is deleterious to disk
- `sectorsSampled` - Number of sectors sampled given the `EPoStsampleRate`.
- `networkPower` - filecoin network’s power - read from the power table, expressed in number of bytes

## ElectionPoSt verification

At a high-level, to validate a winning PoSt Proof:

1. **(Verify post_randomness):**

    - `post_randomness` is appropriate given `chain`, `SPC_LOOKBACK_POST` and `epoch`
    - VRF is generated by miner who submitted the PoSt:
    `{0,1} ?= VRF_miner.Verify(post_randomness)`
    - VRF output is valid given expected inputs:
    `{0,1} ?= VRF_miner.Validate(H(randomness_ticket || chainEpoch))`

2. Rederive eligible sectors in order to verify that winning sectors were appropriately selected (meaning it’s in `S_j`) from
    `{0, 1} ?= sectorID in (HashToSectorID(post_randomness || minerID || j))`

3. **Verify `PartialTicket` returned**

    - Prove that the PartialTicket were appropriately derived from the eligible sectors, by submitting all miner sectors along with the wanted number of tickets and verifying that the outputted PartialTicket match.
    - Verify that no duplicate PartialTickets were submitted, that is there are no two tickets with the same `challengeIndex` (challengeIndices are unique across challenged sectors). Tickets challenging the same sector at different indices are valid).

4. **Derive and validate the `ChallengeTicket` from the PartialTickets**

    - Prove the derived ChallengeTicket is below the target
    `{0, 1} ?= (ChallengeTicket < target)`

5. **Verify the `postProof` using the `PartialTickets`**

    - Verify that the postProof was appropriately derived from the PartialTickets.
    - The PoSt proof will verify the correctness of any PartialTickets passed to it as public inputs. In order to do this, it also needs the sector number along with various on-chain data (randomness, CommR, miner ID, etc.).

As it stands, the proofs caller passes a list of all of the replicas (with filesystem paths and auxiliary metadata) which could be challenged and a wanted number of tickets. This code structure forces the miner architecture to look like a single machine with a large number of disks presented as a single filesystem, as shown in spec. That won’t change for now. VFS systems like NFS may be used to distribute disks among multiple machines virtualizing the FileStore.

There is no requirement to persist the witnesses (list of merkle proofs) for failure recovery. If something fails, we don’t expect recovery to be fast enough to allow the miner to participate in that round anyway.

## Surprise PoSt cleanup

But while this means a miner will never win blocks from faked power, they won’t be penalized either when storage is lost. How do we ensure that the Power Table is accurate? Likewise, how do we onboard new power since it will not win a block unless it has at least X TB or Y % of the network.

This is why we need PoSt surprise. The process is largely the same as ElectionPoSt save for a few differences:

- A miner will only be surprised if they have not submitted a PoSt in some time (more than `SURPRISE_NO_CHALLENGE_PERIOD` epochs)
- A miner will have `CHALLENGE_DURATION` epochs to reply
- The PoSt will be submitted as a message on chain rather than part of the block header
- After being challenged, a miner will no longer be eligible to submit an ElectionPoSt

If a miner fails to respond to the challenge past the `CHALLENGE_DURATION`, they will lose all their power and a portion of their pledge collateral. This is considered a `DetectedFault` and all sectors in the `ProvingSet` will be marked as `Failing`. No deal collateral will be slashed if miners can recover within the next three proving periods. Note that the exact amount of slashed pledge collateral and pledge collateral function itself are subject to change.

Thereafter, the miner will have three more ProvingPeriods (specifically three challenges, so three ProvingPeriods on expectation) to recover their power by submitting a SurprisePoSt. If the miner does not do so, their sectors are permanently terminated and their storage deal collateral slashed (see the StorageMinerActor in the spec).

## PoStSurprise Challenge

Miners are selected at random to be challenged at every clockTick by the storage power actor (with `_selectMinersToSurprise`), with `len(PowerTable)/ProvingPeriod` miners selected per round (for one challenge per miner per proving period on expectation).

If the chosen `StorageMinerActor` is already challenged (`IsChallenged` is True) or they've submitted a post (election or surprise) in the last `SURPRISE_NO_CHALLENGE_PERIOD` epochs, they will not be challenged again. That is if `ShouldChallenge` is False, then the SurprisePoSt notification will be ignored and return success. 

Sampling which miners to surprise works as follows:
```text
// A number of challenged miners is chosen at every round
challNumber = NumMiners / ProvingPeriod

// Using the ticket to seed randomness, a miner is picked from the power table for each challenge
sampledMiners = []
while len(sampledMiners) < challNumber:
    ranHash = H(ticket, i)
    ranIndex = HashToInt(ranHash) mod len(PowerTable)
    chosenMiner = PowerTable[ranIndex].address
    // a miner should only be challenged if
    // - they have not submitted a post in SURPRISE_NO_CHALLENGE_PERIOD epochs 
    // - they are not currently challenged 
    // - they are not already chosen in this epoch
    if !sampledMiners.contains(chosenMiner) &&chosenMiner.shouldChallenge():
        sampledMiners.append(chosenMiner)
```

The surprise process described above is triggered by the cron actor in the `Storage Power Actor`. A miner should only be chosen if:

- they are not currently being challenged,
- they are not already selected for a challenge this round,
- their last PoSt (election or surprise) is older than `SURPRISE_NO_CHALLENGE_PERIOD` epochs.

Once issued a surprise, a miner will have to challenge the sectors in their `ProvingSet`, generating a number of partial tickets determined by a sampling rate (`SPOST_SAMPLE_RATE`). In the system here, we use the same system calls as above:

- `GenerateCandidates` which takes in the miner's sectors and a wanted number of Challenge Tickets and generates a number of inclusion proofs for a number of challenged sectors chosen randomly in proportion to the requested number of challenged tickets.
- `GeneratePoSt` takes a set of ChallengeTickets and generates a __*Proof of Spacetime*__ for them, proving the miner storage as a whole.

## PoStSurprise Response

The miner's response must be a PoSt proof over the PartialTickets generated by the challenge which clear a given target `SURPRISE_TARGET` to be defined.

```
submitSurpriseTickets(partialTickets):
    eligibleTickets = []
    for tix in partialTickets:
        challTix = H(tix) // finalization
        if challTix / 2^len(H) < surprise_target:
            eligibleTickets += tix
    return eligibleTickets
```

If the miner responds to the challenge with a SurprisePoSt within `CHALLENGE_DURATION` epochs, they will keep/recover their power.

## Faults

## Fault Detection

Fault detection happens over the course of the life time of a sector. When the sector is unavailable for some reason, the miner must submit the known `faults` in order for that sector not to be tested before the PoSt challenge begins.
Only faults which have been reported at challenge time, will be accounted for. If any other faults have occured the miner will likely fail to submit a valid PoSt for this proving period. Moreover:

- Faults cannot be declared during a ChallengePeriod
- Faults cannot be recovered during a ChallengePeriod

Accordingly, if the miner does not respond to the challenge, they will lose all their Power and a portion of their pledge collateral. This is considered a `DetectedFault` and all sectors in the `ProvingSet` will be marked as `Failing`. The miner will get challenged again in the next proving period. If the miner does not provide a valid response to `MAX_CONSECUTIVE_FAULTS` challenges in a row, their pledge collateral is slashed and their sectors are permanently terminated. Their storage deal collateral is slashed accordingly (see [Storage Deal States](\missing-link) for more).

Any faulted sectors will not count toward miner power in [Expected Consensus](\missing-link). Through these `Detected` and `Declared` faults, the power table should closely track power in the network.

## Fault Recovery

In order to recover from faults (and make the faulted sectors active once more), a miner must mark the faults as `recovering` and then submit a PoSt proving the recovering sectors.
When such a PoSt proof is successfully submitted all faults are reset and assumed to be recovered. A miner must either resolve a failing sector and accept challenges against it in the next proof submission or fail to recover the failing sector within a proving period and the FaultCount on the sector will be incremented. Sectors that have been in the Failing state after `MAX_CONSECUTIVE_FAULTS` challenges will be terminated and result in a `TerminatedFault`.

**Note**: It is important that all faults are known (i.e submitted to the chain) prior to challenge generation, because otherwise it would be possible to know the challenge set, before the actual challenge time. This would allow a miner to report only faults on challenged sectors, with a gurantee that other faulty sectors would not be detected.

## Fault Penalization

Each reported fault carries a penality with it.

**TODO**: Define the exact penalty structure for this.

# Miner Onboarding

Storage Power Consensus participants are subject to a [Minimum Miner Size](\missing-link), meaning miners smaller than `MIN_MINER_SIZE_STOR` of active (or in-deal) storage cannot produce valid electionPoSts. 

These miners' power does not count as part of the total network power, nor are they able to sumit electionPoSts but they can still run and transmit SurprisePosts as messages to be added on-chain. These miners can also be faulted as usual for lacking to prove their power after a challenge. 

New miners are expected to be onboarded through SurprisePoSts. Once they reach the requisite size, their power will be included in total power and they will be able to submit new blocks with ElectionPoSts as well.
