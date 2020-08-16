---
title: Gas Costs
dashboardWeight: 1
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# VM Gas Cost Constants
---

Every operation that triggers computation or storage on the Filecoin VM incurs payment in terms of `GasCharge` (`VirtualCompute`, `VirtualStorage`). Charges can also be incurred by actors (`ComputeGas`, `StorageGas`). The `GasCharge struct` includes the `Name` of the cost to be incurred, which is taken by the `Pricelist`.

```go
type GasCharge struct {
	Name  string
	Extra interface{}

	ComputeGas int64
	StorageGas int64

	VirtualCompute int64
	VirtualStorage int64
}
```

The total `GasCharge` is computed by the `Total()` function, which also takes into account a multiplier in case either compute or storage are considered more expensive. In the current implementation they both have equal weight.


```go
const (
	GasStorageMulti = 1
	GasComputeMulti = 1
)

func (g GasCharge) Total() int64 {
	return g.ComputeGas*GasComputeMulti + g.StorageGas*GasStorageMulti
}
```

Any `newGasCharge` triggers a function to gather costs incurred by actors, while those costs incurred by the VM are captured by a separate `WithVirtual` function. Both of those functions are called by the VM or the runtime interface to make sure all costs are captured.


The `Pricelist interface` provides prices for operations in the VM. This interface should be APPEND ONLY since last chain checkpoint. It includes the following charges.

- `OnChainMessage(msgSize int) GasCharge`: returns the gas charged to the originator of a message for storing a message of a given size on chain.
- `OnChainReturnValue(dataSize int) GasCharge`: returns the gas used for storing the response of a message in the chain.
- `OnMethodInvocation(value abi.TokenAmount, methodNum abi.MethodNum) GasCharge`: returns the gas used when invoking a method.
- `OnIpldGet(dataSize int) GasCharge`: returns the gas used for getting an object
- `OnIpldPut(dataSize int) GasCharge`: returns the gas used for storing an object. The put operation incurs higher charges than the get operation as the put operation includes storage of on-chain data.
- `OnCreateActor() GasCharge`: returns the gas used for creating an actor
- `OnDeleteActor() GasCharge`: returns the gas used for deleting an actor
- `OnVerifySignature(sigType crypto.SigType, planTextSize int) (GasCharge, error)`: returns the gas cost for verifying a signature is valid for an actor address
- `OnHashing(dataSize int) GasCharge`: returns the gas cost of hashing input data using `HashBlake2b`
- `OnComputeUnsealedSectorCid(proofType abi.RegisteredSealProof, pieces []abi.PieceInfo) GasCharge`: returns the gas cost for computing the unsealed sector CID from the piece CIDs.
- `OnVerifySeal(info abi.SealVerifyInfo) GasCharge`: returns the gas cost for verifying a sealed sector
- `OnVerifyPost(info abi.WindowPoStVerifyInfo) GasCharge`: returns the gas cost for verifying a proof of SpaceTime submission
- `OnVerifyConsensusFault() GasCharge`: returns the gas cost for verifying a consensus fault. The corresponding function should verify that two block headers provide proof of a consensus fault for the following cases:
	- both headers mined by the same actor
	- headers are different
	- first header is of the same or lower epoch as the second
	- at least one of the headers appears in the current chain at or after epoch `earliest`
	- the headers provide evidence of a fault (see the spec for the different fault types).
	The corresponding function is defined as: `func (ps pricedSyscalls) VerifyConsensusFault(h1 []byte, h2 []byte, extra []byte) (*runtime.ConsensusFault, error) {}`. The parameters are all serialized block headers. The third "extra" parameter is 	consulted only for the "parent grinding fault", in which case it must be the sibling of h1 (same parent tipset) and one of the blocks in the parent of h2 (i.e. h2's grandparent).
	The function returns nil and an error if the headers don't prove a fault.



```go
type Pricelist interface {

	OnChainMessage(msgSize int) GasCharge
	OnChainReturnValue(dataSize int) GasCharge
	OnMethodInvocation(value abi.TokenAmount, methodNum abi.MethodNum) GasCharge
	OnIpldGet(dataSize int) GasCharge
	OnIpldPut(dataSize int) GasCharge
	OnCreateActor() GasCharge
	OnDeleteActor() GasCharge

	OnVerifySignature(sigType crypto.SigType, planTextSize int) (GasCharge, error)
	OnHashing(dataSize int) GasCharge
	OnComputeUnsealedSectorCid(proofType abi.RegisteredSealProof, pieces []abi.PieceInfo) GasCharge
	OnVerifySeal(info abi.SealVerifyInfo) GasCharge
	OnVerifyPost(info abi.WindowPoStVerifyInfo) GasCharge
	OnVerifyConsensusFault() GasCharge
}
```
