type Actor struct {
    // CID of the code object for this actor
    code CID

    // Reference to the root of this actor's state
    head &ActorState

    // Counter of the number of messages this actor has sent
    nonce UInt

    // Current Filecoin balance of this actor
    balance UInt
}
