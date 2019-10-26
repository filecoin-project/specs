package sysactors

import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"

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

	cron := &CronActorCode_I{}

	vmr.RegisterActor(vmr.InitActorCodeCID, &InitActorCode_I{})
	vmr.RegisterActor(vmr.CronActorCodeCID, cron)

	// wire in CRON actions.
	// TODO: there's probably a better place to put this, but for now, do it here.
	cron.Actors_ = append(cron.Actors_, vmr.StoragePowerActorAddr)
	cron.Actors_ = append(cron.Actors_, vmr.StorageMarketActorAddr)
}
