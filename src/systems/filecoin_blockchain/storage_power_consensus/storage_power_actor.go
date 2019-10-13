package storage_power_actor

func (spa *StoragePowerActor) UpdatePower(address addr.Address, newPower BytesAmount) {
  spa.Miners()[addr.Address] = spa.Miners()[addr.Address] + newPower
}
