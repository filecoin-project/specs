package filproofs

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (fps *FilecoinProofsSubsystem_I) VerifySeal(sealVerifyInfo sector.SealVerifyInfo) bool {
	registeredProof := sealVerifyInfo.OnChain().RegisteredProof()
	sdr := WinSDRParams(registeredProof)
	return sdr.VerifySeal(sealVerifyInfo)
}
