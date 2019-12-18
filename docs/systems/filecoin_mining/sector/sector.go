package sector

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

// If a sector PreCommit appear at epoch T, then the corresponding ProveCommit
// must appear between epochs
//   (T + MIN_PROVE_COMMIT_SECTOR_EPOCH, T + MAX_PROVE_COMMIT_SECTOR_EPOCH)
// inclusive.
// TODO: placeholder epoch values -- will be set later
const MIN_PROVE_COMMIT_SECTOR_EPOCH = block.ChainEpoch(5)
const MAX_PROVE_COMMIT_SECTOR_EPOCH = block.ChainEpoch(10)

const (
	DeclaredFault StorageFaultType = 1 + iota
	DetectedFault
	TerminatedFault
)
