package storage_mining

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
)

// TODO: placeholder epoch value -- this will be set later
const MAX_PROVE_COMMIT_SECTOR_PERIOD = block.ChainEpoch(3)    // placeholder
const MAX_SURPRISE_POST_RESPONSE_PERIOD = block.ChainEpoch(4) // placeholder
const POST_CHALLENGE_TIME = block.ChainEpoch(1)               // placeholder
// sets the average frequency
const PROVING_PERIOD = block.ChainEpoch(2) // placeholder, 2 days
// how long after a POST before a miner can get challenged again
const SUPRISE_NO_CHALLENGE_PERIOD = block.ChainEpoch(0) // placeholder, 2 hours
// how long miner has to respond to the challenge before it expires
const CHALLENGE_DURATION = block.ChainEpoch(0) // placeholder, 2 hours
// number of detected faults before a miner's sectors are all terminated
const MAX_CONSECUTIVE_FAULTS = 3

const EPOST_SAMPLE_NUM = 1
const EPOST_SAMPLE_DENOM = 25

// FIL deposit per sector precommit in Interactive PoRep
// refunded after ProveCommit but burned if PreCommit expires
const PRECOMMIT_DEPOSIT_PER_BYTE = actor.TokenAmount(0) // placeholder
