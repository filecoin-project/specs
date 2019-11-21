package storage_mining

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

// TODO: placeholder epoch value -- this will be set later
const MAX_PROVE_COMMIT_SECTOR_PERIOD = block.ChainEpoch(3)
const CHALLENGE_CLEANUP_PERIOD = block.ChainEpoch(24 * 60 * 4) // one day
const MAX_SURPRISE_POST_RESPONSE_PERIOD = block.ChainEpoch(4)

const ELECTION_PERIOD_DURATION = block.ChainEpoch(1)
const CLEANUP_PERIOD_DURATION = block.ChainEpoch(1)
