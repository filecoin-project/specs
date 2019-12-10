---
menuTitle: Piece
statusIcon: 🔁
title: Piece - a part of a file
entries:
- piece_store
---


A `Piece` is an object that represents a whole or part of a `File`,
and is used by `Clients` and `Miners` in `Deals`. `Clients` hire `Miners`
to store `Pieces`. 

The piece data structure is designed for proving storage of arbitrary
IPLD graphs and client data. This diagram shows the detailed composition
of a piece and its proving tree, including both full and bandwidth-optimized
piece data structures.

{{< diagram src="pieces.png" title="Pieces, Proving Trees, and Piece Data Structures" >}}

{{< readfile file="piece.id" code="true" lang="go" >}}
