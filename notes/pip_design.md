# Piece Inclusion Proof Design

Miners need to communicate piece inclusion proofs to clients and sometimes put these proofs on chain to demonstrate that they have included a client's piece within a sector. We need a structure for piece inclusion proofs and a specification for how verifiers will verify these proofs. All the proposals listed below assume:

1. The miner constructs a merkle tree from an unsealed sector (sector tree). The root of this tree is CommD. Through commitSector and PoSt messages, the miner publically proves they are in possesion of the sector represented by CommD
2. The client constructs a merkle tree (piece tree) from their piece (the details of this construction may depend on the selected proposal). The root of this tree is CommP.
3. Given CommD, CommP and a PIP (piece inclusion proof), a verifier can verify that data represented by CommP is included within the data included within CommD.

The general strategy for a PIP is to find nodes common to both the piece tree and the sector tree that represent all the piece data. The miner provides a merkle inclusion proof for each of these nodes, this establishing they are in the sector. The after verifying all these proofs, the verifier constructs a partial merkle tree from the nodes. If the tree's root is CommP the proof is valid.

The challege with this style of proof is that it relies on finding common nodes in the two trees and the prevalence of those depend on the alignment of the piece within the sector. In the worst case, only the leaves will be common and the proof will be longer than the data itself. In the best case, the data will perfectly align with a subtree of the sector tree, CommP will be a node in the sector tree, and the miner only needs to provide a single merkle inclusion proof. The miner will have always have the ability to reduce the proof size over the worst case by padding pieces to align or partially align with a subtree. If a miner were required to perfectly align all their pieces, however, we would expect them to waste 1/4 of the sector on average and 1/2 of the sector in the worst case.

The proposals below show ways to optimize the tradeoff between worst case proof size and storage efficiency (each with their own tradeoffs). The question is: how much optimization do we support in the structure of the PIPs to support these options?

## The Shape of PIPs

![PIP Shape Sketch](pip_shape.jpg?raw=true "PIP Shape")

The diagram above illustrates the shape of the various options for PIPs. Every proof is formed from hashes that are either common nodes in the piece and sector trees or just nodes in the sector tree needed to complete the proof to CommD. 

### cover section

The __cover__ section of the proof is a standard merkle inclusion proof to get to the root of the smallest sector subtree that entirely contains the piece. The cover height is the height to the bottom of the cover section, which is log(piece size) +/- 1 based on alignment.

### aligment section

The __alignment height__ is the height of the largest sector subtree such that start of the piece occupies its leftmost node. The larger the alignment height, the more common nodes between trees and the smaller the proof size (in all option) The alignment height of a piece in a tightly packed sector varies greatly depending on the sizes of the pieces that came before. By padding pieces, we can increase the alignment height at the cost of wasted space.

### common section

(the following appeals to intuition, sorry for the lack of rigour)

Any proof that is not fully aligned (cover height + alignment height < log(sector size)), will have a middle section that represents merkle inclusion proofs for all the common nodes. The common nodes will be the roots of a set of subtrees that span the piece, so combining these nodes with __exterior__ nodes to the right and left will allow us to recover the sector tree node at the base of the cover, and we can then use the cover proof to verify against CommD. Combining the common nodes without the exterior information will allow us to recover CommP.

Depending on the construction, the common nodes may form a line at the bottom of the triangle (__suboptimal interior__), or form a line up and down the left and right sides (__optimal interior__). The difference between these sizes is a difference of order and can make a huge difference (the scale of the diagram is deceiving here).

### contributions to proof size

We will be optimizing for the number of hashes, since the orgainizational data needed to detrmine the placement should be small comparatively.

Let:

- tree height = _log(sector size)_
- cover height = _log(piece size)_
- common height = _cover height - alignment height_

Roughly:

* __cover__ = _tree height - cover height_ hashes
* __exterior__ = _2 * common height_ hashes
* __optimal interior__ = _2 * common height_ hashes
* __suboptimal interior__ = _2 ^ common height_ hashes <- note the exponent!

### Design options

| name | size | wasted space (avg) | requires | interactive | in-order | advantages | disadvantages |
|------|------|--------------|----------|-------------|----------|------------|---------------|
| fully aligned | log(sector size / piece size) | 1/4 sector | cover | no | yes | very simple, minimal proof size | maximal space wastage, out-of-order piece inclusion |
| partially aligned | ~2 ^ common height | .5 * pieces * alignment | cover and common | no | yes | maximal worst-case proof size | space wastage tunable, simple |
| interactive | ~4*common height | .5 * pieces * alignment | cover and common | yes | yes | good proof size, low wastage | interactive constructions require significant changes |
| out-of-order | ~4*common height | .5 * pieces * alignment | cover, common, maybe ordering | no | no | good proof size, low wastage | out of order pieces require ordering information or verifier to have knowledge of specific proof construction. |

## Constructions

### Fully Aligned
