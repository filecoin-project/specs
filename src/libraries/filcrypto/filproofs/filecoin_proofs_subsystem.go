package filproofs

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (fps *FilecoinProofsSubsystem_I) VerifySeal(sealVerifyInfo sector.SealVerifyInfo) bool {
	sealCfg := sealVerifyInfo.SealCfg()
	cfg := SDRCfg_I{
		SealCfg_: sealCfg,
	}

	sdr := WinSDRParams(&cfg)
	return sdr.VerifySeal(sealVerifyInfo)
}
