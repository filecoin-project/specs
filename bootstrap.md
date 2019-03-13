# Filecoin Bootstrapping Routine

This spec describes how to implement the protocol in general, for Filecoin-specific processes, see:

- [Mining Blocks](mining.md#mining-blocks) on how consensus is used.
- [Faults](./faults.md) on slashing.
- [Storage Market](./storage-market.md#the-power-table) on how the power table is maintained.
- [Block data structure](data-structures.md#block) for details on fields and encoding. 

Expected Consensus is a probabilistic Byzantine fault-tolerant consensus protocol. At a high
level, it operates by running a leader election every round in which, on expectation, one 
participant may be eligible to submit a block. All valid blocks submitted in a given round form a `TipSet`. Every block in a TipSet adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on this, see [Chain Selection](#chain-selection).

The basic algorithm can be broken down by looking in turn at its two major components: 

- Leader Election
- Chain Selection

## Connect

If a node has never connected to the network before, They will need to connect to initial set of trusted peers. This will generally be a list of bootstrap peers from the config file, though each individual could come up with a set of peers they explicitly trust.

On successive connects, we can augment the ‘initial connect set’ with randomly selected peers remembered from the previous session.

## Peer Set Expansion

Once connected into some portion of the network, a node will need to expand their peer set in order to get a good sample of the network. This can be accomplished many different ways, the simplest for us will be to use the DHT to do random ‘FindPeer’ requests.

Once the node is connected to a sufficiently diverse set of peers, they should then move on to chain selection.

## Chain Selection

Sample what each of your peers thinks of as the latest chain head. If there are no significant disagreements (an insignificant disagreement would be if two peers have different heads, but one is the parent of the other, or some other similar ancestry chain) then that chain should be given to the syncer routine.

If there are disagreements, then the node should watch both chains for some time, syncing enough data to validate them moving forward, but not historically yet. After a ‘trial period’ has elapsed, the node should select the chain that has included valid proofs for the most storage. The other chain(s) should be marked as bad. The selected chain should be passed to the syncer, and the ‘initial sync’ should begin.