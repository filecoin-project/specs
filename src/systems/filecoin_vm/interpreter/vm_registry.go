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

// CodeCIDs for system actors
var (
	InitActorCodeCID           = actor.CodeCID("filecoin/1.0/InitActor")
	CronActorCodeCID           = actor.CodeCID("filecoin/1.0/CronActor")
	AccountActorCodeCID        = actor.CodeCID("filecoin/1.0/AccountActor")
	StoragePowerActorCodeCID   = actor.CodeCID("filecoin/1.0/StoragePowerActor")
	StorageMinerActorCodeCID   = actor.CodeCID("filecoin/1.0/StorageMinerActor")
	StorageMarketActorCodeCID  = actor.CodeCID("filecoin/1.0/StorageMarketActor")
	PaymentChannelActorCodeCID = actor.CodeCID("filecoin/1.0/PaymentChannelActor")
)

var staticActorCodeRegistry = &actorCodeRegistry{}

type actorCodeRegistry struct {
	code map[actor.CodeCID]vmr.ActorCode
}

func (r *actorCodeRegistry) _registerActor(cid actor.CodeCID, actor vmr.ActorCode) {
	r.code[cid] = actor
}

func (r *actorCodeRegistry) _loadActor(cid actor.CodeCID) (vmr.ActorCode, error) {
	a, ok := r.code[cid]
	if !ok {
		return nil, ErrActorNotFound
	}
	return a, nil
}

func RegisterActor(cid actor.CodeCID, actor vmr.ActorCode) {
	staticActorCodeRegistry._registerActor(cid, actor)
}

func LoadActor(cid actor.CodeCID) (vmr.ActorCode, error) {
	return staticActorCodeRegistry._loadActor(cid)
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

	RegisterActor(InitActorCodeCID, &sysactors.InitActorCode_I{})
	RegisterActor(CronActorCodeCID, cron)
	RegisterActor(StoragePowerActorCodeCID, &spc.StoragePowerActor_I{})
	RegisterActor(StorageMarketActorCodeCID, &market.StorageMarketActor_I{})

	// wire in CRON actions.
	// TODO: there's probably a better place to put this, but for now, do it here.
	cron.Actors_ = append(cron.Actors_, addr.StoragePowerActorAddr)
	cron.Actors_ = append(cron.Actors_, addr.StorageMarketActorAddr)
}
