---
title: Election PoSt
---

This document describes Election-PoSt, the Proof-of-Spacetime used in Filecoin.

# High Level API

Election PoSt couples the PoSt process with block production, meaning that in order to produce a block, the miner must produce a valid PoSt proof (snark output). Specifically, a subset of non-faulted sectors the miner is storing (i.e. eligible sectors) allows them to attempt a leader election using a ChallengeTicket any of which could yield a valid ElectionProof. The probability of scratching a winning ChallengeTicket is dependent on sector size and total storage. Miners are rewarded in proportion to the quantity of winning ChallengeTickets they generate in a given epoch, thereby incentivizing a miner to check as much of their storage as allowed in order to express their full power in a leader election.

This does mean that, in a given round, a lucky miner may succeed in generating a block by proving only a single sector, but again on expectation, a miner will have to prove their sectors at every single round in order for their full power to contribute to block generation. In the event one of the miner’s sectors cannot be proven (i.e. the miner does not have access to the nodes from that sector), no tickets will be returned. In order to prevent this, a miner can declare faults on their faulty sectors to avoid having to include them in the eligible sector set. Their power will be reduced accordingly. 

# ElectionPoSt

To enable short PoSt response time, miners are required submit a PoSt when they are elected to mine a block, hence PoSt on election or ElectionPoSt. When miners win a block, they need to immediately generate a PoStProof and submit that along with the ElectionProof. Both the ElectionProof and PoStProof are checked at Block Validation by `StoragePowerConsenusSubsystem`. When a block is included, a special message is added that calls `SubmitElectionPoSt` which will process sector updates in the same way as successful `SubmitSurprisePoSt` do.

# SurprisePoSt

The chain keeps track of the last time that a miner received and submitted a Challenge (Election or Surprise). It randomly challenges a miner to submit a surprisePoSt once per `ProvingPeriod` in the latter half of their `ProvingPeriod`. For every miner challenged, a `NotifyOfPoStSurpriseChallenge` is issued and sent as an on-chain message to surprised `StorageMinerActor`. However, if the `StorageMinerActor` is already proving a SurprisePoStChallenge (`IsChallenged` is True) or the `StorageMinerActor` has received a challenge (by Election or Surprise) within the `MIN_CHALLENGE_PERIOD` (`ShouldChallenge` is False), then the PoSt surprise notification will be ignored and return success.

For more on these components, see {{<sref election-post>}}. At a high-level though both are needed for different reasons:
- By coupling leader election and PoSt, `ElectionPoSt` ensures that miners must do the work to prove their sectors at every round in order to earn block rewards.
- Small miners may not win on a regular basis however, `SurprisePoSt` thus comes in as a lower-bound to how often miners must PoSt and helps ensure the Power Table does not grow stale for its long-tail of smaller miners.

## ElectionPoSt Generation

Filecoin's ElectionPoSt process makes use of two calls to the system library:
- `GenerateCandidates` which takes in the miner's sectors and a wanted number of Challenge Tickets and generates a number of inclusion proofs for a number of challenged sectors chosen randomly in proportion to the requested number of challenged tickets.
- `GeneratePoSt` takes a set of ChallengeTickets and generates a __*Proof of Spacetime*__ for them, proving the miner storage as a whole.

As stated above, a miner is incentivized to repeat this process at every block time in order to check whether they were elected leaders (see <<sref expected_consensus>>). The rationality assumption made by ElectionPoSt is thus that storing files continuously and earning block rewards accordingly will be more profitable to miners than regenerating data at various epochs to sporadically participate in leader election.

At every epoch, each miner will challenge a portion of sectors at random proportional to sectorsSampled, with each sector being issued K PoSt challenges (coverage may not be perfect).

By proving access to the challenged range of nodes (i.e. merkle tree leaf from the committed sector) in the sector, the miner can generate a set of valid ChallengeTickets in order to check them as part of leader election in EC (in order to find the winning tickets, or ElectionProofs). The winning tickets will be stored on the block and used to generate a PoSt (using a SNARK). A block header will thus contain a number of “winning” ChallengeTickets (each containing a SectorID and other elements, used to derive ElectionProofs) and a PostProof generated from the ChallengeTickets.

