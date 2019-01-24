# Payments

### What are payments

### What payments affect

### Dependencies

## Miners Claiming Earnings

Storage miners get paid entirely through payment channels. Payment from a client to a storage miner comes in the form of a set of channel updates that get created when proposing the deal. These updates are each time-locked, and can only be cashed out if the storage miner has not been slashed for the storage that is being paid for. (TODO: working on a multi-lane payment channel construction that should make this all pretty easy, only requiring a single on-chain channel construction between each client and storage miner).

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

### Simple Payment Channel Actor

We need to implement an actor to allow the creation of payment channels between users. The interface for that should look something like this:

```go
type ChannelID *big.Int
type BlockHeight *big.Int
type Signature []byte

type SpendVoucher struct {
    Channel ChannelID
    Amount *TokenAmount
    Sig Signature
}

type PaymentBroker interface {
    // CreateChannel creates a new payment channel from the caller to the target.
    // The value attached to the invocation is used as the deposit, and the channel
    // will expire and return all of its money to the owner after the given block height.
    CreateChannel(target Address, eol BlockHeight) ChannelID
    
    // Update updates the payment channel with the given amounts, and sends the current
    // committed amount to the target. This is useful when you want to checkpoint the 
    // value in a payment, but continue to use the channel afterwards.
    Update(channel ChannelID, amt *TokenAmount, sig Signature)
    
    // Close is called by the target of a payment channel to cash out and close out
    // the payment channel. This is really a courtesy call, as the channel will
    // eventually time out and close on its own.
    Close(channel ChannelID, amt *TokenAmount, sig Signature)
    
    // Extend can be used by the owner of a channel to add more funds to it and
    // extend the channels lifespan.
    Extend(target Address, channel ChannelID, eol BlockHeight)
    
    // Reclaim is used by the owner of a channel to reclaim unspent funds in timed
    // out payment channels they own.
    Reclaim(target Address, channel ChannelID)
}

// MakeSpendVoucher is used by the owner of a channel to create an offline payment
// for the target. Note that any amount may be given, but the target gets to select
// which of the vouchers you give them to cash out, and therefore any rational actor 
// will only ever keep the one with the largest amount. After calling this function,
// you should send the returned SpendVoucher to the target out of band.
func MakeSpendVoucher(ch ChannelID, amt *TokenAmount, sk PrivateKey) *SpendVoucher {
    data := concatBytes(ch, amt)
    sig := sk.Sign(data)
    return &SpendVoucher{
        Channel: ch,
        Amount: amt,
        Sig: sig,
    }
}
```

Channel IDs should be memory efficient and generated is such a way that reordering does not change the channelID. This may be solved by first indexing by the target address internally and then using the nonce of the message that invoked the create channel method.



### Multi-Lane Payment Channel (WIP)

The filecoin storage market may require a way to do incremental payments between two parties, over time, for multiple different transactions. The primary motivating usecase for this is to provide payment for file storage over time, for each file stored. An additional requirement is the ability to have less than one message on chain per transaction 'lane', meaning that payments for multiple files should be aggregateable (Note: its okay if this aggregation is an interactive process).

Let's say that `A` wants to make such an arrangement with `B`. `A` should create the payment channel with enough funds to cover all potential transactions. Then `A` decides to start the first transaction, so they send a signed voucher for the payment channel on 'lane 1', for 2 FIL. They can then send more updates on lane 1 as needed. Then, at some point `A` decides to start another independent transaction to `B`, so they send a voucher on 'lane 2'. The voucher for lane 2 can be cashed out independently of lane 1. However, `B` can ask `A` to 'reconcile' the two payment channels for them into a single update. This update could contain a value, and a list of lanes to close. Cashing out that reconciled update would invalidate the other lanes, meaning `B` couldnt also cash in those. The single update would be much smaller, and therefore cheaper to close out.

Lane state can be easily tracked on-chain with a compact bitfield.

