package storage_power_consensus

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"

func (spa *StoragePowerActor_I) UpdatePower(address addr.Address, newPower BytesAmount) {
	panic("TODO")
	// spa.Miners()[addr.Address] = spa.Miners()[addr.Address] + newPower
}
