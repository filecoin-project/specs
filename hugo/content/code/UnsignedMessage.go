type UnsignedMessage struct {
    to Address
    from Address

    // When receiving a message from a user account the nonce in the message must match
    // the expected nonce in the "from" actor. This prevents replay attacks.
    nonce UInt
    value UInt

    gasPrice UInt
    gasLimit UInt

    method Uint
    params Bytes  // Serialized parameters to the method.
} // representation tuple
