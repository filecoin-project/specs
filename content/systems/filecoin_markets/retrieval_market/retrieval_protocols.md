---
title: "Retrieval Protocols"
weight: 2
---

# Retrieval Protocols
---

The `retrieval market` will initially be implemented as two `libp2p` services.

{{<hint info >}}
**Name**: Query Protocol
**Protocol ID V0**: `/fil/<network-name>/retrieval/qry/0.0.1`  
**Protocol ID V1**: `/fil/<network-name>/retrieval/qry/1.0.0`
{{</hint>}}

Request: CBOR Encoded RetrievalQuery Data Structure
Response: CBOR Encoded RetrievalQueryResponse Data Structure

{{<hint info>}}
**Name**: Retrieval Protocol  
**Protocol ID V0**: `/fil/<network-name>/retrieval/0.0.1` 
{{</hint>}}

V0:
Request: CBOR Encoded RetrievalDealProposal Data Structure
Response: CBOR Encoded RetrievalDealResponse Data Structure  
-- Following  
Request: CBOR Encoded RetrievalPayment Data Structure
Response: CBOR Encoded RetrievalDealResponse Data Structure w/ Blocks

V1: Protocol does not exist in this version