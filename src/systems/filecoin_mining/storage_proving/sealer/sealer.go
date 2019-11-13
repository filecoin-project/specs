package sealer

import "errors"

import util "github.com/filecoin-project/specs/util"
import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	cfg := &filproofs.SDRCfg_I{
		SealCfg_: si.SealCfg(),
	}
	sdr := filproofs.SDRParams(cfg)

	sid := si.SectorID()
	subsectorCount := int(si.SealCfg().SubsectorCount())
	sectorSize := int(si.SealCfg().SectorSize())
	subsectorSize := sectorSize / subsectorCount

	if len(si.UnsealedPaths()) != subsectorCount {
		return SectorSealer_SealSector_FunRet_Make_err(
			errors.New("Wrong number of subsector files."),
		).Impl()
	}

	var subsectorData [][]byte
	for _, unsealedPath := range si.UnsealedPaths() {
		data := make(util.Bytes, si.SealCfg().SectorSize())
		in := file.FromPath(unsealedPath)
		length, err := in.Read(data)

		if err != nil {
			return SectorSealer_SealSector_FunRet_Make_err(err).Impl()
		}

		subsectorData = append(subsectorData, data)

		if length != subsectorSize {
			return SectorSealer_SealSector_FunRet_Make_err(
				errors.New("Subsector file is wrong size"),
			).Impl()
		}
	}

	sealArtifacts := sdr.Seal(sid, subsectorData, si.RandomSeed())
	sealedPaths := si.SealedPaths()

	for i, data := range subsectorData {
		out := file.FromPath(sealedPaths[i])
		length, err := out.Write(data)

		if err != nil {
			return SectorSealer_SealSector_FunRet_Make_err(err).Impl()
		}

		if length != subsectorSize {
			return SectorSealer_SealSector_FunRet_Make_err(
				errors.New("Wrote wrong sealed subsector size"),
			).Impl()
		}

	}

	return SectorSealer_SealSector_FunRet_Make_so(
		SectorSealer_SealSector_FunRet_so(
			&SealOutputs_I{
				ProofAuxTmp_: &sector.ProofAuxTmp_I{
					PersistentAux_: &sector.ProofAux_I{
						CommC_:             sealArtifacts.CommC(),
						CommRLast_:         sealArtifacts.CommRLast(),
						CommRLastTreePath_: sealArtifacts.CommRLastTreePath(),
					},
					CommD_:          sealArtifacts.CommD(),
					CommR_:          sealArtifacts.CommR(),
					CommDTreePaths_: sealArtifacts.CommDTreePaths(),
					CommCTreePath_:  sealArtifacts.CommCTreePath(),
					Seeds_:          sealArtifacts.Seeds(),
					SubsectorData_:  subsectorData,
					KeyLayers_:      sealArtifacts.KeyLayers(),
					Replicas_:       sealArtifacts.Replicas(),
				}})).Impl()
}

func (s *SectorSealer_I) CreateSealProof(si CreateSealProofInputs) *SectorSealer_CreateSealProof_FunRet_I {
	sid := si.SectorID()
	randomSeed := si.InteractiveRandomSeed()
	auxTmp := si.SealOutputs().ProofAuxTmp()
	aux := auxTmp.PersistentAux()

	cfg := &filproofs.SDRCfg_I{
		SealCfg_: si.SealCfg(),
	}

	sdr := filproofs.SDRParams(cfg)
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
