---
title: AccountActor
weight: 3
dashboardWeight: 2
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# AccountActor
---

The `AccountActor` is responsible for user accounts. Account actors are not created by the `InitActor`, but they constructor is called by the system. Account actors are created by sending a message to a public-key style address. The address must be `BLS` or `SECP`.

```go
func (a Actor) Constructor(rt vmr.Runtime, address *addr.Address) *adt.EmptyValue {

	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	switch address.Protocol() {
	case addr.SECP256K1:
	case addr.BLS:
		break // ok
	default:
		rt.Abortf(exitcode.ErrIllegalArgument, "address must use BLS or SECP protocol, got %v", address.Protocol())
	}
	st := State{Address: *address}
	rt.State().Create(&st)
	return nil
}
```
