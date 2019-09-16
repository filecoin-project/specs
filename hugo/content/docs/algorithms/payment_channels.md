---
title: "Payment Channels"
---

## Payment Channels

In order for the Filecoin Markets to work in a timely manner, we need to be able to have off-chain payments. This is a solved problem (at least, for our purposes in v0). Payment channels have been implemented and used in bitcoin, ethereum and many other networks.

The basic premise is this: User A wants to be able to send many small payments to user B. So user A locks up money in a contract that says "this money will only go to user B, and the unclaimed amount will be returned to user A after a set time period". Once that money is locked up, user A can send user B signed transactions that user B can cash out at any time.

For example:

- User A locks up 10 FIL to B
- User B does something for A
- User A sends `SignedVoucher{Channel, 1 FIL}` to B
- User B does something for user A
- User A sends `SignedVoucher{Channel, 2 FIL}` to B

At this point, B has two signed messages from A, but the contract is set up such that it can only be cashed out once. So if B decided to cash out, they would obviously select the message with the higher value. Also, once B cashes out, they must not accept any more payments from A on that same channel.

### Multi-Lane Payment Channel

The filecoin storage market may require a way to do incremental payments between two parties, over time, for multiple different transactions. The primary motivating usecase for this is to provide payment for file storage over time, for each file stored. An additional requirement is the ability to have less than one message on chain per transaction 'lane', meaning that payments for multiple files should be aggregateable (Note: its okay if this aggregation is an interactive process).

Let's say that `A` wants to make such an arrangement with `B`. `A` should create the payment channel with enough funds to cover all potential transactions. Then `A` decides to start the first transaction, so they send a signed voucher for the payment channel on 'lane 1', for 2 FIL. They can then send more updates on lane 1 as needed. Then, at some point `A` decides to start another independent transaction to `B`, so they send a voucher on 'lane 2'. The voucher for lane 2 can be cashed out independently of lane 1. However, `B` can ask `A` to 'reconcile' the two payment channels for them into a single update. This update could contain a value, and a list of lanes to close. Cashing out that reconciled update would invalidate the other lanes, meaning `B` couldnt also cash in those. The single update would be much smaller, and therefore cheaper to close out.

Lane state can be easily tracked on-chain with a compact bitfield.

### Payment Channel Reconciliation

In a situation where peers A and B  have several different payment channels between them, the scenario may frequently come up where A has multiple payment channel updates from B to apply. Submitting each of these individually would cost a noticeable amount in fees, and put excess unnecessary load on the chain. To remedy this, A can contact B and ask them for a single payment channel update for the combined value of all the updates they have (minus some fee to incent B to actually want to do this). This aggregated update would contain a list of the IDs of the other payment channels that it is superceding so that A cannot also cash out on the originals.

# Payment Reconciliation

The filecoin storage market will (likely) have many independent payments between the same parties. These payments will be secured through payment channels, set up initially on chain, but utilized almost entirely off-chain. The point at which they need to touch the chain is when miners wish to cash out their earnings. A naive solution to this problem would have miners perform one on-chain action per file stored for a particular client. This would not scale well. Instead, we need a system where the miner and client can have some additional off-chain communication and end up with the miner submitting only a single message to the chain.

To accomplish this, we introduce the Payment Reconciliation Protocol.

This is a libp2p service run by all participants wanting to participate in payment reconciliation. When Alice has a set of payments from Bob that she is ready to cash out, Alice can send a `ReconcileRequest` to Bob, containing the following information:

```sh
type ReconcileRequest struct {
	vouchers [Vouchers]
	reqVal TokenAmount
}
```

The Vouchers should all be valid vouchers from Bob to Alice, on the same payment channel, and they should all be ready to be cashed in. `ReqVal` is a token amount less than or equal to the sum of all the values in the given vouchers. Generally, this value will be between the total sum of the vouchers, and that total sum minus the fees it would cost to submit them all to the chain.

Bob receives this request, and checks that all the fields are correct, and then ensures that the difference between ReqVal and the vouchers sum is sufficient (this is a parameter that the client can set).  Then, he sends back a response which either contains the requested voucher, or an error status and message.

```sh
type ReconcileResponse struct {
	combined Voucher
	status  Status
	message optional String
}

## TODO: what are the possible status cases?
type Status enum {
    | Success
    | Failure
}
```

Open Questions:

- In a number of usecases, this protocol will require the miner look up and connect to a client to propose reconciliation. How does a miner look up and connect to a client over libp2p given only their filecoin address?
- Without repair miners, this protocol will likely not be used that much. Should that be made clear? Should there be other considerations added to compensate?
