---
title: IPLD
bookCollapseSection: true
weight: 3
dashboardWeight: 1
dashboardState: stable
dashboardAudit: missing
dashboardTests: 0
---

# IPLD - InterPlanetary Linked Data

[IPLD](https://ipld.io/) is the data model of the content-addressable web. It provides standards and formats to build Merkle-DAG data-structures, like those that represent a filesystem. IPLD allows us to treat all hash-linked data structures as subsets of a unified information space, unifying all data models that link data with hashes as instances of IPLD. This means that data can be linked and referenced from totally different data structures. This is a very useful feature that is used in Filecoin.

IPLD introduces several concepts and protocols, such as the concept of content addressing itself, codecs, such as DAG-CBOR, file types, such as Content Addressable aRchives (CARs) and protocols, such as GraphSync.

Please refer to the [IPLD specification](https://github.com/ipld/specs) for more information.

## DAG-CBOR encoding

All system datastructures are stored using DAG-CBOR ( an IPLD codec). DAG-CBOR is a more strict subset of CBOR with a predefined tagging scheme, designed for storage, retrieval and traversal of hash-linked data DAGs.
Files and data stored on the Filecoin network are also stored using various IPLD codecs (not necessarily DAG-CBOR).
IPLD provides a consistent and coherent abstraction above data that allows Filecoin to build and interact with complex, multi-block data structures, such as HAMT and AMT. Filecoin uses the DAG-CBOR codec for the serialization and deserialization of its data structures and interacts with that data using the IPLD Data Model, upon which various tools are built. IPLD Selectors can also used to address specific nodes within a linked data structure.


## GraphSync


