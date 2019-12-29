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

- A payment channel actor (See {{<sref payment_channel_actor>}} for details)
- 2 `libp2p` services - 
   - a protocol for making queries
   - a protocol for negotiating and carrying out retrieval deals (will disappear as a seperate libp2p protocol as retrieval markets evolve)
- A chain-based content routing interface
- A client module to query retrieval miners and initiate deals for retrieval
- A provider module to respond to queries and deal proposals

# VO & V1

V0 of the protocol has participants send data over the retrieval protocol itself in a series of Blocks encoded in Bitswap format and verify received blocks manually. It will only support fetching the payload CID which is at the root of PieceCID's Car File, and will only support fetching the whole DAG.

In V1, the retrieval markets will evolve to support sending arbitrary payload CID's & selectors within a piece (V1). Further, it will piggy back on the Data Transfer system and Graphsync to handle transfer and verification, to support arbitrary selectors, and to reduce roundtrips.
The Data Transfer System will accordingly be augmented to support pausing/resuming and sending intermediate vouchers to facilitate this.
V1 will also include additional mechanisms for timeouts and cancellations. (to be specified)

Though the underlying protocols will change, the API interfaces for the client & provider will not change from V0 to V1.

# Deal Flow (V0)

{{< diagram src="retrieval_flow_v0.mmd.svg" title="Retrieval Flow - V0" >}}

The baseline version of proposing and accepting a deal will work as follows:

- The client finds a provider of a given piece with `FindProviders()`.
- The client queries a provider to see if it meets its retrieval criteria (via Query Protocol)
- The client sends a RetrievalDealProposal to the retrieval miner. (via RetrievalProtocol)
- The provider validates the proposal and rejects it if it is invalid
- If the request is valid, the provider responds to it with an accept message
- The client creates a payment channel as neccesary and a lane, ensures there are free funds in the channel
- The provider unseals the sector as neccesary
- The provider sends blocks over the protocol until it requires payment
- The client consumes blocks over the retrieval protocol and manually verifies them
- When the provider requires payment to proceed, it sends payment request and does not send any more blocks
- The client puts a payment voucher on chain
- The client responds to the provider with a reference to the payment voucher
- The provider redeems the payment voucher on the chain
- The provider resumes sending blocks
- The client consumes blocks until payment is required again
- The process continues until the end of the query

# Deal Flow (V1)

{{< diagram src="retrieval_flow_v1.mmd.svg" title="Retrieval Flow - V1" >}}

The evolved protocol for proposing and accepting a deal will work as follows:

- The client finds a provider of a given piece with `FindProviders()`.
- The client queries a provider to see if it meets its retrieval criteria (via Query Protocol)
- The client schedules a `Data Transfer Pull Request` passing the `RetrievalDealProposal` as a voucher.
- The provider validates the proposal and rejects it if it is invalid
- If the proposal is valid, the provider responds with an accept message and begins monitoring the data transfer process
- The client creates a payment channel as neccesary and a lane, ensures there are free funds in the channel
- The provider unseals the sector as neccesary
- The provider monitors data transfer as it sends blocks over the protocol, until it requires payment
- When the provider requires payment, it pauses the data transfer and sends a request for payment as an intermediate voucher
- The client receives the request for payment
- The client puts a payment voucher on chain
- The client responds to provider with a reference to the payment voucher, sent as an intermediate voucher
- The provider redeems the payment voucher on the chain
- The provider unpauses the request and data resumes sending 
- The process continues until the end of the query

# Bootstrapping Trust

Neither the client nor the provider have any specific reason to trust the other. Therefore, payment for a retrieval deal is done in pieces, sending vouchers as bytes are sent and verified.

The trust process is as follows:
- When the deal is created, client & provider agree to a "payment interval" in bytes, which is the _minimum_ amount of data the provider will send before each required increment
- They also agree to a "payment interval increase" -- the interval will increase by this value after each successful transfer and payment, as trust develops
- Example:
   - If my "payment interval" is 1000, and my "payment interval increase" is 300:
   - The provider must send at least 1000 bytes before it requires any payment (it may send more cause block boundaries are uneven)
   - The client must pay for all bytes sent when the provider requests payment, if the provider has sent at least 1000 bytes
   - The provider now must send at least 1300 bytes before it requests payment again
   - The client must pay for all bytes it hasn't yet paid for when the provider
   requests payment, assuming it's received at least 1300 bytes since last payment
   - The process continues till the end of the retrieval, when the last payment will simply be for the remainder of bytes
- Additional trust mechanisms in the V1 version of the protocol will include agreed upon timeouts and cancellation fees

# Common Data Types

{{< readfile file="types.id" code="true" lang="go" >}}
