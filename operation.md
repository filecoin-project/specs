# Filecoin Node Operation

Running a Filecoin `full node` requires running many different processes and protocols simultaneously. This section describes the set of things you need to do in order to run a fully validating Filecoin node.

(TODO: elaborate on all this, obviously)


- Run [chain validation](validation.md)
  - Listen for new blocks on the blocks pubsub channel
  - As new blocks come in, run through the block validation process
    - Keep track of valid blocks, keep track of the current 'best' block (according to EC rules)
  - Rebroadcast valid blocks
- Listen for messages on the messages pubsub channel
  - Validate each message, rebroadcast valid ones
- Run the DHT protocol for aiding node discovery
  - libp2p-kad-dht, with just 'find node' RPCs enabled
- For each node you connect to, run the ['hello' protocol](network-protocols.md#hello-handshake) with them
  - If response shows that you and the other node have different genesis blocks, disconnect from them
    - (Note: maybe this is not the best thing to do?)
  - If the other node gives a valid chain head that is farther ahead than you, 'sync' the chain and maybe switch to it if it's better.



## Mining

### Storage

To be a Filecoin storage miner means run several processes in addition to running a 'Full Node'.

#### Accept Deals

A Filecoin storage miner should listen for, decide upon, and accept storage deals from clients. The data being stored for each of these deals should be placed into a sector, and sealed up once that sector is full. Sealed sectors are then submitted to the chain via `CommitSector`. This logic is generally abstracted away by the 'Sector Sealing Subsystem'.

To accept deals, miners should run the ['storage deal'](network-protocols.md#storage-deal) service.

#### Prove Storage

Once miners have submitted sealed sectors to the chain, they will be on the hook for proving the data over time. Every `proving period`, miners should take their current `proving set` and call `post.GeneratePost` on it. This process will take a fairly long amount of time (TODO: either put specific parameters here, or link to them) and result in a compact Proof of SpaceTime, which must then be submitted to the chain via `SubmitPoSt`

#### Extend the Blockchain

A storage miner is also responsible for producing blocks to extend the blockchain. At every round, storage miners check to see if they are the leader, and if they are, they submit a new block to the network, earning a reward for doing so.

The responsibilities of storage miners are documented in more detail in the [mining document](mining.md)

### Retrieval

To be a filecoin retrieval miner: (todo)

