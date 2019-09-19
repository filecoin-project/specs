---
title: Proof Generator
---

Block miners should constantly be performing Proofs of SpaceTime, and also checking if they have a winning `ticket` to propose a block at each height/in each round. Rounds are currently set to take around 30 seconds, in order to account for network propagation around the world. The details of both processes are defined here.


## The Miner Actor

After successfully calling `CreateStorageMiner`, a miner actor will be created on-chain, and registered in the storage market. This miner, like all other Filecoin State Machine actors, has a fixed set of methods that can be used to interact with or control it.

For details on the methods on the miner actor, see its entry in the [actors spec](actors.md#storage-miner-actor).

### Owner Worker distinction

The miner actor has two distinct 'controller' addresses. One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key.

### Storage Mining Cycle

Storage miners must continually produce Proofs of SpaceTime over their storage to convince the network that they are actually storing the sectors that they have committed to. Each PoSt covers a miner's entire storage.

#### Step 0: Registration

To initially become a miner, a miner first register a new miner actor on-chain. This is done through the storage market actor's [`CreateStorageMiner`](actors.md#createstorageminer) method. The call will then create a new miner actor instance and return its address.

The next step is to place one or more storage market asks on the market. This is done through the storage markets [`AddAsk`](actors.md#addask) method. A miner may create a single ask for their entire storage, or partition their storage up in some way with multiple asks (at potentially different prices).

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the section on [deal](storage-market.md#deal).

When they have a full sector, they should seal it. This is done by invoking [`PoRep.Seal`](proofs.md#seal) on the sector.

#### Step 1: Commit

When the miner has completed their first seal, they should post it on-chain using [CommitSector](actors.md#commitsector). If the miner had zero committed sectors prior to this call, this begins their proving period.

The proving period is a fixed amount of time in which the miner must submit a Proof of Space Time to the network.

During this period, the miner may also commit to new sectors, but they will not be included in proofs of space time until the next proving period starts.
For example, if a miner currently PoSts for 10 sectors, and commits to 20 more sectors. The next PoSt they submit (i.e. the one they're currently proving) will be for 10 sectors again, the subsequent one will be for 30.

TODO: sectors need to be globally unique. This can be done either by having the seal proof prove the sector is unique to this miner in some way, or by having a giant global map on-chain is checked against on each submission. As the system moves towards sector aggregation, the latter option will become unworkable, so more thought needs to go into how that proof statement could work.

#### Step 2: Proving Storage (PoSt creation)

```go
func ProveStorage(sectorSize BytesAmount, sectors []commR) PoStProof {
    challengeBlockHeight := miner.ProvingPeriodEnd - POST_CHALLENGE_TIME

    // Faults to be used are the currentFaultSet for the miner.
    faults := miner.currentFaultSet
    seed := GetRandFromBlock(challengeBlockHeight)
    return GeneratePoSt(sectorSize, sectors, seed, faults)
}
```

Note: See ['Proof of Space Time'](proof-of-spacetime.md) for more details.

The proving set remains consistent during the proving period. Any sectors added in the meantime will be included in the next proving set, at the beginning of the next proving period.

#### Step 3: PoSt Submission

When the miner has completed their PoSt, they must submit it to the network by calling [SubmitPoSt](actors.md#submitpost). There are two different times that this *could* be done.

1. **Standard Submission**: A standard submission is one that makes it on-chain before the end of the proving period. The length of time it takes to compute the PoSts is set such that there is a grace period between then and the actual end of the proving period, so that the effects of network congestion on typical miner actions is minimized.
2. **Penalized Submission**: A penalized submission is one that makes it on-chain after the end of the proving period, but before the generation attack threshold. These submissions count as valid PoSt submissions, but the miner must pay a penalty for their late submission. (See '[Faults](faults.md)' for more information)
   - Note: In this case, the next PoSt should still be started at the beginning of the proving period, even if the current one is not yet complete. Miners must submit one PoSt per proving period.

Along with the PoSt submission, miners may also submit a set of sectors that they wish to remove from their proving set. This is done by selecting the sectors in the 'done' bitfield passed to `SubmitPoSt`.

### Stop Mining

In order to stop mining, a miner must complete all of its storage contracts, and remove them from their proving set during a PoSt submission. A miner may then call [`DePledge()`](actors.md#depledge) to retrieve their collateral. `DePledge` must be called twice, once to start the cooldown, and once again after the cooldown to reclaim the funds. The cooldown period is to allow clients whose files have been dropped by a miner to slash them before they get their money back and get away with it.

### Faults

Faults are described in the [faults document](faults.md).

### On Being Slashed (WIP, needs discussion)

If a miner is slashed for failing to submit their PoSt on time, they currently lose all their pledge collateral. They do not necessarily lose their storage collateral. Storage collateral is lost when a miner's clients slash them for no longer having the data. Missing a PoSt does not necessarily imply that a miner no longer has the data. There should be an additional timeout here where the miner can submit a PoSt, along with 'refilling' their pledge collateral. If a miner does this, they can continue mining, their mining power will be reinstated, and clients can be assured that their data is still there.

TODO: disambiguate the two collaterals across the entire spec

Review Discussion Note: Taking all of a miners collateral for going over the deadline for PoSt submission is really really painful, and is likely to dissuade people from even mining filecoin in the first place (If my internet going out could cause me to lose a very large amount of money, that leads to some pretty hard decisions around profitability). One potential strategy could be to only penalize miners for the amount of sectors they could have generated in that timeframe.

### Future Work

There are many ideas for improving upon the storage miner, here are ideas that may be potentially implemented in the future.

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, and will not be implemented in the initial release of Filecoin, but could provide interesting capabilities down the road.
