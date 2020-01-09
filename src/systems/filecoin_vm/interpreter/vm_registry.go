package interpreter

import "errors"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import abi "github.com/filecoin-project/specs/actors/abi"
import builtin "github.com/filecoin-project/specs/actors/builtin"
import accact "github.com/filecoin-project/specs/actors/builtin/account"
import cronact "github.com/filecoin-project/specs/actors/builtin/cron"
import initact "github.com/filecoin-project/specs/actors/builtin/init"
import spowact "github.com/filecoin-project/specs/actors/builtin/storage_power"
import smarkact "github.com/filecoin-project/specs/actors/builtin/storage_market"

var (
	ErrActorNotFound = errors.New("Actor Not Found")
)

var staticActorCodeRegistry = &actorCodeRegistry{}

type actorCodeRegistry struct {
	code map[abi.ActorCodeID]vmr.ActorCode
}

func (r *actorCodeRegistry) _registerActor(id abi.ActorCodeID, actor vmr.ActorCode) {
	r.code[id] = actor
}

func (r *actorCodeRegistry) _loadActor(id abi.ActorCodeID) (vmr.ActorCode, error) {
	a, ok := r.code[id]
	if !ok {
		return nil, ErrActorNotFound
	}
	return a, nil
}

func RegisterActor(id abi.ActorCodeID, actor vmr.ActorCode) {
	staticActorCodeRegistry._registerActor(id, actor)
}

func LoadActor(id abi.ActorCodeID) (vmr.ActorCode, error) {
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

	RegisterActor(builtin.InitActorCodeID, &initact.InitActorCode_I{})
	RegisterActor(builtin.CronActorCodeID, cron)
	RegisterActor(builtin.AccountActorCodeID, &accact.AccountActorCode_I{})
	RegisterActor(builtin.StoragePowerActorCodeID, &spowact.StoragePowerActorCode_I{})
	RegisterActor(builtin.StorageMarketActorCodeID, &smarkact.StorageMarketActorCode_I{})

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