In order to simplify implementation and bound block header size, we can set a maximum number of possible election proofs for any block. For instance, for EC.E=5, we can cap challengeTicket submissions at 16 per block, which would cover more than 99.99% of cases (using Chernoff bounds) encountered by a 50% miner (i.e. much more in practice).

The epoch time is divided into three non-overlapping portions: (1) checking sectors for small challenge tickets, (2) computing the SNARK to produce a PoSt certifying the winning challenge tickets, and (3) block production/propagation. The target duration for each is still to be tuned. Proof parameter adjustments might allow a longer block propagation window while maintaining PoSt security (but trades off against sealing time and other things).

If no one has found a winning ticket in this epoch, increment epoch value as part of the post_randomness sampling and try again.

```md
At every round:
1. **(sample randomness)**
The miner draws a randomness ticket from the randomness chain from a given epoch SPC.post_lookback back and concats it with the minerID, and epoch for use as post_randomness (this is as it used to happen previously, except the PoSt_lookback is now also the EC randomness_lookback, which is small here):
    - `post_randomness = VRF_miner(ChainRandomness(currentBlockHeight - SPC.post_lookback))`

2. **(select eligible sectors)**
The miner calls `GenerateCandidates` from proofs with their non-faulted (declared or detected) sectors, meaning those in their `proving set` (from the `StorageMinerActor`) along with a chosen number of `PartialTickets` (`sectorsSampled*eligibleSectors`).
    - Select a subset of sectors of size `sectorsSampled*minerSectors` [details omitted]

    Note also that even with `challengeTicketNum == numSectors`, this process may not sample all of a miner’s sectors in any given epoch, due to collisions in the pseudorandom sector ID selection process.

3. **(generate Partial Ticket(s))** for each selected sector
    - first, generate K `PoStChallenges` (C_i), sampled at random from the chosen sector.
    There will be a chosen challenge-range size (a power of 2), called `ChallengeRangeSize`. Challenge ranges must be subtree-aligned and thus divide the sector into fixed-size data blocks. So challenge ranges can be indexed because allowable challenge ranges do not overlap. C_i is selected from the valid indexes into the challenge ranges.
    `C_i = HashToBlockNumber(post_randomness || S_j || i)`
    C_i corresponds to a given index of a known-sized data block (the sector)
    - for each C_i, miner reads the specified range of data from disk, which is of size `ChallengeRangeSize : C_i_output`
    - Miner generates a PartialTicket using all the PoSt witnesses
    - `PartialTicket = H(post_randomness || minerID ||S_ j || C_1_Output || … || C_K_Output)`

4. **(check Challenge Ticket(s) for winners)**
Given a returned PartialTicket, miner checks it is a winning ticket. Specifically, they do the following:
    - `ChallengeTicket = Finalize(PartialTicket).  = H(ChallengeTicket) / 2^len(H)`
    - Check that `ChallengeTicket < Target`
    - If yes, it is a winning ticket and can be used to submit a block
    - In either case, try again with next sector to increase rewards

5. **(generate a `PoStProof`)** for inclusion in the block header
    - Using the winning tickets from step 4 (there may be multiple tickets from the same `SectorNumber`), call `GeneratePoSt` from proofs to generate a single PoSt for the block

If no one has found a winning ticket in this epoch, increment epoch value as part of the post_randomness sampling and try again.
```

