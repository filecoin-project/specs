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
},
	util.UInt(RegisteredProof_WinStackedDRG32GiBPoSt): &ProofInstance_I{
		ID_:        RegisteredProof_WinStackedDRG32GiBPoSt,
		Type_:      ProofType_PoStProof,
		Algorithm_: ProofAlgorithm_WinStackedDRGPoSt,
	},
	util.UInt(RegisteredProof_StackedDRG32GiBSeal): &ProofInstance_I{
		ID_:        RegisteredProof_StackedDRG32GiBSeal,
		Type_:      ProofType_SealProof,
		Algorithm_: ProofAlgorithm_StackedDRGSeal,
	},
	util.UInt(RegisteredProof_StackedDRG32GiBPoSt): &ProofInstance_I{
		ID_:        RegisteredProof_StackedDRG32GiBPoSt,
		Type_:      ProofType_PoStProof,
		Algorithm_: ProofAlgorithm_StackedDRGPoSt,
	},
})
