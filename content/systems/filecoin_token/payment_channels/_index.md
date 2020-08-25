---
title: Payment Channels
weight: 2
bookCollapseSection: true
dashboardWeight: 1
dashboardState: incorrect
dashboardAudit: 0
dashboardTests: 0
---

# Payment Channels
---

Payment channels are generally used as an mechanism to increase the scalability of blockchains and enable users to transact without involving (i.e., publishing their transactions on) the blockchain, which: i) increases the load of the system, and ii) incurs gas costs for the user. Payment channels generally use a smart contract as an agreement between the two participants. In the Filecoin blockchain Payment Channels are realised by the `paych actor`.

The goal of the Payment Channel Actor specified here is to enable a series of off-chain microtransactions for applications built on top of Filecoin to be reconciled on-chain at a later time with fewer messages that involve the blockchain. Payment channels are already used in the Retrieval Market of the Filecoin Network, but their applicability is not constrained within this use-case only. Hence, here, we provide a detailed description of Payment Channels in the Filecoin network and then describe how Payment Channels are used in the specific case of the Filecoin Retrieval Market.

The payment channel actor can be used to open long-lived, flexible payment channels between users. Filecoin payment channels are _uni-directional_ and can be funded by adding to their balance. Given the context of _uni-directional_ payment channels, we define the sender as the party that receives some service, creates the channel, deposits funds and sends payments. The recipient, on the other hand is defined as the party that provides services and receives payment for the services delivered. The fact that payment channels are uni-directional means that only the sender can add funds and the recipient can receive funds. Payment channels are identified by a unique address, as is the case with all Filecoin actors.

The payment channel state structure looks like this:

```go
// A given payment channel actor is established by From (normally the client/receipent of data)
// to enable off-chain microtransactions to To (normally provider/sender of data) to be reconciled
// and tallied on chain.
type State struct {
	// Channel owner, who has created and funded the actor - the sender
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

- Voucher: a signed message from the creator of the channel (sender) to the recipient that updates the channel balance.
- Redeeming a voucher: the voucher is submitted on-chain by the channel's recipient, but does not trigger movement of funds from the channel to the recipient's account. This action does not incur transaction/gas costs.
- `UpdateChannelState`: this is the process by which a voucher is redeemed, i.e., a voucher is submitted (but not cashed-out) on-chain.
- `Settle`: this process starts closing the channel. It can be called by either the channel creator (sender) or the channel recipient.
- `Collect`: with this process funds are eventually transferred from the sender to the recipient. This process incurs transaction/gas costs.

Throughout the interaction between the two parties, the channel sender (`From` address) is sending vouchers to the recipient (`To` address). The `Value` indicated in the voucher, from the sender to the recipient is progressively increasing (with every new voucher issued) to indicate the value available for the provider to _redeem_. The `Value` is based on the service that the recipient has provided to the sender of the channel. The recipient can `Update` the balance of the channel and the balance `ToSend` to them. The recipient can choose to `Collect` this balance at any time incurring the corresponding gas cost.

## Vouchers

Traditionally, in order to transact through a Payment Channel parties send to each other signed messages that update the balance of the channel. In Filecoin, these signed messages are called _vouchers_.

Redeeming a voucher is not transferring funds from the payment channel to the recipient's account. Instead, redeeming a voucher denotes the fact that some service worth of `voucher_val` has been provided by the recipient of the payment channel to the sender. It is not until the voucher is _collected_ that the funds are dispatched to the provider's account.

This is the structure of the voucher:

```go
// A voucher is sent by `From` to `To` off-chain in order to enable
// `To` to redeem payments on-chain in the future
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

Over the course of a transaction cycle, each participant in the payment channel can send `voucher`s to the other participant. The payment channel's sender (`From`) account holder will send a signed voucher with a given nonce to the receipient (`To`) account holder.

For instance, if the sender (`From` address) has sent to the recipient (`To` address) the following three vouchers (voucher_val, voucher_nonce) for a lane with 100 FIL to be redeemed: (10, 1), (20, 2), (30, 3), then the recipient could choose to redeem (30, 3) bringing the lane's value to 70 (100 - 30) and cancelling the preceding vouchers, i.e., they would not be able to redeem (10, 1) or (20, 2) anymore. However, they could redeem (20, 2), that is, 20 FIL, and then follow up with (30, 3) to redeem the remaining 10 FIL later.

