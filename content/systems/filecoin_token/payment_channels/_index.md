---
title: Payment Channels
weight: 2
bookCollapseSection: true
dashboardAudit: 1
dashboardState: Stable
dashboardInterface: wip
---

# Payment Channels
---

Payment Channels are used in the Filecoin [Retrieval Market](retrieval_market) to enable efficient off-chain payments and accounting between parties for what is expected to be a series of microtransactions, as these occur during data retrieval.

In particular, given that there is no proving method provided for the act of sending data from a provider (miner) to a client, there is no trust anchor between the two. Therefore, in order to avoid mis-behaviour, Filecoin is making use of payment channels in order to realise a step-wise "data transfer <-> payment" relationship between the provider and the recipient of data. Clients issue requests for data that miners are responding to. The miner is entitled to ask for interim payments, the volume-oriented interval for which is agreed in the Deal phase. In order to facilitate this process, the Filecoin client is creating a payment channel once the provider has agreed on the proposed deal. The client should also lock in the payment channel monetary value equal to the one needed for retrieval of the entire block of data requested. Every time a provider is completing transfer of the pre-specified amount of data, they can request a payment. The client is responding to this payment with a voucher which the provider can redeem (immediately or later).


Redeeming a voucher is not transferring funds from the payment channel to the provider's account, that is, it is not an on-chain operation. Instead, redeeming a voucher denotes the fact that data volume worth of `voucher_val` has been transferred to the client and the corresponding value is locked on the provider's side of payment channel. It is not until the voucher is submitted on-chain through `Collect` that the funds are dispatched to the provider's account.

The following provides a high-level overview of payment channels as these are implemented in Filecoin and an accompanying interface. The lotus implementation of [vouchers](https://github.com/filecoin-project/lotus/blob/master/chain/types/voucher.go) and [payment channels](https://github.com/filecoin-project/lotus/blob/master/paychmgr/paych.go) are also good references.

You can also read more about the [Filecoin payment channel actor interface](payment_channel_actor).

In short, the payment channel actor can be used to open long-lived, flexible payment channels between users. Each channel can be funded by adding to their balance. 
The goal of the payment channel actor is to enable a series of off-chain microtransactions to be reconciled on-chain at a later time with fewer messages. 

The structure of the voucher looks like this:

```text
From	: address of recipient of data (Client)
To		: address of sender of data (Provider)
Value	: value of service (agreed on Deal)
Nonce	: lane- and voucher-specific identifier
```

The client (`From`) is sending vouchers to the provider (`To`). The `Value` indicated in the voucher is progressively increasing together with the `Nonce` to indicate the value available to the provider to redeem, based on the data that the latter has sent. The provider can `Update` the balance of the channel and the balance `ToSend` to them. The provider can choose to `Collect` this balance at any time by submitting the redeemed vouchers on chain and incurring the corresponding gas cost.


In addition, payment channels in Filecoin can be split into `lane`s created as part of updating the channel state with a payment `voucher`. Each lane has an associated `nonce` and amount of tokens it can be `redeemed` for. These lanes allow for a lot of accounting between parties to be done off-chain and reconciled via single updates to the payment channel.

One example of such accounting is *merging of lanes*. When a pair of client-provider nodes have a payment channel established between them with many lanes, the provider will have to pay gas cost for each one of the lanes in order to `Collect` funds. Alternatively, the provider can send a "merge" request to the client to ask it to merge some of the lanes and consolidate the funds. As an incentive for the client to accept the merge lane request, the provider can ask for lower fees to balance out the gas cost. For example, if the provider has collected vouchers worth of 10 monetary units from two lanes, say 5 from each, and the gas cost is 2, then it can ask for 9 from the client if the latter accepts to merge the two lanes.

Over the course of a transaction cycle, each participant in the payment channel can send `voucher`s to other participants. The payment channel's client (`From`) account holder will send a signed voucher with a given nonce to the provider (`To`) account holder.

For instance if the client (`From` address) has sent to the provider (`To` address) the following three vouchers (voucher_val, voucher_nonce) for a lane with 100 monetary units to be redeemed: (10, 1), (20, 2), (30, 3), then the provider could choose to redeem (30, 3) bringing the lane's value to 70 (100 - 30) and cancelling the preceding vouchers, i.e., they would not be able to redeem (10, 1) or (20, 2) anymore. However, they could redeem (20, 2), that is, 20 monetary units, and then follow up with (30, 3) to redeem the remaining 10 monetary units later.

The multiple lanes enable two parties to use a single payment channel to adjudicate multiple independent sets of payments.

Vouchers are signed by the sender and authenticated using a `Secret`, `PreImage` pair provided by the paying party. If the `PreImage` is indeed a pre-image of the `Secret` when used as input to some given algorithm (typically a one-way function like a hash), the `Voucher` is valid. The `Voucher` itself contains the `PreImage` but not the `Secret` (communicated separately to the receiving party). This enables multi-hop payments since an intermediary cannot redeem a voucher on their own. Vouchers can also be used to update the minimum height at which a channel will be settled (i.e., closed), or have `TimeLock`s to prevent recipients (`To`) of the voucher from using them (redeeming them) too early. A channel can also have a `MinCloseHeight` to prevent it being closed prematurely (e.g. before the recipient has collected funds) by the sender.

Once their transactions have completed, either party can choose to `Settle` (i.e., close) the channel. There is a two-hour period after `Settle` during which the recipient can submit any outstanding vouchers. Once the vouchers are submitted, the recipient can then `Collect` the `ToPay` amount from the channel. `From` will be refunded the remaining balance in the channel (if any).

Summarising we have the following set of actions and their relation to the chain:

- \[off-chain\] - Two parties agree to a series of transactions (for instance as part of file retrieval) with party **A** paying party **B** up to some **total** sum of Filecoin over time.
- \[on-chain\] - The [Payment Channel Actor](payment_channel_actor.md) is used, called by A, to open a payment channel `from` A `to` B and a lane is opened to increase the `balance` of the channel, triggering a transaction between A and the payment channel actor.
At any time, A can open new lanes to increase the total balance available in the channel (e.g. if A and B choose to do more transactions together).
- \[off-chain\] - Throughout the transaction cycle (e.g. on every piece of data sent via a retrieval deal), party A sends a voucher to party B enabling B to redeem payment from the payment lanes, and incentivizing B to continue providing a service (e.g. sending more data along).
- \[on-chain\] - At regular intervals, B can redeem vouchers and `Update` the payment channel balance available `ToSend` with the vouchers received (past their `timeLock`), decreasing the remaining `Value` of the payment channel.
- \[on-chain\] - At the end of the cycle, past the `MinCloseHeight`, A can choose to `Settle` the payment channel.
- \[off-chain\] - B has a 2hr period to submit any outstanding vouchers after the channel has been `Settled`, after which period B will lose any monetary value that corresponds to non-submitted vouchers.
- \[on-chain\] - B can choose to `Collect` the amount `ToSend` triggering a payment between the payment channel actor and B.
