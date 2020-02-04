---
title: Payment Channel Actor
---
{{<label payment_channel_actor>}}

{{< readfile file="payment_channel_actor.id" code="true" lang="go" >}}

The goal of the payment channel actor is to enable a series of off-chain microtransactions to be reconciled on-chain at a later time with fewer messages. Accordingly, the expectation is `From` will send to `To` vouchers of successively greater `Value` and increasing `Nonce`s. When they choose to, `To` can `Update` the channel to update the balance available `ToSend` to them in the channel, and can choose to `Collect` this balance at any time (incurring a gas cost).
The channel is split into `lane`s created as part of updating the channel state with a payment `voucher`. Each lane has an associated `nonce` and amount of tokens it can be `redeemed` for. These lanes allow for a lot of accounting between parties to be done off chain and reconciled via single updates to the payment channel, merging these lanes to arrive at a desired outcome.

Over the course of a transaction cycle, each party to the payment channel can send the other `voucher`s. The payment channel's  `From` account holder will send a signed voucher with a given nonce to the `To` account holder. The latter can use the voucher to `redeem` part of the lane's value, merging other lanes into it as needed.

For instance if `From` sends `To` the following vouchers (voucher_val, voucher_nonce) for a lane with 100 to be redeemed: (10, 1), (20, 2), (30, 3), then `To` could choose to redeem (30, 3) bringing the lane's value to 70 (100 - 30). They could not redeem (10, 1) or (20, 2) thereafter. They could however redeem (20, 2) for 20, and then (30, 3) for 10 (30 - 20) thereafter.

The multiple lanes enable two parties to use a single payment channel to adjudicate multiple independent sets of payments.

Vouchers are signed by both parties (i.e. explicitly by the sender and implicitly by the recipient submitting it on chain) and authenticated using a `Secret`, `PreImage` pair provided by the paying party. If the `PreImage` is indeed a pre-image of the `Secret` when used as input to some given algorithm (typically a one-way function like a hash), the `Voucher` is valid. The `Voucher` itself contains the `PreImage` but not the `Secret` (communicated separately to the receiving party). This enables multi-hop payments since an intermediary cannot redeem a voucher on their own. They can also be used to update the minimum height at which a channel will be settled. Likewise, vouchers can have `TimeLock`s to prevent their being used too early, likewise a channel can have a `MinSettleHeight` to prevent it being settled prematurely (e.g. before the recipient has collected funds) by the sender.

Once their transactions have completed, either party can choose to `Settle` the channel, the recipient can then `Collect` the `ToPay` amount from the channel. `From` will be refunded the remaining balance in the channel.