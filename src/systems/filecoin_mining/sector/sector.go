package sector

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	util "github.com/filecoin-project/specs/util"
)

var IMPL_FINISH = util.IMPL_FINISH

type Serialization = util.Serialization

const (
	DeclaredFault StorageFaultType = 1 + iota
	DetectedFault
	TerminatedFault
)

var PROOFS ProofRegistry = ProofRegistry(map[abi.RegisteredProof]ProofInstance{abi.RegisteredProof_WinStackedDRG32GiBSeal: &ProofInstance_I{
	ID_:        abi.RegisteredProof_WinStackedDRG32GiBSeal,
	Type_:      ProofType_SealProof,
	Algorithm_: abi.ProofAlgorithm_WinStackedDRGSeal,
	CircuitType_: &ConcreteCircuit_I{
		Name_: "HASHOFCIRCUITDEFINITION1",
	},
},
	abi.RegisteredProof_WinStackedDRG32GiBPoSt: &ProofInstance_I{
		ID_:        abi.RegisteredProof_WinStackedDRG32GiBPoSt,
		Type_:      ProofType_PoStProof,
		Algorithm_: abi.ProofAlgorithm_WinStackedDRGPoSt,
		CircuitType_: &ConcreteCircuit_I{
			Name_: "HASHOFCIRCUITDEFINITION2",
		},
	},
	abi.RegisteredProof_StackedDRG32GiBSeal: &ProofInstance_I{
		ID_:        abi.RegisteredProof_StackedDRG32GiBSeal,
		Type_:      ProofType_SealProof,
		Algorithm_: abi.ProofAlgorithm_StackedDRGSeal,
		CircuitType_: &ConcreteCircuit_I{
			Name_: "HASHOFCIRCUITDEFINITION3",
		},
	},
	abi.RegisteredProof_StackedDRG32GiBPoSt: &ProofInstance_I{
		ID_:        abi.RegisteredProof_StackedDRG32GiBPoSt,
		Type_:      ProofType_PoStProof,
		Algorithm_: abi.ProofAlgorithm_StackedDRGPoSt,
		CircuitType_: &ConcreteCircuit_I{
			Name_: "HASHOFCIRCUITDEFINITION4",
		},
	},
})

func RegisteredProofInstance(r abi.RegisteredProof) ProofInstance {
	return PROOFS[r]
}

func (c *ConcreteCircuit_I) GrothParameterFileName() string {
	return c.Name() + ".params"
}

func (c *ConcreteCircuit_I) VerifyingKeyFileName() string {
	return c.Name() + ".vk"
}
