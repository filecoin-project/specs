# Repair Market Notes

!! This is an sketch of the Filecoin Repair Protocol, it is not final and will not be implemented as described in mainnet !!

The `repair market` is the network where clients can initiate and form `deals` for `repair-backed storage` with `repair miners`. We define `repair-backed storage` to be storage that has been added to the network via a `repair miner`. `repair-backed storage` can be thought of as storage that has been replicated and distributed to ensure that the underlying stored data is retrievable even when several faults have occurred. We discuss `faults` [here](add link). We explore `repair miners`, `storage miners`, and their relationship [here](add link).



Interface Additions:

```go
type RepairMarket interface {	
	// AddRepairAsk adds a repair backed storage ask to the orderbook. Must be
    // called indirectly by a repair miner actor created by this storage
    // market actor.
    AddRepairAsk(price *TokenAmount, size *BytesAmount, expiry uint64) RAskID

    // AddBid adds a bid order to the orderbook. Can be called by anyone. The
    // message must contain the appropriate amount of funds to be locked up for the
    // bid. 
    // 'droot' is a merkleroot of the properly encoded data (see 'Encoding')
    // 'price' is the price per GB per block being paid
    // 'size' is the number of bytes to be stored (post encoding)
    // 'time' is a number of blocks that the data will be stored for
    // 'expiry' is the number of blocks the bid is valid for after it makes it on chain.
    // 'params' specifies the erasure coding parameters for this data
    AddBid(droot *cid.Cid, price *TokenAmount, size *BytesAmount, time, expiry uint64, params CodingParams) BidID
    
    // AddDeal creates a deal from the given ask and bid
    // It must always called by the repair miner in the ask
    // This is for deals between clients and repair miners.
    AddDeal(askID RAskID, bidID BidID, bidOwnerSig Signature) DealID
    
    // CreateRepairMiner creates a new repair miner with the given public key and a pledge
    // of the given size. The pubkey is the one that should be used to sign deals
    // with other entities, the account that calls this method will be set as the 'owner'
    CreateRepairMiner(pubk PublicKey, pledge *BytesAmount, pid libp2p.PeerID) Address
    
    // SlashRepairMiner is used to slash a repair miner who has not correctly posted
    // their proofs of repair in a given timeframe. (TODO: economics)
    SlashRepairMiner(addr Address)
}
```



### The difference between the `storage market` and `repair market`

It is important to understand the difference between storage providers and what services they provide. Note that storage providers are Filecoin miners come in two varieties: storage and repair. A `storage miner` is physically stores data and is the only miner that does so. This miner is responsible for generate Proofs of Space-Time (`PoST`) and managing physical storage in sectors. A `repair miner` does not physically store client data but instead works with `storage miners` to manage fault-tolerate versions of a clients data using `erasure codes`. One might think of a `storage miner` as a hard drive and a `repair miner` as a RAID controller. For example, suppose a client wishes to store its data in a fault-tolerant way, that client would reach out and make a deal with a `repair miner` to store data `client_data`. This `repair miner` would now work with `storage miners` to store parts of the `client_data` (which would be parsed based on an [erasure code](add link)) amongst the `storage miners`. This would ensure that `client_data` would be fault-tolerant up to `x` faults. Note that the mapping defining this parsing would be held by the `repair miner`. We will discuss this further [here](add link).

Note that a client may choose to make `deals` directly with a `storage miner` but there will not be a guarantee of fault tolerance for the stored data.



### Deal failure and data recovery

The Filecoin storage market must ensure that `client_data` remains available in the face of data faults. In order to accomplish this, the Filecoin storage market relies on the functions of `repair miners`. Let's consider the case where deals made directly with a `storage miner`. In this case, a `storage miner` will fault if it fails to prove that it is currently storing the `client_data` via a `PoST`. Such a fault may occur if the `storage miner` goes offline or has disk failures; other types of faults are discussed [here](add link). Thus, `client_data` becomes irretrievable due to a single fault at the `storage miner`.

Now let's consider the case where a client requests that `client_data` be retrievable under some number of faults. In this case we would like to distribute `client_data` to multiple `storage miners` in some efficient and retrievable way; this is where `repair miners` come into play. `repair miners` accept storage deals from clients, similarly to how `storage miners` might accept direct deals from clients. However, unlike `storage miners`, the `repair miners` delegate the actual physical storage of `client data` to `storage miners` according to an `erasure code`. In the event that one of those `storage miners` fails, the `repair miner` will detect this fault, recover the lost data using the erasure code, and communicate the recovered data to a new `storage miner`.



