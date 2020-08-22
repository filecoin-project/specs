---
title: Sector Sealing
dashboardWeight: 2
dashboardState: wip
dashboardAudit: wip
dashboardTests: 0
---

# Sector Sealing
---

{{<embed src="sealing.id" lang="go" >}}

## Drawing randomness for sector commitments

[Tickets](storage_power_consensus#the-ticket-chain-and-drawing-randomness "The Ticket chain and drawing randomness") are used as input to calculation of the ReplicaID in order to tie Proofs-of-Replication to a given chain, thereby preventing long-range attacks (from another miner in the future trying to reuse SEALs).

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

## Picking a Ticket to Seal

When starting to prepare a SEAL in round X, the miner should draw a ticket from X-F with which to compute the SEAL.

## Verifying a Seal's ticket

When verifying a SEAL in round Z, a verifier should ensure that the ticket used to generate the SEAL is found in the range of rounds [Z-T-F-G, Z-T-F+G].

### In Detail

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

### In Practice

The Filecoin protocol will include a `MAX_SEAL_TIME` for each sector size and proof type.