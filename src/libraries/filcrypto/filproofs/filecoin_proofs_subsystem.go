package filproofs

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (fps *FilecoinProofsSubsystem_I) VerifySeal(sealVerifyInfo sector.SealVerifyInfo) bool {
	cfg := sector.RegisteredProofInstance(sealVerifyInfo.RegisteredProof()).Cfg().As_SealCfg()
	sdr := WinSDRParams(cfg)
	return sdr.VerifySeal(sealVerifyInfo)
}
