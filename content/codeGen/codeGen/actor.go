package codeGen

type ActorState interface {
	Code() ActorSubType
	Head() ActorSubState
	Nonce() uint64
	Balance() BigInt
}

type ActorSubState interface{}

type ActorSubTypeID int

const (
	ActorSubTypeID_InitActor ActorSubTypeID = 0
)

type ActorSubType interface {
	ID() ActorSubTypeID
	Code() []byte
}

// InitActor

type ActorID uint64

type AddressMap interface {
	Lookup(addr Address) ActorID
	Update(addr Address, id ActorID) AddressMap
}

type ActorSubState_InitActor interface {
	AddressMap() AddressMap
	NextID() ActorID
}

// TODO: other Actor types
