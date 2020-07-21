---
title: Storage Market
weight: 1
bookCollapseSection: true
dashboardAudit: 1
dashboardState: wip
---

# Storage Market in Filecoin
---

Storage Market subsystem is the data entry point into the network. Storage miners only earn power from data stored in a storage deal and all deals live on the Filecoin network. Specific deal negotiation process happens off chain, clients and miners enter a storage deal after an agreement has been reached and post storage deals on the Filecoin network to earn block rewards and get paid for storing the data in the storage deal. A deal is only valid when it is posted on chain with signatures from both parties and at the time of posting, there are sufficient balances for both parties locked up to honor the deal in terms of deal price and deal collateral.

## Terminology

- **StorageClient** - The party that wants to make a deal to store data
- **StorageProvider** - The party that will store the data in exchange for payment. A storage miner.
- **StorageMarketActor** - The on-chain component of deals. The StorageMarketActor is analogous to an escrow and a ledger for all deals made.
- **StorageAsk** - The current price and parameters a miner is currently offering for storage (analogous to an Ask in a financial market)
- **StorageDealProposal** - A proposal for a storage deal, signed only by the - `Storage client`
- **StorageDeal** - A storage deal proposal with a counter signature from the Provider, which then goes on-chain.

## Deal Flow

The lifecycle for a deal within the storage market contains distinct phases:

1. **Discovery** - Client identifies miners and determines their current asks
2. **Negotiation** (out of band) - Both parties come to an agreement about the terms of the deal, each party commits funds to the deal, data is transferred
4. **Publishing** - Deal is published on chain, making the storage provider publicly accountable for the deal
5. **Handoff** - Once the deal is published, it is handled off the Storage Mining Subsystem. The Storage Mining Subsystem will add the deal to a sector,
seal the sector, and tell the Storage Market Actor that the deal is in a sector, thereby making the deal active.

The remaining parts of the deal lifecycle are handled by the Storage Mining Subsystem, which communicates with the Storage Market Actor in order to process deal payments. See [Storage Mining Subsystem](storage_mining) for more details.

The following diagram outlines the phases of deal flow within the storage market in detail:

{{<svg src="storage_market_flow.mmd.svg" title="Storage Market Deal Flow" >}}

# Discovery

Discovery is the process where a client identifies storage providers (i.e. a miner) who might be able to store their data. There are many ways a client might identify a provider to store their data. This outlines the minimum discovery services a filecoin implementation MUST provide. As the network evolves, third parties may build systems that supplement or enhance these services.

Discovery involves identifying providers and determining their current `StorageAsk`. The steps are as follows:
1. A client queries the chain to gain a list of Storage Miner Actors who have registerd as miners with the StoragePowerActor
2. A client may perform additional queries to each Storage Miner Actor to determine their properties -- worker address, sector size, libp2p Multiaddress etc
3. Once the client identifies potentially suitable providers, it sends a direct libp2p message using the `Storage Query Protocol` to get each potential provider's current `StorageAsk`.
4. Miners respond on the `AskProtocol` with a signed version of their current `StorageAsk`

A `StorageAsk` contains all the properties that a client will need to determine if a given provider will meet its needs for storage at this moment. Providers should update their asks frequently to insure the information they are providing to clients is up to date.

## Negotiation

Negotiation is the out of band process where a storage client and a storage provider come to an agreement about a storage deal and reach the point where a deal is published on chain.

Negotiation begins once a client has discovered a miner whose `StorageAsk` meets their desired criteria.The *recommended* order of operations for negotiating and publishing a deal is as follows:

1. Before sending a proposal to the provider, the `StorageClient` adds funds for a deal, as necessary, to the `StorageMarketActor` (by calling `AddBalance`)
2. In order to propose a storage deal, the `StorageClient` then calculates the piece commitment (`CommP`) for the data it intends to store ahead of time. This is neccesary so that the `StorageProvider` can verify the data the `StorageClient` sends to be stored matches the `CommP` in the `StorageDealProposal`. For more detail about the relationship between payloads, pieces, and `CommP` see [Piece](piece) and [Filproofs](filproofs).
3. The `StorageClient` now creates a `StorageDealProposal` and sends the proposal and the CID for the root of the data payload to be stored to the `StorageProvider` using the `Storage Deal Protocol`

Execution now moves to the `StorageProvider`

