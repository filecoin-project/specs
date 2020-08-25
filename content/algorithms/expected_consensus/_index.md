---
title: "Expected Consensus"
weight: 1
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Expected Consensus
---

TODO: remove stale .id/.go files

## Algorithm

Expected Consensus (EC) is a probabilistic Byzantine fault-tolerant consensus protocol. At a high level, it operates by running a leader election every epoch in which, on expectation, a set number of participants may be eligible to submit a block. EC guarantees that these winners will be anonymous until they reveal themselves by submitting a proof of their election, the `ElectionProof`. Each winning miner can submit one such proof per round and will be rewarded accordingly.

Each proof can be derived from a [properly formatted beacon entry](randomness), as described below. All valid blocks submitted in a given round form a `Tipset`. Every block in a Tipset adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on how to select the heaviest chain, see [Chain Selection](expected_consensus#chain-selection).

While on expectation at least one block will be generated at every round, in cases where no one finds a block in a given round, a miner can simply run leader election again for the next epoch with the appropriate random seed, thereby ensuring liveness in the protocol.

The [Storage Power Consensus](storage_power_consensus) subsystem uses access to EC to use the following facilities:

- Running and verifying [leader election](expected_consensus#secret-leader-election) for block generation.
- Access to a weighting function enabling [Chain Selection](expected_consensus#chain-selection) by the chain manager.
- Access to the most recently [finalized tipset](expected_consensus#finality-in-ec) available to all protocol participants.

{{<embed src="../../systems/filecoin_blockchain/storage_power_consensus/expected_consensus.id"  lang="go" >}} 
{{<embed src="../../systems/filecoin_blockchain/storage_power_consensus/expected_consensus.go"  lang="go" >}}

## Secret Leader Election

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their power. In the case of Filecoin, participants and powers are drawn from the [The Power Table](storage_power_actor#the-power-table), where power is equivalent to storage provided through time.

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of randomness used to run the election. In the case of Filecoin's EC, the blockchain uses [Beacon Entries](storage_power_consensus#beacon-entries) provided by a [drand](drand) beacon. These seeds are used as unbiasable randomness for Leader Election. Every block header contains an `ElectionProof` derived by the miner using the appropriate seed.

### Running a leader election

Design goals here include:

- Miners should be rewarded proportional to their power in the system
- The system should be able to tune how many blocks are put out per epoch on expectation (hence "expected consensus").

At a high-level, leader election works as follows:

- A miner draws an appropriate random seed for this epoch
- They generate an `ElectionProof` from this seed by using a VRF to generate a signature over the seed at a given epoch, as defined in [Randomness](randomness). The `ElectionProof` is the VRF Proof.
  - If the `ElectionProof`'s normalized digest value (i.e. the VRF Digest) is below a miner-specific `ElectionTarget` determined by the miner's power in SPC, it is valid: the miner can craft a block (see [Block Producer](block_producer)). 
    - The `ElectionProof` Digest must be normalized according to both the `ExpectedLeadersPerEpoch` in EC (a network parameter) and the maximum value of the digest (based on the hash used to produce it).
    - The `ElectionTarget` is set as the proportion of the miner's `TotalQualityAdjPower` over the total network quality-adjusted power. Accordingly leader election in EC is proportional to miner power.
  - Otherwise, the miner tries again in the next epoch.

Conceptually, EC yields winners proportionally to their power since it enables each miner to generate a single `ElectionProof` whose digest can be normalized to yield a uniformly-drawn random number between 0 and 1. It is then compared to the miner's power in proportion to total network quality-adjusted power (i.e. also between 0 and 1). The more powerful the miner, the more frequently their `ElectionProof`s will be valid.

We show this below, removing division for ease of implementation:

We have:
```go 
GenerateElectionProof(epoch) {
  electionProofInput := GetRandomness(DomainSeparationTag_ElectionProofProduction, epoch, CBOR_Serialize(miner.address))
  vrfResult := miner.VRFSecretKey.Generate(electionProofInput)

  if IsValidElectionProof(vrfResult.Digest) {
    return vrfResult.Proof
  }
  return nil
}

IsValidElectionProof(proofDigest, miner) {
  // for SHA256, more generally it is 2^len(H)
  const maxDigestSize = 2^256
  // normalizedDigest := proofDigest / maxDigestSize * EC.ExpectedLeadersPerEpoch()
  // ElectionTarget := SPC.GetTotalQualityAdjPower(miner) / SPC.GetNetworkTotalQualityAdjPower()
  // We check that normalizedDigest < ElectionTarget, with
  // For ease of implementation we remove divisions from the above:
  return proofDigest * SPC.GetNetworkTotalQualityAdjPower() * EC.ExpectedLeadersPerEpoch() < maxDigestSize * SPC.GetTotalQualityAdjPower(miner)
}
```

### Election Validation

In order to determine that the mined block was generated by an eligible miner, one must check its `ElectionProof`'s validity, that it was generated using the appropriate beacon entry, and that it is valid according to the miner's `ElectionTarget`, per the above definition.

## Chain Selection

Just as there can be 0 miners win in a round, multiple miners can be elected in a given round. This in turn means multiple blocks can be created in a round, as seen above. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round.

### Chain Weighting

It is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the most storage provided.

In short, the weight at each block is equal to its `ParentWeight` plus that block's delta weight. Details of Filecoin's chain weighting function [are included here](https://observablehq.com/d/3812cd65c054082d).

Delta weight is a term composed of a few elements:

- wPowerFactor: which adds weight to the chain proportional to the total power backing the chain, i.e. accounted for in the chain's power table.
- wBlocksFactor: which adds weight to the chain proportional to the number of tickets mined in a given epoch. It rewards miner cooperation (which will yield more blocks per round on expectation).

The weight should be calculated using big integer arithmetic with order of operations defined above. We use brackets instead of parentheses below for legibility. We have:

```text
w[r+1] = w[r] + (wPowerFactor[r+1] + wBlocksFactor[r+1]) * 2^8
```

For a given tipset `ts` in round `r+1`, we define:

- `wPowerFactor[r+1]  = wFunction(totalPowerAtTipset(ts))`
- `wBlocksFactor[r+1] =  wPowerFactor[r+1] * wRatio * t / e`
  - with `t = |ticketsInTipset(ts)|`
  - `e = expected number of tickets per round in the protocol`
  - and `wRatio in ]0, 1[`
Thus, for stability of weight across implementations, we take:
- `wBlocksFactor[r+1] =  (wPowerFactor[r+1] * b * wRatio_num) / (e * wRatio_den)`

We get:

```text
w[r+1] = w[r] + wFunction(totalPowerAtTipset(ts)) * 2^8 + (wFunction(totalPowerAtTipset(ts)) * len(ts.tickets) * wRatio_num * 2^8) / (e * wRatio_den)
```

Using the 2^8 here to prevent precision loss ahead of the division in the wBlocksFactor.

The exact value for these parameters remain to be determined, but for testing purposes, you may use:

 - `e = 5`
 - `wRatio = .5, or wRatio_num = 1, wRatio_den = 2`
- `wFunction = log2b` with
  - `log2b(X) = floor(log2(x)) = (binary length of X) - 1` and `log2b(0) = 0`. Note that that special case should never be used (given it would mean an empty power table).

{{< hint warning >}}
**Note that if your implementation does not allow for rounding to the fourth decimal**, miners should apply the [tie-breaker below](selecting-between-tipsets-with-equal-weight). Weight changes will be on the order of single digit numbers on expectation, so this should not have an outsized impact on chain consensus across implementations.
{{< /hint >}}

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given Tipset should have
the same `ParentWeight` value) plus the delta weight of each parent. To make the
computation a bit easier, a block's `ParentWeight` is stored in the block itself (otherwise
potentially long chain scans would be required to compute a given block's weight).

### Selecting between Tipsets with equal weight

When selecting between Tipsets of equal weight, a miner chooses the one with the smallest ticket (see [Tickets](storage_power_consensus#tickets)).

In the case where two Tipsets of equal weight have the same min ticket, the miner will compare the next smallest ticket in the Tipset (and select the Tipset with the next smaller ticket). This continues until one Tipset is selected.

The above case may happen in situations under certain block propagation conditions. Assume three blocks B, C, and D have been mined (by miners 1, 2, and 3 respectively) off of block A, with minTicket(B) < minTicket(C) < minTicket(D).

Miner 1 outputs their block B and shuts down. Miners 2 and 3 both receive B but not each others' blocks. We have miner 2 mining a Tipset made of B and C and miner 3 mining a Tipset made of B and D. If both succesfully mine blocks now, other miners in the network will receive new blocks built off of Tipsets with equal weight and the same smallest ticket (that of block B). They should select the block mined atop `[B, C]` since `minTicket(C) < minTicket(D)`.

The probability that two Tipsets with different blocks would have all the same tickets can be considered negligible: this would amount to finding a collision between two 256-bit (or more) collision-resistant hashes.

## Finality in EC
EC enforces a version of soft finality whereby all miners at round N will reject all blocks that fork off prior to round N-F. For illustrative purposes, we can take F to be 900. While strictly speaking EC is a probabilistically final protocol, choosing such an F simplifies miner implementations and enforces a macroeconomically-enforced finality at no cost to liveness in the chain.

## Consensus Faults

Due to the existence of potential forks in EC, a miner can try to unduly influence protocol fairness. This means they may choose to disregard the protocol in order to gain an advantage over the power they should normally get from their storage on the network. A miner should be slashed if they are provably deviating from the honest protocol.

This is detectable when a given miner submits two blocks that satisfy any of the following "consensus faults". In all cases, we must have:

- both blocks were mined by the same miner
- both blocks have valid signatures
- first block's epoch is smaller or equal than second block

### Types of faults

1. **Double-Fork Mining Fault**: two blocks mined at the same epoch (even if they have the same tipset).
   - `B4.Epoch == B5.Epoch` 
![Double-Fork Mining Fault](diagrams/double_fork.dot)


2. **Time-Offset Mining Fault**: two blocks mined off of the same Tipset at different epochs.
   - `B3.Parents == B4.Parents && B3.Epoch != B4.Epoch`
![Time-Offset Mining Fault](diagrams/time_offset.dot)

3. **Parent-Grinding Fault**: one block's parent is a Tipset that provably should have included a given block but does not. While it cannot be proven that a missing block was willfully omitted in general (i.e. network latency could simply mean the miner did not receive a particular block), it can when a miner has successfully mined a block two epochs in a row and omitted one. That is, this condition should be evoked when a miner omits their own prior block.
Specifically, this can be proven with a "witness" block, that is by submitting blocks B2, B3, B4 where B2 is B4's parent and B3's sibling but B3 is not B4's parent.
    - `!B4.Parents.Include(B3) && B4.Parents.Include(B2) && B3.Parents == B2.Parents && B3.Epoch == B2.Epoch`
![Parent-Grinding fault](diagrams/parent_grinding.dot)

### Penalization for faults
A single consensus fault results into:
- miner termination and removal of power from the power table,
- loss of all pledge collateral (which includes the initial pledge and blocks rewards yet to be vested)

### Detection and Reporting

A node that detects and report a consensus fault is called "slasher", any user in Filecoin can be a slasher. They can report consensus faults by calling the `ReportConsensusFault` on the `StorageMinerActor` of the faulty miner. The slasher is rewarded with a portion of the offending miner's [Pledge Collateral](storage_power_actor#pledge-collateral) for notifying the network of the consensu fault.

The reward give to the slasher is a function of some initial share (`SLASHER_INITIAL_SHARE`) and growth rate (`SLASHER_SHARE_GROWTH_RATE`) and it has a maximum `maxReporterShare`. Slasher's share increases exponentially as epoch elapses since the block when the fault is committed (see `RewardForConsensusSlashReport`). Only the first slasher gets their share of the pledge collateral and the remaining pledge collateral is burned. The longer a slasher waits, the higher the likelihood that the slashed collateral will be claimed by another slasher.
