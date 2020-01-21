package builtin

import (
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// The built-in actor code IDs
var SystemActorCodeID abi.ActorCodeID
var InitActorCodeID abi.ActorCodeID
var CronActorCodeID abi.ActorCodeID
var AccountActorCodeID abi.ActorCodeID
var StoragePowerActorCodeID abi.ActorCodeID
var StorageMinerActorCodeID abi.ActorCodeID
var StorageMarketActorCodeID abi.ActorCodeID
var PaymentChannelActorCodeID abi.ActorCodeID
var MultisigActorCodeID abi.ActorCodeID
var RewardActorCodeID abi.ActorCodeID

func init() {
	builder := cid.V1Builder{Codec: cid.Raw, MhType: mh.IDENTITY}
	makeBuiltin := func(s string) abi.ActorCodeID {
		c, err := builder.Sum([]byte(s))
		if err != nil {
			panic(err)
		}
		id := abi.ActorCodeID(c)
		return id
	}

	SystemActorCodeID = makeBuiltin("fil/1/system")
	InitActorCodeID = makeBuiltin("fil/1/init")
	CronActorCodeID = makeBuiltin("fil/1/cron")
	AccountActorCodeID = makeBuiltin("fil/1/account")
	StoragePowerActorCodeID = makeBuiltin("fil/1/storagepower")
	StorageMinerActorCodeID = makeBuiltin("fil/1/storageminer")
	StorageMarketActorCodeID = makeBuiltin("fil/1/storagemarket")
	PaymentChannelActorCodeID = makeBuiltin("fil/1/paymentchannel")
	MultisigActorCodeID = makeBuiltin("fil/1/multisig")
	RewardActorCodeID = makeBuiltin("fil/1/reward")
}
