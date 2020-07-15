package power

import (
	"reflect"

	addr "github.com/filecoin-project/go-address"
	cid "github.com/ipfs/go-cid"
	errors "github.com/pkg/errors"

	abi "github.com/filecoin-project/specs-actors/actors/abi"
	big "github.com/filecoin-project/specs-actors/actors/abi/big"
	. "github.com/filecoin-project/specs-actors/actors/util"
	adt "github.com/filecoin-project/specs-actors/actors/util/adt"
)

type State struct {
	TotalRawBytePower abi.StoragePower
	// TotalBytesCommitted includes claims from miners below min power threshold
	TotalBytesCommitted  abi.StoragePower
	TotalQualityAdjPower abi.StoragePower
	// TotalQABytesCommitted includes claims from miners below min power threshold
	TotalQABytesCommitted abi.StoragePower
	TotalPledgeCollateral abi.TokenAmount
	MinerCount            int64
	// Number of miners having proven the minimum consensus power.
	MinerAboveMinPowerCount int64

	// A queue of events to be triggered by cron, indexed by epoch.
	CronEventQueue cid.Cid // Multimap, (HAMT[ChainEpoch]AMT[CronEvent]

	// Last chain epoch OnEpochTickEnd was called on
	LastEpochTick abi.ChainEpoch

	// Claimed power for each miner.
	Claims cid.Cid // Map, HAMT[address]Claim

	ProofValidationBatch *cid.Cid
}

type Claim struct {
	// Sum of raw byte power for a miner's sectors.
	RawBytePower abi.StoragePower

	// Sum of quality adjusted power for a miner's sectors.
	QualityAdjPower abi.StoragePower
}

type CronEvent struct {
	MinerAddr       addr.Address
	CallbackPayload []byte
}

type AddrKey = adt.AddrKey

func ConstructState(emptyMapCid, emptyMMapCid cid.Cid) *State {
	return &State{
		TotalRawBytePower:       abi.NewStoragePower(0),
		TotalBytesCommitted:     abi.NewStoragePower(0),
		TotalQualityAdjPower:    abi.NewStoragePower(0),
		TotalQABytesCommitted:   abi.NewStoragePower(0),
		TotalPledgeCollateral:   abi.NewTokenAmount(0),
		LastEpochTick:           -1,
		CronEventQueue:          emptyMapCid,
		Claims:                  emptyMapCid,
		MinerCount:              0,
		MinerAboveMinPowerCount: 0,
	}
}

// MinerNominalPowerMeetsConsensusMinimum is used to validate Election PoSt
// winners outside the chain state. If the miner has over a threshold of power
// the miner meets the minimum.  If the network is a below a threshold of
// miners and has power > zero the miner meets the minimum.
func (st *State) MinerNominalPowerMeetsConsensusMinimum(s adt.Store, miner addr.Address) (bool, error) { //nolint:deadcode,unused
	claim, ok, err := st.GetClaim(s, miner)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.Errorf("no claim for actor %v", miner)
	}

	minerNominalPower := claim.QualityAdjPower

	// if miner is larger than min power requirement, we're set
	if minerNominalPower.GreaterThanEqual(ConsensusMinerMinPower) {
		return true, nil
	}

	// otherwise, if ConsensusMinerMinMiners miners meet min power requirement, return false
	if st.MinerAboveMinPowerCount > ConsensusMinerMinMiners {
		return false, nil
	}

	// If fewer than ConsensusMinerMinMiners over threshold miner can win a block with non-zero power
	return minerNominalPower.GreaterThanEqual(abi.NewStoragePower(0)), nil
}

// Parameters may be negative to subtract.
func (st *State) AddToClaim(s adt.Store, miner addr.Address, power abi.StoragePower, qapower abi.StoragePower) error {
	oldClaim, ok, err := st.GetClaim(s, miner)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("no claim for actor %v", miner)
	}

	// TotalBytes always update directly
	st.TotalQABytesCommitted = big.Add(st.TotalQABytesCommitted, qapower)
	st.TotalBytesCommitted = big.Add(st.TotalBytesCommitted, power)

	newClaim := Claim{
		RawBytePower:    big.Add(oldClaim.RawBytePower, power),
		QualityAdjPower: big.Add(oldClaim.QualityAdjPower, qapower),
	}

	prevBelow := oldClaim.QualityAdjPower.LessThan(ConsensusMinerMinPower)
	stillBelow := newClaim.QualityAdjPower.LessThan(ConsensusMinerMinPower)

	if prevBelow && !stillBelow {
		// just passed min miner size
		st.MinerAboveMinPowerCount++
		st.TotalQualityAdjPower = big.Add(st.TotalQualityAdjPower, newClaim.QualityAdjPower)
		st.TotalRawBytePower = big.Add(st.TotalRawBytePower, newClaim.RawBytePower)
	} else if !prevBelow && stillBelow {
		// just went below min miner size
		st.MinerAboveMinPowerCount--
		st.TotalQualityAdjPower = big.Sub(st.TotalQualityAdjPower, oldClaim.QualityAdjPower)
		st.TotalRawBytePower = big.Sub(st.TotalRawBytePower, oldClaim.RawBytePower)
	} else if !prevBelow && !stillBelow {
		// Was above the threshold, still above
		st.TotalQualityAdjPower = big.Add(st.TotalQualityAdjPower, qapower)
		st.TotalRawBytePower = big.Add(st.TotalRawBytePower, power)
	}

	AssertMsg(newClaim.RawBytePower.GreaterThanEqual(big.Zero()), "negative claimed raw byte power: %v", newClaim.RawBytePower)
	AssertMsg(newClaim.QualityAdjPower.GreaterThanEqual(big.Zero()), "negative claimed quality adjusted power: %v", newClaim.QualityAdjPower)
	AssertMsg(st.MinerAboveMinPowerCount >= 0, "negative number of miners larger than min: %v", st.MinerAboveMinPowerCount)
	return st.setClaim(s, miner, &newClaim)
}

