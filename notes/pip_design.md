## Piece Inclusion Proof Design

Miners need to communicate piece inclusion proofs to clients and sometimes put these proofs on chain to demonstrate that they have included a client's piece within a sector. We need a structure for piece inclusion proofs and a specification for how verifiers will verify these proofs. All the proposals listed below assume:

1. The miner constructs a merkle tree from an unsealed sector (sector tree). The root of this tree is CommD. Through commitSector and PoSt messages, the miner publically proves they are in possesion of the sector represented by CommD
2. The client constructs a merkle tree (piece tree) from their piece (the details of this construction may depend on the selected proposal). The root of this tree is CommP.
3. Given CommD, CommP and a PIP (piece inclusion proof), a verifier can verify that data represented by CommP is included within the data included within CommD.

The general strategy for a PIP is to find nodes common to both the piece tree and the sector tree that represent all the piece data. The miner provides a merkle inclusion proof for each of these nodes, this establishing they are in the sector. The after verifying all these proofs, the verifier constructs a partial merkle tree from the nodes. If the tree's root is CommP the proof is valid.

The challege with this style of proof is that it relies on finding common nodes in the two trees and the prevalence of those depend on the alignment of the piece within the sector. In the worst case, only the leaves will be common and the proof will be longer than the data itself. In the best case, the data will perfectly align with a subtree of the sector tree, CommP will be a node in the sector tree, and the miner only needs to provide a single merkle inclusion proof. The miner will have always have the ability to reduce the proof size over the worst case by padding pieces to align or partially align with a subtree. If a miner were required to perfectly align all their pieces, however, we would expect them to waste 1/4 of the sector on average and 1/2 of the sector in the worst case.

The proposals below show ways to optimize the tradeoff between worst case proof size and storage efficiency (each with their own tradeoffs). The question is: how much optimization do we support in the structure of the PIPs to support these options?

### The Shape of PIPs

