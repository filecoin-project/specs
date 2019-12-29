---
menuTitle: "Payment Channels"
statusIcon: 🔁
title: "Payment Channels"
entries:
- payment_channel_actor
---

{{<label payment_channels>}}

Payment Channels are used in the Filecoin {{<sref retrieval_market>}} in order to enable efficient off-chain payments between parties for what are expected to be series of microtransactions, specifically those occuring as part of retrieval market data retrieval.

Note that a lot of the following is based on the [Lotus state channel](https://github.com/filecoin-project/lotus/blob/master/chain/actors/actor_paych.go) and [voucher](https://github.com/filecoin-project/lotus/blob/master/chain/types/voucher.go) implementations.

You can also read more about the {{<sref payment_channel_actor "Filecoin payment channel actor">}}.

Each payment channel has a `Secret` and an associated set of `lanes` associated to it, with each lane having an associated `nonce` and amount of tokens it can be `redeemed` for (their relative values determine a payout schedule, e.g. linear if all amounts are the same).

The `Secret` is generated from an initial value through successive hashing, with each lane standing in for a specific pre-image of the `Secret`.
Throughout the channel's lifecycle, the sender will send successive vouchers to the recipient. Each signed voucher has:

- a `secretPreimage`,
- a set of lanes to be merged,
- a `minCloseHeight` before which the associated channel cannot be closed (serving as a timelock),
- a unique `nonce`.
    
Each voucher encodes a set of lanes to be merged. Each lane merge is a payment, through which part of the payment channel's balance will eventually go to its recipient. Its `secretPreimage` is then used along with a hash function (`SHA256`) to generate the channel's `Secret`. The `Secret` cannot be generated without the `PreImage` given a secure Collision-Resistant Hash Function. If the voucher is thus proven correct, the `amount` redeemed in all of its lanes is added to the payment channel's `ToPay` total.
The sender sends successive vouchers (starting at `nonce = 1`) eachh of which merges more lanes, enabling the recipient to `UpdateChannelState` with the latest received nonce and prove that the sender has enabled it to merge payment channels. This update requires validating:

- current the voucher's preimage, hashed `nonce` times corresponds to the payment channel `Secret`
- the voucher's nonce is not smaller than the current largest nonce seen by the payment channel
- there are enough funds in the channel to enable the voucher's lanes to merge.

On channel `Close`, the recipient can collect the `ToPay` amount from the channel.

So we have:

- \[off-chain\] - Two parties agree to a series of transactions (for instance as part of file retrieval) with party **A** paying party **B** up to some **total** sum of Filecoin over time.
- \[on-chain\] - The {{<sref payment_channel_actor>}} is used called by A to open a payment channel `from` A `to` B and whose `balance` is funded by A with total Filecoin. The channel also tracks:

    - an amount `ToSend`
    - a `closingAt` block height
    - a `minClosingHeight` before which it cannot be closed, determined by the incoming vouchers
    - a map of its `laneStates` keeping track of payment lane states, 
    - a `Secret` against which vouchers are checked for merging lanes
- \[off-chain\] - Throughout the transaction cycle (e.g. on every piece of data sent via a retrieval deal), party A sends a voucher to party B enabling B to redeem more payment from the payment channel, and incentivizing B to continue providing a service (e.g. sending more data along).
- \[on-chain\] - At regular intervals, either A or B can `Update` the state channel with the vouchers to redeem them (merging the payment lanes) and increase the amount `ToSend`.
- \[on-chain\] - At the end of the cycle, either A or B can `Close` the state channel to effectively transfer `ToSend` amount to B and the remainder of the state channel's balance (if any) back to A.
