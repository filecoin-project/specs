package codeGen

type StateTree interface {
	GetActorState() ActorState
	UpdateActorState(a ActorState) StateTree

	CID() CID
}
