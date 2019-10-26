package runtime

import "errors"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"

// import sysactors "github.com/filecoin-project/specs/systems/filecoin_vm/sysactors"

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

// Addresses for singleton system actors
var (
	InitActorAddr           = &addr.Address_I{} // TODO
	CronActorAddr           = &addr.Address_I{} // TODO
	StoragePowerActorAddr   = &addr.Address_I{} // TODO
	StorageMarketActorAddr  = &addr.Address_I{} // TODO
	PaymentChannelActorAddr = &addr.Address_I{} // TODO
	BurntFundsActorAddr     = &addr.Address_I{} // TODO
)

var staticActorCodeRegistry = &actorCodeRegistry{}

type actorCodeRegistry struct {
	code map[actor.CodeCID]ActorCode
}

func (r *actorCodeRegistry) _registerActor(cid actor.CodeCID, actor ActorCode) {
	r.code[cid] = actor
}

func (r *actorCodeRegistry) _loadActor(cid actor.CodeCID) (ActorCode, error) {
	a, ok := r.code[cid]
	if !ok {
		return nil, ErrActorNotFound
	}
	return a, nil
}

func RegisterActor(cid actor.CodeCID, actor ActorCode) {
	staticActorCodeRegistry._registerActor(cid, actor)
}

func LoadActor(cid actor.CodeCID) (ActorCode, error) {
	return staticActorCodeRegistry._loadActor(cid)
}
