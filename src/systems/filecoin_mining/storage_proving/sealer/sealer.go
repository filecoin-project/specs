package sealer

import util "github.com/filecoin-project/specs/util"

import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	sdr := filproofs.SDRParams(si.SealCfg())
	sid := si.SectorID()

	data := make(util.Bytes, si.SealCfg().SectorSize())
	f := file.FromPath(si.SealedPath())
	length, _ := f.Read(data)

	if util.UInt(length) != util.UInt(si.SealCfg().SectorSize()) {
		return &SectorSealer_SealSector_FunRet_I{
			rawValue: "Sector file is wrong size",
			which:    SectorSealer_SealSector_FunRet_Case_err,
		}
	}

	sealArtifacts := sdr.Seal(sid, data, si.RandomSeed())

	return SectorSealer_SealSector_FunRet_Make_so(
		SectorSealer_SealSector_FunRet_so(
			&SealOutputs_I{
				ProofAuxTmp_: &sector.ProofAuxTmp_I{
					PersistentAux_: &sector.ProofAux_I{
						CommC_:             sealArtifacts.CommC(),
						CommRLast_:         sealArtifacts.CommRLast(),
						CommRLastTreePath_: sealArtifacts.CommRLastTreePath(),
					},
					CommD_:         sealArtifacts.CommD(),
					CommR_:         sealArtifacts.CommR(),
					CommDTreePath_: sealArtifacts.CommDTreePath(),
					Data_:          data,
					KeyLayers_:     sealArtifacts.KeyLayers(),
					Replica_:       sealArtifacts.Replica(),
				}})).Impl()
}

func (s *SectorSealer_I) CreateSealProof(si CreateSealProofInputs) *SectorSealer_CreateSealProof_FunRet_I {
	sid := si.SectorID()
	randomSeed := si.RandomSeed()
	auxTmp := si.SealOutputs().ProofAuxTmp()
	aux := auxTmp.PersistentAux()

	sdr := filproofs.SDRParams(si.SealCfg())
	proof := sdr.CreateSealProof(randomSeed, auxTmp)

	onChain := sector.OnChainSealVerifyInfo_I{
		SealedCID_: auxTmp.CommR(),
		// Epoch_:  ? // TODO
		Proof_: proof,
	}

	return SectorSealer_CreateSealProof_FunRet_Make_so(
		SectorSealer_CreateSealProof_FunRet_so(
			&CreateSealProofOutputs_I{
				SealInfo_: &sector.SealVerifyInfo_I{
					SectorID_: sid,
					OnChain_:  &onChain,
				},
				ProofAux_: aux,
			})).Impl()
}

func (s *SectorSealer_I) VerifySeal(sc sector.SealCfg, sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {

	sdr := filproofs.SDRParams(sc)
	result := sdr.VerifySeal(sc, sv)

	return SectorSealer_VerifySeal_FunRet_Make_ok(
		SectorSealer_VerifySeal_FunRet_ok(result)).Impl()
}
