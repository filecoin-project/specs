package sector

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

// TODO: placeholder epoch value -- this will be set later
// ProveCommitSector needs to be submitted within MAX_PROVE_COMMIT_SECTOR_EPOCH after PreCommit
const MAX_PROVE_COMMIT_SECTOR_EPOCH = block.ChainEpoch(3)

const (
	DeclaredFault StorageFaultType = 1 + iota
	DetectedFault
	TerminatedFault
)
