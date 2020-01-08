package interpreter

import "errors"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import accact "github.com/filecoin-project/specs/actors/builtin/account"
import cronact "github.com/filecoin-project/specs/actors/builtin/cron"
import initact "github.com/filecoin-project/specs/actors/builtin/init"
import spowact "github.com/filecoin-project/specs/actors/builtin/storage_power"
import smarkact "github.com/filecoin-project/specs/actors/builtin/storage_market"

var (
	ErrActorNotFound = errors.New("Actor Not Found")
)

// CodeIDs for system actors
var (
	SystemActorCodeID         = actor.CodeID_Make_Builtin(actor.BuiltinActorID_System)
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

	cron := &cronact.CronActorCode_I{}

	RegisterActor(InitActorCodeID, &initact.InitActorCode_I{})
	RegisterActor(CronActorCodeID, cron)
	RegisterActor(AccountActorCodeID, &accact.AccountActorCode_I{})
	RegisterActor(StoragePowerActorCodeID, &spowact.StoragePowerActorCode_I{})
	RegisterActor(StorageMarketActorCodeID, &smarkact.StorageMarketActorCode_I{})

	// wire in CRON actions.
	// TODO: move this to CronActor's constructor method
	cron.Entries_ = append(cron.Entries_, &cronact.CronTableEntry_I{
		ToAddr_:    addr.StoragePowerActorAddr,
		MethodNum_: ai.Method_StoragePowerActor_OnEpochTickEnd,
	})

	cron.Entries_ = append(cron.Entries_, &cronact.CronTableEntry_I{
		ToAddr_:    addr.StorageMarketActorAddr,
		MethodNum_: ai.Method_StorageMarketActor_OnEpochTickEnd,
	})
}
