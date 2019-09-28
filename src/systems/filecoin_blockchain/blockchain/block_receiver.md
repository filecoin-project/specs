---
title: Block Receiver
---

A node must decode and perform syntactic validation for every block received
before passing it on (e.g. in a lipbp2p pubsub validator).

# Syntactic Validation

{{<goFile Block>}}

A syntactically valid block:

- must include a well-formed miner address
- must include at least one well-formed ticket, and if more they form a valid ticket chain
- must include an election proof which is a valid signature by the miner address of the final ticket
- must include at least one parent CID
- must include a positive parent weight
- must include a positive height
- must include well-formed state root, messages, and receipts CIDs
- must include a timestamp not in the future