**Parameters:**
```md
- `Randomness_ticket` --  a ticket drawn from the randomness ticket chain at a prior tipset 
- `Randomness_lookback` -- how far back to draw randomness from the randomness ticket chain - it will be as large as allowed by PoSt security, likely 1 or 2 epochs
- `K (e.g. 20-100s)` - number of  challenges per sector -- must be large enough such that the PoSpace is secure.
- `ChallengeRangeSize` - challenge read size (between 32B and 256KB)  -- based on security analysis.
- `sectorsSampled` - sector sampling fraction (e.g. 1, .10, .04) -- 1 to start-- It should be large enough to make it irrational to fully regenerate sectors. We may choose some subset if cost of verifying all is deleterious to disk
- `Target` -- target value under which PoSt value must be for block creation -- `target = activePowerInSector/networkPower * sectorsSampled * EC.ExpectedLeaders`.
Put another way check `challengeTicket * networkPower * sectorsSampled_denom < activePowerInSector * sectorsSampled_num * EC.ExpectedLeaders`
- `networkPower` - filecoin network’s power - read from the power table, expressed in number of bytes
```

## ElectionPoSt verification

At a high-level, to validate a winning PoSt Proof:
```md
1. **(Verify post_randomness):**
    - `post_randomness` is appropriate given `chain`, `post_lookback` and `epoch`
    - VRF is generated by miner who submitted the PoSt:
    `{0,1} ?= VRF_miner.Verify(post_randomness)`
    - VRF output is valid given expected inputs:
    `{0,1} ?= VRF_miner.Validate(H(randomness_ticket || chainEpoch))`

2. Rederive eligible sectors in order to verify that winning sectors were appropriately selected (meaning it’s in `S_j`) from
    `{0, 1} ?= sectorID in (HashToSectorID(post_randomness || minerID || j))`

3. **Verify `PartialTicket` returned**
    - Prove that the PartialTicket were appropriately derived from the eligible sectors, by submitting all miner sectors along with the wanted number of tickets and verifying that the outputted PartialTicket match.

4. **Derive and validate the `ChallengeTicket` from the PartialTickets**
    - Prove the derived ChallengeTicket is below the target
    `{0, 1} ?= (ChallengeTicket < target)`

5. **Verify the `postProof` using the `PartialTickets`**
    - Verify that the postProof was appropriately derived from the PartialTickets.
    - The PoSt proof will verify the correctness of any PartialTickets passed to it as public inputs. In order to do this, it also needs the sector number along with various on-chain data (randomness, CommR, miner ID, etc.).
```

As it stands, the proofs caller passes a list of all of the replicas (with filesystem paths and auxiliary metadata) which could be challenged and a wanted number of tickets. This code structure forces the miner architecture to look like a single machine with a large number of disks presented as a single filesystem, as shown in spec. That won’t change for now. VFS systems like NFS may be used to distribute disks among multiple machines virtualizing the FileStore.

There is no requirement to persist the witnesses (list of merkle proofs) for failure recovery. If something fails, we don’t expect recovery to be fast enough to allow the miner to participate in that round anyway.

## Surprise PoSt cleanup

But while this means a miner will never win blocks from faked power, they won’t be penalized either when storage is lost. How do we ensure that the Power Table is accurate? Likewise, how do we onboard new power since it will not win a block unless it has at least X TB or Y % of the network.

This is why we need PoSt surprise: 

Atop ElectionPost, a miner will be surprised with a PoSt challenge in every ProvingPeriod (~2 days). The ProvingPeriod resets whenever the miner publishes a PoSt (election or surprise). 
This SurprisePoSt will use a challenge drawn from the chain at the start of this recovery period. Its PoStProof must be a proof over the PartialTickets for all sectors a miner is storing (i.e. a miner must submit a PoStProof made up of all partialTickets for all sectors in the `ProvingSet`, not just the winning ones on sampled sectors). For this reason a miner is incentivized to declare faults in order to successfully generate this PoStProof.

Upon receiving a PoSt surprise challenge, a miner has a given ChallengePeriod (~2 hours) past which they, they will lose all their power and a portion of their pledge collateral if they have not submitted a PoSt on chain. This is considered a `DetectedFault` and all sectors in the `ProvingSet` will be marked as `Failing`. No deal collateral will be slashed if miners can recover within the next three proving periods. Note that the exact amount of slashed pledge collateral and pledge collateral function itself are subject to change.

