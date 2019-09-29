
# Payment Channel Actor (DEPRECATED)

- **Code Cid:** `<codec:raw><mhType:identity><"paych">`

The payment channel actor manages the on-chain state of a point to point payment channel.

```sh
type PaymentChannel struct {
  from Address
  to   Address

  toSend       TokenAmount

  closingAt      UInt
  minCloseHeight UInt

  laneStates {UInt:LaneState}
} representation tuple

type SignedVoucher struct {
  TimeLock BlockHeight
  SecretPreimage Bytes
  Extra ModVerifyParams
  Lane Uint
  Nonce Uint
  Merges []Merge
  Amount TokenAmount
  MinCloseHeight Uint

  Signature Signature
}

type ModVerifyParams struct {
  Actor Address
  Method Uint
  Data Bytes
}

type Merge struct {
  Lane Uint
  Nonce Uint
}

type LaneState struct {
  Closed bool
  Redeemed TokenAmount
  Nonce Uint
}

type PaymentChannelMethod union {
  | PaymentChannelConstructor 0
  | UpdateChannelState 1
  | Close 2
  | Collect 3
} representation keyed
```

## Methods

| Name | Method ID |
|--------|-------------|
| `Constructor` | 1 |
| `UpdateChannelState` | 2 |
| `Close` | 3 |
| `Collect` | 4 |

## `Constructor`

**Parameters**

```sh
type PaymentChannelConstructor struct {
  to Address
}
```

**Algorithm**

{{% notice todo %}}

TODO: Define me

{{% /notice %}}

## `UpdateChannelState`

**Parameters**

```sh
type UpdateChannelState struct {
  sv SignedVoucher
  secret Bytes
  proof Bytes
} representation tuple
```

**Algorithm**

```go
func UpdateChannelState(sv SignedVoucher, secret []byte, proof []byte) {
  if !self.validateSignature(sv) {
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

  if sv.Extra != nil {
    ret := vmctx.Send(sv.Extra.Actor, sv.Extra.Method, sv.Extra.Data, proof)
    if ret != 0 {
      Fatal("spend voucher verification failed")
    }
  }

  ls := self.LaneStates[sv.Lane]
  if ls.Closed {
    Fatal("cannot redeem a voucher on a closed lane")
  }

  if ls.Nonce > sv.Nonce {
    Fatal("voucher has an outdated nonce, cannot redeem")
  }

  var mergeValue TokenAmount
  for _, merge := range sv.Merges {
    if merge.Lane == sv.Lane {
      Fatal("voucher cannot merge its own lane")
    }

    ols := self.LaneStates[merge.Lane]
    if ols.Nonce >= merge.Nonce {
      Fatal("merge in voucher has outdated nonce, cannot redeem")
    }

    mergeValue += ols.Redeemed
    ols.Nonce = merge.Nonce
  }

  ls.Nonce = sv.Nonce
  balanceDelta = sv.Amount - (mergeValue + ls.Redeemed)
  ls.Redeemed = sv.Amount

  newSendBalance = self.ToSend + balanceDelta
  if newSendBalance < 0 {
    // TODO: is this impossible?
    Fatal("voucher would leave channel balance negative")
  }

  if newSendBalance > self.Balance {
    Fatal("not enough funds in channel to cover voucher")
  }

  self.ToSend = newSendBalance

  if sv.MinCloseHeight != 0 {
    if self.ClosingAt != 0 && self.ClosingAt < sv.MinCloseHeight {
      self.ClosingAt = sv.MinCloseHeight
    }
    if self.MinCloseHeight < sv.MinCloseHeight {
      self.MinCloseHeight = sv.MinCloseHeight
    }
  }
}

func Hash(b []byte) []byte {
  return blake2b.Sum(b)
}
```

## `Close`

**Parameters**

```sh
type Close struct {
} representation tuple
```

**Algorithm**

```go
const ChannelClosingDelay = 6 * 60 * 2 // six hours

func Close() {
  if msg.From != self.From && msg.From != self.To {
    Fatal("not authorized to close channel")
  }
  if self.ClosingAt != 0 {
    Fatal("Channel already closing")
  }

  self.ClosingAt = chain.Now() + ChannelClosingDelay
  if self.ClosingAt < self.MinCloseHeight {
    self.ClosingAt = self.MinCloseHeight
  }
}
```

## `Collect`

**Parameters**

```sh
type Collect struct {
} representation tuple
```

**Algorithm**

```go
func Collect() {
  if self.ClosingAt == 0 {
    Fatal("payment channel not closing or closed")
  }

  if chain.Now() < self.ClosingAt {
    Fatal("Payment channel not yet closed")
  }

  TransferFunds(self.From, self.Balance-self.ToSend)
  TransferFunds(self.To, self.ToSend)
  self.ToSend = 0
}
```