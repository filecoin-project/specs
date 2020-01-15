package actors

import (
	cid "github.com/ipfs/go-cid"
)

// CallSeqNum is an invocation (Call) sequence (Seq) number (Num).
// This is a value used for securing against replay attacks:
// each AccountActor (user) invocation must have a unique CallSeqNum
// value. The sequenctiality of the numbers is used to make it
// easy to verify, and to order messages.
//
// Q&A
// - > Does it have to be sequential?
//   No, a random nonce could work against replay attacks, but
//   making it sequential makes it much easier to verify.
// - > Can it be used to order events?
//   Yes, a user may submit N separate messages with increasing
//   sequence number, causing them to execute in order.
//
type CallSeqNum int64

type ActorSubstateCID cid.Cid

func (x ActorSubstateCID) Ref() ActorSubstateCID {
	return x
}
