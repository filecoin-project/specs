---
title: Forest
weight: 3
dashboardWeight: 1
dashboardState: reliable
dashboardAudit: n/a
implRepos:
  - { lang: rust, repo: https://github.com/ChainSafe/forest }
---

# Forest

Forest is an implementation of Filecoin written in Rust. It focuses on performance and low resource usage. It is compatible with most of the JSON-RPC API exposed by Lotus, making it easy to use with existing tools and libraries.

Forest does not support all features of the reference implementation but is a good fit for specific applications, such as:

- generating Filecoin snapshots,
- running a bootstrap node,
- running an RPC node.

Forest does not provide storage provider functionality.

Links:

- [Source code](https://github.com/ChainSafe/forest)
- [Documentation](https://docs.forest.chainsafe.io/)
- [Website](https://forest.chainsafe.io/)

The Forest implementation of Filecoin is maintained by [ChainSafe](https://chainsafe.io/).
