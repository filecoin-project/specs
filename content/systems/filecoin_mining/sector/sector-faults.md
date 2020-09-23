---
title: Sector Faults
weight: 5
dashboardWeight: 2
dashboardState: stable
dashboardAudit: wip
dashboardTests: 0
---

# Sector Faults

It is very important for storage providers to have a strong incentive to both report the failure to the chain and attempt recovery from the fault in order to uphold the storage guarantee for the networkʼs clients. Without this incentive, it is impossible to distinguish an honest minerʼs hardware failure from malicious behavior, which is necessary to treat miners fairly. The size of the fault fees depend on the severity of the failure and the rewards that the miner is expected to earn from the sector to make sure incentives are aligned. The three types of sector storage fault fees are:

* **Sector fault fee:** This fee is paid per sector per day while the sector is in a faulty state. The size of the fee is slightly more than the amount the sector is expected to earn per day in block rewards. If a sector remains faulty for more than 2 consecutive weeks, the sector will pay a termination fee and be _removed from the chain state_. As storage miner reliability increases above a reasonable threshold, the risk posed by these fees decreases rapidly.
* **Sector fault detection fee:** This is a one-time fee paid in the event of a failure if the miner does not report it honestly and instead the unreported failure is caught by the   chain. Given the probabilistic nature of our PoSt checks, this is set to a few days worth of block reward that would be expected to be earned by a particular sector.
* **Sector termination fee:** A sector can be terminated before its expiration through automatic faults or miner decisions. A termination fee is charged that is, in principle, equivalent to how much a sector has earned so far, up to a limit in order to avoid discouraging long sector lifetimes. In an active termination, the miner decides to stop mining and they pay a fee to leave. In a fault termination, a sector is in a faulty state for too long, and the chain terminates the deal, returns unpaid deal fees to the client and penalizes the miner. Termination fee is currently capped at 90 days worth of block reward that a sector will earn. Miners are responsible for deciding to comply with local regulations, and may sometimes need to accept a termination fee for complying with content laws. Many of the concepts and parameters above make use of the notion of “how much a sector would have earned in a day” in order to understand and align incentives for participants. This concept is robustly tracked and extrapolated on chain.