4. The `StorageProvider` inspects the deal to verify that the deal's parameters match its own internal criteria (such as price, piece size, deal duration, etc) -- and rejects the proposal if it doesn't (it sends a rejection to the client over the
`Storage Deal Protocol`)
5. The `StorageProvider` queries the `StorageMarketActor` to verify the `StorageClient` has deposited enough funds to make the deal (i.e. balance > total storage price) and rejects the proposal if it hasn't.
6. If all criteria are met, the `StorageProvider` responds using the `Storage Deal Protocol` indicate an intent to accept the deal.

Execution now moves back to the `StorageClient`

7. The `StorageClient` opens a push request for the payload data using the `Data Transfer Module`, and sends the request to the provider along with a voucher containing the CID for the `StorageDealProposal` 
8. The `StorageProvider` checks the voucher and verifies the CID matches a storage deal proposal it has received and verified but not put on chain already. If so, it accepts the data transfer request from the `StorageClient`.
9. The `Data Transfer Module` now transfers for the payload data to be stored from the  `StorageClient` to the `StorageProvider` using `GraphSync`.
10. Once complete, the `Data Transfer Module` notifies the `StorageProvider`
11. The `StorageProvider` recalculates the piece commitment (`CommP`) from the data that was transferred and verifies it matches the piece commitment in the `StorageDealProposal`

## Publishing

Data is now transferred, both parties have agreed, and it's time to publish the deal. Given that the counter signature on a deal proposal is a standard message signature by the provider and the signed deal is an on-chain message, it is usually the `StorageProvider` that publishes the deal. However, if `StorageProvider` decides to send this signed on-chain message to the client before calling `PublishStorageDeal` then client can publish the deal on chain. Client's funds are not locked until the deal is published and a published deal that is not activated within some window will result in an on-chain penalty.

12. First, the `StorageProvider` adds collateral for the deal as needed to the `StorageMarketActor` (using `AddBalance`)
13. Now, the `StorageProvider` prepares and signs the on-chain `StorageDeal` message with the `StorageDealProposal` signed by the client and its own signature. It can now either send this message back to the client or call `PublishStorageDeals` on the `StorageMarketActor` to publish the deal. It is recommended for `StorageProvider` to send back the signed message before `PublishStorageDeals` is called.
14. After calling `PublishStorageDeals`, `StorageProvider` sends a message to the `StorageClient` on the `Storage Deal Protocol` with the CID of the message that it is putting on chain for convenience.
15. If all goes well, the `StorageMarketActor` responds with an on chain `DealID` for the published deal.


Finally, the client verifies the deal.

16. The `StorageClient` queries the node for the CID of the message published on chain (sent by the provider). It then inspects the message parameters to make sure they match the expected deal.

## Handoff

Now that a deal is published, it needs to be stored, sealed, and proven in order for the provider to be paid. (see [Storage Deal](storage_deal) for more information about how deal payments are made). These later stages of a deal are handled by the [Storage Mining Subsystem](storage_mining). So the final task for the Storage Market is to handoff to the Storage Mining Subsystem

1.  The `StorageProvider` writes the serialized, padded piece to a shared [Filestore](filestore). 
2.  The `StorageProvider` calls `HandleStorageDeal` on the `StorageMiner` with the published `StorageDeal` and filestore path.

A note re order of operations: the only requirement to publish a storage deal with the `StorageMarketActor` is that the `StorageDealProposal` is signed by the `StorageClient`, the publish message is signed by the `StorageProvider`, and both parties have deposited adequate funds/collateral in the `StorageMarketActor`. As such, it's not required that the steps listed above happen in exactly the order listed above. However, the above order is *recommended* because it generally minimizes the ability of either party to act maliciously.

## Data Types

The following data types unique to the Storage Market:

{{<embed src="storage_market.id" lang="go" >}}
{{<embed src="storage_market.go" lang="go" >}}

Details about `StorageDealProposal` and `StorageDeal` (which are used in the Storage Market and elsewhere) specifically can be found in [Storage Deal](storage_deal).

## Protocols
{{<hint info >}}
**Name**: Storage Query Protocol  
**Protocol ID**: `/fil/<network-name>/storage/ask/1.0.1`
{{</hint >}}

Request: CBOR Encoded AskProtocolRequest Data Structure
Response: CBOR Encoded AskProtocolResponse Data Structure

{{<hint info >}}
**Name**: Storage Deal Protocol  
**Protocol ID**: `/fil/<network-name>/storage/mk/1.0.1`
{{</hint >}}

Request: CBOR Encoded DealProtocolRequest Data Structure
Response: CBOR Encoded DealProtocolResponse Data Structure
