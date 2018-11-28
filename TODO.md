
- [ ] process on how to sync the blockchain
- [ ] hello protocol
- [ ] message processing (aka, how do i do state transitions?)
- [ ] retrieval market not yet done
- [x] nothing about proofs yet
- [ ] 'how to run a node' should be documented
- [ ] cbor CHAMP documentation
- [ ] need link to 'cbor-ipld' spec
  - [ ] https://github.com/ipld/specs/pull/73
- [ ] 'id' stuff
- [ ] message signatures
- [ ] payment channel usage



----



- split proofs up into abstract and concrete
  - naming is key, don't call the implementation just 'Proof of Replication'



------------

Node Operation : 

- [ ] When you start describing how you sync to the longest chain when connecting to a full node, I'm wondering whether this is describing expected consensus? Because that's the consensus process 

Storage market: 

- [ ] For the storage market, is it always miners who put the ask out, and clients are takers? Or is there asks as well as client bids, which miners can accept?
- [ ] Why does the client run a "local merkle translation" before the deal?
- [ ] Does the client "look at" asks, or simply select the lowest? Is this automated node behavior, or a human decision?
- [ ] Same with miner decision to accept, is it automated or human?
- [ ] "proofs of spacetime" needs backticks to indicate it's a special word
- [ ] Is the process of miners sending in proofs every proving period scalable on-chain? Or is this an off-chain thing?



-------

libp2p spec process:

Base wire formats:

- [ ] multistream select
- [ ] multiplex
- [ ] yamux
- [ ] secio

Then a doc describing exactly how these are tied together, a spec for the 'swarm'.

A document describing how dialing works.

A document describing the 'libp2p api'. Which is not language specific, but a broader description of the capabilities of each subsystem in terms of what it does.

Libp2p 'builtin' protocols

- [ ] Identify
- [ ] Ping
- [ ] pubsub
- [ ] DHT