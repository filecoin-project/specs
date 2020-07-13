---
title: Message Syncer
weight: 1
---

# Message Syncer
---

{{< hint warning >}}
TODO:

- explain message syncer works
- include the message syncer code
{{< /hint >}}


## Message Propagation

Messages are propagated over the libp2p pubsub channel `/fil/messages`. On this channel, every serialised `SignedMessage` is announced (see [Message](\missing-link)).

Upon receiving the message, its validity must be checked: the signature must be valid, and the account in question must have enough funds to cover the actions specified. If the message is not valid it should be dropped and must not be forwarded.

**TODO:** discuss checking signatures and account balances, some tricky bits that need consideration. Does the fund check cause improper dropping? E.g. I have a message sending funds then use the newly constructed account to send funds, as long as the previous wasn't executed the second will be considered "invalid" ... though it won't be at the time of execution.
