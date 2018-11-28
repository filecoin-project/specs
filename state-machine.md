# The Filecoin State Machine

The majority of Filecoin's user facing functionality (payments, storage market, power table, etc) is managed through the Filecoin State Machine. The network generates a series of blocks, and agrees which 'chain' of blocks is the correct one. Each block contains a series of state transitions called `messages`, and a checkpoint of the current `global state` after the application of those `messages`. 

The `global state` here consists of a set of `actors`, each with their own private `state`.

An `actor` is the Filecoin equivalent of Ethereum's smart contracts, it is essentially an 'object' in the filecoin network with state and a set of methods that can be used to interact with it. Every actor has a Filecoin balance attributed to it, a `state` pointer, and a `nonce` which tracks the number of messages sent by this actor. (TODO: the nonce is really only needed for external user interface actors, AKA `account actors`. Maybe we should find a way to clean that up?)

### Method Invocation
There are two routes to calling a method on an `actor`.

First, to call a method as an external participant of the system (aka, a normal user with Filecoin) you must send a signed `message` to the network, and pay a fee to the miner that includes your `message`.  The signature on the message must match the key associated with an account with sufficient Filecoin to pay for the messages execution. The fee here is equivalent to transaction fees in Bitcoin and Ethereum, where it is proportional to the work that is done to process the message (Bitcoin prices messages per byte, Ethereum uses the concept of 'gas'. We will likely also use 'gas').

Second, an `actor` may call a method on another actor.  However, the only time this may happen is as a result of some actor being invoked by an external users message (note: an actor called by a user may call another actor that then calls another actor, as many layers deep as the execution can afford to run for).

### State Representation

The `global state` is modeled as a map of actor addresses to actor structs. This map is implemented by an ipld HAMT (TODO: link to spec for our HAMT). Each actor's `state` is an ipld pointer to a graph that can be entirely defined by the actor.

### Execution

Message execution currently relies entirely on 'built-in' code, with a common external interface. All method invocations have a set of parameters, which are a cbor encoded array of the parameters set by the function definition.

TODO: expand on message execution.