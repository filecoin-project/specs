package reward

import (
	"math"

	addr "github.com/filecoin-project/go-address"
	actor "github.com/filecoin-project/specs/actors"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	serde "github.com/filecoin-project/specs/actors/serde"
	autil "github.com/filecoin-project/specs/actors/util"
	cid "github.com/ipfs/go-cid"
)

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime

var IMPL_FINISH = autil.IMPL_FINISH
var IMPL_TODO = autil.IMPL_TODO
var TODO = autil.TODO

type VestingFunction int

const (
	None VestingFunction = iota
	Linear
	// TODO: potential options
	// PieceWise
	// Quadratic
	// Exponential
)

type Reward struct {
	VestingFunction
	StartEpoch      abi.ChainEpoch
	EndEpoch        abi.ChainEpoch
	Value           abi.TokenAmount
	AmountWithdrawn abi.TokenAmount
}

func (r *Reward) AmountVested(elapsedEpoch abi.ChainEpoch) abi.TokenAmount {
	switch r.VestingFunction {
	case None:
		return r.Value
	case Linear:
		TODO() // BigInt
		vestedProportion := math.Max(1.0, float64(elapsedEpoch)/float64(r.StartEpoch-r.EndEpoch))
		return abi.TokenAmount(uint64(r.Value) * uint64(vestedProportion))
	default:
		return abi.TokenAmount(0)
	}
}

// ownerAddr to a collection of Reward
type RewardBalanceAMT map[addr.Address][]Reward

type RewardActorState struct {
	RewardMap RewardBalanceAMT
}

func (st *RewardActorState) CID() cid.Cid {
	panic("TODO")
}

func (st *RewardActorState) _withdrawReward(rt vmr.Runtime, ownerAddr addr.Address) abi.TokenAmount {
	rewards, found := st.RewardMap[ownerAddr]
	if !found {
		rt.AbortStateMsg("ra._withdrawReward: ownerAddr not found in RewardMap.")
	}

	rewardToWithdrawTotal := abi.TokenAmount(0)
	indicesToRemove := make([]int, len(rewards))

	for i, r := range rewards {
		elapsedEpoch := rt.CurrEpoch() - r.StartEpoch
		unlockedReward := r.AmountVested(elapsedEpoch)
		withdrawableReward := unlockedReward - r.AmountWithdrawn

		if withdrawableReward < 0 {
			rt.AbortStateMsg("ra._withdrawReward: negative withdrawableReward.")
		}

		r.AmountWithdrawn = unlockedReward // modify rewards in place
		rewardToWithdrawTotal += withdrawableReward

		if r.AmountWithdrawn == r.Value {
			indicesToRemove = append(indicesToRemove, i)
		}
	}

	updatedRewards := removeIndices(rewards, indicesToRemove)
	st.RewardMap[ownerAddr] = updatedRewards

	return rewardToWithdrawTotal
}

type RewardActor struct{}

func (a *RewardActor) Constructor(rt vmr.Runtime) InvocOutput {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	// initialize Reward Map with investor accounts
	panic("TODO")
}

func (a *RewardActor) State(rt Runtime) (vmr.ActorStateHandle, RewardActorState) {
	h := rt.AcquireState()
	stateCID := cid.Cid(h.Take())
	var state RewardActorState
	if !rt.IpldGet(stateCID, &state) {
		rt.AbortAPI("state not found")
	}
	return h, state
}

func (a *RewardActor) WithdrawReward(rt vmr.Runtime) {
	vmr.RT_ValidateImmediateCallerIsSignable(rt)
	ownerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	// withdraw available funds from RewardMap
	withdrawableReward := st._withdrawReward(rt, ownerAddr)
	UpdateReleaseRewardActorState(rt, h, st)

	rt.SendFunds(ownerAddr, withdrawableReward)
}

func (a *RewardActor) AwardBlockReward(
	rt vmr.Runtime,
	miner addr.Address,
	penalty abi.TokenAmount,
	minerNominalPower abi.StoragePower,
	currPledge abi.TokenAmount,
) {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)

	inds := rt.CurrIndices()
	pledgeReq := inds.PledgeCollateralReq(minerNominalPower)
	currReward := inds.GetCurrBlockRewardForMiner(minerNominalPower, currPledge)
	TODO()                                                                              // BigInt
	underPledge := math.Max(float64(abi.TokenAmount(0)), float64(pledgeReq-currPledge)) // 0 if over collateralized
	rewardToGarnish := math.Min(float64(currReward), float64(underPledge))

	TODO()
	// handle penalty here
	// also handle penalty greater than reward
	actualReward := currReward - abi.TokenAmount(rewardToGarnish)
	if rewardToGarnish > 0 {
		// Send fund to SPA for collateral
		rt.Send(
			builtin.StoragePowerActorAddr,
			builtin.Method_StoragePowerActor_AddBalance,
			serde.MustSerializeParams(miner),
			abi.TokenAmount(rewardToGarnish),
		)
	}

	h, st := a.State(rt)
	if actualReward > 0 {
		// put Reward into RewardMap
		newReward := &Reward{
			StartEpoch:      rt.CurrEpoch(),
			EndEpoch:        rt.CurrEpoch(),
			Value:           actualReward,
			AmountWithdrawn: abi.TokenAmount(0),
			VestingFunction: None,
		}
		rewards, found := st.RewardMap[miner]
		if !found {
			rewards = make([]Reward, 0)
		}
		rewards = append(rewards, *newReward)
		st.RewardMap[miner] = rewards
	}
	UpdateReleaseRewardActorState(rt, h, st)
}

func UpdateReleaseRewardActorState(rt Runtime, h vmr.ActorStateHandle, st RewardActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(&st))
	h.UpdateRelease(newCID)
}

func removeIndices(rewards []Reward, indices []int) []Reward {
	// remove fully paid out Rewards by indices
	panic("TODO")
}
