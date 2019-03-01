# Filecoin Data Propagation

The filecoin network needs to broadcast blocks and messages to all peers in the network. This document details how that process works.

Both blocks and messages are propagated using the gossipsub libp2p pubsub router. The pubsub messages are authenticated. For blocks, the pubsub hop validation function is set to check that the block is valid before re-propagating. For messages, a similar validity check is run, the signature must be valid, and the account in question must have enough funds to cover the actions specified.

### Links

- [Gossipsub Spec](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub)
- [Block Validity Check](mining.md#chain-validation)
- TODO: Link to message validity check function

## Block Propagation

Blocks are propagated over the libp2p pubsub channel `/fil/blocks`. The block is [serialized](data-structures.md#block) and the raw bytes are sent as the content of the pubsub message. No messages, computed state, or other additional information is sent along with the block message.

Each Filecoin node sets a 'validation function' for the blocks topic that checks that the block is properly constructed, its ticket is valid, the block signature is valid, the miner is a valid miner, and the block is a child of a known good tipset. (TODO: clarify which of these checks are needed, any slowness here impacts propagation time significantly, this is not a full validity check) If an invalid block is received, the peer it was received from should be marked as potentially bad (TODO: we could blacklist peers who send bad blocks, maybe need support from libp2p for this?)

TODO: we should likely be smarter here and track which messages we could send along with each block to improve propagation time

## Message Propagation

Messages are propagated over the libp2p pubsub channel `/fil/messages`. The message is [serialized](data-structures.md#messages) and the raw bytes are sent as the content of the pubsub message.

The pubsub validation function for messages checks that the content of each pubsub message on this topic is, first, under the maximum size limit for a message, and then that it is a properly constructed message. (TODO: discuss checking signatures and account balances, some tricky bits that need consideration). If an invalid message is received from a peer, that peer should be marked as potentially bad.
