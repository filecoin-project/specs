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
| binary aligned | log(sector size / piece size) | 1/4 sector | cover | no | yes | very simple, minimal proof size | maximal space wastage, out-of-order piece inclusion |
| partially aligned | ~2 ^ common height | .5 * pieces * alignment | cover and common | no | yes space wastage tunable, simple | | maximal worst-case proof size | 
| interactive | ~4*common height | .5 * pieces * alignment | cover and common | yes | yes | good proof size, low wastage | interactive constructions require significant changes |
| out-of-order | ~4*common height | .5 * pieces * alignment | cover, common, maybe ordering | no | no | good proof size, low wastage | out of order pieces require ordering information or verifier to have knowledge of specific proof construction. |

## Constructions

### Binary Aligned Proofs

Binary aligned proofs require miners to only store binary-aligned pieces (ie. pieces whose size are the sector size over a power of 2). This has the effect of turning the alignment all the way up to the cover, so there is no common section. This is the only way we can guarantee CommP is found within the sector tree. This gives us minimally small proofs. The proof is a simple merkle inclusion proof so it is also easiest to encode and to verify.

The downsides are that there is no way for miners to tune the amount of wasted space in the sector. A miner could completely fill a sector with pieces and yet only store half a sector's worth of usable data (the other half will be padding). To get even that storage efficiency, a miner will need to store pieces out of order. For example, if the pieces come in 1/4, 1/2, and 1/4 of the sector size, the miner will place the first 1/4 in the first 1/4 of the sector. The second piece will only fit in the second 1/2 of the sector. To fill the sector, the miner will need to put the second 1/4 in between. 

### Partially Aligned Proofs

With partially aligned proofs, the miner will either pad pieces to a fixed power of two (e.g. round to the next Mb) or pad to some local alignment minimum. The distinction is that miners are not required to pad so that entire pieces are binary aligned. The miner then finds the largest subtree of the left side of the piece tree that overlaps with a subtree in the sector tree and places it. The miner repeats this process until the piece is fully place. The piece will be contiguous and in-order in the sector.

The advantage of this is that the miner gets to tune exactly how much padding they can live with to get better proof sizes. With this and the following constructions, the miner can still binary align the pieces for optimal proof sizes at their discretion. While this construction requires us to be able to store the common part (a tree of inclusion proofs) and decode it for verification, it doesn't require any special knowlege on the verifier's part about how the proof was constructed.

The disadvantage is that all the subtrees the miner places will be no bigger than the first, and all of them will be less than or equal to the alignment height. So we require piece size / 2 ^ (alignment height) (plus side hashes and cover). So alignment makes a big difference, but the proof is still linear in piece size. If we require Mb alignment and store a 100Mb piece, our PIP will contain over 100 hashes.

### Interactive Proofs

Interactive proofs require the client to get involved to improve piece size. The client begins a storage deal by asking to place a piece of a given size. The miner responds with a location within the sector where they intend to store the piece. The client then computes the smallest fully aligned subtree of the sector tree that will fully contain their pieces and generates a merkle tree of the same size with the same piece placement to compute CommP. This means most of the interior nodes of the piece tree will be found within the sector tree. The miner constructs a proof as before, but this time they are able to find much taller overlapping subtrees.

The resulting proof will have a set of subtrees of descending size on the left, and another set of descending subtrees going to the right. The result is that they will have to include two times the height of the common area hashes (plus exterior hashes and the cover). This is logarithmic in piece size. The miner can align the piece to reduce the proof size further, but this now has a linear rather than exponential effect on proof size.

This is a dramatic improvement, but the interaction is a high cost. In addition to the additionaly complexity in the proof flow, it requires that the miner determine piece placement at the start of the deal rather than just prior to sealing. More importantly, it probably requires miners to put a lease on the allocate space while the deal progresses so they don't attempt to put two deals in the same space.

### Out-of-order Proofs

It's possible to acheive the same proof sizes as interactive proofs without interaction, but it requires storing parts of the piece out of order. The idea is to binary align the largest subtrees of the piece tree with the largest subtrees in the sector tree where the piece will reside. The algorithm looks something like this:

1. Put the piece tree in a list of subtrees
2. If a piece in the list binary aligns at the sector location, place it, move to sector location after the subtree, and repeat this step.
3. Else if the list contains a subtree that is too big to binary align, pick the smallest subtree that is too big, split it so that one piece binary aligns, put them all in the list and goto step 2.
4. Else place the largest piece. Repeat step 4.

This creates the same triangular shape to the subtree sizes in the interactive construction and gets roughly the same size proof.

The first downside of this construction is that the hashes in the proof are out-of-order and will need to be reordered before a verifier can recover CommP. This could mean including ordering information with the proof. It might be able to recover the odering if the verifier understands how it was contructed, but this will require labeling the proof as an out-of-order proof with some sort of metadata. Out-of-order pieces will also be a problem for retrieval, but the retrieval miner could use the same mechanism as the verifier to recover the original piece.

