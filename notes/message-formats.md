# Example message formats

The following code is from the [go-filcoin repo](https://github.com/filecoin-project/go-filecoin/blob/master/types/block.go):


## Blocks on the blockchain

```go
//Example structure of a block in the FIL blockchain
type block struct {
  
  // Unique identity of the miner `actor` who is this block's leader
  miner [size?]type
  
  // StateRoot is a cid pointer to the state tree after application
  // of the transactions state transitions.
  stateRoot [size?]byte
  
  // Description of messgeRoot
  messageRoot: [size?]byte
  
  // ...
  receiptsRoot: [size?]byte
  
  // ...
  parent: [size?]type
  
  // ...
  weight: [size?]type
  
  // ???
  signature: [size? 65?]byte
  
  // ...
  height: [size?]type
  
  // ...
  ticket: [size?]type
  
  // A list of all messages included in this block. Messages can 
  // be of any type supported by FIL. This includes all `actors` as
  // well as transactions such as deals and SEALs. There can be up
  // to X messages in a block.
  messages []message
}
```

## Messages in a block

```go
type msg struct {
  // List of all sub-messages associated with this message. An example
  // where there are multiple messages is X.
  messages []message
}
```

## Items in a message

```go
type message struct {

  // Size of the cryptographic signature
  signature [65]byte 
  
  // Target address of the message
  toAddress: [size?]byte
    
  // Nonce is a temporary field used to differentiate blocks 
  // for testing
  nonce: [size?]byte
    
  /// ???
  value: [size?]type
    
  // Limit in the amoubnt of work associated with the message
  gasLimit: [size?]type
    
  // Cost of work as of the instatiation of this message
  gasPrice: [size?]type
    
  // ???
  methodLabel: [size?]type
}
```

## 