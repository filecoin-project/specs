package sector

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

// TODO: placeholder epoch value -- this will be set later
// ProveCommitSector needs to be submitted within MAX_PROVE_COMMIT_SECTOR_EPOCH after PreCommit
const MAX_PROVE_COMMIT_SECTOR_EPOCH = block.ChainEpoch(3)

func (r *FaultReport_I) GetDeclaredFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (r *FaultReport_I) GetDetectedFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (r *FaultReport_I) GetTerminatedFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

type StorageDealPaymentAction = int

const (
	ExpireStorageDeals StorageDealPaymentAction = 0
	CreditStorageDeals StorageDealPaymentAction = 1
)

type StorageDealSlashAction = int

const (
	SlashDeclaredFaults   StorageDealSlashAction = 0
	SlashDetectedFaults   StorageDealSlashAction = 1
	SlashTerminatedFaults StorageDealSlashAction = 2
)
