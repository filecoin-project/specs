package interpreter

import "errors"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import sysactors "github.com/filecoin-project/specs/systems/filecoin_vm/sysactors"

var (
	ErrActorNotFound = errors.New("Actor Not Found")
)

var staticActorCodeRegistry = &actorCodeRegistry{}

// CodeCIDs for system actors
var (
	InitActorCodeCID           = actor.CodeCID("filecoin/1.0/InitActor")
	CronActorCodeCID           = actor.CodeCID("filecoin/1.0/CronActor")
	AccountActorCodeCID        = actor.CodeCID("filecoin/1.0/AccountActor")
	StoragePowerActorCodeCID   = actor.CodeCID("filecoin/1.0/StoragePowerActor")
	StorageMarketActorCodeCID  = actor.CodeCID("filecoin/1.0/StorageMarketActor")
	PaymentChannelActorCodeCID = actor.CodeCID("filecoin/1.0/PaymentChannelActor")
)

// Addresses for singleton system actors
var (
	InitActorAddr           = &addr.Address_I{} // TODO
	CronActorAddr           = &addr.Address_I{} // TODO
	StoragePowerActorAddr   = &addr.Address_I{} // TODO
	StorageMarketActorAddr  = &addr.Address_I{} // TODO
	PaymentChannelActorAddr = &addr.Address_I{} // TODO
)

// init is called in Go during initialization of a program.
// this is an idiomatic way to do this. Implementations should approach this
// howevery they wish. The point is to initialize a static registry with
// built in pure types that have the code for each actor. Once we have
// a way to load code from the StateTree, use that instead.
func init() {
	registerBuiltinActors(staticActorCodeRegistry)
}

func registerBuiltinActors(r *actorCodeRegistry) {
	// TODO

	cron := &sysactors.CronActorCode_I{}

	r.RegisterActor(InitActorCodeCID, &sysactors.InitActorCode_I{})
	r.RegisterActor(CronActorCodeCID, cron)

	// wire in CRON actions.
	// TODO: there's probably a better place to put this, but for now, do it here.
	cron.Actors_ = append(cron.Actors_, StoragePowerActorAddr)
	cron.Actors_ = append(cron.Actors_, StorageMarketActorAddr)
}

// ActorCode is the interface that all actor code types should satisfy.
// It is merely a method dispatch interface.
type ActorCode interface {
	InvokeMethod(input vmr.InvocInput, method actor.MethodNum, params actor.MethodParams) vmr.InvocOutput
}

type actorCodeRegistry struct {
	code map[actor.CodeCID]ActorCode
}

func (r *actorCodeRegistry) RegisterActor(cid actor.CodeCID, actor ActorCode) {
	r.code[cid] = actor
}

func (r *actorCodeRegistry) LoadActor(cid actor.CodeCID) (ActorCode, error) {
	a, ok := r.code[cid]
	if !ok {
		return nil, ErrActorNotFound
	}
	return a, nil
}
