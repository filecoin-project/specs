---
title: Gas Fee
weight: 6
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: coming
dashboardTests: 0
---

# Gas Fees

## Summary

As is traditionally the case with many blockchains, Gas is a unit of measure of how much storage and/or compute resource an on-chain message operation consumes in order to be executed. At a high level, it works as follows: the message sender specifies the maximum amount they are willing to pay in order for their message to be executed and included in a block. This is specified both in terms of total number of units of gas (`GasLimit`), which is generally expected to be higher than the actual `GasUsed` and in terms of the price (or fee) per unit of gas (`GasFeeCap`).

Traditionally, `GasUsed * GasFeeCap` goes to the block producing miner as a reward. The result of this product is treated as the priority fee for message inclusion, that is, messages are ordered in decreasing sequence and those with the highest `GasUsed * GasFeeCap` are prioritised, given that they return more profit to the miner.

However, it has been observed that this tactic (of paying `GasUsed * GasFee`) is problematic for block producing miners for a few reasons. Firstly, a block producing miner may include a very expensive message (in terms of chain resources required) for free in which case the chain itself needs to bear the cost. Secondly, message senders can set arbitrarily high prices but for low-cost messages (again, in terms of chain resources), leading to a DoS vulnerability.

In order to overcome this situation, the Filecoin blockchain defines a `BaseFee`, which is burnt for every message. The rationale is that given that Gas is a measure of on-chain resource consumption, it makes sense for it to be burned, as compared to be rewarded to miners. This way, fee manipulation from miners is avoided. The `BaseFee` is dynamic, adjusted automatically according to network congestion. This fact, makes the network resilient against spam attacks. Given that the network load increases during SPAM attacks, maintaining full blocks of SPAM messages for an extended period of time is impossible for an attacker due to the increasing `BaseFee`.

Finally, `GasPremium` is the priority fee included by senders to incentivize miners to pick the most profitable messages. In other words, if a message sender wants its message to be included more quickly, they can set a higher `GasPremium`.

## Parameters

- `GasUsed` is a measure of the amount of resources (or units of gas) consumed, in order to execute a message. Each unit of gas is measured in attoFIL and therefore, `GasUsed` is a number that represents the units of energy consumed. `GasUsed` is independent of whether a message was executed correctly or failed.
- `BaseFee` is the set price per unit of gas (measured in attoFIL/gas unit) to be burned (sent to an unrecoverable address) for every message execution. The value of the `BaseFee` is dynamic and adjusts according to current network congestion parameters. For example, when the network exceeds 5B gas limit usage, the `BaseFee` increases and the opposite happens when gas limit usage falls below 5B. The `BaseFee` applied to each block should be included in the block itself. It should be possible to get the value of the current `BaseFee` from the head of the chain. The `BaseFee` applies per unit of `GasUsed` and therefore, the total amount of gas burned for a message is `BaseFee * GasUsed`. Note that the `BaseFee` is incurred for every message, but its value is the same for all messages in the same block.
- `GasLimit` is measured in units of gas and set by the message sender. It imposes a hard limit on the amount of gas (i.e., number of units of gas) that a message’s execution should be allowed to consume on chain. A message consumes gas for every fundamental operation it triggers, and a message that runs out of gas fails. When a message fails, every modification to the state that happened as a result of this message's execution is reverted back to its previous state. Independently of whether a message execution was successful or not, the miner will receive a reward for the resources they consumed to execute the message (see `GasPremium` below).
- `GasFeeCap` is the maximum price that the message sender is willing to pay per unit of gas (measured in attoFIL/gas unit). Together with the `GasLimit`, the `GasFeeCap` is setting the maximum amount of FIL that a sender will pay for a message: a sender is guaranteed that a message will never cost them more than `GasLimit * GasFeeCap` attoFIL (not including any Premium that the message includes for its recipient).
- `GasPremium` is the price per unit of gas (measured in attoFIL/gas) that the message sender is willing to pay (on top of the `BaseFee`) to "tip" the miner that will include this message in a block. A message typically earns its miner `GasLimit * GasPremium` attoFIL, where effectively `GasPremium = GasFeeCap - BaseFee`. Note that `GasPremium` is applied on `GasLimit`, as opposed to `GasUsed`, in order to make message selection for miners more straightforward.

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/vm/burn.go"  lang="go" symbol="ComputeGasOverestimationBurn">}}

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/store/basefee.go"  lang="go" symbol="ComputeNextBaseFee">}}

## Notes & Implications

- The `GasFeeCap` should always be higher than the network's `BaseFee`. If a message’s `GasFeeCap` is lower than the `BaseFee`, then the remainder comes from the miner (as a penalty). This penalty is applied to the miner because they have selected a message that pays less than the network `BaseFee` (i.e., does not cover the network costs). However, a miner might want to choose a message whose `GasFeeCap` is smaller than the `BaseFee` if the same sender has another message in the message pool whose `GasFeeCap` is much bigger than the `BaseFee`. Recall, that a miner should pick all the messages of a sender from the message pool, if more than one exists. The justification is that the increased fee of the second message will pay off the loss from the first.

- If `BaseFee + GasPremium` > `GasFeeCap`, then the miner might not earn the entire `GasLimit * GasPremium` as their reward.

- A message is hard-constrained to spending no more than `GasFeeCap * GasLimit`. From this amount, the network `BaseFee` is paid (burnt) first. After that, up to `GasLimit * GasPremium` will be given to the miner as a reward.

- A message that runs out of gas fails with an "out of gas" exit code. `GasUsed * BaseFee` will still be burned (in this case `GasUsed = GasLimit`), and the miner will still be rewarded `GasLimit * GasPremium`. This assumes that `GasFeeCap > BaseFee + GasPremium`.

- A low value for the `GasFeeCap` will likely cause the message to be stuck in the message pool, as it will not be attractive-enough in terms of profit for any miner to pick it and include it in a block. When this happens, there is a procedure to update the `GasFeeCap` so that the message becomes more attractive to miners. The sender can push a new message into the message pool (which, by default, will propagate to other miners' message pool) where: i) the identifier of the old and new messages is the same (e.g., same `Nonce`) and ii) the `GasPremium` is updated and increased by at least 25% of the previous value.
