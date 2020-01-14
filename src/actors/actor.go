package actors

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	util "github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
)

// CallSeqNum is an invocation (Call) sequence (Seq) number (Num).
// This is a value used for securing against replay attacks:
// each AccountActor (user) invocation must have a unique CallSeqNum
// value. The sequenctiality of the numbers is used to make it
// easy to verify, and to order messages.
//
// Q&A
// - > Does it have to be sequential?
//   No, a random nonce could work against replay attacks, but
//   making it sequential makes it much easier to verify.
// - > Can it be used to order events?
//   Yes, a user may submit N separate messages with increasing
//   sequence number, causing them to execute in order.
//
type CallSeqNum util.UVarint

// Code is a serialized object that contains the code for an Actor.
// Until we accept external user-provided contracts, this is the
// serialized code for the actor in the Filecoin Specification.
type Code abi.Bytes

type ActorSystemStateCID cid.Cid

// Actor is a base computation object in the Filecoin VM. Similar
// to Actors in the Actor Model (programming), or Objects in Object-
// Oriented Programming, or Ethereum Contracts in the EVM.
//
// ActorState represents the on-chain storage all actors keep.
type ActorState struct {
	// Identifies the code this actor executes.
	CodeID_ abi.ActorCodeID
	// CID of the root of optional actor-specific sub-state.
	State_ ActorSubstateCID
	// Balance of tokens held by this actor.
	Balance_ abi.TokenAmount
	// Expected sequence number of the next message sent by this actor.
	// Initially zero, incremented when an account actor originates a top-level message.
	// Always zero for other abi.
	CallSeqNum_ CallSeqNum
}

func (st *ActorState) CID() cid.Cid {
	panic("TODO")
}

func (a *ActorState) CodeID() abi.ActorCodeID {
	return a.CodeID_
}

func (a *ActorState) State() ActorSubstateCID {
	return a.State_
}

func (a *ActorState) Balance() abi.TokenAmount {
	return a.Balance_
}

func (a *ActorState) CallSeqNum() CallSeqNum {
	return a.CallSeqNum_
}

type ActorSubstateCID cid.Cid

func (x ActorSubstateCID) Ref() *ActorSubstateCID {
	return &x
}

// ActorState represents the on-chain storage actors keep. This type is a
// union of concrete types, for each of the Actors:
// - InitActor
// - CronActor
// - AccountActor
// - PaymentChannelActor
// - StoragePowerActor
// - StorageMinerActor
// - StroageMarketActor
//
// TODO: move this into a directory inside the VM that patches in all
// the actors from across the system. this will be where we declare/mount
// all actors in the VM.
// type ActorState union {
//     Init struct {
//         AddressMap  {addr.Address: ActorID}
//         NextID      ActorID
//     }
// }
