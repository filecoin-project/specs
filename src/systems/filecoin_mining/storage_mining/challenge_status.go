package storage_mining

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
)

func (cs *ChallengeStatus_I) OnNewChallenge(currEpoch block.ChainEpoch) ChallengeStatus {
	cs.LastChallengeEpoch_ = currEpoch
	return cs
}

// Call by _onSuccessfulPoSt or _onMissedSurprisePoSt
func (cs *ChallengeStatus_I) OnChallengeResponse(currEpoch block.ChainEpoch) ChallengeStatus {
	cs.LastChallengeEndEpoch_ = currEpoch
	return cs
}

func (cs *ChallengeStatus_I) IsChallenged() bool {
	// true (isChallenged) when LastChallengeEpoch is later than LastChallengeEndEpoch
	return cs.LastChallengeEpoch() > cs.LastChallengeEndEpoch()
}

func (cs *ChallengeStatus_I) ShouldChallenge(currEpoch block.ChainEpoch, minChallengePeriod block.ChainEpoch) bool {
	return currEpoch > (cs.LastChallengeEpoch()+minChallengePeriod) && !cs.IsChallenged()
}
