---
title: Message Pool
entries:
- message_syncer
---

The Message Pool is stores and propagates uncommitted messages in Filecoin.

{{< readfile file="message_pool_subsystem.id" code="true" lang="go" >}}


TODOs:

- discuss how messages are meant to propagate slowly/async
- explain algorithms for choosing profitable txns
