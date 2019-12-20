package sysactors

import (
	"math"

	ipld "github.com/filecoin-project/specs/libraries/ipld"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////

var TODO = util.TODO

func (a *RewardActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, RewardActorState) {
	h := rt.AcquireState()
	stateCID := ipld.CID(h.Take())
	var state RewardActorState_I
	if !rt.IpldGet(stateCID, &state) {
		rt.AbortAPI("state not found")
	}
	return h, &state
}
func UpdateReleaseRewardActorState(rt Runtime, h vmr.ActorStateHandle, st RewardActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *RewardActorState_I) CID() ipld.CID {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (r *Reward_I) AmountVested(elapsedEpoch block.ChainEpoch) actor.TokenAmount {
	switch r.VestingFunction() {
	case VestingFunction_None:
		return r.Value()
	case VestingFunction_Linear:
		TODO() // BigInt
		vestedProportion := math.Max(1.0, float64(elapsedEpoch)/float64(r.StartEpoch()-r.EndEpoch()))
		return actor.TokenAmount(uint64(r.Value()) * uint64(vestedProportion))
	default:
		return actor.TokenAmount(0)
	}
}

func (st *RewardActorState_I) _withdrawReward(rt vmr.Runtime, ownerAddr addr.Address) actor.TokenAmount {

	rewards, found := st.RewardMap()[ownerAddr]
	if !found {
		rt.AbortStateMsg("ra._withdrawReward: ownerAddr not found in RewardMap.")
	}

	rewardToWithdrawTotal := actor.TokenAmount(0)
	indicesToRemove := make([]int, len(rewards))

	for i, r := range rewards {
		elapsedEpoch := rt.CurrEpoch() - r.StartEpoch()
		unlockedReward := r.AmountVested(elapsedEpoch)
		withdrawableReward := unlockedReward - r.AmountWithdrawn()

		if withdrawableReward < 0 {
			rt.AbortStateMsg("ra._withdrawReward: negative withdrawableReward.")
		}

		r.Impl().AmountWithdrawn_ = unlockedReward // modify rewards in place
		rewardToWithdrawTotal += withdrawableReward

		if r.AmountWithdrawn() == r.Value() {
			indicesToRemove = append(indicesToRemove, i)
		}
	}

	updatedRewards := removeIndices(rewards, indicesToRemove)
	st.RewardMap()[ownerAddr] = updatedRewards

	return rewardToWithdrawTotal
}

func (a *RewardActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// initialize Reward Map with investor accounts
	panic("TODO")
}

func (a *RewardActorCode_I) AwardBlockReward(
	rt vmr.Runtime,
	miner addr.Address,
	penalty actor.TokenAmount,
	minerStoragePower block.StoragePower,
	minerActiveSectorWeight block.SectorWeight,
	minerInactiveSectorWeight block.SectorWeight,
	currPledge actor.TokenAmount,
) {
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)

	inds := rt.CurrIndices()
	pledgeReq := inds.BlockReward_PledgeCollateralReq(minerActiveSectorWeight, minerInactiveSectorWeight, currPledge)
	currReward := inds.BlockReward_GetCurrRewardForMiner(minerStoragePower, currPledge)
	TODO()                                                                                // BigInt
	underPledge := math.Max(float64(actor.TokenAmount(0)), float64(pledgeReq-currPledge)) // 0 if over collateralized
	rewardToGarnish := math.Min(float64(currReward), float64(underPledge))

	TODO()
	// handle penalty here
	// also handle penalty greater than reward
	actualReward := currReward - actor.TokenAmount(rewardToGarnish)
	if rewardToGarnish > 0 {
		// Send fund to SPA for collateral
		rt.Send(
			addr.StoragePowerActorAddr,
			ai.Method_StoragePowerActor_AddBalance,
			[]util.Serialization{
				addr.Serialize_Address(miner),
			},
			actor.TokenAmount(rewardToGarnish),
		)
	}

	h, st := a.State(rt)
	if actualReward > 0 {
		// put Reward into RewardMap
		newReward := &Reward_I{
			StartEpoch_:      rt.CurrEpoch(),
			EndEpoch_:        rt.CurrEpoch(),
			Value_:           actualReward,
			AmountWithdrawn_: actor.TokenAmount(0),
			VestingFunction_: VestingFunction_None,
		}
		rewards, found := st.RewardMap()[miner]
		if !found {
			rewards = make([]Reward, 0)
		}
		rewards = append(rewards, newReward)
		st.Impl().RewardMap_[miner] = rewards
	}
	UpdateReleaseRewardActorState(rt, h, st)
}

// called by ownerAddress
func (a *RewardActorCode_I) WithdrawReward(rt vmr.Runtime) {
	// withdraw available funds from RewardMap
	h, st := a.State(rt)

	ownerAddr := rt.ImmediateCaller()
	withdrawableReward := st._withdrawReward(rt, ownerAddr)
	UpdateReleaseRewardActorState(rt, h, st)

	// send funds to owner
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    ownerAddr,
		Value_: withdrawableReward,
	})
}

func removeIndices(rewards []Reward, indices []int) []Reward {
	// remove fully paid out Rewards by indices
	panic("TODO")
}
