---
title: Data Transfer
weight: 3
dashboardWeight: 1
dashboardState: stable
dashboardAudit: missing
dashboardTests: 0
---

# Data Transfer in Filecoin

The _Data Transfer Protocol_ is a protocol for transferring all or part of a `Piece` across the network when a deal is made. The overall goal for the data transfer module is for it to be an abstraction of the underlying transport medium over which data is transferred between different parties in the Filecoin network. Currently, the underlying medium or protocol used to actually do the data transfer is GraphSync. As such, the Data Transfer Protocol can be thought of as a negotiation protocol.

The Data Transfer Protocol is used both for Storage and for Retrieval Deals. In both cases, the data transfer request is initiated by the client. The primary reason for this is that clients will more often than not be behind NATs and therefore, it is more convenient to start any data transfer from their side. In the case of Storage Deals the data transfer request is initiated as a _push request_ to send data to the storage provider. In the case of Retrieval Deals the data transfer request is initiated as a _pull request_ to retrieve data by the storage provider.

The request to initiate a data transfer includes a voucher or token (none to be confused with the Payment Channel voucher) that points to a specific deal that the two parties have agreed to before. This is so that the storage provider can identify and link the request to a deal it has agreed to and not disregard the request. As described below the case might be slightly different for retrieval deals, where both a deal proposal and a data transfer request can be sent at once.

## Modules

This diagram shows how Data Transfer and its modules fit into the picture with the Storage and Retrieval Markets.
In particular, note how the Data Transfer Request Validators from the markets are plugged into the Data Transfer module,
but their code belongs in the Markets system.

![Data Transfer](data-transfer-modules.png)

## Terminology

- **Push Request**: A request to send data to the other party - normally initiated by the client and primarily in case of a Storage Deal.
- **Pull Request**: A request to have the other party send data - normally initiated by the client and primarily in case of a Retrieval Deal.
- **Requestor**: The party that initiates the data transfer request (whether Push or Pull) - normally the client, at least as currently implemented in Filecoin, to overcome NAT-traversal problems.
- **Responder**: The party that receives the data transfer request - normally the storage provider.
- **Data Transfer Voucher or Token**: A wrapper around storage- or retrieval-related data that can identify and validate the transfer request to the other party.
- **Request Validator**: The data transfer module only initiates a transfer when the responder can validate that the request is tied directly to either an existing storage or retrieval deal. Validation is not performed by the data transfer module itself. Instead, a request validator inspects the data transfer voucher to determine whether to respond to the request or disregard the request.
- **Transporter**:  Once a request is negotiated and validated, the actual transfer is managed by a transporter on both sides. The transporter is part of the data transfer module but is isolated from the negotiation process. It has access to an underlying verifiable transport protocol and uses it to send data and track progress.
- **Subscriber**: An external component that monitors progress of a data transfer by subscribing to data transfer events, such as progress or completion.
- **GraphSync**: The default underlying transport protocol used by the Transporter. The full graphsync specification can be found [here](https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md)

## Request Phases

There are two basic phases to any data transfer:

1. Negotiation: the requestor and responder agree to the transfer by validating it with the data transfer voucher.
2. Transfer: once the negotiation phase is complete, the data is actually transferred. The default protocol used to do the transfer is Graphsync.

Note that the Negotiation and Transfer stages can occur in separate round trips,
or potentially the same round trip, where the requesting party implicitly agrees by sending the request, and the responding party can agree and immediately send or receive data. Whether the process is taking place in a single or multiple round-trips depends in part on whether the request is a push request (storage deal) or a pull request (retrieval deal), and on whether the data transfer negotiation process is able to piggy back on the underlying transport mechanism. 
In the case of GraphSync as transport mechanism, data transfer requests can piggy back as an extension to the GraphSync protocol using [GraphSync's built-in extensibility](https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md#extensions). So, only a single round trip is required for Pull Requests. However, because Graphsync is a request/response protocol with no direct support for `push` type requests, in the Push case, negotiation happens in a seperate request over data transfer's own libp2p protocol `/fil/datatransfer/1.0.0`. Other future transport mechinisms might handle both Push and Pull, either, or neither as a single round trip.
Upon receiving a data transfer request, the data transfer module does the decoding the voucher and delivers it to the request validators. In storage deals, the request validator checks if the deal included is one that the recipient has agreed to before. For retrieval deals the request includes the proposal for the retrieval deal itself. As long as request validator accepts the deal proposal, everything is done at once as a single round-trip.

It is worth noting that in the case of retrieval the provider can accept the deal and the data transfer request, but then pause the retrieval itself in order to carry out the unsealing process. The storage provider has to unseal all of the requested data before initiating the actual data transfer. Furthermore, the storage provider has the option of pausing the retrieval flow before starting the unsealing process in order to ask for an unsealing payment request. Storage providers have the option to request for this payment in order to cover unsealing computation costs and avoid falling victims of misbehaving clients.

## Example Flows

### Push Flow

![Data Transfer - Push Flow](push-flow.mmd)

1. A requestor initiates a Push transfer when it wants to send data to another party.
2. The requestors' data transfer module will send a push request to the responder along with the data transfer voucher.
3. The responder's data transfer module validates the data transfer request via the Validator provided as a dependency by the responder.
4. The responder's data transfer module initiates the transfer by making a GraphSync request.
5. The requestor receives the GraphSync request, verifies that it recognises the data transfer and begins sending data.
6. The responder receives data and can produce an indication of progress.
7. The responder completes receiving data, and notifies any listeners.

The push flow is ideal for storage deals, where the client initiates the data transfer straightaway
once the provider indicates their intent to accept and publish the client's deal proposal.


## Pull Flow - Single Round Trip

![Data Transfer - Single Round Trip Pull Flow](alternate-pull-flow.mmd)

1. A requestor initiates a Pull transfer when it wants to receive data from another party.
2. The requestor’s data transfer module initiates the transfer by making a pull request embedded in the GraphSync request to the responder. The request includes the data transfer voucher.
3. The responder receives the GraphSync request, and forwards the data transfer request to the data transfer module.
4. The responder's data transfer module validates the data transfer request via a PullValidator provided as a dependency by the responder.
5. The responder accepts the GraphSync request and sends the accepted response along with the data transfer level acceptance response.
6. The requestor receives data and can produce an indication of progress. This timing of this step comes later in time, after the storage provider has finished unsealing the data.
7. The requestor completes receiving data, and notifies any listeners.

## Protocol

A data transfer CAN be negotiated over the network via the Data Transfer Protocol, a libp2p protocol type.

Using the Data Transfer Protocol as an independent libp2p communciation mechanism is not a hard requirement -- as long as both parties have an implementation of the Data Transfer Subsystem that can talk to the other, any
transport mechanism (including offline mechanisms) is acceptable.

## Data Structures


{{<embed src="github:filecoin-project/go-data-transfer/types.go"  lang="go" title="Data Transfer Types">}}


{{<embed src="github:filecoin-project/go-data-transfer/statuses.go"  lang="go" title="Data Transfer Statuses">}}


{{<embed src="github:filecoin-project/go-data-transfer/manager.go"  lang="go" symbol="Manager" title="Data Transfer Manager">}}
