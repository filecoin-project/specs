package block

import (
	filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
	util "github.com/filecoin-project/specs/util"
)

func SmallerBytes(a, b util.Bytes) util.Bytes {
	if util.CompareBytesStrict(a, b) > 0 {
		return b
	}
	return a
}

func (block *BlockHeader_I) ExtractElectionSeed() ElectionSeed {
	panic("TODO")
	// var ret []byte

	// for _, currBlock := range lookbackTipset.Blocks() {
	// 	for _, currTicket := range currBlock.Tickets() {

	// 		currSeed := Hash(
	// 			HashRole_ElectionSeedFromVRFOutput,
	// 			currTicket.VRFResult().bytes(),
	// 		)
	// 		if ret == nil {
	// 			ret = currSeed
	// 		} else {
	//             ret = SmallerBytes(currSeed, ret)
	//         }
	// 	}
	// }

	// Assert(ret != nil)
	// return ElectionSeed.FromBytesInternal(nil, ret)
}

// func GenerateElectionTicket(k VRFKeyPair, seed ElectionSeed) Ticket {
// 	var vrfResult VRFResult = VRFEval(k, seed.ToBytesInternal())

// 	var vdfInput []byte = Hash(
// 		HashRole_TicketVDFInputFromVRFOutput,
// 		vrfResult.ToBytesInternal(),
// 	)
// 	var vdfResult VDFResult = VDFEval(vdfInput)

// 	return &TicketI{
// 		vrfResult,
// 		vdfResult,
// 	}
// }

func (self *BlockHeader_I) ValidateTickets(pubKey filcrypto.PubKey) bool {
	for _, tix := range self.Tickets_ {
		panic("TODO")
		panic(tix)
		// tix.Validate()
	}

	return true
}
