---
title: Filecoin CryptoEconomics
weight: 8
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# Filecoin CryptoEconomics

The Filecoin network is a complex multi-agent economic system. The CryptoEconomics of Filecoin touch on most parts of the system. As such related mechanisms and details have been included in many places across this specification. This section aims to explain in more detail the mechanisms and parameters of the system that contribute to the overall network-level goals.

Next, we provide the parameters of the cryptoeconomic model. It is advised that the reader refers to the following sections that are closely related to the Filecoin CryptoEconomic model.

- [Miner Collaterals](miner_collaterals)
- The Minting Model
- Token Allocation

## Initial Parameter Recommendation

Economic analyses and models were developed to design, validate, and parameterize the mechanisms described in the sections listed above. Cryptoeconomics is a young field, where global expertise is both sparse and shallow. Developing these recommendations is advancing the state of the art, not only in the field of decentralized storage networks, but also of cryptoeconomic mechanism design as a wider discipline.

The following table summarizes initial parameter recommendations for Filecoin. Monitoring, testing, validation and recommendations will continue to evolve and adapt. When changes to these parameters are due they will be announced and applied through FIPs.


| **Parameter**  | **Value**   |
| :------------- | :---------- |
| Baseline Initial Value | 1 EiB | 
| Baseline Annual Growth Rate  | 200% |
| Percent simple minting vs baseline minting | 30% / 70% |
| Reward delay and linear vesting period | 20 days |
| Linear vesting period | 180 days |
| Sector quality multipliers | Committed Capacity: 1x <br> Regular Deals: 1x <br> Verified Client Deals: 10x |
| Initial pledge function | 20 days worth of block reward + <br> share of 30% FIL circulating supply target | 
| Minimum sector lifetime | 180 days |
| Maximum sector lifetime | 540 days |
| Minimum deal duration | 180 days |
| Maximum deal duration | 540 days |
| Sector Fault Fee | 2.14 days |
| Sector Fault Detection Fee | 5 days worth of estimated block reward |
| Sector Termination Fee | Estimated number of days of block reward that a sector has earned; capped at 90 days |
| Network Transaction Fee | Dynamic fee structure based on network congestion |

## Design Principles Justification

**Baseline Minting:** Filecoin tokens are a limited resource. The rate at which tokens are deployed into the network should be controlled to maximize their net benefit to the community, just like the consumption of any exhaustible common-pool resource. The purpose of baseline minting is to: (a) reward participants proportionally to the storage they provide rather than exponentially, based on the time when they joined the network, and (b) to adjust the minting rate based on approximated network utility in order to maintain a relatively steady flow of block rewards over longer time periods.

**Initial Pledge:** The justification for having an initial pledge is as follows: firstly, having an initial pledge forces miners to behave responsibly on their sector commitments and holds them accountable for not keeping up to their promise, even before they earn any block reward. Secondly, requiring a pledge of stake in the network supports and enhances the security of the consensus mechanism.

**Block Reward Vesting:** In order to reduce the initial pledge requirement of a sector, the network considers all vesting block rewards as collateral. However, tracking block rewards on a per-sector level is not scalable. Instead, the protocol tracks rewards at a per-miner level and linearly vests block rewards over a fixed duration.

**Minimum Sector Lifetime:** The justification for a minimum sector lifetime is as follows. Committing a sector to the Filecoin Network currently requires a moderately computationally-expensive "sealing" operation up-front, whose amortized cost is lower if the sector's lifetime is longer. In addition, a sector commitment will involve on-chain transactions, for which gas fees will be paid. The net effect of these transaction costs will be subsidized by the block reward, but only for sectors that will contribute to the network and earn rewards for a sufficiently long duration. Under current constraints, short-lived sectors would reduce the overall capacity of the network to deliver useful storage over time.

**Sector Fault Fee:** If stored sectors are withdrawn from the network only temporarily, a substantial fraction of those sectors' value may be recovered in case the data storage is quickly restored to normal operation — this means that the network need not levy a termination fee immediately. However, even temporary interruptions can be disruptive, and also damage confidence about whether the sector is recoverable or permanently lost. In order to account for this situation, the network charges a much smaller fee per day that a sector is not being proven as promised (until enough days have passed that the network writes off that sector as terminated).

**Sector Fault Detection Fee:** If a sector is temporarily damaged, storage miners are expected to proactively detect, report, and repair the fault. An unannounced interruption in service is both more disruptive for clients and more of a signal that the fault may not have been caught early enough to fully recover. Finally, dishonest storage miners may have some chance of briefly evading detection and earning rewards despite being out of service. For all these reasons, a higher penalty is applied when the network detects an undeclared fault.

**Sector Termination Fee:** The ultimate goal of the Filecoin Network is to provide useful data storage. Use-cases for unreliable data storage, which may vanish without warning, are much rarer than use-cases for reliable data storage, which is guaranteed in advance to be maintained for a given duration. So to the extent that committed sectors disappear from the network, most of the value provided by those sectors is canceled out, in most cases. If storage miners had little to lose by terminating active sectors compared to their realized gains, then this would be a negative externality that fails to be effectively managed by the storage market; termination fees internalize this cost.
