---
title: Payment Channels
weight: 3
bookCollapseSection: true
dashboardWeight: 1
dashboardState: stable
dashboardAudit: wip
dashboardTests: 0
---

# Payment Channels

Payment channels are generally used as a mechanism to increase the scalability of blockchains and enable users to transact without involving (i.e., publishing their transactions on) the blockchain, which: i) increases the load of the system, and ii) incurs gas costs for the user. Payment channels generally use a smart contract as an agreement between the two participants. In the Filecoin blockchain Payment Channels are realised by the `paychActor`.

The goal of the Payment Channel Actor specified here is to enable a series of off-chain microtransactions for applications built on top of Filecoin to be reconciled on-chain at a later time with fewer messages that involve the blockchain. Payment channels are already used in the Retrieval Market of the Filecoin Network, but their applicability is not constrained within this use-case only. Hence, here, we provide a detailed description of Payment Channels in the Filecoin network and then describe how Payment Channels are used in the specific case of the Filecoin Retrieval Market.

The payment channel actor can be used to open long-lived, flexible payment channels between users. Filecoin payment channels are _uni-directional_ and can be funded by adding to their balance. Given the context of _uni-directional_ payment channels, we define the **payment channel sender** as the party that receives some service, creates the channel, deposits funds and _sends_ payments (hence the term _payment channel sender_). The _payment channel recipient_, on the other hand is defined as the party that provides services and _receives payment_ for the services delivered (hence the term _payment channel recipient_). The fact that payment channels are uni-directional means that only the payment channel sender can add funds and only the recipient can receive funds. Payment channels are identified by a unique address, as is the case with all Filecoin actors.

The payment channel state structure looks like this:

```go
// A given payment channel actor is established by From (the receipent of a service)
// to enable off-chain microtransactions to To (the provider of a service) to be reconciled
// and tallied on chain.
type State struct {
	// Channel owner, who has created and funded the actor - the channel sender
	From addr.Address
	// Recipient of payouts from channel
	To addr.Address

	// Amount successfully redeemed through the payment channel, paid out on `Collect()`
	ToSend abi.TokenAmount

	// Height at which the channel can be `Collected`
	SettlingAt abi.ChainEpoch
	// Height before which the channel `ToSend` cannot be collected
	MinSettleHeight abi.ChainEpoch

	// Collections of lane states for the channel, maintained in ID order.
	LaneStates []*LaneState
}
```

Before continuing with the details of the Payment Channel and its components and features, it is worth defining a few terms.

- Voucher: a signed message created by either of the two channel parties that updates the channel balance. To differentiate from the **payment channel sender/recipient**, we refer to the voucher parties as **voucher sender/recipient**, who might or might not be the same as the payment channel ones (i.e., the voucher sender might be either the payment channel recipient or the payment channel sender).
- Redeeming a voucher: the voucher MUST be submitted on-chain by the opposite party from the one that created it. Redeeming a voucher does not trigger movement of funds from the channel to the recipient's account, but it does incur message/gas costs. A voucher can be redeemed at any time up to `Collect` (see below), as long as it has a higher `Nonce` than previously submitted vouchers.
- `UpdateChannelState`: this is the process by which a voucher is redeemed, i.e., a voucher is submitted (but not cashed-out) on-chain.
- `Settle`: this process starts closing the channel. It can be called by either the channel creator (sender) or the channel recipient.
- `Collect`: with this process funds are eventually transferred from the payment channel sender to the payment channel recipient. This process incurs message/gas costs.

## Vouchers

Traditionally, in order to transact through a Payment Channel, the payment channel parties send to each other signed messages that update the balance of the channel. In Filecoin, these signed messages are called _vouchers_.

Throughout the interaction between the two parties, the channel sender (`From` address) is sending vouchers to the recipient (`To` address). The `Value` included in the voucher indicates the value available for the receiving party to _redeem_. The `Value` is based on the service that the _payment channel recipient_ has provided to the _payment channel sender_. Either the _payment channel recipient_ or the _payment channel sender_ can `Update` the balance of the channel and the balance `ToSend` to the _payment channel recipient_ (using a voucher), but the `Update` (i.e., the voucher) has to be accepted by the other party before funds can be collected. Furthermore, the voucher has to be redeemed by the opposite party from the one that issued the voucher. The _payment channel recipient_ can choose to `Collect` this balance at any time incurring the corresponding gas cost.

