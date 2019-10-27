package filproofs

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (fps *FilecoinProofsSubsystem_I) VerifySeal(sealVerifyInfo sector.SealVerifyInfo) bool {
	sealCfg := sealVerifyInfo.SealCfg()

	sdr := SDRParams(sealCfg, nil)
	return sdr.VerifySeal(sealVerifyInfo)
}
