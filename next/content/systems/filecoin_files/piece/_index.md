---
title: Piece
weight: 2
bookCollapseSection: true
dashboardAudit: 1
dashboardState: permanent
dashboardInterface: stable
---

# Piece - Part of a file
---

A `Piece` is an object that represents a whole or part of a `File`,
and is used by `Clients` and `Miners` in `Deals`. `Clients` hire `Miners`
to store `Pieces`. 

The piece data structure is designed for proving storage of arbitrary
IPLD graphs and client data. This diagram shows the detailed composition
of a piece and its proving tree, including both full and bandwidth-optimized
piece data structures.


{{< figure src="diagrams/pieces.png" title="Pieces, Proving Trees, and Piece Data Structures" zoom="true">}}


{{< embed src="piece.id" lang="go" >}}
