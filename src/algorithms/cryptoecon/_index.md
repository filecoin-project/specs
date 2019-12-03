---
menuTitle: Cryptoecon
title: Cryptoecon -- Placeholders
---

The Filecoin network is a complex multi-agent economic system. This section aims to explain some mechanisms and parameters in the system that can help achieve network-level goals. For now, just lists some key mechanisms and parameters that are still subject to changes during testnet but that must be resolved before mainnet launch. Note that this is a list of economic levers that are in consideration and not all of them will be used in mainnet. Some may be added or changed after mainnet launch, through the FIP process, requiring a **network upgrade**.

1. **Block reward minting function** - parameters may change over time, including the exponent of the block reward decay function.
2. **Block reward vesting function** - block reward earned from mining may be required to vest over some period of time.
3. **Pledge collateral function and slashing** - pledge collateral needs to satisfy certain security constraints for consensus and hence needs to be re-evaluated.
4. **Deal collateral requirement** - minimum deal collaterals may be adjusted to achieve appropriate incentive structure.
5.**Interactive PoRep slashing** - some penalties may be introduced for failing `ProveCommit` in interactive `PoRep`.
6. **Network transaction fee** - burned as a network fee during proof or txn submission.
7. **Reward for pledged but unused storage** - explicit reward for available but unused storage.
8. **Minimum miner size** - the minimum size of sectors a storage miner must have to produce blocks. There are several security and scalability parameters that depend on this.

(This list is incomplete, you can help by expanding it.)
