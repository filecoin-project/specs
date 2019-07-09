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


**Definitions**

- `PoStProvingPeriod: UInt = 24 * 60 * 60 / 30 = 2880` (1 day, given `30s` block time)
- `PoStTimeout: UInt` (less than 1 ProvingPeriod)
- `SectorFailureTimout: UInt` (more than 1 Proving Period)

TODO: define the above constants

#### PoSt Late Submission

- **After `ProvingPeriodEnd`:**
  - **Condition**: If the miner posts their PoSt after the proving period ends, but before `PoStTimeout`
  - **Reporting:** The miner submits their PoSt as usual, but includes the `LateSubmissionFee`.
  - **Penalization:**
    - *Economic penalization*: To determine the penalty amount, `ComputeLateSubmissionFee(minerPower)` is called.
    - *Power penalization*: The miners' power is reduced to `0`.
  - **Recovery**:
    - **Conditions for Recovery:**
      - the submission is within the `PoStTimeout`,
      - the fee provided matches the required `LateSubmissionFee` for the miner
    - *Economic*: There is no recovery of the late submission fee.
    - *Power*: The miners' power is reset matching the submitted PoSt.
- **After `PoStTimeout`:**
  - **Condition**: If the miner posts their PoSt after the proving period ends, but before `SectorFailureTimeout`
  - **Reporting:** The miner submits their PoSt as usual, but includes the `LateSubmissionFee + Lost Storage Collateral`.
  - **Penalization:**
    - *Economic penalization*: To determine the penalty amount, `ComputeLateSubmissionFee(minerPower)` is called, in addition all storage collateral is lost for the sectors in current `ProvingSet`
    - *Power penalization*: The miners' power is reduced to `0`.
  - **Recovery**:
    - **Conditions for Recovery:**
      - the submission is within the `SectorFailureTimeout`,
      - the fee provided matches the required `LateSubmissionFee` for the miner and the required storage collateral
    - *Economic*: There is no recovery of the fees.
    - *Power*: The miners' power is reset matching the submitted PoSt.
- **After `SectorFailureTimeout`:**
  - **Condition**: If the miner posts their PoSt after the proving period ends, or not at all.
  - **Reporting:** The same as "Unreported storage fault slashing"
  - **Penalization:**
    - *Economic penalization*: All storage collateral is lost.
    - *Power penalization*: The miners' power is reduced to `0`.
  - **Recovery**:
    - Only resubmission of the sectors can lead to recovery.

#### Unreported storage fault slashing

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

#### Reported storage fault penalty

- **Condition:** The miner submits their PoSt with a non-empty set of 'missing sectors'.
- **Reporting:** The miner can specify some sectors that they failed to prove during the proving period.
  - Note: These faults are output by the `ProveStorage` routine, and are posted on-chain when posting the proof. This occurs when the miner (for example) has a disk failure, or other local data corruption.
- **Check:** The chain checks that the proof verifies with the missing sectors.
- **Penalization:** The miner is penalized for collateral and power proportional to the number of missing sectors. The sectors are also removed from the miners proving set.
  - TODO: should the collateral lost here be proportional to the remaining time?
  - TODO(nicola): check if the time between posting two proofs allows for a generation attack if it does not then we might reconsider the sector not being lost
- Note: if a sector is missed here, and they are recovered after the fact, the miner could simple 're-commit' the sector. They still have to pay the collateral, but the data can be quickly re-introduced into the system to avoid clients calling them out for breach of contract (this would only work because the sector commD/commR is the same)
- Note: In the case where a miner is temporarily unable to prove some of their data, they can simply wait for the temporary unavailability to recover, and then continue proving, submitting the proofs a bit late if necessary (paying appropriate fees, as described above).



#### Breach of contract dispute

- **Condition:** A client who has stored data with a miner, and the miner removes the sector containing that data before the end of the agreed upon time period.
- **Reporting:** The client invokes `ArbitrateDeal` on the offending miner actor with a signed deal from that miner for the storage in question. Note: the reporting must happen within one proving period of the miner removing the storage erroneously.
- **Check:** The chain checks that the deal was correctly signed by the miner in question, that the deal has not yet expired, and that the sector referenced by the deal is no longer in the miners proving set.
- **Penalization:** The miner is penalized an amount proportional to the incorrectly removed sector. This penalty is taken from their pledged collateral .
- Note: This implies that miners cannot re-seal data into different sectors. We could come up with a protocol where the client gives the miner explicit consent to re-seal, but that is more complicated and can be done later.
