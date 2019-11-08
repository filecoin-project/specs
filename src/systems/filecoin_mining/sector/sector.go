package sector

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

func (r *FaultReport_I) GetDeclaredFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (r *FaultReport_I) GetDetectedFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (r *FaultReport_I) GetTerminatedFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}
