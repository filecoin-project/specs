package sector

import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

// TODO
func (css *CompactSectorSet) SectorsOn() []SectorNumber {
	var sectorNo []SectorNumber
	return sectorNo
}

// TODO
func (css *CompactSectorSet) SectorsOff() []SectorNumber {
	var sectorNo []SectorNumber
	return sectorNo
}

// TODO
func (css *CompactSectorSet) Add(sectorNo SectorNumber) CompactSectorSet {
	var newCompactSectorSet CompactSectorSet
	return newCompactSectorSet
}

// TODO
func (css *CompactSectorSet) Remove(sectorNo SectorNumber) CompactSectorSet {
	var newCompactSectorSet CompactSectorSet
	return newCompactSectorSet
}

// TODO
func (css *CompactSectorSet) DeepCopy() CompactSectorSet {
	var newCompactSectorSet CompactSectorSet
	return newCompactSectorSet
}

// TODO
func (css *CompactSectorSet) Extend(css2 CompactSectorSet) CompactSectorSet {
	var newCompactSectorSet CompactSectorSet
	return newCompactSectorSet
}

// TODO
func (css *CompactSectorSet) Contain(sectorNo SectorNumber) bool {
	return false
}

/////////////////////////////////////////////////////////////////////////

// TODO
func (css *CompactDealSet) DealsOn() []deal.DealID {
	var dealId []deal.DealID
	return dealId
}

// TODO
func (css *CompactDealSet) DealsOff() []deal.DealID {
	var dealId []deal.DealID
	return dealId
}

// TODO
func (css *CompactDealSet) Add(dealId deal.DealID) CompactDealSet {
	var newCompactDealSet CompactDealSet
	return newCompactDealSet
}

// TODO
func (css *CompactDealSet) Remove(dealId deal.DealID) CompactDealSet {
	var newCompactDealSet CompactDealSet
	return newCompactDealSet
}

// TODO
func (css *CompactDealSet) Extend(css2 CompactDealSet) CompactDealSet {
	var newCompactDealSet CompactDealSet
	return newCompactDealSet
}
