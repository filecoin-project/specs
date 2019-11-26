package sector

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

// TODO: placeholder epoch value -- this will be set later
// ProveCommitSector needs to be submitted within MAX_PROVE_COMMIT_SECTOR_EPOCH after PreCommit
const MAX_PROVE_COMMIT_SECTOR_EPOCH = block.ChainEpoch(3)

type StorageFaultType = int

const (
	DeclaredFault   StorageFaultType = 0
	DetectedFault   StorageFaultType = 1
	TerminatedFault StorageFaultType = 2
)
