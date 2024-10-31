---
title: Expected Consensus
weight: 1
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Expected Consensus

## Algorithm

Expected Consensus (EC) is a probabilistic Byzantine fault-tolerant consensus protocol. At a high level, it operates by running a leader election every epoch in which, on expectation, a set number of participants may be eligible to submit a block. EC guarantees that these winners will be anonymous until they reveal themselves by submitting a proof that they have been elected, the `ElectionProof`. Each winning miner can submit one such proof per round and will be rewarded proportionally to its power. From this point on, each winning miner also creates a proof of storage (aka Winning PoSt). Each proof can be derived from a [properly formatted beacon entry](randomness), as described below.

All valid blocks submitted in a given round form a `Tipset`. Every block in a Tipset adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on how to select the heaviest chain, see [Chain Selection](expected_consensus#chain-selection). While on expectation at least one block will be generated at every round, in cases where no one finds a block in a given round, a miner can simply run leader election again for the next epoch with the appropriate random seed, thereby ensuring liveness in the protocol.

The randomness used in the proofs is generated from [DRAND](https://drand.love), an unbiasable randomness generator, through a beacon. When the miner wants to publish a new block, they invoke the `getRandomness` function providing the chain height (i.e., epoch) as input. The randomness value is returned through the DRAND beacon and included in the block. For the details of DRAND and its implementation, please consult the project's [documentation](https://drand.love/docs/overview/) and [specification](https://drand.love/docs/specification/).

The [Storage Power Consensus](storage_power_consensus) subsystem uses access to EC to use the following facilities:

- Access to verifiable randomness for the protocol, derived from [Tickets](storage_power_consensus#tickets).
- Running and verifying [leader election](expected_consensus#secret-leader-election) for block generation.
- Access to a weighting function enabling [Chain Selection](expected_consensus#chain-selection) by the chain manager.
- Access to the most recently [finalized tipset](expected_consensus#finality-in-ec) available to all protocol participants.

## Tickets in EC

There are two kinds of tickets:

1. ElectionProof ticket: which is the VRF that runs based on DRAND input. In particular, the miner gets the DRAND randomness beacon and gives it as input to the VRF together with the miner's worker's key.
2. the ticket is generated using the VRF as above, but the input includes the concatenation of the previous ticket. This means that the new ticket is generated running the VRF on the old ticket concatenated with the new DRAND value (and the key as before).

```text
## Get Randomness value from  DRAND beacon, by giving the current epoch as input.
Drand_value = GetRandmness(current epoch)

## Create ElectionProof ticket by calling VRF and  giving the secret key of the miner's worker and the DRAND value obtained in the previous step
Election_Proof = VRF(sk, drand_value)

## Extend the VRF ticket chain by concatenating the previous proof/ticket with the current one by following the same process as above (i.e., call VRF function with the secret key of the miner's worker and the DRAND value of the current epoch).
VRF chain: new_ticket = VRF(sk, drand_value || previous ticket)
```

Within Storage Power Consensus (SPC), a miner generates a new ticket for every block on which they run a leader election. This means that the ticket chain is always as long as the blockchain.

Through the use of VRFs and thanks to the unbiasable design of DRAND, we achieve the following two properties.

- _Ensure leader secrecy:_ meaning a block producer will not be known until they release their block to the network.
- _Prove leader election:_ meaning a block producer can be verified by any participant in the network.

## Secret Leader Election

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their power. In the case of Filecoin, participants and powers are drawn from the [The Power Table](storage_power_actor#the-power-table), where power is equivalent to storage provided over time.

Leader Election in Expected Consensus must be _Secret, Fair and Verifiable_. This is achieved through the use of randomness used to run the election. In the case of Filecoin's EC, the blockchain uses [Beacon Entries](storage_power_consensus#beacon-entries) provided by a [drand](drand) beacon. These seeds are used as unbiasable randomness for Leader Election. Every block header contains an `ElectionProof` derived by the miner using the appropriate seed. As noted earlier, there are two ways through which randomness can be used in the Filecoin EC: i) through the ElectionProof ticket, and ii) through the VRF ticket chain.

### Running a leader election

The miner whose block has been submitted must be checked to verify that they are eligible to mine a block in this round, i.e., they have not been slashed in the previous round.

Design goals here include:

- Miners should be rewarded proportional to their power in the system
- The system should be able to tune how many blocks are put out per epoch _on expectation_ (hence "expected consensus").

A miner will use the ElectionProof ticket to uniformly draw a value from 0 to 1 when crafting a block.

#### Winning a block

**Step 1:** Check for leader election

A miner checks if they are elected for the current epoch by running `GenerateElectionProof`.

Recall that a miner is elected proportionally to their quality adjusted power at `ElectionPowerTableLookback`.

A requirement for setting `ElectionPowerTableLookback` is that it must be larger than finality. This is because if `ElectionPowerTableLookback` is shorter, a malicious miner could create sybils with different VRF keys to increase the chances of election and then fork the chain to assign power to those keys.

The steps of this well-known attack in Proof of Stake systems would be:

1. The miner generates keys used in the VRF part of the election until they find a key that would allow them to win.
2. The miner forks the chain and creates a miner with the winning key.

This is generally a problem in Proof of Stake systems where the stake table is read from the past to make sure that no staker can do a transfer of stake to a new key that they found to be winning.

**Step 2:** Generate a storage proof (WinningPoSt)

An elected miner gets the randomness value through the DRAND randomness generator based on the current epoch and uses it to generate WinningPoSt.

WinningPoSt uses the randomness to select a sector for which the miner must generate a proof. If the miner is not able to generate this proof within some predefined amount of time, then they will not be able to create a block. The sector is chosen from the power table `WinningPoStSectorSetLookback` epochs in the past.

Similarly to `ElectionPowerTableLookback`, a requirement for setting `WinningPoStSectorSetLookback` is that it must be larger than finality. This is to enforce that a miner cannot play with the power table and change which sector is challenged for WinningPoSt (i.e., set the challenged sector to one of their preference).

If `WinningPoStSectorSetLookback` is not longer than finality, a miner could try to create forks to change their sectors allocation to get a more favourable sector to be challenged for example. A simple attack could unfold as follows:

- The power table at epoch `X` shows that the attacker has sectors 1, 2, 3, 4, 5.
- The miner decides to not store sector 5.
- The miner wins the election at epoch `X`.
  - Main fork: Miner is asked a WinningPoSt for sector 5 for which they won't be able to provide a proof.
  - The miner creates a fork, terminates sector 5 in epochs before `X`.
  - At `X`, the miner is now challenged a different sector (not 5).

Note that there are variants of this attack in which the miner introduces a new sector to change which sector will be challenged.

> What happens if a sector expired after `WinningPoStSectorSetLookback`?

An expired sector will not be challenged during WindowPoSt (hence not penalized after its expiration). However, an edge case around the fact that `WinningPoStSectorSetLookback` is longer than finality is that due to the lookback, a miner can be challenged for an expired sector between `expirationEpoch` and `expirationEpoch + WinningPoStSectorSetLookback - 1`. Therefore, it is important that miners keep an expired sector for `WinningPoStSectorSetLookback` more epochs after expiration, or they will not be able to generate a WinningPoSt (and get the corresponding reward).

Example:

- At epoch `X`:
  - Sector expires and miner deletes the sector.
- At epoch `X+WinningPoStSectorSetLookback-1`:
  - The expired sector gets selected for WinningPoSt
  - The miner will not be able to generate the WinningPoSt and they will not win the sector.

**Step 3:** Block creation and propagation

If the above is successful, miners build a block and propagate it.

### GenerateElectionProof

`GenerateElectionProof` outputs the result of whether a miner won the block or not as well as the quality of the block mined.

The "WinCount" is an integer that is used for weight and block reward calculations. For example a WinCount equal to "2" is equivalent as two blocks of quality "1".

#### High level algorithm

- Get the percentage of power at block `ElectionPowerTableLookback`
  - Get the power of the miner at block `ElectionPowerTableLookback`
  - Get the total network power at block `ElectionPowerTableLookback`
- Get randomness for the current epoch using `GetRandomness`.
- Generate a VRF and compute its hash
  - The storage miner's `workerKey` is given as input in the VRF process
- Compute WinCount: The smaller the hash is, the higher the WinCount will be
  - Compute the probability distribution of winning `k` blocks (i.e. Poisson distribution see below for details)
  - Let `h_n` be the normalised VRF, i.e. `h_n = vrf.Proof/2^256`, where `vrf.Proof` is the ElectionProof ticket.
  - The probability of winning one block is `1-P[X=0]`, where `X` is a Poisson random variable following Poisson distribution with parameter `lambda = MinerPowerShare*ExpectedLeadersPerEpoch`. Thus, if `h_n` is less than `1-P[X=0]`, the miner wins at least one block.
  - Similarly if `h_n` is less than `1-P[X=0]-P[X=1]` we have at least two blocks and so on.
  - While it is not permitted for a single miner to publish two distinct blocks, in this case, the miner produces a single block which earns two block rewards

#### Explanations - Poisson Sortition

Filecoin is building on the principle that a miner possessing _X%_ of network power should win as many times as X miners with 1% of network power in the election algorithm.

A straightforward solution to model the situation is using a Binomial distribution with parameter`p=MinerPower/TotalPower` and `n=ExpectedLeadersPerEpoch`. However, given that we effectively want every miner to roll an uncorrelated/independent dice and want to be invariant to miner pooling, it turns out that Poisson is the ideal distribution for our case.

Despite this finding, we wanted to assess the difference between the two distributions in terms of the probability mass function.

Using `lambda = MinerPower*ExepectedLeader` as the parameter for the Poisson distribution, and assuming `TotalPower =10000`, `minerPower = 3300` and `ExpectedLeaderPerEpoch = 5`, we find (see table) that the probability mass function for the Binomial and the Poisson distributions do not differ much anyway.

| k   | Binomial | Poisson |
| --- | -------- | ------- |
| 0   | 0.19197  | 0.19205 |
| 1   | 0.31691  | 0.31688 |
| 2   | 0.26150  | 0.26143 |
| 3   | 0.14381  | 0.14379 |
| 4   | 0.05930  | 0.05931 |
| 5   | 0.01955  | 0.01957 |
| 6   | 0.00537  | 0.00538 |
| 7   | 0.00126  | 0.00127 |
| 8   | 0.00026  | 0.00026 |
| 9   | 0.00005  | 0.00005 |

**Justification for the need of _WinCount_**

It should not be possible for a miner to split their power into multiple identities and have more chances of winning more blocks than keeping their power under one identity. In particular, Strategy 2 below should not be possible to achieve.

- Strategy 1: A miner with X% can run a single election and win a single block.
- Strategy 2: The miner splits its power in multiple sybil miners (with the sum still equal to X%), running multiple elections to win more blocks.

WinCount guarantees that a lucky single block will earn the same reward as the reward that the miner would earn had they split their power into multiple sybils.

**Alternative Options for the Distribution/Sortition**

Bernoulli, Binomial and Poisson distributions have been considered for the _WinCount_ of a miner with power `p` out of a total network power of `N`. There are the following options:

- Option 1: WinCount(p,N) ~ Bernoulli(pE/N)
- Option 2: WinCount(p,N) ~ Binomial(E, p/N)
- Option 3: WinCount(p,N) ~ Binomial(p, E/N)
- Option 4: WinCount(p,N) ~ Binomial(p/M, ME/N)
- Option 5: WinCount(p,N) ~ Poisson(pE/N)

Note that in Options 2-5 the expectation of the win-count grows linearly with the miner's power `p`. That is, `𝔼[WinCount(p,N)] = pE/N`. For Option 1 this property does not hold when p/N > 1/E.

Furthermore, in Options 1, 3 and 5 the \_WinCount distribution is invariant to the number of Sybils in the system. In particular WinCount(p,N)=2WinCount(p/2,N), which is a desirable property.

In Option 5 (the one used in Filecoin Leader Election), the ticket targets for each \_WinCount k that range from 1 to mE (with m=2 or 3) shall approximate the upside-down CDF of a Poisson distribution with rate λ=pE/N, or explicitly, 2²⁵⁶(1-exp(-pE/N)∑ᵏ⁻¹ᵢ₌₀(pE/N)ⁱ/(i!)).

**Rationale for the Poisson Sortition choice**

- Option 1 - Bernoulli(pE/N): this option is easy to implement, but comes with a drawback: if the miner's power exceeds 1/E, the miner's WinCount is always 1, but never higher than 1.
- Option 2 - Binomial(E, p/N): the expectation of WinCount stays the same irrespectively of whether the miner splits their power into more than one Sybil nodes, but the variance increases if they choose to Sybil. Risk-seeking miners will prefer to Sybil, while risk-averse miners will prefer to pool, none of which is a behaviour the protocol should encourage. This option is not computationally-expensive as it would involve calculation of factorials and fixed-point multiplications (or small integer exponents) only.
- Option 3 - Binomial(p, E/N): this option is computationally inefficient. It involves very large integer exponents.
- Option 4 - Binomial(p/M, ME/N): the complexity of this option depends on the value of M. A small M results in high computational cost, similarly to Option 3. A large M, on the other hand, leads to a situation similar to that of Option 2, where a risk-seeking miner is incentivized to Sybil. Clearly none of these are desirable properties.
- Option 5 - Poisson(pE/N): the chosen option presents the implementation difficulty of having to hard-code the co-efficients (see below), but overcomes all of the problems of the previous options. Furthermore, the expensive part, that is calculating exp(λ), or exp(-pE/N) has to be calculated only once.

**Coefficient Approximation**

We have used the Horner rule with 128-bit fixed-point coefficients in decimal, in order to approximate the co-efficients of `exp(-x)`. The coefficients are given below:

```text
(x * (x * (x * (x * (x * (x * (x * (
-648770010757830093818553637600
*2^(-128)) +
67469480939593786226847644286976
*2^(-128)) +
-3197587544499098424029388939001856
*2^(-128)) +
89244641121992890118377641805348864
*2^(-128)) +
-1579656163641440567800982336819953664
*2^(-128)) +
17685496037279256458459817590917169152
*2^(-128)) +
-115682590513835356866803355398940131328
*2^(-128))
+ 1) /
(x * (x * (x * (x * (x * (x * (x * (x * (x * (x * (x * (x * (x * (
1225524182432722209606361
*2^(-128)) +
114095592300906098243859450
*2^(-128)) +
5665570424063336070530214243
*2^(-128)) +
194450132448609991765137938448
*2^(-128)) +
5068267641632683791026134915072
*2^(-128)) +
104716890604972796896895427629056
*2^(-128)) +
1748338658439454459487681798864896
*2^(-128)) +
23704654329841312470660182937960448
*2^(-128)) +
259380097567996910282699886670381056
*2^(-128)) +
2250336698853390384720606936038375424
*2^(-128)) +
14978272436876548034486263159246028800
*2^(-128)) +
72144088983913131323343765784380833792
*2^(-128)) +
224599776407103106596571252037123047424
*2^(-128))
+ 1)
```

#### Implementation Guidelines

- The `ElectionProof` ticket struct in the block header has two fields:
  - `vrf.Proof`, the output of the VRF, or `ElectionProof` ticket.
  - `WinCount` that corresponds to the result of the Poisson Sortition.
- `WinCount` needs to be `> 0` for winning blocks.
- `WinCount` is included in the tipset weight function. The sum of `WinCount`s of a tipset replaces the size of tipset factor in the weight function.
- `WinCount` is passed to Reward actor to increase the reward for blocks winning multiple times.

```go
GenerateElectionProof(epoch) {
  electionProofInput := GetRandomness(DomainSeparationTag_ElectionProofProduction, epoch, CBOR_Serialize(miner.address))
  vrfResult := miner.VRFSecretKey.Generate(electioinProofInput)

  if GetWinCount(vrfResult.Digest,minerID,epoch)>0 {
    return vrfResult.Proof, GetWinCount(vrfResult.Digest,minerID,epoch)
  }
  return nil
}

GetWinCount(proofDigest, minerID,epoch) {
  // for SHA256, more generally it is 2^len(H)
    const maxDigestSize = 2^256
    minerPower = GetminerPower(minerID, epoch-PowerTableLookback)
    TotalPower = GetTotalPower(epoch-PowerTableLookback)
    if minerPower = 0 {
        return 0
    }
    lambda = minerPower/totalPower*ExpectedLeaderPerEpoch
    h = hash(proofDigest)/maxDigestSize
    rhs = 1 - PoissPmf(lambda, 0)

    WinCount = 0
    for h < rhs {
      WinCount++
      rhs -= PoissPmf(lambda, WinCount)
    }
    return WinCount
}
```

### Leader Election Verification

In order to verify that the leader ElectionProof ticket in a block is correct, miners perform the following checks:

- Verify that the randomness is correct by checking `GetRandomness(epoch)`
- Use this randomness to verify the VRF correctness `Verify_VRF(vrf.Proof,beacon,public_key)`, where vrf.Proof is the ElectionProof ticket.
- Verify ElectionProof.WinCount > 0 by checking `GetWinCount(vrf.Proof, miner,epoch)`, where `vrf.Proof` is the ElectionProof ticket.

## Chain Selection

Just as there can be 0 miners win in a round, there can equally be multiple miners elected in a given round. This in turn means multiple blocks can be created in a round, as seen above. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round.

### Chain Weighting

It is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the biggest amount of committed storage.

In short, the weight at each block is equal to its `ParentWeight` plus that block's delta weight. Details of Filecoin's chain weighting function [are included here](https://observablehq.com/d/3812cd65c054082d).

Delta weight is a term composed of a few elements:

- wPowerFactor: which adds weight to the chain proportional to the total power backing the chain, i.e. accounted for in the chain's power table.
- wBlocksFactor: which adds weight to the chain proportional to the number of tickets mined in a given epoch. It rewards miner cooperation (which will yield more blocks per round on expectation).

The weight should be calculated using big integer arithmetic with order of operations defined above. We use brackets instead of parentheses below for legibility. We have:

```text
w[r+1] = w[r] + (wPowerFactor[r+1] + wBlocksFactor[r+1]) * 2^8
```

For a given tipset `ts` in round `r+1`, we define:

- `wPowerFactor[r+1] = wFunction(totalPowerAtTipset(ts))`
- `wBlocksFactor[r+1] = wPowerFactor[r+1] * wRatio * t / e`
  - with `t = |ticketsInTipset(ts)|`
  - `e = expected number of tickets per round in the protocol`
  - and `wRatio in ]0, 1[`
    Thus, for stability of weight across implementations, we take:
- `wBlocksFactor[r+1] = (wPowerFactor[r+1] * b * wRatio_num) / (e * wRatio_den)`

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
**Note that if your implementation does not allow for rounding to the fourth decimal**, miners should apply the [tie-breaker below](#section-algorithms.expected_consensus.selecting-between-tipsets-with-equal-weight). Weight changes will be on the order of single digit numbers on expectation, so this should not have an outsized impact on chain consensus across implementations.
{{< /hint >}}

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given Tipset should have
the same `ParentWeight` value) plus the delta weight of each parent. To make the
computation a bit easier, a block's `ParentWeight` is stored in the block itself (otherwise
potentially long chain scans would be required to compute a given block's weight).

### Selecting between Tipsets with equal weight

When selecting between Tipsets of equal weight, a miner chooses the one with the smallest final ElectionProof ticket.

In the case where two Tipsets of equal weight have the same minimum VRF output, the miner will compare the next smallest ticket in the Tipset (and select the Tipset with the next smaller ticket). This continues until one Tipset is selected.

The above case may happen in situations under certain block propagation conditions. Assume three blocks B, C, and D have been mined (by miners 1, 2, and 3 respectively) off of block A, with minTicket(B) < minTicket(C) < minTicket(D).

Miner 1 outputs their block B and shuts down. Miners 2 and 3 both receive B but not each others' blocks. We have miner 2 mining a Tipset made of B and C and miner 3 mining a Tipset made of B and D. If both succesfully mine blocks now, other miners in the network will receive new blocks built off of Tipsets with equal weight and the same smallest VRF output (that of block B). They should select the block mined atop `[B, C]` since minVRF(C) < minVRF(D).

The probability that two Tipsets with different blocks would have all the same VRF output can be considered negligible: this would amount to finding a collision between two 256-bit (or more) collision-resistant hashes. Behaviour is explicitly left unspecified in this case.

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
   Specifically, this can be proven with a "witness" block, that is by submitting blocks B2, B3, B4 where B2 is B4's parent and B3's sibling but B3 is not B4's parent. - `!B4.Parents.Include(B3) && B4.Parents.Include(B2) && B3.Parents == B2.Parents && B3.Epoch == B2.Epoch`
   ![Parent-Grinding fault](diagrams/parent_grinding.dot)

### Penalization for faults

A single consensus fault results into:

- miner suspension
- loss of all pledge collateral (which includes the initial pledge and blocks rewards yet to be vested)

### Detection and Reporting

A node that detects and reports a consensus fault is called "slasher". Any user in Filecoin can be a slasher. They can report consensus faults by calling the `ReportConsensusFault` on the `StorageMinerActor` of the faulty miner. The slasher is rewarded with a portion of the penalty paid by the offending miner's `ConsensusFaultPenalty` for notifying the network of the consensus fault. Note that some slashers might not get the full reward because of the low balance of the offending miners. However rational honest miners are still incentivised to notify the network about consensus faults.

The reward given to the slasher is a function of some initial share (`SLASHER_INITIAL_SHARE`) and growth rate (`SLASHER_SHARE_GROWTH_RATE`) and it has a maximum `maxReporterShare`. Slasher's share increases exponentially as epoch elapses since the block when the fault is committed (see `RewardForConsensusSlashReport`). Only the first slasher gets their share of the pledge collateral and the remaining pledge collateral is burned. The longer a slasher waits, the higher the likelihood that the slashed collateral will be claimed by another slasher.
