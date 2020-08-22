---
title: CID
description: CIDs - Content IDentifiers
dashboardWeight: 1
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# CIDs - Content IDentifiers
---

For most objects referenced by Filecoin, a Content Identifier (CID for short) is used. Any pointer inclusions in the Filecoin spec `id` files (e.g. `&Object`) denotes the CID of said object. Some objects explicitly name a CID field. The spec treats these notations interchangeably.
This is effectively a hash value, prefixed with its hash function (multihash) as well as extra labels to inform applications about how to deserialize the given data.

For a more detailed specification, we refer the reader to the
[IPLD repository](https://github.com/ipld/cid).
