package sealer

import . "github.com/filecoin-project/specs/util"

import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	sid := si.SectorID()
	commD := sector.UnsealedSectorCID(s.ComputeDataCommitment(si.UnsealedPath()).As_commD())

	data := make(Bytes, si.SealCfg().SectorSize())
	f := file.FromPath(si.SealedPath())
	length, _ := f.Read(data)

	if UInt(length) != UInt(si.SealCfg().SectorSize()) {
		return &SectorSealer_SealSector_FunRet_I{
			rawValue: "Sector file is wrong size",
			which:    SectorSealer_SealSector_FunRet_Case_err,
		}
	}

	sdr := filproofs.SDRParams()
	sealOutputs := sdr.Seal(sid, commD, data)

	return &SectorSealer_SealSector_FunRet_I{
		rawValue: &SealOutputs_I{
			SealInfo_: &sector.SealVerifyInfo_I{
				SectorID_: sid,
				OnChain_: &sector.OnChainSealVerifyInfo_I{
					SealedCID_: sealOutputs.CommR,
					Proof_:     sealOutputs.Proof,
				},
			},
			ProofAuxTmp_: &sector.ProofAuxTmp_I{
				PersistentAux_: &sector.ProofAux_I{
					CommC_:                sealOutputs.CommC,
					CommRLast_:            sealOutputs.CommRLast,
					CachedMerkleTreePath_: sealOutputs.TreePath,
				},
				CommD_: commD,
			},
		},
		which: SectorSealer_SealSector_FunRet_Case_so,
	}
}

func (s *SectorSealer_I) CreateSealProof(si CreateSealProofInputs) *SectorSealer_CreateSealProof_FunRet_I {
	sid := si.SectorID()
	randomSeed := si.SealOutputs().SealInfo().OnChain().RandomSeed()
	layers := si.SealOutputs().ProofAuxTmp().Layers()

	sdr := filproofs.SDRParams()
	sealedCID, cachedMerkleTreePath := sdr.CreateSealProof(layers)

	var proof sector.SealProof

	return &SectorSealer_CreateSealProof_FunRet_I{
		rawValue: &CreateSealProofOutputs_I{
			SealInfo_: &sector.SealVerifyInfo_I{
				SectorID_: sid,
				OnChain_: &sector.OnChainSealVerifyInfo_I{
					SealedCID_:  sealedCID,
					RandomSeed_: randomSeed,
					Proof_:      proof,
				},
			},
			ProofAux_: &sector.ProofAux_I{
				CommRLast_:            sector.Commitment{},
				CommC_:                sector.Commitment{},
				CachedMerkleTreePath_: cachedMerkleTreePath,
			},
		},
	}
}

func (s *SectorSealer_I) VerifySeal(sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {
	return &SectorSealer_VerifySeal_FunRet_I{}
}

func (s *SectorSealer_I) ComputeDataCommitment(unsealedPath file.Path) *SectorSealer_ComputeDataCommitment_FunRet_I {
	// TODO: Generate merkle tree using appropriate hash.
	return &SectorSealer_ComputeDataCommitment_FunRet_I{}
}
