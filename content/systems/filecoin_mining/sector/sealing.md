---
title: Sector Sealing
weight: 4
dashboardWeight: 2
dashboardState: stable
dashboardAudit: wip
dashboardTests: 0
---

# Sector Sealing

Before a Sector can be used, the Miner must _seal_ the Sector: encode the data in the Sector to prepare it for the proving process.

- **Unsealed Sector**: A Sector of raw data.
  - **UnsealedCID (CommD)**: The root hash of the Unsealed Sector's merkle tree. Also called CommD, or "data commitment."
- **Sealed Sector**: A Sector that has been encoded to prepare it for the proving process.
  - **SealedCID (CommR)**: The root hash of the Sealed Sector's merkle tree. Also called CommR, or "replica commitment."

Sealing a sector through Proof-of-Replication (PoRep) is a computation-intensive process that results in a unique encoding of the sector. Once data is sealed, storage miners: generate a proof; run a SNARK on the proof to compress it; and finally, submit the result of the compression to the blockchain as a certification of the storage commitment. Depending on the PoRep algorithm and protocol security parameters, cost profiles and performance characteristics vary and tradeoffs have to be made among sealing cost, security, onchain footprint, retrieval latency and so on. However, sectors can be sealed with commercial hardware and sealing cost is expected to decrease over time. The Filecoin Protocol will launch with Stacked Depth Robust (SDR) PoRep with a planned upgrade to Narrow Stacked Expander (NSE) PoRep with improvement in both cost and retrieval latency.

The Lotus-specific set of functions applied to the sealing of a sector can be found [here](https://github.com/filecoin-project/lotus/blob/master/cmd/lotus-miner/sealing.go).

## Randomness

Randomness is an important attribute that helps the network verify the integrity of Miners' stored data. Filecoin's block creation process includes two types of randomness:

- [DRAND](drand): Values pulled from a distributed random beacon
- VRF: The output of a _Verifiable Random Function_ (VRF), which takes the previous block's VRF value and produces the current block's VRF value.

Each block produced in Filecoin includes values pulled from these two sources of randomness.

When Miners submit proofs about their stored data, the proofs incorporate references to randomness added at specific epochs. Assuming these values were not able to be predicted ahead of time, this helps ensure that Miners generated proofs at a specific point in time.

There are two proof types. Each uses one of the two sources of randomness:

- Windowed PoSt: Uses Drand values
- Proof of Replication (PoRep): Uses VRF values

## Drawing randomness for sector commitments

Tickets are used as input to calculation of the ReplicaID in order to tie Proofs-of-Replication to a given chain, thereby preventing long-range attacks (from another miner in the future trying to reuse SEALs).

The ticket has to be drawn from a finalized block in order to prevent the miner from potential losing storage (in case of a chain reorg) even though their storage is intact.

Verification should ensure that the ticket was drawn no farther back than necessary by the miner. We note that tickets can uniquely be associated with a given round in the protocol (lest a hash collision be found), but that the round number is explicited by the miner in `commitSector`.

We present precisely how ticket selection and verification should work. In the below, we use the following notation:

- `F`-- Finality (number of rounds)
- `X`-- round in which SEALing starts
- `Z`-- round in which the SEAL appears (in a block)
- `Y`-- round announced in the SEAL `commitSector` (should be X, but a miner could use any Y <= X), denoted by the ticket selection
- `T`-- estimated time for SEAL, dependent on sector size
- `G = T + variance`-- necessary flexibility to account for network delay and SEAL-time variance.

We expect Filecoin will be able to produce estimates for sector commitment time based on sector sizes, e.g.:
`(estimate, variance) <--- SEALTime(sectors)`
G and T will be selected using these.

**Picking a Ticket to Seal:** When starting to prepare a SEAL in round X, the miner should draw a ticket from X-F with which to compute the SEAL.

**Verifying a Seal's ticket:** When verifying a SEAL in round Z, a verifier should ensure that the ticket used to generate the SEAL is found in the range of rounds `[Z-T-F-G, Z-T-F+G]`.

```text
                               Prover
           ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
          │

          ▼
         X-F ◀───────F────────▶ X ◀──────────T─────────▶ Z
     -G   .  +G                 .                        .
  ───(┌───────┐)───────────────( )──────────────────────( )────────▶
      └───────┘                 '                        '        time
 [Z-T-F-G, Z-T-F+G]
          ▲

          └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
                              Verifier
```

Note that the prover here is submitting a message on chain (i.e. the SEAL). Using an older ticket than necessary to generate the SEAL is something the miner may do to gain more confidence about finality (since we are in a probabilistically final system). However it has a cost in terms of securing the chain in the face of long-range attacks (specifically, by mixing in chain randomness here, we ensure that an attacker going back a month in time to try and create their own chain would have to completely regenerate any and all sectors drawing randomness since to use for their fork's power).

We break this down as follows:

- The miner should draw from `X-F`.
- The verifier wants to find what `X-F` should have been (to ensure the miner is not drawing from farther back) even though Y (i.e. the round of the ticket actually used) is an unverifiable value.
- Thus, the verifier will need to make an inference about what `X-F` is likely to have been based on:
  - (known) round in which the message is received (Z)
  - (known) finality value (F)
  - (approximate) SEAL time (T)
- Because T is an approximate value, and to account for network delay and variance in SEAL time across miners, the verifier allows for G offset from the assumed value of `X-F`: `Z-T-F`, hence verifying that the ticket is drawn from the range `[Z-T-F-G, Z-T-F+G]`.

In Practice, the Filecoin protocol will include a `MAX_SEAL_TIME` for each sector size and proof type.
