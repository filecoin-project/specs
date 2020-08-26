---
title: Lotus
weight: 1
implRepos: 
  - { lang: go,   repo: https://github.com/filecoin-project/lotus, auditState: done, auditDate: '2020-07-31' }
  - { lang: go,   repo: https://github.com/filecoin-project/go-fil-markets, auditState: done, auditDate: '2020-07-31' }
  - { lang: go,   repo: https://github.com/filecoin-project/specs-actors, auditState: done, auditDate: '2020-07-31' }
  - { lang: rust, repo: https://github.com/filecoin-project/rust-fil-proofs }
---

# Lotus
---

Lotus is an implementation of the Filecoin Distributed Storage Network. Lotus is written in Go and it is designed to be modular and interoperable with other implementations of Filecoin.

You can run the Lotus software client to join the Filecoin Testnet. Lotus can run on MacOS and Linux. Windows is not supported yet.

The two main components of Lotus are:
1. **The Lotus Node** can sync the blockchain, validating all blocks, transfers, and deals along the way. It can also facilitate the creation of new storage deals. Running this type of node is ideal for users that do not wish to contribute storage to the network, produce new blocks and extend the blockchain.
2. **The Lotus Storage Miner** can register as a miner in the network, register storage, accept deals and store data. The Lotus Storage Miner can produce blocks, extend the blockchain and receive rewards for new blocks added to the network.

You can find the Lotus codebase [here](https://github.com/filecoin-project/lotus) and further documentation, how-to guides and a list of FAQs in at [lotu.sh](https://lotu.sh).

The Lotus implementation of Filecoin is supported by [Protocol Labs](https://protocol.ai/).
