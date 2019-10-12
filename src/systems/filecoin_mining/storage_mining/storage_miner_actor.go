package storage_mining

import base_mining "github.com/filecoin-project/specs/systems/filecoin_mining"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import sealing "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (sma *StorageMinerActor_I) SubmitPoSt(postProof base_mining.PoStProof, sectorStateSets sector.SectorStateSets) {
	panic("TODO")
}

func (sma *StorageMinerActor_I) CommitSector(onChainInfo sealing.OnChainSealVerifyInfo) {
	currentSectorContent, found := sma.Sectors()[onChainInfo.SectorNumber()]
	if found {
		newSectorContent := append(currentSectorContent, onChainInfo)
		sma.Sectors()[onChainInfo.SectorNumber()] = newSectorContent
	} else {
		sma.Sectors()[onChainInfo.SectorNumber()] = []sealing.OnChainSealVerifyInfo{onChainInfo}
	}

	// TODO broadcast message on chain
}
