package sector

import util "github.com/filecoin-project/specs/util"

// NOTE: It's fairly unclear how any of this should interface/cooperate with filcrypto/filproofs.
// Leaving now to preserve some historical intent for later refactoring.

// TODO: add SHA256 to filcrypto
// TODO: import SHA256 from filcrypto
var SHA256 = func([]byte) []byte { return nil }

var SealSeedHash = SHA256

// This is superseded (heh) by fliproofs.computeSealSeed. Should it live here?
// func GenSealSeed(m MinerID, s SectorNumber, r SealRandomness, cid UnsealedSectorCID) SealSeed {
// 	var buf []byte
// 	// TODO: buf := m || s || r || cid
// 	h := SealSeedHash(buf)
// 	return SealSeed(h)
// }

func (x PieceInfo_I) Ref() *PieceInfo_I {
	return &x
}

func (svi *OnChainSealVerifyInfo_I) IsValidAtSealEpoch() bool {
	// We can just hardcode logic for the range of epochs at which each circuit type is valid.
	switch PROOFS[util.UInt(svi.RegisteredProof())].CircuitType() {
	}
	panic("TODO")
}

func (cfg *SealInstanceCfg_I) SectorSize() SectorSize {
	switch cfg.Which() {
	case SealInstanceCfg_Case_WinStackedDRGCfgV1:
		{
			return cfg.As_WinStackedDRGCfgV1().SectorSize()
		}
	}
	panic("TODO")
}