```go
type SpendVoucher struct {
    // Amount is the amount of FIL that this voucher can be redeemed for
    Amount TokenAmount
    
    // Nonce is a number that sets the ordering of vouchers. If you try to redeem
    // a voucher with an equal or lower nonce, the operation will fail. Nonces are
    // per lane.
    Nonce uint64
    
    // Lane specifies which 'lane' of the payment channel this voucher is for.
    // Lanes may be either open or closed, a voucher for a closed lane may not be redeemed
    Lane uint64

    // Merges specifies a list of lane-nonce pairs that this voucher will close.
    // This voucher may not be redeemed if any of the lanes specified here are already
    // closed, or their nonce specified here is lower than the nonce of the lane on-chain.
    Merges []MergePair

    TimeLock uint64

    SecretPreimage []byte

    RequiredSector []byte

    MinCloseHeight uint64
    
    Sig Signature
}

type MergePair struct {
    Lane uint64
    Nonce uint64
}
```

```go
type PaymentChannel struct {
    From Address
    To Address

    ChannelTotal TokenAmount
    ToSend TokenAmount

    ClosingAt uint64
    MinCloseHeight uint64

    LaneStates map[uint64]LaneState
}

type LaneState struct {
    Nonce uint64
    Redeemed TokenAmount
}

func (paych *PaymentChannel) validateSignature(sv SpendVoucher) {
    if msg.From == paych.From {
        ValidateSignature(sv.SerializeNoSig(), sv.Signature, paych.To)
    } else if msg.From == paych.To {
        ValidateSignature(sv.SerializeNoSig(), sv.Signature, paych.From)
    } else {
        Fatal("bad programmer")
    }
}

func (paych *PaymentChannel) UpdateChannelState(sv SpendVoucher, secret []byte) {
    if !paych.validateSignature(sv) {
        Fatal("Signature Invalid")
    }

    if chain.Now() < sv.TimeLock {
        Fatal("cannot use this voucher yet!")
    }

    if sv.SecretPreimage != nil {
        if Hash(secret) != sv.SecretPreimage {
            Fatal("Incorrect secret!")
        }
    }

    if sv.RequiredSector != nil {
        miner, found := GetMiner(msg.From)
        if !found {
            Fatal("Redeemer is not a miner")
        }

        if !miner.HasSector(sv.RequiredSector) {
            Fatal("miner does not have sector, cannot redeem payment")
        }
    }

    ls := paych.LaneStates[sv.Lane]
    if ls.Closed {
        Fatal("cannot redeem a voucher on a closed lane")
    }

    if ls.Nonce > sv.Nonce {
        Fatal("voucher has an outdated nonce, cannot redeem")
    }

    var mergeValue TokenAmount
    for _, merge := range sv.Merges {
        ols := paych.LaneStates[merge.Lane]

        if ols.Nonce >= merge.Nonce {
            Fatal("merge in voucher has outdated nonce, cannot redeem")
        }

        mergeValue += ols.Redeemed
        ols.Nonce = merge.Nonce
    }

    ls.Nonce = sv.Nonce
    balanceDelta = sv.Amount - (mergeValue + ls.Redeemed)
    ls.Redeemed = sv.Amount

    newSendBalance = paych.ToSend + balanceDelta
    if newSendBalance < 0 {
        // TODO: is this impossible?
        Fatal("voucher would leave channel balance negative")
    }

    if newSendBalance > paych.ChannelTotal {
        Fatal("not enough funds in channel to cover voucher")
    }

    paych.ToSend = newSendBalance

    if sv.MinCloseHeight != 0 {
        if paych.ClosingAt < sv.MinCloseHeight {
            paych.ClosingAt = sv.MinCloseHeight
        }
        if paych.MinCloseHeight < sv.MinCloseHeight {
            paych.MinCloseHeight = sv.MinCloseHeight
        }
    }
}

func (paych *PaymentChannel) Withdraw(upTo TokenAmount) {
    // TODO: this ones tricky, withdraw funds without closing it out entirely...
}

func (paych *PaymentChannel) Close() {
    if msg.From != paych.From && msg.From != paych.To {
        Fatal("not authorized to close channel")
    }
    if paych.ClosingAt != 0 {
        Fatal("Channel already closing")
    }

    paych.ClosingAt = chain.Now() + ChannelClosingDelay
    if paych.ClosingAt < paych.MinCloseHeight {
        paych.ClosingAt = paych.MinCloseHeight
    }
}

func (paych *PaymentChannel) Collect() {
    if paych.ClosingAt == 0 {
        Fatal("payment channel not closing or closed")
    }
    if chain.Now() < paych.ClosingAt {
        Fatal("Payment channel not yet closed")
    }
    Transfer(paych.ChannelTotal - paych.ToSend, paych.From)
    Transfer(paych.ToSend, paych.To)
}
```



