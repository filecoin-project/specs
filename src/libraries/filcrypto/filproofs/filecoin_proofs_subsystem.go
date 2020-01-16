package filproofs

import abi "github.com/filecoin-project/specs/actors/abi"

func (fps *FilecoinProofsSubsystem_I) VerifySeal(sealVerifyInfo abi.SealVerifyInfo) bool {
	registeredProof := sealVerifyInfo.OnChain.RegisteredProof
	sdr := WinSDRParams(registeredProof)
	return sdr.VerifySeal(sealVerifyInfo)
}
