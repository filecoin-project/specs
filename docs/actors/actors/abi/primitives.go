package abi

import cid "github.com/ipfs/go-cid"

// The abi package contains definitions of all types that cross the VM boundary and are used
// within actor code.
//
// Primitive types include numerics and opaque array types.

// Epoch number of the chain state, which acts as a proxy for time within the VM.
type ChainEpoch int64

// ActorCodeID identifies an actor's code (either one of the builtin actors,
// or, in the future, a CID of VM bytecode for a user-defined actor).
type ActorCodeID cid.Cid

// ActorID is a sequential number assigned to actors in the state tree.
// ActorIDs are assigned by the InitActor when an actor is created.
type ActorID int64

// MethodNum is an integer that represents a particular method
// in an actor's function table. These numbers are used to compress
// invocation of actor code, and to decouple human language concerns
// about method names from the ability to uniquely refer to a particular
// method.
//
// Consider MethodNum numbers to be similar in concerns as for
// offsets in function tables (in programming languages), and for
// tags in ProtocolBuffer fields. Tags in ProtocolBuffers recommend
// assigning a unique tag to a field and never reusing that tag.
// If a field is no longer used, the field name may change but should
// still remain defined in the code to ensure the tag number is not
// reused accidentally. The same should apply to the MethodNum
// associated with methods in Filecoin VM Actors.
type MethodNum int64

// Method params are the CBOR-serialization of a heterogenous array of values.
type MethodParams []byte

// TODO: remove this alias after actor types are realized from .id files.
type Bytes []byte

// TokenAmount is an amount of Filecoin tokens. This type is used within
// the VM in message execution, to account movement of tokens, payment
// of VM gas, and more.
type TokenAmount int64 // TODO bigint

// The randomness seed is a string of byte, distinguished from Randomness
// for expressiveness: it hasn't been given the needed entropy
type RandomnessSeed []byte

// Randomness is a string of random bytes
type Randomness []byte
