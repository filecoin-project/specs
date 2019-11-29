package storage_mining

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
)

func (st *StorageMinerActorState_I) _isChallenged() bool {
	return st.ChallengeStatus().IsChallenged()
}

func (st *StorageMinerActorState_I) _canBeElected(epoch block.ChainEpoch) bool {
	return st.ChallengeStatus().CanBeElected(epoch)
}

func (st *StorageMinerActorState_I) _challengeHasExpired(epoch block.ChainEpoch) bool {
	return st.ChallengeStatus().ChallengeHasExpired(epoch)
}

func (st *StorageMinerActorState_I) ShouldChallenge(currEpoch block.ChainEpoch, challengeFreePeriod block.ChainEpoch) bool {
	return st.ChallengeStatus().ShouldChallenge(currEpoch, challengeFreePeriod)
}

func (st *StorageMinerActorState_I) _processStagedCommittedSectors(rt Runtime) {
	for sectorNo, stagedInfo := range st.StagedCommittedSectors() {
		st.Sectors()[sectorNo] = stagedInfo.Sector()
		st.Impl().ProvingSet_.Add(sectorNo)
		st.SectorTable().Impl().CommittedSectors_.Add(sectorNo)
		st.SectorUtilization()[sectorNo] = stagedInfo.Utilization()
	}

	// empty StagedCommittedSectors
	st.Impl().StagedCommittedSectors_ = make(map[sector.SectorNumber]StagedCommittedSectorInfo)
}

func (st *StorageMinerActorState_I) _updateSectorUtilization(rt Runtime, lastPoStResponse block.ChainEpoch) []deal.DealID {
	// TODO: verify if we should update Sector utilization for failing sectors
	// depends on decision around collateral requirement for inactive power
	// and what happens when a failing sector expires

	ret := make([]deal.DealID, 0)

	for _, sectorNo := range st.Impl().ProvingSet_.SectorsOn() {

		utilizationInfo := st._getUtilizationInfo(rt, sectorNo)
		newUtilization := utilizationInfo.CurrUtilization()

		currEpoch := rt.CurrEpoch()
		newExpiredDealIDs := make([]deal.DealID, 0)
		newExpiredDeals := utilizationInfo.DealExpirationAMT().Impl().ExpiredDealsInRange(lastPoStResponse, currEpoch)

		for _, expiredDeal := range newExpiredDeals {
			expiredPower := expiredDeal.Power()
			newUtilization -= expiredPower
			newExpiredDealIDs = append(newExpiredDealIDs, expiredDeal.ID())

		}

		st.SectorUtilization()[sectorNo].Impl().CurrUtilization_ = newUtilization
		ret = append(ret, newExpiredDealIDs...)
	}

	return ret

}

func (st *StorageMinerActorState_I) _getActivePower(rt Runtime) block.StoragePower {
	var activePower = block.StoragePower(0)

	for _, sectorNo := range st.SectorTable().Impl().ActiveSectors_.SectorsOn() {
		utilizationInfo, found := st.SectorUtilization()[sectorNo]
		if !found {
			rt.Abort("sm._getActivePower: sectorNo not found in SectorUtilization")
		}
		activePower += utilizationInfo.CurrUtilization()
	}

	return activePower
}

func (st *StorageMinerActorState_I) _getInactivePower(rt Runtime) block.StoragePower {
	var inactivePower = block.StoragePower(0)

	// iterate over sectorNo in CommittedSectors, RecoveringSectors, and FailingSectors
	inactiveProvingSet := st.SectorTable().Impl().CommittedSectors_.Extend(st.SectorTable().RecoveringSectors())
	inactiveSectorSet := inactiveProvingSet.Extend(st.SectorTable().FailingSectors())

	for _, sectorNo := range inactiveSectorSet.SectorsOn() {
		utilizationInfo := st._getUtilizationInfo(rt, sectorNo)
		inactivePower += utilizationInfo.CurrUtilization()
	}

	return inactivePower
}

// move Sector from Active/Failing
// into Cleared State which means deleting the Sector from state
// remove SectorNumber from all states on chain
// update SectorTable
func (st *StorageMinerActorState_I) _updateClearSector(rt Runtime, sectorNo sector.SectorNumber) {
	sectorState := st.Sectors()[sectorNo].State()
	switch sectorState.StateNumber {
	case SectorActiveSN:
		// expiration case
		st.SectorTable().Impl().ActiveSectors_.Remove(sectorNo)
	case SectorFailingSN:
		// expiration and termination cases
		st.SectorTable().Impl().FailingSectors_.Remove(sectorNo)
	default:
		// Committed and Recovering should not go to Cleared directly
		rt.Abort("invalid state in clearSector")
	}

	delete(st.Sectors(), sectorNo)
	delete(st.SectorUtilization(), sectorNo)
	st.ProvingSet_.Remove(sectorNo)
	st.SectorExpirationQueue().Remove(sectorNo)

	// Send message to SMA
}

// move Sector from Committed/Recovering into Active State
// reset FaultCount to zero
// update SectorTable
func (st *StorageMinerActorState_I) _updateActivateSector(rt Runtime, sectorNo sector.SectorNumber) {
	sectorState := st.Sectors()[sectorNo].State()
	switch sectorState.StateNumber {
	case SectorCommittedSN:
		st.SectorTable().Impl().CommittedSectors_.Remove(sectorNo)
	case SectorRecoveringSN:
		st.SectorTable().Impl().RecoveringSectors_.Remove(sectorNo)
	default:
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("sm._updateActivateSector: invalid state in activateSector")
	}

	st.Sectors()[sectorNo].Impl().State_ = SectorActive()
	st.SectorTable().Impl().ActiveSectors_.Add(sectorNo)
}