Thereafter, the miner will have three more ProvingPeriods (specifically three challenges, so three ProvingPeriods on expectation) to recover their power by submitting a SurprisePoSt. If the miner does not do so, their sectors are permanently terminated and their storage deal collateral slashed (see the StorageMinerActor in the spec).

Faults are largely independent from this process. The only rules are:
- Faults cannot be declared during a ChallengePeriod
- Faults cannot be recovered during a ChallengePeriod

Note that this PoStSurprise will frequently be used by small miners who do not win blocks, and by miners as they are onboarding power to the network, since they will not be able to win any blocks (they have no power) and hence gain power through ElectionPoSt.

Miners earn no reward from submitting PoStSurprise messages. This mechanism does forces a miner to announce their faults and clean up the Power Table accordingly.

Surprise PoSt works as follows:
```text
// A number of challenged miners is chosen at every round
challNumber = CHALLENGE_FREQUENCY*NumMiners / ProvingPeriod

// Using the ticket to seed randomness, a miner is picked from the power table for each challenge
sampledMiners = []
For i=0; i < challNumber; i++:
    ranHash = H(ticket, i)
    ranIndex = HashToInt(ranHash) mod len(PowerTable)
    chosenMiner = PowerTable[ranIndex].address
    // a miner should only be challenged if they have not submitted a post in ProvingPeriod/CHALLENGE_FREQUENCY epochs and are not currently challenged
    if chosenMiner.shouldChallenge(ProvingPeriod/CHALLENGE_FREQUENCY):
        sampledMiners.append(chosenMiner)
```

The surprise process described above is triggered by the cron actor in the storage_power_actor (through which the power table is searched for challengeable miners). A miner should be getting randomly sampled twice per proving period on expectation, but would only be sampled if they are in the latter half of their proving period leading to one challenge per proving period on expectation.
This is done as follows: if there are M miners in the power table and a Proving Period of length P, 2M/P challenges will be issued at each epoch. Miners are sampled using a randomness ticket from the chain and will only be challenged if they have not submitted a PoSt in at least PP/2 epochs and are not currently being challenged (this is checked using the storage_miner_actor).

An alternative approach would be to assign a probability of being challenged to each miner which grows at every epoch to be 1 PP epochs from the last challenge (but this would require more computation since every miner would have to be checked at every epoch).

## Fault Detection

Fault detection happens over the course of the life time of a sector. When the sector is unavailable for some reason, the miner must submit the known `faults` in order for that sector not to be tested before the PoSt challenge begins.
Only faults which have been reported at challenge time, will be accounted for. If any other faults have occured the miner will likely fail to submit a valid PoSt for this proving period.

Any faulted sectors will not count toward miner power in {{<sref expected_consensus>}}.

## Fault Recovery

In order to recover from faults (and make the faulted sectors active once more), a miner must mark the faults as `recovering` and then submit a PoSt proving the recovering sectors.
When such a PoSt proof is successfully submitted all faults are reset and assumed to be recovered. A miner must either resolve a failing sector and accept challenges against it in the next proof submission or fail to recover the failing sector within a proving period and the FaultCount on the sector will be incremented. Sectors that have been in the Failing state for more than `MaxFaultCount` consecutive epochs will be terminated and result in a `TerminatedFault`.

{{% notice note %}}
**Note**: It is important that all faults are known (i.e submitted to the chain) prior to challenge generation, because otherwise it would be possible to know the challenge set, before the actual challenge time. This would allow a miner to report only faults on challenged sectors, with a gurantee that other faulty sectors would not be detected.
{{% /notice %}}

## Fault Penalization

Each reported fault carries a penality with it.

{{% notice todo %}}
**TODO**: Define the exact penality structure for this.
{{% /notice %}}

Note that surprise PoSt will frequently be used by small miners who do not win blocks, and by miners as they are onboarding power to the network, since they will not be able to win any blocks (they have no power) and hence gain power through ElectionPoSt.

While it resets a miner’s ElectionWindow, miners earn no reward from submitting SurprisePoSt messages. This mechanism does forces a miner to announce their faults and clean up the Power Table accordingly.
