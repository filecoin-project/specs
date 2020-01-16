package builtin

import abi "github.com/filecoin-project/specs/actors/abi"

/////////////////////////////////////////////////////////////
// PoSt
/////////////////////////////////////////////////////////////

// how long miner has to respond to the challenge before it expires
const CHALLENGE_DURATION = abi.ChainEpoch(4) // placeholder, 2 hours
// sets the average frequency
const PROVING_PERIOD = abi.ChainEpoch(2) // placeholder, 2 days
// how long after a POST challenge before a miner can get challenged again
const SURPRISE_NO_CHALLENGE_PERIOD = abi.ChainEpoch(0) // placeholder, 2 hours
// number of detected faults before a miner's sectors are all terminated
const MAX_CONSECUTIVE_FAULTS = 3

// Time between when a temporary sector fault is declared, and when it becomes
// effective for purposes of reducing the active proving set for PoSts.
const DECLARED_FAULT_EFFECTIVE_DELAY = abi.ChainEpoch(20) // placeholder

/////////////////////////////////////////////////////////////
// Storage mining
/////////////////////////////////////////////////////////////

// If a sector PreCommit appear at epoch T, then the corresponding ProveCommit
// must appear between epochs
//   (T + MIN_PROVE_COMMIT_SECTOR_EPOCH, T + MAX_PROVE_COMMIT_SECTOR_EPOCH)
// inclusive.
// TODO: placeholder epoch values -- will be set later
const MIN_PROVE_COMMIT_SECTOR_EPOCH = abi.ChainEpoch(5)
const MAX_PROVE_COMMIT_SECTOR_EPOCH = abi.ChainEpoch(10)

const SPC_LOOKBACK_POST = 1   // cheap to generate, should be set as close to current TS as possible
const SPC_LOOKBACK_SEAL = 500 // should be approximately the same as finality

/////////////////////////////////////////////////////////////
// Storage power
/////////////////////////////////////////////////////////////

const MIN_MINER_SIZE_STOR = 1 << 40 // placeholder, 100 TB
const MIN_MINER_SIZE_TARG = 3       // placeholder

/////////////////////////////////////////////////////////////
// Faults and slashing
/////////////////////////////////////////////////////////////

// FIL deposit per sector precommit in Interactive PoRep
// refunded after ProveCommit but burned if PreCommit expires
const PRECOMMIT_DEPOSIT_PER_BYTE = abi.TokenAmount(0) // placeholder
const FAULT_SLASH_PERC_DECLARED = 1                   // placeholder
const FAULT_SLASH_PERC_DETECTED = 10                  // placeholder
const FAULT_SLASH_PERC_TERMINATED = 100               // placeholder

const SLASHER_INITIAL_SHARE_NUM = 1            // placeholder
const SLASHER_INITIAL_SHARE_DENOM = 1000       // placeholder
const SLASHER_SHARE_GROWTH_RATE_NUM = 102813   // placeholder
const SLASHER_SHARE_GROWTH_RATE_DENOM = 100000 // placeholder
