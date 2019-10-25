package runtime

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import ipld "github.com/filecoin-project/specs/libraries/ipld"
import st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import gascost "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/gascost"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import util "github.com/filecoin-project/specs/util"

// InvocInput represents inputs to a particular actor invocation.
type InvocInput struct {
	Runtime     Runtime
	InTree      st.StateTree
	OriginActor actor.Actor
	CallSeqNum  actor.CallSeqNum
	FromActor   actor.Actor
	ToActor     actor.Actor
	Method      actor.MethodNum
	Params      actor.MethodParams
	Value       actor.TokenAmount
	GasLimit    msg.GasAmount
	GasUsed     msg.GasAmount
	// GasPrice    GasPrice
}

type InvocOutput struct {
	OutTree     st.StateTree
	ExitCode    msg.ExitCode
	ReturnValue []byte
	GasUsed     msg.GasAmount
}

func ErrorInvocOutput(tree st.StateTree, ec msg.ExitCode) InvocOutput {
	return InvocOutput{
		OutTree:     tree,
		GasUsed:     gascost.CodeLookupFail,
		ExitCode:    exitcode.InvalidMethod,
		ReturnValue: nil,
	}
}

func (r *Runtime_I) Fatal(string) Runtime_Fatal_FunRet {
	panic("TODO")
}

func (r *Runtime_I) Send(to addr.Address, method actor.MethodNum, params actor.MethodParams, value actor.TokenAmount) msg.MessageReceipt {
	panic("TODO")
}

func (s *VMState_I) Epoch() st.Epoch {
	panic("TODO")
}

func (s *VMState_I) Balance(id actor.ActorID) actor.TokenAmount {
	panic("TODO")
}

func (s *VMState_I) Randomness(e st.Epoch, offset uint64) Randomness {
	panic("TODO")
}

func (s *VMState_I) ComputeActorAddress(creator addr.Address, nonce actor.CallSeqNum) addr.Address {
	panic("TODO")
}

func (s *VMStorage_I) Put(o IpldObject) VMStorage_Put_FunRet {
	panic("TODO")
}

func (s *VMStorage_I) Get(c ipld.CID) VMStorage_Get_FunRet {
	panic("TODO")
}

func (s *VMStorage_I) Commit(old ipld.CID, new ipld.CID) error {
	panic("TODO")
}

func (s *VMStorage_I) Head() ipld.CID {
	panic("TODO")
}

func (rt *Runtime_I) CurrEpoch() block.ChainEpoch {
	panic("TODO")
}

func (rt *Runtime_I) ReadState() util.Any {
	panic("TODO")
}

func (rt *Runtime_I) AcquireState() ActorStateHandle {
	panic("TODO")
}

func (rt *Runtime_I) ValidateCallerIs(addr.Address) Runtime_ValidateCallerIs_FunRet {
	panic("TODO")
}

func (rt *Runtime_I) CronActorAddress() addr.Address {
	panic("TODO")
}

func (rt *Runtime_I) StoragePowerActorAddress() addr.Address {
	panic("TODO")
}
