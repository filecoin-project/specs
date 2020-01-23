---
menuTitle: "Payment Channels"
statusIcon: 🔁
title: "Payment Channels"
entries:
- payment_channel_actor
---

{{<label payment_channels>}}

Payment Channels are used in the Filecoin {{<sref retrieval_market>}} in order to enable efficient off-chain payments between parties for what are expected to be series of microtransactions, specifically those occuring as part of retrieval market data retrieval.

Note that the following provides a high-level overview of payment channels and an accompanying interface. We use the lotus implementation as a reference implementation for vouchers in Filecoin. See their [payment channel actor](https://github.com/filecoin-project/lotus/blob/master/chain/actors/actor_paych.go), [voucher](https://github.com/filecoin-project/lotus/blob/master/chain/types/voucher.go) and[payment channel](https://github.com/filecoin-project/lotus/tree/master/paych) implementations for more details.

You can also read more about the {{<sref payment_channel_actor "Filecoin payment channel actor interface">}}.

In short, the payment channel actor can be used to open payment channels between users. Each channel can be funded multiple times by opening payment `lane`s, adding to the channel `Value`. Each lane has an associated `nonce` and amount of tokens it can be `redeemed` for.

Payments are redeemed using `voucher`s. The `From` account holder will send a signed voucher with a given nonce to the `To` account holder. The latter can use the voucher to `redeem` part of the lane's value. The goal of the payment channel actor is to enable a series of off-chain microtransactions to be reconciled on-chain at a later time with fewer messages. Accordingly, the expectation is `From` will send to `To` vouchers of successively greater `Value` and increasing `Nonce`s. When they choose to, `To` can `Update` the channel to update the balance available `ToSend` to them in the channel, and can choose to `Collect` this balance at any time (incurring a gas cost).

For instance if `From` sends `To` the following vouchers (voucher_val, voucher_nonce) for a lane with 100 to be redeemed: (10, 1), (20, 2), (30, 3), then `To` could choose to redeem (30, 3) bringing the lane's value to 70 (100 - 30). They could not redeem (10, 1) or (20, 2) thereafter. They could however redeem (20, 2) for 20, and then (30, 3) for 10 (30 - 20) thereafter.

The multiple lanes enable two parties to use a single payment channel to adjudicate multiple independent sets of payments.

Vouchers are signed by the sender and authenticated using a `Secret`, `PreImage` pair provided by the paying party. If the `PreImage` is indeed a pre-image of the `Secret` when used as input to some given algorithm (typically a one-way function like a hash), the `Voucher` is valid. The `Voucher` itself contains the `PreImage` but not the `Secret` (communicated separately to the receiving party). This enables multi-hop payments since an intermediary cannot redeem a voucher on their own.

Once their transactions have completed, either party can choose to `Close` the channel, the recipient can collect the `ToPay` amount from the channel.

Vouchers can have `TimeLock`s to prevent their being used too early, likewise a channel can have a `MinCloseHeight` to prevent it being closed prematurely (e.g. before the recipient has collected funds) by the sender.

So we have:

- \[off-chain\] - Two parties agree to a series of transactions (for instance as part of file retrieval) with party **A** paying party **B** up to some **total** sum of Filecoin over time.
- \[on-chain\] - The {{<sref payment_channel_actor>}} is used called by A to open a payment channel `from` A `to` B and a lane is opened to increase the `balance` of the channel, triggering a transaction between A and the payment channel actor.
At any time, A can open new lanes to increase the total balance available in the channel (e.g. if A and B choose to do more transactions together).
- \[off-chain\] - Throughout the transaction cycle (e.g. on every piece of data sent via a retrieval deal), party A sends a voucher to party B enabling B to redeem more payment from the payment lanes, and incentivizing B to continue providing a service (e.g. sending more data along).
- \[on-chain\] - At regular intervals, B can `Update` the payment channel balance available `ToSend` with the vouchers received (past their `timeLock`), decreasing the remaining `Value` of the payment channel.
- \[on-chain\] - B can choose to `Collect` the amount `ToSend` triggering a payment between the payment channel actor and B.
- \[on-chain\] - At the end of the cycle, past the `MinCloseHeight`, A can choose to `Close` the payment channel.