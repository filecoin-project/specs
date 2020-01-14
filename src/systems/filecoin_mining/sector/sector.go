package sector

import (
	util "github.com/filecoin-project/specs/util"
)

var IMPL_FINISH = util.IMPL_FINISH

type Serialization = util.Serialization

const (
	DeclaredFault StorageFaultType = 1 + iota
	DetectedFault
	TerminatedFault
)

var PROOFS ProofRegistry = ProofRegistry(map[util.UInt]ProofInstance{util.UInt(RegisteredProof_WinStackedDRG32GiBSeal): &ProofInstance_I{
	ID_:        RegisteredProof_WinStackedDRG32GiBSeal,
	Type_:      ProofType_SealProof,
	Algorithm_: ProofAlgorithm_WinStackedDRGSeal,
	CircuitType_: &ConcreteCircuit_I{
		Name_: "HASHOFCIRCUITDEFINITION1",
	},
},
	util.UInt(RegisteredProof_WinStackedDRG32GiBPoSt): &ProofInstance_I{
		ID_:        RegisteredProof_WinStackedDRG32GiBPoSt,
		Type_:      ProofType_PoStProof,
		Algorithm_: ProofAlgorithm_WinStackedDRGPoSt,
		CircuitType_: &ConcreteCircuit_I{
			Name_: "HASHOFCIRCUITDEFINITION2",
		},
	},
	util.UInt(RegisteredProof_StackedDRG32GiBSeal): &ProofInstance_I{
		ID_:        RegisteredProof_StackedDRG32GiBSeal,
		Type_:      ProofType_SealProof,
		Algorithm_: ProofAlgorithm_StackedDRGSeal,
		CircuitType_: &ConcreteCircuit_I{
			Name_: "HASHOFCIRCUITDEFINITION3",
		},
	},
	util.UInt(RegisteredProof_StackedDRG32GiBPoSt): &ProofInstance_I{
		ID_:        RegisteredProof_StackedDRG32GiBPoSt,
		Type_:      ProofType_PoStProof,
		Algorithm_: ProofAlgorithm_StackedDRGPoSt,
		CircuitType_: &ConcreteCircuit_I{
			Name_: "HASHOFCIRCUITDEFINITION4",
		},
	},
})

func RegisteredProofInstance(r RegisteredProof) ProofInstance {
	return PROOFS[util.UInt(r)]
}

func (c *ConcreteCircuit_I) GrothParameterFileName() string {
	return c.Name() + ".params"
}

func (c *ConcreteCircuit_I) VerifyingKeyFileName() string {
	return c.Name() + ".vk"
}
