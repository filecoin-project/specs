---
title: Retrieval Market
bookCollapseSection: true
weight: 2
---

# Retrieval Market in Filecoin
---

## Components

The `retrieval market` refers to the process of negotiating deals for a provider to serve stored data to a client. It should be highlighted that the negotiation process for the retrieval happens primarily off-chain. It is only some parts of it (mostly relating to redeeming vouchers from payment channels) that involve interaction with the blockchain.

The main components are as follows:

- A payment channel actor (see {{<link payment_channel_actor>}} for details)
- Two `libp2p` services - 
   - a protocol for making queries
   - a protocol for negotiating and carrying out retrieval deals. The intention for this protocol is to become a separate libp2p protocol as the retrieval market design evolves in the near future.
- A chain-based content routing interface
- A client module to query retrieval miners and initiate deals for retrieval
- A provider module to respond to queries and deal proposals

## VO & V1

V0 of the protocol has participants send data over the retrieval protocol itself in a series of Blocks encoded in Bitswap format and verify received blocks manually. It will only support fetching the payload CID which is at the root of PieceCID's `.car` File, and will only support fetching the whole DAG.

In V1, the retrieval market has evolved to support sending arbitrary payload CIDs & selectors within a piece. Further, it piggybacks on the Data Transfer system and Graphsync to handle transfer and verification, to support arbitrary selectors, and to reduce round trips.
The Data Transfer System is augmented accordingly to support pausing/resuming and sending intermediate vouchers to facilitate this.
V1 will also include additional mechanisms for timeouts and cancellations (to be specified).

Although the underlying protocols will change, the API interfaces for the client & provider have not changed from V0 to V1.

## Deal Flow (V0)

NOTE: V0 is now obsolete, as V1 has been implemented and is being used.

{{<svg src="retrieval_flow_v0.mmd.svg" title="Retrieval Flow - V0" >}}

The baseline version of proposing and accepting a deal in V0 of the retrieval market works as follows:

- The client finds a provider of a given piece with `FindProviders()`.
- The client queries a provider to see if it meets its retrieval criteria (via Query Protocol)
- The client sends a RetrievalDealProposal to the retrieval miner. (via RetrievalProtocol)
- The provider validates the proposal and rejects it if it is invalid
- If the request is valid, the provider responds to it with an accept message
- The client creates a payment channel as neccesary and a "lane" (see [Payment Channels](payment_channels) for more details) and ensures there are enough funds in the channel
- The provider unseals the sector as necessary
- The provider sends blocks over the protocol until it requires payment
- The client consumes blocks over the retrieval protocol and manually verifies them
- When the provider requires payment to proceed, it sends payment request and does not send any more blocks
- The client creates and stores a payment voucher off-chain
- The client responds to the provider with a reference to the payment voucher
- The provider redeems the payment voucher off-chain
- The provider resumes sending blocks
- The client consumes blocks until payment is required again
- The process continues until the end of the query

## Deal Flow (V1)

{{<svg src="retrieval_flow_v1.mmd.svg" title="Retrieval Flow - V1" >}}

The evolved Filecoin Retrieval Market protocol, currently in use, for proposing and accepting a deal works as follows:

- The client finds a provider of a given piece with `FindProviders()`.
- The client queries a provider to see if it meets its retrieval criteria (via Query Protocol)
- The client schedules a `Data Transfer Pull Request` passing the `RetrievalDealProposal` as a voucher.
- The provider validates the proposal and rejects it if it is invalid
- If the proposal is valid, the provider responds with an accept message and begins monitoring the data transfer process
- The client creates a payment channel as necessary and a "lane" and ensures there are enough funds in the channel
- The provider unseals the sector as necessary
- The provider monitors data transfer as it sends blocks over the protocol, until it requires payment
- When the provider requires payment, it pauses the data transfer and sends a request for payment as an intermediate voucher
- The client receives the request for payment
- The client creates and stores payment voucher off-chain
- The client responds to the provider with a reference to the payment voucher, sent as an intermediate voucher
- The provider redeems the payment voucher off-chain
- The provider resumes both the request and sending data
- The process continues until the end of the query

Some extra notes worth making with regard to the above process are as follows:

- The payment channel is created by the client.
- The payment channel is not created until the provider accepts the deal.
- The vouchers are also created by the client and (a reference/identifier to these vouchers is) sent to the provider.
- The payment indicated in the voucher is not taken out of the payment channel funds upon creation and exchange of vouchers between the client and the provider.
- In order for money to be taken out of the payment channel, the provider has to *redeem* the voucher.
- Once the data transfer is complete, there is a 2hr *redeem period* within which the provider has to redeem the voucher, otherwise, the client is free to close the channel. In this case, the provider has not received the funds from the service they provided.
- The provider can redeem the vouchers that they have collected during the transfer or at the end of it.
- The provider can ask for a small payment ahead of the transfer, before they start unsealing data to send to the client. The payment is meant to support the providers' computational cost of unsealing the first chunk of data (where chunk is the agreed step-wise data transfer). This process is needed in order to avoid clients from carrying out a DoS attack, according to which they start several deals and cause the provider to engage a large amount of computational resources.

## Bootstrapping Trust

Neither the client nor the provider have any specific reason to trust the other. Therefore, payment for a retrieval deal is done incrementally, sending vouchers as bytes are sent and verified.

The trust process is as follows:
- When the deal is created, client & provider agree to a "payment interval" in bytes, which is the _minimum_ amount of data the provider will send before each required increment
- They also agree to a "payment interval increase" -- the interval will increase by this value after each successful transfer and payment, as trust develops
- Example:
   - If my "payment interval" is 1000, and my "payment interval increase" is 300, then:
   - The provider must send at least 1000 bytes before they require any payment (they may end up sending slightly more because block boundaries are uneven).
   - The client must pay (i.e., issue a voucher) for all bytes sent when the provider requests payment, provided that the provider has sent at least 1000 bytes.
   - The provider now must send at least 1300 bytes before they request payment again.
   - The client must pay (i.e., issue subsequent vouchers) for all bytes it has not yet paid for when the provider requests payment, assuming it has received at least 1300 bytes since last payment.
   - The process continues until the end of the retrieval, when the last payment will simply be for the remainder of bytes.

## Common Data Types

{{<embed src="types.id" lang="go" >}}
