
# Multisig Account Actor (DEPRECATED)

- **Code Cid**: `<codec:raw><mhType:identity><"multisig">`

A basic multisig account actor. Allows sending of messages like a normal account actor, but with the requirement of M of N parties agreeing to the operation. Completed and/or cancelled operations stick around in the actors state until explicitly cleared out. Proposers may cancel transactions they propose, or transactions by proposers who are no longer approved signers.

Self modification methods (add/remove signer, change requirement) are called by
doing a multisig transaction invoking the desired method on the contract itself. This means the 'signature
threshold' logic only needs to be implemented once, in one place.

The [init actor](#init-actor) is used to create new instances of the multisig.

```sh
type MultisigActorState struct {
    signers [Address]
    required UInt
    nextTxId UInt
    initialBalance UInt
    startingBlock UInt
    unlockDuration UInt
    transactions {UInt:Transaction}
}

type Transaction struct {
    txID UInt
    to Address
    value TokenAmount
    method &ActorMethod
    approved [Address]
    completed Bool
    canceled Bool
    retcode UInt
}
```

## Methods

| Name | Method ID |
|--------|-------------|
| `MultisigConstructor` | 1 |
| `Propose` | 2 |
| `Approve` | 3 |
| `Cancel` | 4 |
| `ClearCompleted` | 5 |
| `AddSigner` | 6 |
| `RemoveSigner` | 7 |
| `SwapSigner` | 8 |
| `ChangeRequirement` | 9 |


## `Constructor`

This method sets up the initial state for the multisig account

**Parameters**

```sh
type MultisigConstructor struct {
    ## The addresses that will be the signatories of this wallet.
    signers [Address]
    ## The number of signatories required to perform a transaction.
    required UInt
    ## Unlock time (in blocks) of initial filecoin balance of this wallet. Unlocking is linear.
    unlockDuration UInt
} representation tuple
```

**Algorithm**

```go
func Multisig(signers []Address, required UInt, unlockDuration UInt) {
  self.Signers = signers
  self.Required = required
  self.initialBalance = msg.Value
  self.unlockDuration = unlockDuration
  self.startingBlock = VM.CurrentBlockHeight()
}
```

## `Propose`

Propose is used to propose a new transaction to be sent by this multisig. The proposer must be a signer, and the proposal also serves as implicit approval from the proposer. If only a single signature is required, then the transaction is executed immediately.

**Parameters**


```sh
type Propose struct {
    ## The address of the target of the proposed transaction.
    to Address
    ## The amount of funds to send with the proposed transaction.
    value TokenAmount
    ## The method and parameters that will be invoked on the proposed transactions target.
    method &ActorMethod
} representation tuple
```

**Algorithm**

```go
func Propose(to Address, value TokenAmount, method String, params Bytes) UInt {
  if !isSigner(msg.From) {
    Fatal("not authorized")
  }

  txid := self.NextTxID
  self.NextTxID++

  tx := Transaction{
    TxID:     txid,
    To:       to,
    Value:    value,
    Method:   method,
    Params:   params,
    Approved: []Address{msg.From},
  }

  self.Transactions.Append(tx)

  if self.Required == 1 {
    if !self.canSpend(tx.value) {
      Fatal("transaction amount exceeds available")
    }
    tx.RetCode = vm.Send(tx.To, tx.Value, tx.Method, tx.Params)
    tx.Complete = true
  }

  return txid
}
```

## `Approve`

Approve is called by a signer to approve a given transaction. If their approval pushes the approvals for this transaction over the threshold, the transaction is executed.

**Parameters**

```sh
type Approve struct {
    ## The ID of the transaction to approve.
    txid UInt
} representation tuple
```

**Algorithm**

```go
func Approve(txid UInt) {
  if !self.isSigner(msg.From) {
    Fatal("not authorized")
  }

  tx := self.getTransaction(txid)
  if tx.Complete {
    Fatal("transaction already completed")
  }
  if tx.Canceled {
    Fatal("transaction canceled")
  }

  for _, signer := range tx.Approved {
    if signer == msg.From {
      Fatal("already signed this message")
    }
  }

  tx.Approved.Append(msg.From)

  if len(tx.Approved) >= self.Required {
    if !self.canSpend(tx.Value) {
      Fatal("transaction amount exceeds available")
    }
    tx.RetCode = vm.Send(tx.To, tx.Value, tx.Method, tx.Params)
    tx.Complete = true
  }
}
```

## `Cancel`

**Parameters**

```sh
type Cancel struct {
    txid UInt
} representation tuple
```

**Algorithm**

```go
func Cancel(txid UInt) {
  if !self.isSigner(msg.From) {
    Fatal("not authorized")
  }

  tx := self.getTransaction(txid)
  if tx.Complete {
    Fatal("cannot cancel completed transaction")
  }
  if tx.Canceled {
    Fatal("transaction already canceled")
  }

  proposer := tx.Approved[0]
  if proposer != msg.From && isSigner(proposer) {
    Fatal("cannot cancel another signers transaction")
  }

  tx.Canceled = true
}
```

## `ClearCompleted`

**Parameters**

```sh
type ClearCompleted struct {
} representation tuple
```

**Algorithm**

```go
func ClearCompleted() {
  if !self.isSigner(msg.From) {
    Fatal("not authorized")
  }

  for tx := range self.Transactions {
    if tx.Completed || tx.Canceled {
      self.Transactions.Remove(tx)
    }
  }
}
```

## `AddSigner`

**Parameters**

```sh
type AddSigner struct {
    signer Address
    increaseReq bool
} representation tuple
```

**Algorithm**

```go
func AddSigner(signer Address, increaseReq bool) {
  if msg.From != self.Address {
    Fatal("add signer must be called by wallet itself")
  }
  if self.isSigner(signer) {
    Fatal("new address is already a signer")
  }
  if increaseReq {
    self.Required = self.Required + 1
  }

  self.Signers.Append(signer)
}
```

## `RemoveSigner`

**Parameters**

```sh
type RemoveSigner struct {
    signer Address
    decreaseReq bool
} representation tuple
```

**Algorithm**

```go
func RemoveSigner(signer Address, decreaseReq bool) {
  if msg.From != self.Address {
    Fatal("remove signer must be called by wallet itself")
  }
  if !self.isSigner(signer) {
    Fatal("given address was not a signer")
  }
  if decreaseReq || len(self.Signers)-1 < self.Required {
    // Reduce Required outherwise the wallet is locked out
    self.Required = self.Required - 1
  }

  self.Signers.Remove(signer)
}
```

## `SwapSigner`

**Parameters**

```sh
type SwapSigner struct {
    old Address
    new Address
} representation tuple
```

**Algorithm**

```go
func SwapSigner(old Address, new Address) {
  if msg.From != self.Address {
    Fatal("swap signer must be called by wallet itself")
  }
  if !self.isSigner(old) {
    Fatal("given old address was not a signer")
  }
  if self.isSigner(new) {
    Fatal("given new address was already a signer")
  }

  self.Signers.Remove(old)
  self.Signers.Append(new)
}
```

## `ChangeRequirement`

**Parameters**

```sh
type ChangeRequirement struct {
    requirement UInt
} representation tuple
```

**Algorithm**

```go
func ChangeRequirement(req UInt) {
  if msg.From != self.Address {
    Fatal("change requirement must be called by wallet itself")
  }
  if req < 1 {
    Fatal("requirement must be at least 1")
  }
  if req > len(self.Signers) {
    Fatal("requirement must be less than number of signers")
  }

  self.Required = req
}
```
