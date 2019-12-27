package chain

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

type Randomness = util.Randomness
type Serialization = util.Serialization
type DomainSeparationTag int

const (
	DomainSeparationTag_TicketDrawing DomainSeparationTag = 1 + iota
	DomainSeparationTag_TicketProduction
	DomainSeparationTag_PoStChallengeSeed
	// SurprisePoStSelectMiners
	// SurprisePoStSampleSectors
	// SurprisePoStVRFRandomnessInput
)

// Returns the tipset at or immediately prior to `epoch`.
func (chain *Chain_I) TipsetAtEpoch(epoch abi.ChainEpoch) Tipset {
	current := chain.HeadTipset()
	for current.Epoch() > epoch {
		current = current.Parents()
	}

	return current
}

// Draws randomness from the tipset at or immediately prior to `epoch`.
func (chain *Chain_I) RandomnessAtEpoch(epoch abi.ChainEpoch) util.Bytes {
	ts := chain.TipsetAtEpoch(epoch)
	return ts.MinTicket().DrawRandomness(epoch)
}

func (chain *Chain_I) PreparePoStChallengeSeed(randomness util.Randomness, minerAddr addr.Address) util.Randomness {

	randInput := Serialize_PoStChallengeSeedInput(&PoStChallengeSeedInput_I{
		ticket_:    randomness,
		minerAddr_: minerAddr,
	})
	input := DomainSeparationTag_PoStChallengeSeed.DeriveRand(randInput)
	return input
}

func (chain *Chain_I) GetTicketProductionRand(epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_TICKET)
}

func (chain *Chain_I) GetSealRand(epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_SEAL)
}

func (chain *Chain_I) GetPoStChallengeRand(epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_POST)
}

// Derive a random byte string from a domain separation tag and an arbitrary
// serializable object.
//
// Note: to produce values of type Serialization, use the auto-generated method prototypes
//   Serialize_T(T) Serialization
// for each type T defined in .id.
//
// In order to derive randomness from a collection of objects, rather than just a single
// object, define a struct at the .id level that contains those objects as member fields.
// This will then cause a Serialize_*() method to be generated for the struct type.
func (tag DomainSeparationTag) DeriveRand(s Serialization) Randomness {
	return _deriveRandInternal(tag, s, -1)
}

func _deriveRandInternal(tag DomainSeparationTag, s Serialization, index int) Randomness {
	buffer := []byte{}
	buffer = append(buffer, filcrypto.LittleEndianBytesFromInt(int(tag))...)
	buffer = append(buffer, filcrypto.LittleEndianBytesFromInt(int(index))...)
	buffer = append(buffer, util.Bytes(s)...)
	ret := filcrypto.SHA256(buffer)
	return Randomness(ret)
}
