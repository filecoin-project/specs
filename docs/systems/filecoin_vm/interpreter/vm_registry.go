package interpreter

import (
	"errors"

	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	accact "github.com/filecoin-project/specs/actors/builtin/account"
	cronact "github.com/filecoin-project/specs/actors/builtin/cron"
	initact "github.com/filecoin-project/specs/actors/builtin/init"
	smarkact "github.com/filecoin-project/specs/actors/builtin/storage_market"
	spowact "github.com/filecoin-project/specs/actors/builtin/storage_power"
	vmr "github.com/filecoin-project/specs/actors/runtime"
)

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
		ToAddr_:    builtin.StoragePowerActorAddr,
		MethodNum_: builtin.Method_StoragePowerActor_OnEpochTickEnd,
	})

	cron.Entries_ = append(cron.Entries_, &cronact.CronTableEntry_I{
		ToAddr_:    builtin.StorageMarketActorAddr,
		MethodNum_: builtin.Method_StorageMarketActor_OnEpochTickEnd,
	})
}
