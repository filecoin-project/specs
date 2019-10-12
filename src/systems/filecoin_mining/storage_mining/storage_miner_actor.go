package storage_mining

import base_mining "github.com/filecoin-project/specs/systems/filecoin_mining"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import sealing "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (sma *StorageMinerActor_I) SubmitPoSt(postProof base_mining.PoStProof, sectorStateSets sector.SectorStateSets) {
	panic("TODO")
}

func (sma *StorageMinerActor_I) CommitSector(onChainInfo sealing.OnChainSealVerifyInfo) {
	currSector, found := sma.Sectors()[onChainInfo.SectorNumber()]
	if !found {
		sma.Sectors()[onChainInfo.SectorNumber()] = []sealing.OnChainSealVerifyInfo{onChainInfo}
	} else {
		newSector := append(currSector, onChainInfo)
		sma.Sectors()[onChainInfo.SectorNumber()] = newSector
	}

	// TODO broadcast message on chain
}