// failSector moves Sector from Active/Committed/Recovering into Failing State
// and increments FaultCount if asked to do so (DeclareFaults does not increment faultCount)
// move Sector from Failing to Cleared State if increment results in faultCount exceeds MaxFaultCount
// update SectorTable
// remove from ProvingSet
func (st *StorageMinerActorState_I) _updateFailSector(rt Runtime, sectorNo sector.SectorNumber, increment bool) {
	newFaultCount := st.Sectors()[sectorNo].State().FaultCount

	if increment {
		newFaultCount += 1
	}

	state := st.Sectors()[sectorNo].State()
	switch state.StateNumber {
	case SectorActiveSN:
		// wont be terminated from Active
		st.SectorTable().Impl().ActiveSectors_.Remove(sectorNo)
		st.SectorTable().Impl().FailingSectors_.Add(sectorNo)
		st.ProvingSet_.Remove(sectorNo)
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	case SectorCommittedSN:
		st.SectorTable().Impl().CommittedSectors_.Remove(sectorNo)
		st.SectorTable().Impl().FailingSectors_.Add(sectorNo)
		st.ProvingSet_.Remove(sectorNo)
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	case SectorRecoveringSN:
		st.SectorTable().Impl().RecoveringSectors_.Remove(sectorNo)
		st.SectorTable().Impl().FailingSectors_.Add(sectorNo)
		st.ProvingSet_.Remove(sectorNo)
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	case SectorFailingSN:
		// no change to SectorTable but increase in FaultCount
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	default:
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("Invalid sector state in CronAction")
	}

	if newFaultCount > MAX_CONSECUTIVE_FAULTS {
		// TODO: heavy penalization: slash pledge collateral and delete sector
		// TODO: SendMessage(SPA.SlashPledgeCollateral)

		st._updateClearSector(rt, sectorNo)
		st.SectorTable().Impl().TerminatedFaults_.Add(sectorNo)
	}
}

func (st *StorageMinerActorState_I) _updateExpireSectors(rt Runtime) {
	currEpoch := rt.CurrEpoch()

	queue := st.SectorExpirationQueue()
	for queue.Peek().Expiration() <= currEpoch {
		expiredSectorNo := queue.Pop().SectorNumber()

		state := st.Sectors()[expiredSectorNo].State()

		switch state.StateNumber {
		case SectorActiveSN:
			// Note: in order to verify if something was stored in the past, one must
			// scan the chain. SectorNumber can be re-used.

			// do nothing about deal payment
			// it will be evaluated after _updateSectorUtilization
			st._updateClearSector(rt, expiredSectorNo)
		case SectorFailingSN:
			// TODO: check if there is any fault that we should handle here
			// If a SectorFailing Expires, return remaining StorageDealCollateral and remove sector
			// SendMessage(sma.SettleExpiredDeals(sc.DealIDs()))

			// a failing sector expires, no change to FaultCount
			st._updateClearSector(rt, expiredSectorNo)
		default:
			// Note: SectorCommittedSN, SectorRecoveringSN transition first to SectorFailingSN, then expire
			rt.Abort("Invalid sector state in SectorExpirationQueue")
		}
	}

	// Return PledgeCollateral for active expirations
	// SendMessage(spa.Depledge) // TODO
	rt.Abort("TODO: refactor use of this method in order for caller to send this message")
}

func (st *StorageMinerActorState_I) _assertSectorDidNotExist(rt Runtime, sectorNo sector.SectorNumber) {
	_, found := st.Sectors()[sectorNo]
	if found {
		rt.Abort("sm._assertSectorDidNotExist: sector already exists.")
	}
}

func (st *StorageMinerActorState_I) _getUtilizationInfo(rt Runtime, sectorNo sector.SectorNumber) sector.SectorUtilizationInfo {
	utilizationInfo, found := st.SectorUtilization()[sectorNo]

	if !found {
		rt.Abort("sm._getUtilizationInfo: utilization info not found.")
	}

	return utilizationInfo
}

func (st *StorageMinerActorState_I) _initializeUtilizationInfo(rt Runtime, deals []deal.OnChainDeal) sector.SectorUtilizationInfo {

	var maxUtilization block.StoragePower
	var dealExpirationAMT deal.DealExpirationAMT

	for _, d := range deals {
		dealID := d.ID()
		dealExpiration := d.Deal().Proposal().EndEpoch()
		// TODO: verify what counts towards power here
		// There is PayloadSize, OverheadSize, and Total, see piece.id
		dealPayloadPower := block.StoragePower(d.Deal().Proposal().PieceSize().PayloadSize())

		expirationValue := &deal.DealExpirationValue_I{
			ID_:    dealID,
			Power_: dealPayloadPower,
		}
		dealExpirationAMT.Impl().Add(dealExpiration, expirationValue)

		maxUtilization += dealPayloadPower

	}

	initialUtilizationInfo := &sector.SectorUtilizationInfo_I{
		DealExpirationAMT_: dealExpirationAMT,
		MaxUtilization_:    maxUtilization,
		CurrUtilization_:   maxUtilization,
	}

	return initialUtilizationInfo

}
