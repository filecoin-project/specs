---
title: Filecoin Parameters
weight: 3
dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
---

# Filecoin Parameters
---

## `SectorMaximumLifetimeSDR`

actors/abi/sector.go

```go
const SectorMaximumLifetimeSDR = ChainEpoch(1_262_277 * 5)
```

**Description:** This parameter is denoting the maximum duration (i.e., from activation to expiration) of a sector sealed with SDR.

**Motivation:** The setting guarantees that SDR is secure in the _cost model_ for WindowPoSt and in the _time model_ for WinningPoSt. The setting is based on estimation of hardware latency improvement and hardware and software cost reduction over time.

---

---

actors/builtin/market/policy.go

```go
const DealUpdatesInterval = builtin.EpochsInDay
```
**Description:** The number of blocks between payouts for deals

---

```go
var ProvCollateralPercentSupplyNum = big.NewInt(5)
var ProvCollateralPercentSupplyDenom = big.NewInt(100)
```

**Description:** The percentage of normalized circulating supply that must be covered by provider collateral in a deal.

---

```go
var DealMinDuration = abi.ChainEpoch(180 * builtin.EpochsInDay)
var DealMaxDuration = abi.ChainEpoch(540 * builtin.EpochsInDay)
```

**Description:** Minimum & Maximum Deal Duration, set at 100 and 540 days, respectively.

---

---

actors/builtin/miner/monies.go


```go
var PreCommitDepositFactor = 20
```

**Description:** Amount of deposit for PreCommitting a sector. This deposit is lost if a PreCommit is not followed up by a ProveCommit, within a predefined time period.

---

```go
var InitialPledgeFactor = 20
var PreCommitDepositProjectionPeriod = abi.ChainEpoch(PreCommitDepositFactor) * builtin.EpochsInDay
var InitialPledgeProjectionPeriod = abi.ChainEpoch(InitialPledgeFactor) * builtin.EpochsInDay
```

**Description:** Amount of Pledge collateral to be deposited per sector (in expected block rewards per day)
// IP = IPBase(precommit time) + AdditionalIP(precommit time)
// IPBase(t) = BR(t, InitialPledgeProjectionPeriod)
// AdditionalIP(t) = LockTarget(t)*PledgeShare(t)
// LockTarget = (LockTargetFactorNum / LockTargetFactorDenom) * FILCirculatingSupply(t)
// PledgeShare(t) = sectorQAPower / max(BaselinePower(t), NetworkQAPower(t))
// PARAM_FINISH

```go
var LockTargetFactorNum = big.NewInt(3)
var LockTargetFactorDenom = big.NewInt(10)
```

**Description:** Fraction of available supply that the AdditionalIP targets to lock.








Some of these parameters are used around the code in the Filecoin subsystems and ABI. Others are used as part of the proofs libraries.

