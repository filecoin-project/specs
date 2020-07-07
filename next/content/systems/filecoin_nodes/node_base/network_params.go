package node_base

import (
	addr "github.com/filecoin-project/go-address"
)

// Parameters for on-chain calculations are in actors/builtin/network_params.go

/////////////////////////////////////////////////////////////
// Global
/////////////////////////////////////////////////////////////

const NETWORK = addr.Testnet

// how many sectors should be challenged in surprise post (if miner has fewer, will get dup challenges)
const SURPRISE_CHALLENGE_COUNT = 200 // placeholder

const EPOST_SAMPLE_RATE_NUM = 1    // placeholder
const EPOST_SAMPLE_RATE_DENOM = 25 // placeholder
const SPOST_SAMPLE_RATE_NUM = 1    // placeholder
const SPOST_SAMPLE_RATE_DENOM = 50 // placeholder

/////////////////////////////////////////////////////////////
// Consensus
/////////////////////////////////////////////////////////////

const FINALITY = 500          // placeholder
const SPC_LOOKBACK_TICKET = 1 // we chain blocks together one after the other
