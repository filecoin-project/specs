package storage_power_consensus

func (spa *StoragePowerActor) UpdatePower(address addr.Address, newPower BytesAmount) {
	spa.Miners()[addr.Address] = spa.Miners()[addr.Address] + newPower
}
