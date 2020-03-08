---
menuTitle: "Payment Channels"
statusIcon: 🔁
title: "Payment Channels"
entries:
- payment_channel_actor
---

{{<label payment_channels>}}

Payment Channels are used in the Filecoin {{<sref retrieval_market>}} to enable efficient off-chain payments and accounting between parties for what is expected to be series of microtransactions, specifically those occurring as part of retrieval market data retrieval.

Note that the following provides a high-level overview of payment channels and an accompanying interface. The lotus implementation of [vouchers](https://github.com/filecoin-project/lotus/blob/master/chain/types/voucher.go) and [payment channels](https://github.com/filecoin-project/lotus/tree/master/paych) are also good references.

You can also read more about the {{<sref payment_channel_actor "Filecoin payment channel actor interface">}}.

In short, the payment channel actor can be used to open long-lived, flexible payment channels between users. Each channel can be funded by adding to their balance. 
The goal of the payment channel actor is to enable a series of off-chain microtransactions to be reconciled on-chain at a later time with fewer messages. Accordingly, the expectation is `From` will send to `To` vouchers of successively greater `Value` and increasing `Nonce`s. When they choose to, `To` can `Update` the channel to update the balance available `ToSend` to them in the channel, and can choose to `Collect` this balance at any time (incurring a gas cost).
The channel is split into `lane`s created as part of updating the channel state with a payment `voucher`. Each lane has an associated `nonce` and amount of tokens it can be `redeemed` for. These lanes allow for a lot of accounting between parties to be done off chain and reconciled via single updates to the payment channel, merging these lanes to arrive at a desired outcome.

Over the course of a transaction cycle, each party to the payment channel can send the other `voucher`s. The payment channel's  `From` account holder will send a signed voucher with a given nonce to the `To` account holder. The latter can use the voucher to `redeem` part of the lane's value, merging other lanes into it as needed.

For instance if `From` sends `To` the following vouchers (voucher_val, voucher_nonce) for a lane with 100 to be redeemed: (10, 1), (20, 2), (30, 3), then `To` could choose to redeem (30, 3) bringing the lane's value to 70 (100 - 30). They could not redeem (10, 1) or (20, 2) thereafter. They could however redeem (20, 2) for 20, and then (30, 3) for 10 (30 - 20) thereafter.

The multiple lanes enable two parties to use a single payment channel to adjudicate multiple independent sets of payments.

Vouchers are signed by the sender and authenticated using a `Secret`, `PreImage` pair provided by the paying party. If the `PreImage` is indeed a pre-image of the `Secret` when used as input to some given algorithm (typically a one-way function like a hash), the `Voucher` is valid. The `Voucher` itself contains the `PreImage` but not the `Secret` (communicated separately to the receiving party). This enables multi-hop payments since an intermediary cannot redeem a voucher on their own. They can also be used to update the minimum height at which a channel will be closed. Likewise, vouchers can have `TimeLock`s to prevent they are being used too early, likewise a channel can have a `MinCloseHeight` to prevent it being closed prematurely (e.g. before the recipient has collected funds) by the sender.

Once their transactions have completed, either party can choose to `Close` the channel, the recipient can then `Collect` the `ToPay` amount from the channel. `From` will be refunded the remaining balance in the channel.

So we have:

- \[off-chain\] - Two parties agree to a series of transactions (for instance as part of file retrieval) with party **A** paying party **B** up to some **total** sum of Filecoin over time.
- \[on-chain\] - The {{<sref payment_channel_actor>}} is used called by A to open a payment channel `from` A `to` B and a lane is opened to increase the `balance` of the channel, triggering a transaction between A and the payment channel actor.
At any time, A can open new lanes to increase the total balance available in the channel (e.g. if A and B choose to do more transactions together).
- \[off-chain\] - Throughout the transaction cycle (e.g. on every piece of data sent via a retrieval deal), party A sends a voucher to party B enabling B to redeem more payment from the payment lanes, and incentivizing B to continue providing a service (e.g. sending more data along).
- \[on-chain\] - At regular intervals, B can `Update` the payment channel balance available `ToSend` with the vouchers received (past their `timeLock`), decreasing the remaining `Value` of the payment channel.
- \[on-chain\] - At the end of the cycle, past the `MinCloseHeight`, A can choose to `Close` the payment channel.
- \[on-chain\] - B can choose to `Collect` the amount `ToSend` triggering a payment between the payment channel actor and B.