Most are generated/finalized using the [orient framework](https://github.com/filecoin-project/orient). It is used to modelize the Filecoin network.

{{<hint warning>}}
⚠️ **WARNING:** Filecoin is not yet launched, and we are finishing protocol spec and implementations. Parameters are set here as placeholders and highly likely to change to fit product and security requirements.
{{</hint>}}

## Code parameters

{{<embed src="../systems/filecoin_nodes/node_base/network_params.go" lang="go" >}}

## Orient parameters

| LAMBDA | SPACEGAP | BLOCK-SIZE-KIB | SECTOR-SIZE-GIB |
|--------|----------|----------------|-----------------|
| 10     | 0.03     | 2.6084006      | 1024            |
| 10     | 0.03     | 2.9687543      | 1024            |
| 10     | 0.03     | 4.60544        | 256             |
| 10     | 0.03     | 6.9628344      | 256             |
| 10     | 0.03     | 7.195217       | 128             |
| 10     | 0.03     | 12.142387      | 128             |
| 10     | 0.03     | 15.2998495     | 1024            |
| 10     | 0.03     | 22.186821      | 32              |
| 10     | 0.03     | 42.125595      | 32              |
| 10     | 0.03     | 55.240646      | 256             |
| 10     | 0.03     | 107.03619      | 128             |
| 10     | 0.03     | 406.86823      | 32              |
| 10     | 0.06     | 2.3094485      | 1024            |
| 10     | 0.06     | 2.37085        | 1024            |
| 10     | 0.06     | 3.4674127      | 256             |
| 10     | 0.06     | 4.686779       | 256             |
| 10     | 0.06     | 4.9769444      | 128             |
| 10     | 0.06     | 7.705842       | 128             |
| 10     | 0.06     | 9.3208065      | 1024            |
| 10     | 0.06     | 13.775977      | 32              |
| 10     | 0.06     | 25.303907      | 32              |
| 10     | 0.06     | 32.48009       | 256             |
| 10     | 0.06     | 62.670723      | 128             |
| 10     | 0.06     | 238.65137      | 32              |
| 10     | 0.1      | 2.1490319      | 1024            |
| 10     | 0.1      | 2.1985393      | 1024            |
| 10     | 0.1      | 3.0452213      | 256             |
| 10     | 0.1      | 3.8423958      | 256             |
| 10     | 0.1      | 4.1540065      | 128             |
| 10     | 0.1      | 6.059966       | 128             |
| 10     | 0.1      | 7.102623       | 1024            |
| 10     | 0.1      | 10.6557865     | 32              |
| 10     | 0.1      | 19.063526      | 32              |
| 10     | 0.1      | 24.036263      | 256             |
| 10     | 0.1      | 46.211964      | 128             |
| 10     | 0.1      | 176.24756      | 32              |
| 10     | 0.2      | 1.9889219      | 1024            |
| 10     | 0.2      | 2.1184843      | 1024            |
| 10     | 0.2      | 2.7405148      | 256             |
| 10     | 0.2      | 3.2329829      | 256             |
| 10     | 0.2      | 3.5601068      | 128             |
| 10     | 0.2      | 4.8721666      | 128             |
| 10     | 0.2      | 5.501524       | 1024            |
| 10     | 0.2      | 8.404295       | 32              |
| 10     | 0.2      | 14.560543      | 32              |
| 10     | 0.2      | 17.942131      | 256             |
| 10     | 0.2      | 34.33397       | 128             |
| 10     | 0.2      | 131.21773      | 32              |
| 80     | 0.03     | 6.5753794      | 1024            |
| 80     | 0.03     | 10.902712      | 1024            |
| 80     | 0.03     | 19.707468      | 256             |
| 80     | 0.03     | 36.63338       | 128             |
| 80     | 0.03     | 37.16689       | 256             |
| 80     | 0.03     | 71.018715      | 128             |
| 80     | 0.03     | 94.63942       | 1024            |
| 80     | 0.03     | 133.81236      | 32              |
| 80     | 0.03     | 265.37668      | 32              |
| 80     | 0.03     | 357.2812       | 256             |
| 80     | 0.03     | 695.79944      | 128             |
| 80     | 0.03     | 2639.3792      | 32              |
| 80     | 0.06     | 4.183762       | 1024            |
| 80     | 0.06     | 6.1194773      | 1024            |
| 80     | 0.06     | 10.603248      | 256             |
| 80     | 0.06     | 18.887196      | 128             |
| 80     | 0.06     | 18.958448      | 256             |
| 80     | 0.06     | 35.526344      | 128             |
| 80     | 0.06     | 46.80707       | 1024            |
| 80     | 0.06     | 66.525635      | 32              |
| 80     | 0.06     | 130.80322      | 32              |
| 80     | 0.06     | 175.19678      | 256             |
| 80     | 0.06     | 340.8757       | 128             |
| 80     | 0.06     | 1293.6443      | 32              |
| 80     | 0.1      | 3.2964888      | 1024            |
| 80     | 0.1      | 4.3449306      | 1024            |
| 80     | 0.1      | 7.2257156      | 256             |
| 80     | 0.1      | 12.203384      | 256             |
| 80     | 0.1      | 12.303692      | 128             |
| 80     | 0.1      | 22.359337      | 128             |
| 80     | 0.1      | 29.061607      | 1024            |
| 80     | 0.1      | 41.564106      | 32              |
| 80     | 0.1      | 80.880165      | 32              |
| 80     | 0.1      | 107.64613      | 256             |
| 80     | 0.1      | 209.20566      | 128             |
| 80     | 0.1      | 794.4138       | 32              |
| 80     | 0.2      | 2.6560488      | 1024            |
| 80     | 0.2      | 3.0640512      | 1024            |
| 80     | 0.2      | 4.7880635      | 256             |
| 80     | 0.2      | 7.32808        | 256             |
| 80     | 0.2      | 7.552495       | 128             |
| 80     | 0.2      | 12.856943      | 128             |
| 80     | 0.2      | 16.252815      | 1024            |
| 80     | 0.2      | 23.55217       | 32              |
| 80     | 0.2      | 44.856293      | 32              |
| 80     | 0.2      | 58.89311       | 256             |
| 80     | 0.2      | 114.18173      | 128             |
| 80     | 0.2      | 434.17523      | 32              |
