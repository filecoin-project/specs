package filproofs

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (fps *FilecoinProofsSubsystem_I) VerifySeal(sealVerifyInfo sector.SealVerifyInfo) bool {
	cfg := sealVerifyInfo.SealCfg()
	sdr := WinSDRParams(cfg)
	return sdr.VerifySeal(sealVerifyInfo)
}
