package sealer

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
     return &SectorSealer_SealSector_FunRet_I {}
}

func (s *SectorSealer_I) VerifySeal(sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {
     return &SectorSealer_VerifySeal_FunRet_I {}
}

