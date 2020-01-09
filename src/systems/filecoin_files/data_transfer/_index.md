---
menuTitle: Data Transfer
statusIcon: 🔁
title: "Data Transfer in Filecoin"
---

_Data Transfer_ is a system for transferring all or part of a `Piece` across the network when a deal is made.

# Modules

This diagram shows how Data Transfer and its modules fit into the picture with the Storage and Retrieval Markets.
In particular, note how the Data Transfer Request Validators from the markets are plugged into the Data Transfer module,
but their code belongs in the Markets system.

{{< diagram src="data-transfer-modules.png" title="Data Transfer - Push Flow" >}}


# Terminology

- **Push Request**: A request to send data to the other party
- **Pull Request**: A request to have the other party send data
- **Requestor**: The party that initiates the data transfer request (whether Push or Pull)
- **Responder**: The party that receives the data transfer request
- **Data Transfer Voucher**: A wrapper around storage or retrieval data that can identify and validate the transfer request to the other party
- **Request Validator**: The data transfer module only initiates a transfer when the responder can validate that the request is tied directly to either an existing storage deal or retrieval deal. Validation is not performed by the data transfer module itself. Instead, a request validator inspects the data transfer voucher to determine whether to respond to the request.
- **Scheduler**:  Once a request is negotiated and validated, actual transfer is managed by a scheduler on both sides. The scheduler is part of the data transfer module but is isolated from the negotiation process. It has access to an underlying verifiable transport protocol and uses it to send data and track progress.
- **Subscriber**: An external component that monitors progress of a data transfer by subscribing to data transfer events, such as progress or completion.
- **GraphSync**: The default underlying transfer protocol used by the Scheduler. The full graphsync specification can be found at [https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md](https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md)

# Request Phases

There are two basic phases to any data transfer:

1. Negotiation - the requestor and responder agree to the transfer by validating with the data transfer voucher
2. Transfer - Once both parties have negotiated and agreed upon, the data is actually transferred. The default protocol used to do the transfer is Graphsync

Note that the Negotiation and Transfer stages can occur in seperate round trips,
or potentially the same round trip, where the requesting party implicitly agrees by sending the request, and the responding party can agree and immediately send or receive data.

# Example Flows

## Push Flow

{{< diagram src="push-flow.mmd.svg" title="Data Transfer - Push Flow" >}}

1. A requestor initiates a Push transfer when it wants to send data to another party.
2. The requestors' data transfer module will send a push request to the responder along with the data transfer voucher. It also puts the data transfer in the scheduler queue, meaning it expects the responder to initiate a transfer once the request is verified
3. The responder's data transfer module validates the data transfer request via the Validator provided as a dependency by the responder
4. The responder's data transfer module schedules the transfer
5. The responder makes a GraphSync request for the data
6. The requestor receives the graphsync request, verifies it's in the scheduler and begins sending data
7. The responder receives data and can produce an indication of progress
8. The responder completes receiving data, and notifies any listeners

The push flow is ideal for storage deals, where the client initiates the push
once it verifies the the deal is signed and on chain

## Pull Flow

{{< diagram src="pull-flow.mmd.svg" title="Data Transfer - Pull Flow" >}}

1. A requestor initiates a Pull transfer when it wants to receive data from another party.
2. The requestors' data transfer module will send a pull request to the responder along with the data transfer voucher.
3. The responder's data transfer module validates the data transfer request via a PullValidator provided as a dependency by the responder
4. The responder's data transfer module schedules the transfer (meaning it is expecting the requestor to initiate the actual transfer)
5. The responder's data transfer module sends a response to the requestor saying it has accepted the transfer and is waiting for the requestor to initiate the transfer
6. The requestor schedules the data transfer
7. The requestor makes a GraphSync request for the data
8. The responder receives the graphsync request, verifies it's in the scheduler and begins sending data
9. The requestor receives data and can produce an indication of progress
10. The requestor completes receiving data, and notifies any listeners

The pull flow is ideal for retrieval deals, where the client initiates the pull when the deal is agreed upon.

# Alternater Pull Flow - Single Round Trip

{{< diagram src="alternate-pull-flow.mmd.svg" title="Data Transfer - Single Round Trip Pull Flow" >}}

1. A requestor initiates a Pull transfer when it wants to receive data from another party.
2. The requestor’s DTM schedules the data transfer
3. The requestor makes a Graphsync request to the responder with a data transfer request
4. The responder receives the graphsync request, and forwards the data transfer request to the data transfer module
5. The requestors' data transfer module will send a pull request to the responder along with the data transfer voucher.
6. The responder's data transfer module validates the data transfer request via a PullValidator provided as a dependency by the responder
7. The responder's data transfer module schedules the transfer
8. The responder sends a graphsync response along with a data transfer accepted response piggypacked
9. The requestor receives data and can produce an indication of progress
10. The requestor completes receiving data, and notifies any listeners

# Protocol

A data transfer CAN be negotiated over the network via the {{<sref data_transfer_protocol "Data Transfer Protocol">}}, a Libp2p protocol type

A Pull request expects a response. The requestor does not initiate the transfer
until they know the request is accepted.

The responder should send a response to a push request as well so the requestor can release the resources (if not accepted). However, if the Responder accepts the request they can immediately initiate the transfer

Using the Data Transfer Protocol as an independent libp2p communciation mechanism is not a hard requirement -- as long as both parties have an implementation of the Data Transfer Subsystem that can talk to the other, any
transport mechanism (including offline mechanisms) is acceptable.

# Data Structures

{{< readfile file="data_transfer_subsystem.id" code="true" lang="go" >}}
