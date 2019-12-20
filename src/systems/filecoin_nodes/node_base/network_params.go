package node_base

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
)

/////////////////////////////////////////////////////////////
// Commitment
/////////////////////////////////////////////////////////////

// TODO: placeholder epoch value -- this will be set later
const MAX_PROVE_COMMIT_SECTOR_PERIOD = block.ChainEpoch(3) // placeholder

/////////////////////////////////////////////////////////////
// PoSt
/////////////////////////////////////////////////////////////

const MAX_SURPRISE_POST_RESPONSE_PERIOD = block.ChainEpoch(4) // placeholder
const POST_CHALLENGE_TIME = block.ChainEpoch(1)               // placeholder
// sets the average frequency
const PROVING_PERIOD = block.ChainEpoch(2) // placeholder, 2 days
// how long after a POST challenge before a miner can get challenged again
const SUPRISE_NO_CHALLENGE_PERIOD = block.ChainEpoch(0) // placeholder, 2 hours
// how many sectors should be challenged in surprise post (if miner has fewer, will get dup challenges)
const SURPRISE_CHALLENGE_COUNT = 200 // placeholder
// how long miner has to respond to the challenge before it expires
const CHALLENGE_DURATION = block.ChainEpoch(0) // placeholder, 2 hours
// number of detected faults before a miner's sectors are all terminated
const MAX_CONSECUTIVE_FAULTS = 3
const EPOST_SAMPLE_RATE_NUM = 1    // placeholder
const EPOST_SAMPLE_RATE_DENOM = 25 // placeholder
const SPOST_SAMPLE_RATE_NUM = 1    // placeholder
const SPOST_SAMPLE_RATE_DENOM = 50 // placeholder

/////////////////////////////////////////////////////////////
// Consensus
/////////////////////////////////////////////////////////////

const MIN_MINER_SIZE_STOR = 1 << 40 // placeholder, 100 TB
const MIN_MINER_SIZE_TARG = 3       // placeholder
const FINALITY = 500                // placeholder
const SPC_LOOKBACK_RANDOMNESS = 300 // this is EC.K maybe move it there. TODO
const SPC_LOOKBACK_TICKET = 1       // we chain blocks together one after the other
const SPC_LOOKBACK_POST = 1         // cheap to generate, should be set as close to current TS as possible
const SPC_LOOKBACK_SEAL = FINALITY  // should be set to finality

/////////////////////////////////////////////////////////////
// Cryptoecon
/////////////////////////////////////////////////////////////

// FIL deposit per sector precommit in Interactive PoRep
// refunded after ProveCommit but burned if PreCommit expires
const PRECOMMIT_DEPOSIT_PER_BYTE = actor.TokenAmount(0) // placeholder
const FAULT_SLASH_PERC_DECLARED = 1                     // placeholder
const FAULT_SLASH_PERC_DETECTED = 10                    // placeholder
const FAULT_SLASH_PERC_TERMINATED = 100                 // placeholder
