---
title: Indices
weight: 3
dashboardWeight: 1
dashboardState: incorrect
dashboardAudit: 0
dashboardTests: 0
---

# Macroeconomic Indices
---

Indices are a set of global economic indicators computed from State Tree and a collection of pure functions to compute policy output based on user state/action. Indices are used to compute and implement economic mechanisms and policies for the system. There are no persistent states in Indicies. Neither can Indices introduce any state mutation. Note that where indices should live is a design decision. It is possible to break Indices into multiple files or place indices in different actors once all economic mechanisms have been decided on. Temporarily, Indices is a holding file for all potential macroeconomic indicators that the system needs to be aware of.