### Payment Channel Reconciliation

In a situation where peers A and B  have several different payment channels between them, the scenario may frequently come up where A has multiple payment channel updates from B to apply. Submitting each of these individually would cost a noticeable amount in fees, and put excess unnecessary load on the chain. To remedy this, A can contact B and ask them for a single payment channel update for the combined value of all the updates they have (minus some fee to incent B to actually want to do this). This aggregated update would contain a list of the IDs of the other payment channels that it is superceding so that A cannot also cash out on the originals.

# Payment Reconciliation
The filecoin storage market will (likely) have many independent payments between the same parties. These payments will be secured through payment channels, set up initially on chain, but utilized almost entirely off-chain. The point at which they need to touch the chain is when miners wish to cash out their earnings. A naive solution to this problem would have miners perform one on-chain action per file stored for a particular client. This would not scale well. Instead, we need a system where the miner and client can have some additional off-chain communication and end up with the miner submitting only a single message to the chain.

To accomplish this, we introduce the Payment Reconciliation Protocol.

This is a libp2p service run by all participants wanting to participate in payment reconciliation. When Alice has a set of payments from Bob that she is ready to cash out, Alice can send a `ReconcileRequest` to Bob, containing the following information:

```go
type ReconcileRequest struct {
  Vouchers []Vouchers

  ReqVal TokenAmount
}
```

The Vouchers should all be valid vouchers from Bob to Alice, on the same payment channel, and they should all be ready to be cashed in. `ReqVal` is a token amount less than or equal to the sum of all the values in the given vouchers. Generally, this value will be between the total sum of the vouchers, and that total sum minus the fees it would cost to submit them all to the chain.

Bob receives this request, and checks that all the fields are correct, and then ensures that the difference between ReqVal and the vouchers sum is sufficient (this is a parameter that the client can set).  Then, he sends back a response which either contains the requested voucher, or an error status and message.

```go
type ReconcileResponse struct {
  Combined Voucher

  Status StatusCode
  Message string
}
```

Open Questions:

- In a number of usecases, this protocol will require the miner look up and connect to a client to propose reconciliation. How does a miner look up and connect to a client over libp2p given only their filecoin address?
- Without repair miners, this protocol will likely not be used that much. Should that be made clear? Should there be other considerations added to compensate?

## Storage Miner Payments

TODO: these bits were pulled out of a different doc, and describe strategies by which client payments to a miner might happen. We need to organize 'clients paying miners' better, unclear if it should be the same doc that talks about payment channel constructions.

1. **Updates Contingent on Inclusion Proof**
   - In this case, the miner must provide an inclusion proof that shows the client data is contained in one of the miners sectors on chain, and submit that along with the payment channel update.
   - This can be pretty expensive for smaller files, and ideally, we make it to one of the latter two options
   - This option does however allow clients to upload their files and leave.
2. **Update Contingent on CommD Existence**
   - For this, the client needs to wait around until the miner finishes packing a sector, and computing its commD. The client then signs a set of payment channel updates that are contingent on the given commD existing on chain.
   - This route makes it difficult for miners to re-seal smaller files (really, small files just suck)
3. **Reconciled Payment**
   - In either of the above cases, the miner may go back to the client and say "Look, these payment channel updates you gave me are able to be cashed in right now, could you take them all and give me back a single update for a slightly smaller amount?".
   - The slightly smaller amount could be the difference in transaction fees, meaning the client saves money, and the miner gets the same amount.
