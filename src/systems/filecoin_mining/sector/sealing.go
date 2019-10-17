package sector

// TODO: add blake2s_256 to filcrypto
// TODO: import blake2s_256 from filcrypto
var Blake2s_256 = func([]byte) []byte { return nil }

var SealSeedHash = Blake2s_256

func GenSealSeed(m MinerID, s SectorNumber, r SealRandomness, cid UnsealedSectorCID) SealSeed {
	var buf []byte
	// TODO: buf := m || s || r || cid
	h := SealSeedHash(buf)
	return SealSeed(h)
}
