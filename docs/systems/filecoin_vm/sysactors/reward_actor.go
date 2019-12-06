package sysactors

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import ipld "github.com/filecoin-project/specs/libraries/ipld"

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////

func (a *RewardActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, RewardActorState) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.AbortAPI("IPLD lookup error")
	}
	state := DeserializeState(stateBytes.As_Bytes())
	return h, state
}
func UpdateReleaseRewardActorState(rt Runtime, h vmr.ActorStateHandle, st RewardActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *RewardActorState_I) CID() ipld.CID {
	panic("TODO")
}
func DeserializeState(x Bytes) RewardActorState {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (st *RewardActorState_I) _withdrawReward(rt vmr.Runtime, ownerAddr addr.Address) actor.TokenAmount {

	rewards, found := st.RewardMap()[ownerAddr]
	if !found {
		rt.AbortStateMsg("ra._withdrawReward: ownerAddr not found in RewardMap.")
	}

	rewardToWithdrawTotal := actor.TokenAmount(0)
	indicesToRemove := make([]int, len(rewards))

	for i, r := range rewards {
		elapsedEpoch := rt.CurrEpoch() - r.StartEpoch()
		unlockedReward := actor.TokenAmount(uint64(r.ReleaseRate()) * uint64(elapsedEpoch))
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

func (a *RewardActorCode_I) MintReward(rt vmr.Runtime) {
	// block reward function should live here
	// put Reward into RewardMap
	panic("TODO")
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
