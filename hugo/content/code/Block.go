type Block struct {
    // Miner is the address of the miner actor that mined this block.
    miner Address

    // Tickets is a chain (possibly singleton) of tickets ending with a winning ticket
    // submitted with this block.
    tickets [Ticket]

    // ElectionProof is generated from a past ticket and proves this miner is a leader
    // in this block's round.
    electionProof ElectionProof

    // Parents is an array of distinct CIDs of parents on which this block was based.
    // Typically one, but can be several in the case where there were multiple winning
    // ticket-holders for a given round. The order of parent CIDs is not defined.
    parents [&Block]

    // ParentWeight is the aggregate chain weight of the parent set.
    parentWeight UInt

    // Height is the chain height of this block.
    height UInt

    // StateRoot is a CID pointer to the VM state tree after application of the state
    // transitions corresponding to this block's messages.
    stateRoot &StateTree

    // Messages is the set of messages included in this block. This field is the CID
    // of the TxMeta object that contains the bls and secpk signed message trees.
    messages &TxMeta

    // BLSAggregate is an aggregated BLS signature for all the messages in this block
    // that were signed using BLS signatures.
    blsAggregate Signature

    // MessageReceipts is a set of receipts matching to the sending of the `Messages`.
    // This field is the CID of the root of a sharray of MessageReceipts.
    messageReceipts &[MessageReceipt]

    // The block Timestamp is used to enforce a form of block delay by honest miners.
    // Unix time UTC timestamp (in seconds) stored as an unsigned integer.
    timestamp Timestamp

    // BlockSig is a signature over the hash of the entire block with the miners
    // worker key to ensure that it is not tampered with after creation
    blockSig Signature
}

type TxMeta struct {
    blsMessages &[&Message]<Sharray>
    secpkMessages &[&SignedMessage]<Sharray>
}
