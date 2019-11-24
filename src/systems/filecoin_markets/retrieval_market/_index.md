---
menuTitle: Retrieval Market
statusIcon: 🔁
title: "Retrieval Market in Filecoin"
entries:
- retrieval_peer_resolver
- retrieval_protocols
- retrieval_client
- retrieval_provider
---

# Components

The `retrieval market` refers to the process of negotiating deals for a provider to serve stored data to a client.

The main components are as follows:

- A payment channel actor (See {{<sref payment_channel_actor "Payment Channel Actor">}} for details)
- 2 `libp2p` services - 
   - a protocol for making queries
   - a protocol for negotiating and carrying out retrieval deals (will disappear as a seperate libp2p protocol as retrieval markets evolve)
- A chain-based content routing interface
- A client module to query retieval miners and initiate deals for retrieval
- A provider module to respond to queries and deal proposals

# VO & V1

The retrieval markets will evolve to support sending arbitrary payload CID's & selectors within a piece (V1). The first version (V0) will only support fetching the payload CID which is at the root of PieceCID's Car File, and will only support fetching the whole DAG.

The first version will send data over the retrieval protocol itself a series of Blocks encoded in Bitswap format and verify received blocks manually. 

The protocol will evolve to piggy back on the Data Transfer system and Graphsync to handle transfer and verification, and to support arbitrary selectors, and to reduce roundtrips. The Data Transfer System will be augmented to support pausing/resuming and sending intermediate vouchers to facilitate this

The underlying protocols will change, but the API interfaces for the client & provider will not change

# Deal Flow (V0)

{{< diagram src="retrieval_flow_v0.mmd.svg" title="Retrieval Flow - V0" >}}

The baseline version of proposing and accepting a deal will work as follows:

- The client finds a provider of a given piece with FindProviders.
- The client queries a provider to see if it meets its retrieval criteria (via Query Protocol)
- The client creates a payment channel as neccesary and a lane, ensures there are free funds in the channel
- The client sends a RetrievalDealProposal to the retrieval miner. (via RetrievalProtocol)
- The providers validates the proposal and rejects it if it is invalid: invalid price per byte, invalid payment channel are reasons for rejection
- If the request is valid, the provider responds to with an accept message
- The provider may also send a small number of blocks over the protocol before requiring payment
- The client consumes blocks over the retrieval protocol and manually verifies them
- When the requires payment to proceed, it sends a request for payment and does not send any more blocks
- The client responds with a payment (or cancels the request)
- The provider resumes sending blocks
- The client consumes blocks until payment is required again
- The process continues till the end of the query

# Deal Flow (V1)

{{< diagram src="retrieval_flow_v1.mmd.svg" title="Retrieval Flow - V1" >}}

The evolved protocol for proposing and accepting a deal will work as follows:

- The client finds a provider of a given piece with FindProviders.
- The client queries a provider to see if it meets its retrieval criteria (via Query Protocol)
- The client creates a payment channel as neccesary and a lane, ensures there are free funds in the channel
- The client schedules a Data Transfer Pull Request passing the RetrievalDealProposal as a voucher.
- The providers validates the proposal and rejects it if it is invalid: invalid price per byte, invalid payment channel are reasons for rejection
- If the proposal is valid, the provider responds with an accept message and begins monitoring the data transfer process
- The provider may begin sending a small amount of bytes before first payment
- When the provider requires payment, it pauses the data transfer and sends a request for payment as an intermediate voucher
- The client receives the request for payment responds with a payment, also sent as an intermediate voucher (or cancels the request)
- The provider unpauses the request and data resumes sending and the process continues till the end of the query

### Bootstrapping Trust

Neither the client nor the provider have any specific reason to trust the other. Therefore, payment for a retrieval deal is done in pieces, sending vouchers as bytes are sent and verified, or potentially requiring prepayment to proceed. To minimize the ability for either side to fault, the miner can choose to require payment upfront before sending some data, or send data before notifying the client it needs more. The client can choose to cancel the rest of a request if the payment terms become unfriendly. Generally, initial increments should be small and then as trust is built, the exchange can shift to working in larger pre-payment or post-payment increments. Different miners may choose to operate on prepayment or postpayment and may switch during the course of a deal. (prepayment may make particular sense towards the end of a deal)