Redeeming a voucher does not transfer funds from the payment channel to the recipient's account. Instead, redeeming a voucher attests that some service of worth `Value` has been provided by the payment channel recipient to the payment channel sender. It is not until the whole payment channel is _collected_ that the funds are dispatched to the recipient's account.

This is the structure of the voucher:

```go
// A voucher can be created and sent by any of the two parties. The `To` payment channel address can redeem the voucher and then `Collect` the funds.
type SignedVoucher struct {
	// ChannelAddr is the address of the payment channel this signed voucher is valid for
	ChannelAddr addr.Address
	// TimeLockMin sets a min epoch before which the voucher cannot be redeemed
	TimeLockMin abi.ChainEpoch
	// TimeLockMax sets a max epoch beyond which the voucher cannot be redeemed
	// TimeLockMax set to 0 means no timeout
	TimeLockMax abi.ChainEpoch
	// (optional) The SecretPreImage is used by `To` to validate
	SecretPreimage []byte
	// (optional) Extra can be specified by `From` to add a verification method to the voucher
	Extra *ModVerifyParams
	// Specifies which lane the Voucher is added to (will be created if does not exist)
	Lane uint64
	// Nonce is set by `From` to prevent redemption of stale vouchers on a lane
	Nonce uint64
	// Amount voucher can be redeemed for
	Amount big.Int
	// (optional) MinSettleHeight can extend channel MinSettleHeight if needed
	MinSettleHeight abi.ChainEpoch

	// (optional) Set of lanes to be merged into `Lane`
	Merges []Merge

	// Sender's signature over the voucher
	Signature *crypto.Signature
}
```

Over the course of a transaction cycle, each participant in the payment channel can send `Voucher`s to the other participant.

For instance, if the payment channel sender (`From` address) has sent to the payment channel recipient (`To` address) the following three vouchers `(voucher_val, voucher_nonce)` for a lane with 100 FIL to be redeemed: (10, 1), (20, 2), (30, 3), then the recipient could choose to redeem (30, 3) bringing the lane's value to 70 (100 - 30) and cancelling the preceding vouchers, i.e., they would not be able to redeem (10, 1) or (20, 2) anymore. However, they could redeem (20, 2), that is, 20 FIL, and then follow up with (30, 3) to redeem the remaining 10 FIL later.

It is worth highlighting that while the `Nonce` is a strictly increasing value to denote the sequence of vouchers issued within the remit of a payment channel, the `Value` is not a strictly increasing value. Decreasing `Value` (although expected rarely) can be realized in cases of refunds that need to flow in the direction from the payment channel recipient to the payment channel sender. This can be the case when some bits arrive corrupted in the case of file retrieval, for instance.

Vouchers are signed by the party that creates them and are authenticated using a (`Secret`, `PreImage`) pair provided by the paying party (channel sender). If the `PreImage` is indeed a pre-image of the `Secret` when used as input to some given algorithm (typically a one-way function like a hash), the `Voucher` is valid. The `Voucher` itself contains the `PreImage` but not the `Secret` (communicated separately to the receiving party). This enables multi-hop payments since an intermediary cannot redeem a voucher on their own. Vouchers can also be used to update the minimum height at which a channel will be settled (i.e., closed), or have `TimeLock`s to prevent voucher recipients from redeeming them too early. A channel can also have a `MinCloseHeight` to prevent it being closed prematurely (e.g. before the payment channel recipient has collected funds) by the payment channel creator/sender.

Once their transactions have completed, either party can choose to `Settle` (i.e., close) the channel. There is a 12hr period after `Settle` during which either party can submit any outstanding vouchers. Once the vouchers are submitted, either party can then call `Collect`. This will send the payment channel recipient the `ToPay` amount from the channel, and the channel sender (`From` address) will be refunded the remaining balance in the channel (if any).

## Lanes

