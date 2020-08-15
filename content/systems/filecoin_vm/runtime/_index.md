---
title: Runtime
weight: 5
bookCollapseSection: true
dashboardWeight: 1
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# VM Runtime Environment (Inside the VM)
---

## Receipts

A `MessageReceipt` contains the result of a top-level message execution. Every syntactically valid and correctly signed message can be included in a block and will produce a receipt from execution. 

A syntactically valid receipt has:

- a non-negative `ExitCode`,
- a non empty `Return` value only if the exit code is zero, and
- a non-negative `GasUsed`.

```go
type MessageReceipt struct {
	ExitCode exitcode.ExitCode
	Return   []byte
	GasUsed  int64
}
```

## `vm/runtime` interface

The `runtime` is the VM's internal runtime object. What is included in this interface is everything that is accessible to the actors.

The `Runtime interface` includes the following:

- `Message()`: Information related to the current message being executed. When an actor invokes a method on another actor as a sub-call, these values reflect the sub-call context, rather than the top-level context.
- `CurrEpoch()`: The current chain epoch number. Note that the genesis block has epoch zero.
- `ValidateImmediateCallerAcceptAny()`:	Satisfies the requirement that every exported actor method must invoke at least one caller validation method before returning, without making any assertions about the caller.
- `	ValidateImmediateCallerIs(addrs ...addr.Address)`: Validates that the immediate caller's address matches exactly one of a set of expected addresses and aborts if it does not. The caller address is always normalized to an ID address, so expected addresses must be ID addresses to have any expectation of passing validation.
- `	ValidateImmediateCallerType(types ...cid.Cid)` Validates that the immediate caller is an actor with code CID matching one of a set of expected CIDs, aborting if it does not.
- `	CurrentBalance() abi.TokenAmount`: Returns the balance of the receiver. The value should always be greater or equal to zero.
- `ResolveAddress(address addr.Address) (addr.Address, bool)`: Resolves the address of any protocol to an ID address (via the Init actor's table). This allows resolution of externally-provided SECP, BLS, or actor addresses to the canonical form. If the argument is an ID address it is returned directly.
- `	GetActorCodeCID(addr addr.Address) (ret cid.Cid, ok bool)`: Look up the code ID of an actor address. The address will be resolved via ResolveAddress, if necessary, so it does not need to be an ID-address.
- `	GetRandomnessFromBeacon(personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) abi.Randomness`: Returns a (pseudo)random byte array drawing from a random beacon at a prior epoch. The beacon value is combined with the personalization tag, epoch number, and explicitly provided entropy. The personalization tag may be any int64 value. The epoch must be less than the current epoch. The epoch may be negative, in which case it addresses the beacon value from genesis block. The entropy may be any byte array, or nil.
- `	GetRandomnessFromTickets(personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) abi.Randomness`: Samples randomness from the ticket chain. Randomess sampled through this method is unique per potential fork, and as a result, processes relying on this randomness are tied to whichever fork they choose. See GetRandomnessFromBeacon (above) for notes about the personalization tag, epoch, and entropy.
- `	State() StateHandle`: Provides a handle for the actor's state object.
- `	Send(toAddr addr.Address, methodNum abi.MethodNum, params CBORMarshaler, value abi.TokenAmount) (SendReturn, exitcode.ExitCode)`: Sends a message to another actor, returning the exit code and return value envelope. If the invoked method does not return successfully, its state changes (and that of any messages it sent in turn) will be rolled back. The result is never a bare nil, but may be (a wrapper of) `adt.Empty`.
- `	Abortf(errExitCode exitcode.ExitCode, msg string, args ...interface{})`: Halts execution upon an error from which the receiver cannot recover. The caller will receive the exitcode and an empty return value. State changes made within this call will be rolled back. This method does not return. The provided exit code must be `>= exitcode.FirstActorExitCode`. The message and args are for diagnostic purposes and do not persist on chain. They should be suitable for passing to `fmt.Errorf(msg, args...)`.
- `	NewActorAddress() addr.Address`: Computes an address for a new actor. The returned address is intended to uniquely refer to the actor even in the event of a chain re-org (whereas an ID-address might refer to a different actor after messages are re-ordered).
- `	CreateActor(codeId cid.Cid, address addr.Address)`: Creates an actor with code `codeID` and address `address`, with empty state. Can only be called by `InitActor`. It aborts if the provided address has previously been created/already exists.
- `	DeleteActor(beneficiary addr.Address)`: Deletes the executing actor from the state tree, transferring any balance to beneficiary. It aborts if the beneficiary does not exist or is the calling actor. Can only be called by the actor itself.
- `	Syscalls() Syscalls`: Provides the system call interface.
- `	TotalFilCircSupply() abi.TokenAmount`: Returns the total token supply in circulation at the beginning of the current epoch. The circulating supply is the sum of:
	- rewards emitted by the reward actor,
	- funds vested from lock-ups in the genesis state,
	 less the sum of:
	- funds burnt,
	- pledge collateral locked in storage miner actors (recorded in the storage power actor),
	- deal collateral locked by the storage market actor.
- `	Context() context.Context`: Provides a Go context for use by HAMT, etc. The VM is intended to provide an idealised machine abstraction, with infinite storage, so this context must not be used by actor code directly.
- `	StartSpan(name string) TraceSpan`: Starts a new tracing span. The span must be `End()`'ed explicitly, typically with a deferred invocation.
- `	ChargeGas(name string, gas int64, virtual int64)`: Charges specified amount of `gas` for execution.
	- `name` provides information about gas charging point.
	- `virtual` sets virtual amount of gas to charge, this amount is not counted toward execution cost. This functionality is used for observing global changes in total gas charged if amount of gas charged was to be changed.
- `	Log(level LogLevel, msg string, args ...interface{})`: Note events that may make debugging easier.


```go
type Runtime interface {

	Message() Message
	CurrEpoch() abi.ChainEpoch
	ValidateImmediateCallerAcceptAny()
	ValidateImmediateCallerIs(addrs ...addr.Address)
	ValidateImmediateCallerType(types ...cid.Cid)
	CurrentBalance() abi.TokenAmount
	ResolveAddress(address addr.Address) (addr.Address, bool)
	GetActorCodeCID(addr addr.Address) (ret cid.Cid, ok bool)
	GetRandomnessFromBeacon(personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) abi.Randomness
	GetRandomnessFromTickets(personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) abi.Randomness
	State() StateHandle
	Store() Store
	Send(toAddr addr.Address, methodNum abi.MethodNum, params CBORMarshaler, value abi.TokenAmount) (SendReturn, exitcode.ExitCode)
	Abortf(errExitCode exitcode.ExitCode, msg string, args ...interface{})
	NewActorAddress() addr.Address
	CreateActor(codeId cid.Cid, address addr.Address)
	DeleteActor(beneficiary addr.Address)
	Syscalls() Syscalls
	TotalFilCircSupply() abi.TokenAmount
	Context() context.Context
	StartSpan(name string) TraceSpan
	ChargeGas(name string, gas int64, virtual int64)
	Log(level LogLevel, msg string, args ...interface{})
}
```

The `Store interface` defines the storage module exposed to actors. Retrieves and deserializes an object from the store into `o`. Returns whether successful. Serializes and stores an object, returning its CID.

```go
type Store interface {
	Get(c cid.Cid, o CBORUnmarshaler) bool
	Put(x CBORMarshaler) cid.Cid
}
```

The `Message interface` contains information available to the actor about the executing message. These values are fixed for the duration of an invocation.

It includes:
- the address of the immediate calling actor `Caller()`, which is always an ID-address. If an actor invokes its own method, then `Caller() == Receiver()`.
- The address of the actor receiving the message `Receiver()`, which is always an ID-address.
- The value attached to the message being processed, implicitly added to `CurrentBalance()` of the `Receiver()` before method invocation. This value comes from `Caller()`.

```go
type Message interface {

	Caller() addr.Address
	Receiver() addr.Address
	ValueReceived() abi.TokenAmount
}
```

The `Syscalls interfacce` includes pure functions implemented as primitives by the runtime. These functions include:

- `	VerifySignature(signature crypto.Signature, signer addr.Address, plaintext []byte) error`: Verifies that a signature is valid for an address and plaintext. If the address is a public-key type address, it is used directly. If it's an ID-address, the actor is looked up in state. The ID-address must belong to an account actor, and the public key is obtained from it's state.
- `	HashBlake2b(data []byte) [32]byte`: Hashes input data using blake2b with 256 bit output.
- `	ComputeUnsealedSectorCID(reg abi.RegisteredSealProof, pieces []abi.PieceInfo) (cid.Cid, error)`: Computes an unsealed sector CID (CommD) from its constituent piece CIDs (CommPs) and sizes.
- `VerifySeal(vi abi.SealVerifyInfo) error`: Verifies a sector seal proof.
- `VerifyPoSt(vi abi.WindowPoStVerifyInfo) error`: Verifies a proof of spacetime.
- `	VerifyConsensusFault(h1, h2, extra []byte) (*ConsensusFault, error)`: Verifies that two block headers provide proof of a consensus fault:
	- both headers mined by the same actor
	- headers are different
	- first header is of the same or lower epoch as the second
	- the headers provide evidence of a fault (see the spec for the different fault types).
	The parameters are all serialized block headers. The third "extra" parameter is consulted only for the "parent grinding fault", in which case it must be the sibling of h1 (same parent tipset) and one of the blocks in an ancestor of h2. Returns nil and an error if the headers don't prove a fault.



```go
type Syscalls interface {

	VerifySignature(signature crypto.Signature, signer addr.Address, plaintext []byte) error
	HashBlake2b(data []byte) [32]byte
	ComputeUnsealedSectorCID(reg abi.RegisteredSealProof, pieces []abi.PieceInfo) (cid.Cid, error)
	VerifySeal(vi abi.SealVerifyInfo) error

	BatchVerifySeals(vis map[address.Address][]abi.SealVerifyInfo) (map[address.Address][]bool, error)

	VerifyPoSt(vi abi.WindowPoStVerifyInfo) error
	VerifyConsensusFault(h1, h2, extra []byte) (*ConsensusFault, error)
}
```

The `SendReturn interface` is the return type from a message sent from one actor to another. This abstracts over the internal representation of the return, in particular whether it has been serialized to bytes or just passed through.

```go
type SendReturn interface {
	Into(CBORUnmarshaler) error
}
```

The `TraceSpan interface` provides (minimal) tracing facilities to actor code.

```go
type TraceSpan interface {
	// Ends the span
	End()
}
```

The `StateHandle interface` provides mutable, exclusive access to actor state. 	It creates initializes the state object. This is only valid in a constructor function and when the state has not yet been initialized. The `Readonly` option loads a readonly copy of the state into the argument. Any modification to the state is illegal and will result in an abort.

`Transaction` loads a mutable version of the state into the `obj` argument and protects the execution from side effects (including message send). The second argument is a function which allows the caller to mutate the state. If the state is modified after this function returns, execution will abort. The gas cost of this method is that of a `Store.Put` of the mutated state object. Note: the Go signature is not ideal due to lack of type system power.


```go

type StateHandle interface {

	Create(obj CBORMarshaler)
	Readonly(obj CBORUnmarshaler)

	Transaction(obj CBORer, f func())
}
```

The VM Runtime Interface includes also the `ConsensusFault struct` which returns the result of checking two headers for a consensus fault. It includes the address of the miner at fault (always an ID address) - `Target`. It also includes the `Epoch` of the fault, which is the higher epoch of the two blocks causing it and the `Type` of fault.

```go
type ConsensusFault struct {

	Target addr.Address
	Epoch abi.ChainEpoch
	Type ConsensusFaultType
}
```

## Exit Codes

### Common Runtime Exit Codes

There are some common error codes that are shared by different actors. Apart from those, actors can also define their own codes. These include the following:

- `	ErrIllegalArgument = FirstActorErrorCode + iota`: Indicates a method parameter is invalid.
- `ErrNotFound`: Indicates a requested resource does not exist.
- `ErrForbidden`: Indicates an action is disallowed.
- `ErrInsufficientFunds`: Indicates a balance of funds is insufficient.
- `ErrIllegalState`: Indicates an actor's internal state is invalid.
- `ErrSerialization`: Indicates de/serialization failure within actor code.

```go

const (
	ErrIllegalArgument = FirstActorErrorCode + iota
	ErrNotFound
	ErrForbidden
	ErrInsufficientFunds
	ErrIllegalState
	ErrSerialization

	// Common error codes stop here.  If you define a common error code above
	// this value it will have conflicting interpretations
	FirstActorSpecificExitCode = ExitCode(32)
)
```

### Reserved System Exit Codes

Apart from the actor-specific common error codes above, there are some system error codes, which are reserved for use by the runtime. These system error codes must not be used by actors explicitly. Correspondingly, no runtime invocation should abort with an exit code outside this list.

- `SysErrSenderInvalid = ExitCode(1)`: Indicates that the actor identified as the sender of a message is not valid as a message sender. Possible causes are:
	- not present in the state tree
	- not an account actor (for top-level messages)
	- code CID is not found or invalid
- `	SysErrSenderStateInvalid = ExitCode(2)`: Indicates that the sender of a message is not in a state to send the message. Possible causes are:
	- invocation out of sequence (mismatched CallSeqNum)
	- insufficient funds to cover execution
- `SysErrInvalidMethod = ExitCode(3)`: Indicates failure to find a method in an actor.
- `SysErrInvalidParameters = ExitCode(4)`: Indicates non-decodeable or syntactically invalid parameters for a method.
- `SysErrInvalidReceiver = ExitCode(5)`: Indicates that the receiver of a message is not valid (and cannot be implicitly created).
- `SysErrInsufficientFunds = ExitCode(6)`: Indicates that a message sender has insufficient balance for the value being sent. Note that this is distinct from `SysErrSenderStateInvalid` when a top-level sender can't cover value transfer + gas. This code is only expected to come from inter-actor sends.
- `SysErrOutOfGas = ExitCode(7)`: Indicates message execution (including subcalls) used more gas than the specified limit.
- `SysErrForbidden = ExitCode(8)`: Indicates message execution is forbidden for the caller by runtime caller validation.
- `SysErrorIllegalActor = ExitCode(9)`: Indicates actor code performed a disallowed operation. Disallowed operations include:
	- mutating state outside of a state acquisition block
	- failing to invoke caller validation
	- aborting with a reserved exit code (including success or a system error).
- `SysErrorIllegalArgument = ExitCode(10)`: Indicates an invalid argument passed to a runtime method.
- `SysErrSerialization = ExitCode(11)`: Indicates  an object failed to de/serialize for storage.


```go

	SysErrSenderInvalid = ExitCode(1)
	SysErrSenderStateInvalid = ExitCode(2)
	SysErrInvalidMethod = ExitCode(3)
	SysErrInvalidParameters = ExitCode(4)
	SysErrInvalidReceiver = ExitCode(5)
	SysErrInsufficientFunds = ExitCode(6)
	SysErrOutOfGas = ExitCode(7)
	SysErrForbidden = ExitCode(8)
	SysErrorIllegalActor = ExitCode(9)
	SysErrorIllegalArgument = ExitCode(10)
	SysErrSerialization = ExitCode(11)


	SysErrorReserved3 = ExitCode(12)
	SysErrorReserved4 = ExitCode(13)
	SysErrorReserved5 = ExitCode(14)
	SysErrorReserved6 = ExitCode(15)
)

// The initial range of exit codes is reserved for system errors.
// Actors may define codes starting with this one.
const FirstActorErrorCode = ExitCode(16)

```




