---
title: Retrieval Market
bookCollapseSection: true
weight: 3
dashboardWeight: 2
dashboardState: stable
dashboardAudit: missing
dashboardTests: 0
---

# Retrieval Market in Filecoin

## Components

The `retrieval market` refers to the process of negotiating deals for a provider to serve stored data to a client. It should be highlighted that the negotiation process for the retrieval happens primarily off-chain. It is only some parts of it (mostly relating to redeeming vouchers from payment channels) that involve interaction with the blockchain.

The main components are as follows:

- A [payment channel actor](payment_channels#payment-channel-actor)
- A protocol for making queries
- A Data Transfer subsystem and protocol used to query retrieval miners and initiate retrieval deals
- A chain-based content routing interface
- A client module to query retrieval miners and initiate deals for retrieval
- A provider module to respond to queries and deal proposals

The retrieval market operate by piggybacking on the Data Transfer system and Graphsync to handle transfer and verification, to support arbitrary selectors, and to reduce round trips. The retrieval market can support sending arbitrary payload CIDs & selectors within a piece. 

The Data Transfer System is augmented accordingly to support pausing/resuming and sending intermediate vouchers to facilitate this.


## Deal Flow in the Retrieval Market

![Retrieval Flow](retrieval_flow_v1.mmd)

The Filecoin Retrieval Market protocol for proposing and accepting a deal works as follows:

- The client finds a provider of a given piece with `FindProviders()`.
- The client queries a provider to see if it meets its retrieval criteria (via Query Protocol)
- The client schedules a `Data Transfer Pull Request` passing the `RetrievalDealProposal` as a voucher.
- The provider validates the proposal and rejects it if it is invalid.
- If the proposal is valid, the provider responds with an accept message and begins monitoring the data transfer process.
- The client creates a payment channel as necessary and a "lane" and ensures there are enough funds in the channel.
- The provider unseals the sector as necessary.
- The provider monitors data transfer as it sends blocks over the protocol, until it requires payment.
- When the provider requires payment, it pauses the data transfer and sends a request for payment as an intermediate voucher.
- The client receives the request for payment.
- The client creates and stores a payment voucher off-chain.
- The client responds to the provider with a reference to the payment voucher, sent as an intermediate voucher (i.e., acknowledging receipt of a part of the data and channel or lane value).
- The provider validates the voucher sent by the client and saves it to be redeemed on-chain later
- The provider resumes sending data and requesting intermediate payments.
- The process continues until the end of the data transfer.

Some extra notes worth making with regard to the above process are as follows:

- The payment channel is created by the client.
- The payment channel is created when the provider accepts the deal, unless an open payment channel already exists between the given client and provider.
- The vouchers are also created by the client and (a reference/identifier to these vouchers is) sent to the provider.
- The payment indicated in the voucher is not taken out of the payment channel funds upon creation and exchange of vouchers between the client and the provider.
- In order for money to be transferred to the provider's payment channel side, the provider has to *redeem* the voucher
- In order for money to be taken out of the payment channel, the provider has to submit the voucher on-chain and `Collect` the funds.
- Both redeeming and collecting vouchers/funds can be done at any time during the data transfer, but redeeming vouchers and collecting funds involves the blockchain, which further means that it incurs gas cost.
- Once the data transfer is complete, the client or provider may Settle the channel. There is then a 12hr period within which the provider has to submit the redeemed vouchers on-chain in order to collect the funds. Once the 12hr period is complete, the client may collect any unclaimed funds from the channel, and the provider loses the funds for vouchers they did not submit.
- The provider can ask for a small payment ahead of the transfer, before they start unsealing data. The payment is meant to support the providers' computational cost of unsealing the first chunk of data (where chunk is the agreed step-wise data transfer). This process is needed in order to avoid clients from carrying out a DoS attack, according to which they start several deals and cause the provider to engage a large amount of computational resources.

## Bootstrapping Trust

Neither the client nor the provider have any specific reason to trust each other. Therefore, trust is established indirectly by payments for a retrieval deal done *incrementally*. This is achieved by sending vouchers as the data transfer progresses.

Trust establishment proceeds as follows:
- When the deal is created, client & provider agree to a "payment interval" in bytes, which is the _minimum_ amount of data the provider will send before each required increment.
- They also agree to a "payment interval increment". This means that the interval will increase by this value after each successful transfer and payment, as trust develops between client and provider.
- Example:
   - If my "payment interval" is 1000, and my "payment interval increase" is 300, then:
   - The provider must send at least 1000 bytes before they require any payment (they may end up sending slightly more because block boundaries are uneven).
   - The client must pay (i.e., issue a voucher) for all bytes sent when the provider requests payment, provided that the provider has sent at least 1000 bytes.
   - The provider now must send at least 1300 bytes before they request payment again.
   - The client must pay (i.e., issue subsequent vouchers) for all bytes it has not yet paid for when the provider requests payment, assuming it has received at least 1300 bytes since last payment.
   - The process continues until the end of the retrieval, when the last payment will simply be for the remainder of bytes.

## Data Representation in the Retrieval Market

The retrieval market works based on the Payload CID. The PayloadCID is the hash that represents the root of the IPLD DAG of the UnixFS version of the file. At this stage the file is a raw system file with IPFS-style representation. In order for a client to request  for some data under the retrieval market, they have to know the PayloadCID. It is important to highlight that PayloadCIDs are not stored or registered on-chain.

{{<embed src="github:filecoin-project/go-fil-markets/retrievalmarket/types.go"  lang="go" title="Retrieval Market - Common Data Types">}}