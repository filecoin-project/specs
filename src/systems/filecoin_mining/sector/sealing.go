package sector

import abi "github.com/filecoin-project/specs/actors/abi"

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

func (cfg *SealInstanceCfg_I) SectorSize() abi.SectorSize {
	switch cfg.Which() {
	case SealInstanceCfg_Case_WinStackedDRGCfgV1:
		{
			return cfg.As_WinStackedDRGCfgV1().SectorSize()
		}
	}
	panic("TODO")
}