In addition, payment channels in Filecoin can be split into `lane`s created as part of updating the channel state with a payment `voucher`. Each lane has an associated `nonce` and amount of tokens it can be `redeemed` for. Lanes can be thought of as transactions for several different services provided by the channel recipient to the channel sender. The `nonce` plays the role of a sequence number of vouchers within a given lane, where a voucher with a higher nonce replaces a voucher with a lower nonce.

Payment channel lanes allow for a lot of accounting between parties to be done off-chain and reconciled via single updates to the payment channel. The multiple lanes enable two parties to use a single payment channel to adjudicate multiple independent sets of payments.

One example of such accounting is _merging of lanes_. When a pair of channel sender-recipient nodes have a payment channel established between them with many lanes, the channel recipient will have to pay gas cost for each one of the lanes in order to `Collect` funds. Merging of lanes allow the channel recipient to send a "merge" request to the channel sender to request merging of (some of the) lanes and consolidate the funds. This way, the recipient can reduce the overall gas cost. As an incentive for the channel sender to accept the merge lane request, the channel recipient can ask for a lower total value to balance out the gas cost. For instance, if the recipient has collected vouchers worth of 10 FIL from two lanes, say 5 from each, and the gas cost of submitting the vouchers for these funds is 2, then it can ask for 9 from the creator if the latter accepts to merge the two lanes. This way, the channel sender pays less overall for the services it received and the channel recipient pays less gas cost to submit the voucher for the services they provided.

## Lifecycle of a Payment Channel

Summarising, we have the following sequence:

0. Two parties agree to a series of transactions (for instance as part of file retrieval) with one party paying the other party up to some _total_ sum of Filecoin over time. This is part of the deal-phase, it takes place off-chain and does not (at this stage) involve payment channels.
1. The Payment Channel Actor is used, called by the payment channel sender (who is the recipient of some service, e.g., file in case of file retrieval) to create the payment channel and deposit funds.
2. Any of the two parties can create vouchers to send to the other party.
3. The voucher recipient saves the voucher locally. Each voucher has to be submitted by the opposite party from the one that created the voucher.
4. Either immediately or later, the voucher recipient "redeems" the voucher by submitting it to the chain, calling `UpdateChannelState`
5. The channel sender or the channel recipient `Settle` the payment channel.
6. 12-hour period to close the channel begins.
7. If any of the two parties have outstanding (i.e., non-redeemed) vouchers, they should now submit the vouchers to the chain (there should be the option of this being done automatically). If the channel recipient so desires, they should send a "merge lanes" request to the sender.
8. 12-hour period ends.
9. Either the channel sender or the channel recipient calls `Collect`.
10. Funds are transferred to the channel recipient's account and any unclaimed balance goes back to channel sender.

## Payment Channels as part of the Filecoin Retrieval

Payment Channels are used in the Filecoin [Retrieval Market](retrieval_market) to enable efficient off-chain payments and accounting between parties for what is expected to be a series of microtransactions, as these occur during data retrieval.

In particular, given that there is no proving method provided for the act of sending data from a provider (miner) to a client, there is no trust anchor between the two. Therefore, in order to avoid mis-behaviour, Filecoin is making use of payment channels in order to realise a step-wise "data transfer <-> payment" relationship between the data provider and the client (data receiver). Clients issue requests for data that miners are responding to. The miner is entitled to ask for interim payments, the volume-oriented interval for which is agreed in the Deal phase. In order to facilitate this process, the Filecoin client is creating a payment channel once the provider has agreed on the proposed deal. The client should also lock monetary value in the payment channel equal to the amount needed for retrieval of the entire block of data requested. Every time a provider is completing transfer of the pre-specified amount of data, they can request a payment. The client responds to this request with a voucher which the provider can redeem (immediately or later), as per the process described earlier.

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/paychmgr/paych.go" lang="go" title="Payment Channel Implementation">}}

{{<embed src="https://github.com/filecoin-project/specs-actors/blob/v0.9.12/actors/builtin/paych/paych_actor.go" lang="go" symbol="SignedVoucher">}}

{{<embed src="https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/paych/paych_actor.go" lang="go" title="Payment Channel Actor">}}
