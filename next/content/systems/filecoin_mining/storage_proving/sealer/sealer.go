package sealer

import "errors"

import abi "github.com/filecoin-project/specs-actors/actors/abi"
import util "github.com/filecoin-project/specs/util"
import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	registeredProof := si.RegisteredProof()
	sectorSize := si.SectorSize()
	sdr := filproofs.WinSDRParams(registeredProof)

	sid := si.SectorID()
	sectorSizeInt := int(sectorSize)

	unsealedPath := si.UnsealedPath()

	data := make(util.Bytes, sectorSize)
	in := file.FromPath(unsealedPath)
	inLength, err := in.Read(data)

	if err != nil {
		return SectorSealer_SealSector_FunRet_Make_err(err).Impl()
	}
	if inLength != sectorSizeInt {
		return SectorSealer_SealSector_FunRet_Make_err(
			errors.New("Sector file is wrong size"),
		).Impl()
	}

	sealArtifacts := sdr.Seal(si.RegisteredProof(), sid, data, si.RandomSeed())
	sealedPath := si.SealedPath()

	out := file.FromPath(sealedPath)
	outLength, err := out.Write(data)

	if err != nil {
		return SectorSealer_SealSector_FunRet_Make_err(err).Impl()
	}

	if outLength != sectorSizeInt {
		return SectorSealer_SealSector_FunRet_Make_err(
			errors.New("Wrote wrong sealed sector size"),
		).Impl()
	}

	return SectorSealer_SealSector_FunRet_Make_so(
		SectorSealer_SealSector_FunRet_so(
			&SealOutputs_I{
				ProofAuxTmp_: &sector.ProofAuxTmp_I{
					PersistentAux_: &sector.PersistentProofAux_I{
						CommC_:             sealArtifacts.CommC(),
						CommQ_:             sealArtifacts.CommQ(),
						CommRLast_:         sealArtifacts.CommRLast(),
						CommRLastTreePath_: sealArtifacts.CommRLastTreePath(),
					},
					CommD_:         sealArtifacts.CommD(),
					CommR_:         sealArtifacts.CommR(),
					CommDTreePath_: sealArtifacts.CommDTreePath(),
					CommCTreePath_: sealArtifacts.CommCTreePath(),
					CommQTreePath_: sealArtifacts.CommQTreePath(),
					Seed_:          sealArtifacts.Seed(),
					KeyLayers_:     sealArtifacts.KeyLayers(),
				}})).Impl()
}

func (s *SectorSealer_I) CreateSealProof(si CreateSealProofInputs) *SectorSealer_CreateSealProof_FunRet_I {
	sid := si.SectorID()
	randomSeed := si.InteractiveRandomSeed()
	auxTmp := si.SealOutputs().ProofAuxTmp()
	aux := auxTmp.PersistentAux()

	registeredProof := si.RegisteredProof()

	sdr := filproofs.WinSDRParams(registeredProof)
	proof := sdr.CreateSealProof(randomSeed, auxTmp)

	onChain := abi.OnChainSealVerifyInfo{
		SealedCID: auxTmp.CommR(),
		// Epoch_:  ? // TODO
		Proof: proof,
	}

	return SectorSealer_CreateSealProof_FunRet_Make_so(
		SectorSealer_CreateSealProof_FunRet_so(
			&CreateSealProofOutputs_I{
				SealInfo_: abi.SealVerifyInfo{
					SectorID: sid,
					OnChain:  onChain,
				},
				ProofAux_: aux,
			})).Impl()
}
