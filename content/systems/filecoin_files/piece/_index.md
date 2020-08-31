---
title: Piece
weight: 2
bookCollapseSection: true
dashboardWeight: 1.5
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# Piece - Part of a file

A `Piece` is an object that represents a whole or part of a `File`,
and is used by `Clients` and `Miners` in `Deals`. `Clients` hire `Miners`
to store `Pieces`. 

The piece data structure is designed for proving storage of arbitrary
IPLD graphs and client data. This diagram shows the detailed composition
of a piece and its proving tree, including both full and bandwidth-optimized
piece data structures.


![Pieces, Proving Trees, and Piece Data Structures](pieces.png)


{{< embed src="piece.id" lang="go" >}}