## Market Flow (with repair)

- Client selects erasure coding parameters for their data, and encodes it. They then take the encoded data and put it into a single merkledag, and take the hash of the root as their 'dataroot'
- The client then posts a bid to the storage market containing the dataroot, the erasure coding params, the funds required for storage, and other parameters.
- The client then looks for a repair miner ask on chain that suits their needs
- Once found, they contact the repair miner and send them a deal proposal
- If the repair miner agrees, they respond with a tenative acceptance of the deal
- The client now sends the data to the repair miner
- The repair miner validates that they have received the correct data, and signs a deal for storing the data, which they post to chain *and* send to the client
- Upon receipt of the signed deal, the client can be assured their file will be stored, and may leave
  - TODO: there is a 'double spend' type issue here where the repair miner may be able to push another transaction to the chain before the one they send to the client, causing the deal to become invalid. A mechanism for slashing repair miners who do this should probably be implemented.
- EXIT CLIENT
- Now, the repair miner has a period of time in which they have to collect proofs from storage miners that the data they agree to store is stored
- For each piece of data they agreed to store, they must contact a storage miner, and run the 'make storage deal' protocol. (also described above)
  - Repair Miner sends the Storage Miner a signed storage deal proposal, including the repair deal, and the data to be stored
  - Upon receiving and deciding to accept the proposal, the storage miner responds with an agreement to store the data in a given sector at a specified offset
  - The storage miner then finishes filling up that sector and seals it
  - Once the seal is complete, they provide the repair miner with a proof the data is stored within the sector as agreed.
    - If the storage miner fails to provide the proof, the repair miner may post the agreement on chain and the storage miner must respond with the correct proof or be slashed.
  - Once the repair miner receives the correct proof, they send the storage miner a series of time-locked payment channel updates to pay for the storage. These updates are valid contingent on the absense of any fault claims on-chain made by the repair miner against the storage miner for the file in question.
    - Still TODO, this needs a bit more work. We could make the condition for the payments require a certain sector exist on-chain, and send them up-front, thus removing the risk of the repair miner not sending the data. Alternatively, some settlement protocol can be devised to solve this.





```go
type RepairMinerActor interface {
    AddAsk(price *TokenAmount, size *NumBytes) RAskID
    // TODO...
}
```





## The `repair miner`

In order to provide data repair durability guarantees, the services of a repair miner are required. A repair miners job is to make deals with clients to store their data, and then delegate the actual long term storage of that data to dedicated storage miners. The repair miners hold the data necessary to complete repairs of any data lost by faulty storage miners.

### Repair Miner State

The repair miner is responsible for keeping around all information needed to ensure they can correctly repair the data they have deals for, as well as information required to prove they are doing their job correctly. 

This data includes:

- The 'File Root' object for every deal they have made with a client. This is the object referenced by the hash placed in client bids. It contains pointers to each of the chunks of the file being stored. This is necessary to know what data needs to be fetched when a repair needs to happen.
- 'File Inclusion Proofs' for every piece they store with a storage miner. These proofs are most likely just a merkle proof that shows the file in question is contained in a sector that is committed on-chain by a storage miner.
- Temporarily, the repair miner must store the files they are agreeing to store for clients. Once the data is safely delegated to a storage miner, they may delete it from their local storage. 
- The total amount of storage that the repair miner is provably delegating. This can be used as a 'Power' metric for the repair miner.

### Repair

The filecoin repair protocol is a system for automatically detecting damaged or missing sectors and reintroducing that data back into the network in order to maintain the resiliency of all files stored.

This protocol encompasses a way to detect faults, and a way to repair faults. These responsibilities are handled by dedicated repair miners, who constantly monitor the data they are responsible for failure, and reintroduce any missing data to new miners.

### Fault Detection

In order to detect storage faults, repair miners should watch the miners that they have storage agreements with, and make sure that they do not miss their PoSTs, and also that they do not remove any sectors from their proving set that they should still be storing. Note that since miners must report any faulty sectors along with their PoST, detection of missing data becomes fairly simple.

### Sector Repair

Upon seeing a sector that contains data they are responsible for fail, repair miners should recover the missing data by checking the erasure coding parameters, and using that to recover the correct pieces required to repair the chunks that are gone. Once the data is recovered, they simply need to make another agreement with some other storage miner to store the data, collect those proofs, and then they are done. Note that this means the repair miner is on the hook for any additional storage costs incurred by repairing data.



## Claiming Earnings

The money that clients pay to repair miners for their services is locked up when they initially post a bid. It can be withdrawn over time, proportional to the time of the deal that has elapsed.
