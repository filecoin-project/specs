package interpreter

import "errors"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import market "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market"
import spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
import sysactors "github.com/filecoin-project/specs/systems/filecoin_vm/sysactors"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"

var (
	ErrActorNotFound = errors.New("Actor Not Found")
)

// CodeIDs for system actors
var (
	InitActorCodeID           = actor.CodeID_Make_Builtin(actor.BuiltinActorID_Init)
	CronActorCodeID           = actor.CodeID_Make_Builtin(actor.BuiltinActorID_Cron)
	AccountActorCodeID        = actor.CodeID_Make_Builtin(actor.BuiltinActorID_Account)
	StoragePowerActorCodeID   = actor.CodeID_Make_Builtin(actor.BuiltinActorID_StoragePower)
	StorageMinerActorCodeID   = actor.CodeID_Make_Builtin(actor.BuiltinActorID_StorageMiner)
	StorageMarketActorCodeID  = actor.CodeID_Make_Builtin(actor.BuiltinActorID_StorageMarket)
	PaymentChannelActorCodeID = actor.CodeID_Make_Builtin(actor.BuiltinActorID_PaymentChannel)
)

var staticActorCodeRegistry = &actorCodeRegistry{}

type actorCodeRegistry struct {
	code map[actor.CodeID]vmr.ActorCode
}

func (r *actorCodeRegistry) _registerActor(id actor.CodeID, actor vmr.ActorCode) {
	r.code[id] = actor
}

func (r *actorCodeRegistry) _loadActor(id actor.CodeID) (vmr.ActorCode, error) {
	a, ok := r.code[id]
	if !ok {
		return nil, ErrActorNotFound
	}
	return a, nil
}

func RegisterActor(id actor.CodeID, actor vmr.ActorCode) {
	staticActorCodeRegistry._registerActor(id, actor)
}

func LoadActor(id actor.CodeID) (vmr.ActorCode, error) {
	return staticActorCodeRegistry._loadActor(id)
}

// init is called in Go during initialization of a program.
// this is an idiomatic way to do this. Implementations should approach this
// howevery they wish. The point is to initialize a static registry with
// built in pure types that have the code for each actor. Once we have
// a way to load code from the StateTree, use that instead.
func init() {
	_registerBuiltinActors()
}

func _registerBuiltinActors() {
	// TODO

	cron := &sysactors.CronActorCode_I{}

	RegisterActor(InitActorCodeID, &sysactors.InitActorCode_I{})
	RegisterActor(CronActorCodeID, cron)
	RegisterActor(AccountActorCodeID, &sysactors.AccountActorCode_I{})
	RegisterActor(StoragePowerActorCodeID, &spc.StoragePowerActorCode_I{})
	RegisterActor(StorageMarketActorCodeID, &market.StorageMarketActorCode_I{})

	// wire in CRON actions.
	// TODO: there's probably a better place to put this, but for now, do it here.
	cron.Actors_ = append(cron.Actors_, addr.StoragePowerActorAddr)
	cron.Actors_ = append(cron.Actors_, addr.StorageMarketActorAddr)
}
