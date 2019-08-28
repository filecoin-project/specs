type ChainManager_BlockState interface {
    self Block
    children []Block
    parents []Block
    num_unvalidated_parents UInt
    chain_weight UInt
    priority UInt
    candidates_connected Set<Block>
}