func (st *State) GetClaim(s adt.Store, a addr.Address) (*Claim, bool, error) {
	hm, err := adt.AsMap(s, st.Claims)
	if err != nil {
		return nil, false, err
	}

	var out Claim
	found, err := hm.Get(AddrKey(a), &out)
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to get claim for address %v from store %s", a, st.Claims)
	}
	if !found {
		return nil, false, nil
	}
	return &out, true, nil
}

func (st *State) addPledgeTotal(amount abi.TokenAmount) {
	st.TotalPledgeCollateral = big.Add(st.TotalPledgeCollateral, amount)
	Assert(st.TotalPledgeCollateral.GreaterThanEqual(big.Zero()))
}

func (st *State) appendCronEvent(store adt.Store, epoch abi.ChainEpoch, event *CronEvent) error {
	mmap, err := adt.AsMultimap(store, st.CronEventQueue)
	if err != nil {
		return err
	}

	err = mmap.Add(epochKey(epoch), event)
	if err != nil {
		return errors.Wrapf(err, "failed to store cron event at epoch %v for miner %v", epoch, event)
	}
	st.CronEventQueue, err = mmap.Root()
	if err != nil {
		return err
	}
	return nil
}

func (st *State) loadCronEvents(store adt.Store, epoch abi.ChainEpoch) ([]CronEvent, error) {
	mmap, err := adt.AsMultimap(store, st.CronEventQueue)
	if err != nil {
		return nil, err
	}

	var events []CronEvent
	var ev CronEvent
	err = mmap.ForEach(epochKey(epoch), &ev, func(i int64) error {
		events = append(events, ev)
		return nil
	})
	return events, err
}

func (st *State) clearCronEvents(store adt.Store, epoch abi.ChainEpoch) error {
	mmap, err := adt.AsMultimap(store, st.CronEventQueue)
	if err != nil {
		return err
	}

	err = mmap.RemoveAll(epochKey(epoch))
	if err != nil {
		return errors.Wrapf(err, "failed to clear cron events")
	}
	st.CronEventQueue, err = mmap.Root()
	if err != nil {
		return err
	}
	return nil
}

func (st *State) setClaim(s adt.Store, a addr.Address, claim *Claim) error {
	Assert(claim.RawBytePower.GreaterThanEqual(big.Zero()))
	Assert(claim.QualityAdjPower.GreaterThanEqual(big.Zero()))

	hm, err := adt.AsMap(s, st.Claims)
	if err != nil {
		return err
	}

	if err = hm.Put(AddrKey(a), claim); err != nil {
		return errors.Wrapf(err, "failed to put claim with address %s power %v in store %s", a, claim, st.Claims)
	}

	st.Claims, err = hm.Root()
	if err != nil {
		return err
	}
	return nil
}

func (st *State) deleteClaim(s adt.Store, a addr.Address) error {
	hm, err := adt.AsMap(s, st.Claims)
	if err != nil {
		return err
	}

	if err = hm.Delete(AddrKey(a)); err != nil {
		return errors.Wrapf(err, "failed to delete claim at address %s from store %s", a, st.Claims)
	}
	st.Claims, err = hm.Root()
	if err != nil {
		return err
	}
	return nil
}

// CurrentTotalPower returns current power values accounting for minimum miner
// and minimum power
func CurrentTotalPower(st *State) (abi.StoragePower, abi.StoragePower) {
	if st.MinerAboveMinPowerCount < ConsensusMinerMinMiners {
		return st.TotalBytesCommitted, st.TotalQABytesCommitted
	}
	return st.TotalRawBytePower, st.TotalQualityAdjPower
}

func epochKey(e abi.ChainEpoch) adt.Keyer {
	return adt.IntKey(int64(e))
}

func init() {
	// Check that ChainEpoch is indeed a signed integer to confirm that epochKey is making the right interpretation.
	var e abi.ChainEpoch
	if reflect.TypeOf(e).Kind() != reflect.Int64 {
		panic("incorrect chain epoch encoding")
	}
}
