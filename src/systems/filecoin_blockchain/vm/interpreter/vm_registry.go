package interpreter

import "errors"
import actor "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/actor"
import vmr "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/runtime"

var (
	ErrActorNotFound = errors.New("Actor Not Found")
)

var staticActorCodeRegistry = actorCodeRegistry{}

// init is called in Go during initialization of a program.
// this is an idiomatic way to do this. Implementations should approach this
// howevery they wish. The point is to initialize a static registry with
// built in pure types that have the code for each actor. Once we have
// a way to load code from the StateTree inst
func init() {
	staticRegisterBuiltinActors()
}

func staticRegisterBuiltinActors() {
	// TODO
	// registerActor(InitActorCodeCID, new sysactors.InitActorCode())
}

// ActorCode is the interface that all actor code types should satisfy.
// It is merely a method dispatch interface.
type ActorCode interface {
	CodeCID() actor.CodeCID
	InvokeMethod(input vmr.InvocInput, method actor.MethodNum, params actor.MethodParams) vmr.InvocOutput
}

type actorCodeRegistry struct {
	code map[actor.CodeCID]ActorCode
}

func (r *actorCodeRegistry) RegisterActor(cid actor.CodeCID, actor ActorCode) {
	r.code[cid] = actor
}

func (r *actorCodeRegistry) LoadActor(cid actor.CodeCID) (ActorCode, error) {
	a, ok := r.code[cid]
	if !ok {
		return nil, ErrActorNotFound
	}
	return a, nil
}
