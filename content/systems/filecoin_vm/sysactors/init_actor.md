---
title: InitActor
weight: 1
dashboardWeight: 2
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# InitActor
---

The `InitActor` has the power to create new actors, e.g., those that enter the system. It maintains a table resolving a public key and temporary actor addresses to their canonical ID-addresses.

The actor should be able to construct a canonical ID address for the actor that will persist even in case of a chain re-organisation.

```go
type ExecParams struct {
	CodeCID           cid.Cid `checked:"true"` // invalid CIDs won't get committed to the state tree
	ConstructorParams []byte
}

type ExecReturn struct {
	IDAddress     addr.Address // The canonical ID-based address for the actor.
	RobustAddress addr.Address // A more expensive but re-org-safe address for the newly created actor.
}
```

```go
func (a Actor) Exec(rt runtime.Runtime, params *ExecParams) *ExecReturn {
	rt.ValidateImmediateCallerAcceptAny()
	callerCodeCID, ok := rt.GetActorCodeCID(rt.Message().Caller())
	autil.AssertMsg(ok, "no code for actor at %s", rt.Message().Caller())
	if !canExec(callerCodeCID, params.CodeCID) {
		rt.Abortf(exitcode.ErrForbidden, "caller type %v cannot exec actor type %v", callerCodeCID, params.CodeCID)
	}

	// Compute a re-org-stable address.
	// This address exists for use by messages coming from outside the system, in order to
	// stably address the newly created actor even if a chain re-org causes it to end up with
	// a different ID.
	uniqueAddress := rt.NewActorAddress()

	// Allocate an ID for this actor.
	// Store mapping of pubkey or actor address to actor ID
	var st State
	var idAddr addr.Address
	rt.State().Transaction(&st, func() {
		var err error
		idAddr, err = st.MapAddressToNewID(adt.AsStore(rt), uniqueAddress)
		builtin.RequireNoErr(rt, err, exitcode.ErrIllegalState, "failed to allocate ID address")
	})

	// Create an empty actor.
	rt.CreateActor(params.CodeCID, idAddr)

	// Invoke constructor.
	_, code := rt.Send(idAddr, builtin.MethodConstructor, runtime.CBORBytes(params.ConstructorParams), rt.Message().ValueReceived())
	builtin.RequireSuccess(rt, code, "constructor failed")

	return &ExecReturn{idAddr, uniqueAddress}
}
```
