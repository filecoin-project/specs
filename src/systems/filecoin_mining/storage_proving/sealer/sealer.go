package sealer

import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import util "github.com/filecoin-project/specs/util"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	sid := si.SectorID()

	commD := sector.UnsealedSectorCID(s.ComputeDataCommitment(si.UnsealedPath()).As_commD())

	buf := make(util.Bytes, si.SealCfg().SectorSize())
	f := file.FromPath(si.SealedPath())
	length, _ := f.Read(buf) 

	// TODO: How do we meant to handle errors in implementation methods? This could get tedious fast.
	
	if util.UInt(length) != util.UInt(si.SealCfg().SectorSize()) {
		panic("Sector file is wrong size.")
	}
	
	return &SectorSealer_SealSector_FunRet_I{
		rawValue: Seal(sid, si.RandomSeed(), commD, buf),
	}
}

func (s *SectorSealer_I) VerifySeal(sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {
	return &SectorSealer_VerifySeal_FunRet_I{}
}

func (s *SectorSealer_I) ComputeDataCommitment(unsealedPath file.Path) *SectorSealer_ComputeDataCommitment_FunRet_I {
	return &SectorSealer_ComputeDataCommitment_FunRet_I{}
}

func ComputeReplicaID(sid sector.SectorID, commD sector.UnsealedSectorCID, seed sector.SealRandomSeed) *SectorSealer_ComputeReplicaID_FunRet_I {

	_, _ = sid.MinerID(), (sid.Number())

	return &SectorSealer_ComputeReplicaID_FunRet_I{}
}

// type SealOutputs struct {
//     SealInfo  sector.SealVerifyInfo
//     ProofAux  sector.ProofAux
// }

// type SealVerifyInfo struct {
//     SectorID
//     OnChain OnChainSealVerifyInfo
// }

func Seal(sid sector.SectorID, randomSeed sector.SealRandomSeed, commD sector.UnsealedSectorCID, data util.Bytes) *SealOutputs_I {
	replicaID := ComputeReplicaID(sid, commD, randomSeed).As_replicaID()
	_ = replicaID

	return &SealOutputs_I{}
}
