---
menuTitle: Deals
title: Market Deals
---

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties. All deal negotiation happen off chain and a request-response style storage deal protocol is in place to submit agreed-upon storage deals onto the network with `CommitSector` to gain storage power on chain. Hence, there is a `StorageDealProposal` and a `RetrievalDealProposal` that are half sign contracts submitted by clients to be counter-signed and posted on-chain by the miners. 

{{< readfile file="deal.id" code="true" lang="go" >}}
