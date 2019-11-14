package deal

// TODO
func (css *CompactDealSet) DealsOn() []DealID {
	var dealID []DealID
	return dealID
}

// TODO
func (css *CompactDealSet) DealsOff() []DealID {
	var dealID []DealID
	return dealID
}

// TODO
func (css *CompactDealSet) Add(dealID DealID) CompactDealSet {
	var newCompactDealSet CompactDealSet
	return newCompactDealSet
}

// TODO
func (css *CompactDealSet) Remove(dealID DealID) CompactDealSet {
	var newCompactDealSet CompactDealSet
	return newCompactDealSet
}

// TODO
func (css *CompactDealSet) Extend(css2 CompactDealSet) CompactDealSet {
	var newCompactDealSet CompactDealSet
	return newCompactDealSet
}

type DealState int

const (
	PublishedDealState DealState = 0
	ActiveDealState    DealState = 1
)

// TODO
// get state of the deal corresponding to dealID
func (dss *DealStateSet) GetDealState(dealID DealID) DealState {
	return DealState(-1)
}

// TODO
// get list of dealIDs that are in Published state
func (dss *DealStateSet) PublishedDeals() []DealID {
	var dealID []DealID
	return dealID
}

// TODO
// get list of dealIDs that are in Active state
func (dss *DealStateSet) ActiveDeals() []DealID {
	var dealID []DealID
	return dealID
}

// TODO
// add a deal to DealStateSet and initialize it
// check if deal exists before Publish
func (dss *DealStateSet) Publish(dealID DealID) {
	panic("TODO")
}

// TODO
// activate an existing PublishedDeal
// check if deal is in PublishedState before Activate
func (dss *DealStateSet) Activate(dealID DealID) {
	panic("TODO")
}

// TODO
// remove a deal from DealStateSet
func (dss *DealStateSet) Clear(dealID DealID) {
	panic("TODO")
}
