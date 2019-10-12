package sealer

import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	sid := si.SectorID();
	
	commD := sector.UnsealedSectorCID(s.ComputeDataCommitment(si.UnsealedPath()).As_commD());
	replicaID := s.ComputeReplicaID(sid, commD, si.RandomSeed()).As_replicaID();
	
	_ = replicaID;
	return &SectorSealer_SealSector_FunRet_I {}
}

func (s *SectorSealer_I) VerifySeal(sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {
     return &SectorSealer_VerifySeal_FunRet_I {}
}

func (s *SectorSealer_I) ComputeDataCommitment(unsealedPath file.Path) *SectorSealer_ComputeDataCommitment_FunRet_I {
     return &SectorSealer_ComputeDataCommitment_FunRet_I {}
}

func (s *SectorSealer_I) ComputeReplicaID(sid sector.SectorID, commD sector.UnsealedSectorCID, seed sector.SealRandomSeed) *SectorSealer_ComputeReplicaID_FunRet_I {
	
	_, _ = sid.MinerID(), (sid.Number()); 

	return &SectorSealer_ComputeReplicaID_FunRet_I {}
}