Vouchers are signed by the voucher sender and authenticated using a `Secret`, `PreImage` pair provided by the paying party (channel sender). If the `PreImage` is indeed a pre-image of the `Secret` when used as input to some given algorithm (typically a one-way function like a hash), the `Voucher` is valid. The `Voucher` itself contains the `PreImage` but not the `Secret` (communicated separately to the receiving party). This enables multi-hop payments since an intermediary cannot redeem a voucher on their own. Vouchers can also be used to update the minimum height at which a channel will be settled (i.e., closed), or have `TimeLock`s to prevent recipients (`To`) of the voucher from using them (redeeming them) too early. A channel can also have a `MinCloseHeight` to prevent it being closed prematurely (e.g. before the recipient has collected funds) by the sender.

Once their transactions have completed, either party can choose to `Settle` (i.e., close) the channel. There is a 12hr period after `Settle` during which either party can submit any outstanding vouchers. Once the vouchers are submitted, either party can then call `Collect`. This will send the recipient the `ToPay` amount from the channel, and the channel sender (`From` address) will be refunded the remaining balance in the channel (if any).


## Lanes

In addition, payment channels in Filecoin can be split into `lane`s created as part of updating the channel state with a payment `voucher`. Each lane has an associated `nonce` and amount of tokens it can be `redeemed` for. Lanes can be thought of as transactions for several different services provided by the channel recipient to the channel sender. The `nonce` plays the role of a sequence number of vouchers within a given lane, where a voucher with a higher nonce replaces a voucher with a lower nonce.

Payment channel lanes allow for a lot of accounting between parties to be done off-chain and reconciled via single updates to the payment channel. The multiple lanes enable two parties to use a single payment channel to adjudicate multiple independent sets of payments.

One example of such accounting is *merging of lanes*. When a pair of channel sender-recipient nodes have a payment channel established between them with many lanes, the channel recipient will have to pay gas cost for each one of the lanes in order to `Collect` funds. Merging of lanes allow the recipient to send a "merge" request to the channel sender to request merging of (some of the) lanes and consolidate the funds. This way, the recipient can reduce the overall gas cost. As an incentive for the channel sender to accept the merge lane request, the channel recipient can ask for a lower total value to balance out the gas cost. For instance, if the recipient has collected vouchers worth of 10 FIL from two lanes, say 5 from each, and the gas cost of submitting the vouchers for these funds is 2, then it can ask for 9 from the creator if the latter accepts to merge the two lanes. This way, the channel sender pays less overall for the services it received and the channel recipient pays less gas cost to submit the voucher for the services they provided.

## Lifecycle of a Payment Channel

Summarising, we have the following sequence:

0. Two parties agree to a series of transactions (for instance as part of file retrieval) with one party paying the other party up to some _total_ sum of Filecoin over time. This is part of the deal-phase, it takes place off-chain and does not (at this stage) involve payment channels.
1. The [Payment Channel Actor](payment_channel_actor.md) is used, called the sender to create the payment channel and deposits funds.
2. The sender creates a voucher and sends the voucher to the recipient.
3. The recipient saves the voucher locally.
4. Either immediately or later, the recipient "redeems" the voucher by submitting it to the chain, calling `UpdateChannelState`
5. The sender or the recipient `Settle`.
6. 12-hour period to close the channel begins.
7. If the recipient has not already done so, they should now submit the vouchers to the chain (there should be the option of thing being done automatically). If the recipient so desires, they should send a "merge lanes" request to the sender.
8. 12-hour period ends.
9. Either the sender or the recipient calls `Collect`
10. Funds are transferred to the recipient and any unclaimed balance goes back to sender.


## Payment Channels as part of the Filecoin Retrieval
 
Payment Channels are used in the Filecoin [Retrieval Market](retrieval_market) to enable efficient off-chain payments and accounting between parties for what is expected to be a series of microtransactions, as these occur during data retrieval.

In particular, given that there is no proving method provided for the act of sending data from a provider (miner) to a client, there is no trust anchor between the two. Therefore, in order to avoid mis-behaviour, Filecoin is making use of payment channels in order to realise a step-wise "data transfer <-> payment" relationship between the data provider and the client (data receiver). Clients issue requests for data that miners are responding to. The miner is entitled to ask for interim payments, the volume-oriented interval for which is agreed in the Deal phase. In order to facilitate this process, the Filecoin client is creating a payment channel once the provider has agreed on the proposed deal. The client should also lock monetary value in the payment channel equal to the one needed for retrieval of the entire block of data requested. Every time a provider is completing transfer of the pre-specified amount of data, they can request a payment. The client is responding to this payment with a voucher which the provider can redeem (immediately or later), as per the process described earlier.


The following provides a high-level overview of payment channels as these are implemented in Filecoin and an accompanying interface. The lotus implementation of [vouchers](https://github.com/filecoin-project/lotus/blob/master/chain/types/voucher.go) and [payment channels](https://github.com/filecoin-project/lotus/blob/master/paychmgr/paych.go) are also good references.

You can also read more about the [Filecoin payment channel actor interface](payment_channel_actor).
