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
	commR, commC, commRLast, proof, cachedTreePath := sdr.Seal(sid, commD, data)

	return &SectorSealer_SealSector_FunRet_I{
		rawValue: &SealOutputs_I{
			SealInfo_: &sector.SealVerifyInfo_I{
				SectorID_: sid,
				OnChain_: &sector.OnChainSealVerifyInfo_I{
					SealedCID_:   commR,
					UnsealedCID_: commD,
					Proof_:       proof,
				},
			},
			ProofAuxTmp_: &sector.ProofAuxTmp_I{
				CommC_:                commC,
				CommRLast_:            commRLast,
				CachedMerkleTreePath_: cachedTreePath,
			},
		},
		which: SectorSealer_SealSector_FunRet_Case_so,
	}
}

func (s *SectorSealer_I) CreateSealProof(si CreateSealProofInputs) *SectorSealer_CreateSealProof_FunRet_I {
	sid := si.SectorID()
	commD := si.SealOutputs().SealInfo().OnChain().UnsealedCID()
	layers := si.SealOutputs().ProofAuxTmp().Layers()

	return &SectorSealer_CreateSealProof_FunRet_I{
		rawValue: CreateSealProof(sid, si.RandomSeed(), commD, layers[len(layers)-1]),
	}
}

func (s *SectorSealer_I) VerifySeal(sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {
	return &SectorSealer_VerifySeal_FunRet_I{}
}

func (s *SectorSealer_I) ComputeDataCommitment(unsealedPath file.Path) *SectorSealer_ComputeDataCommitment_FunRet_I {
	// TODO: Generate merkle tree using appropriate hash.
	return &SectorSealer_ComputeDataCommitment_FunRet_I{}
}

func ComputeReplicaID(sid sector.SectorID, commD sector.UnsealedSectorCID) *SectorSealer_ComputeReplicaID_FunRet_I {
	_, _ = sid.MinerID(), (sid.Number())

	// FIXME: Implement
	return &SectorSealer_ComputeReplicaID_FunRet_I{}
}

func UnsealedSectorCID(h filproofs.Blake2sHash) sector.UnsealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SealedSectorCID(h filproofs.PedersenHash) sector.SealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SDRParams() *filproofs.StackedDRG_I {
	return &filproofs.StackedDRG_I{}
}

func CreateSealProof(sid sector.SectorID, randomSeed sector.SealRandomSeed, commD sector.UnsealedSectorCID, replica Bytes) *CreateSealProofOutputs_I {
	var cachedMerkleTreePath file.Path // FIXME: get this

	commR, cachedMerkleTreePath := repHash(replica)

	var proof sector.SealProof

	return &CreateSealProofOutputs_I{
		SealInfo_: &sector.SealVerifyInfo_I{
			SectorID_: sid,
			OnChain_: &sector.OnChainSealVerifyInfo_I{
				SealedCID_:   SealedSectorCID(commR),
				UnsealedCID_: commD,
				RandomSeed_:  randomSeed,
				Proof_:       proof,
			},
		},
		ProofAux_: &sector.ProofAux_I{
			CommRLast_:            sector.Commitment{},
			CommC_:                sector.Commitment{},
			CachedMerkleTreePath_: cachedMerkleTreePath,
		},
	}
}

func repHash(data Bytes) (filproofs.PedersenHash, file.Path) {
	return Bytes{}, file.Path("") // FIXME
}
