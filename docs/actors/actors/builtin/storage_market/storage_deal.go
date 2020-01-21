package storage_market

import (
	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	acrypto "github.com/filecoin-project/specs-actors/actors/crypto"
)

// Note: Deal Collateral is only released and returned to clients and miners
// when the storage deal stops counting towards power. In the current iteration,
// it will be released when the sector containing the storage deals expires,
// even though some storage deals can expire earlier than the sector does.
// Collaterals are denominated in PerEpoch to incur a cost for self dealing or
// minimal deals that last for a long time.
// Note: ClientCollateralPerEpoch may not be needed and removed pending future confirmation.
// There will be a Minimum value for both client and provider deal collateral.
type StorageDealProposal struct {
	PieceCID        abi.PieceCID // CommP
	PieceSize       abi.PieceSize
	Client          addr.Address
	Provider        addr.Address
	ClientSignature acrypto.Signature

	// Nominal start epoch. Deal payment is linear between StartEpoch and EndEpoch,
	// with total amount StoragePricePerEpoch * (EndEpoch - StartEpoch).
	// Storage deal must appear in a sealed (proven) sector no later than StartEpoch,
	// otherwise it is invalid.
	StartEpoch           abi.ChainEpoch
	EndEpoch             abi.ChainEpoch
	StoragePricePerEpoch abi.TokenAmount

	ProviderCollateral abi.TokenAmount
	ClientCollateral   abi.TokenAmount
}

func (p *StorageDealProposal) CID() {
	panic("TODO")
}

func (p *StorageDealProposal) Duration() abi.ChainEpoch {
	return (p.EndEpoch - p.StartEpoch)
}

func (p *StorageDealProposal) TotalStorageFee() abi.TokenAmount {
	return abi.TokenAmount(uint64(p.StoragePricePerEpoch) * uint64(p.Duration()))
}

func (p *StorageDealProposal) ClientBalanceRequirement() abi.TokenAmount {
	return (p.ClientCollateral + p.TotalStorageFee())
}

func (p *StorageDealProposal) ProviderBalanceRequirement() abi.TokenAmount {
	return p.ProviderCollateral
}

// Everything in this struct will go on chain
// Provider's signature is implicit in the message containing this structure.
type StorageDeal struct {
	Proposal StorageDealProposal
}

func (d *StorageDeal) CID() {
	panic("TODO")
}

type OnChainDeal struct {
	ID               abi.DealID
	Deal             StorageDeal
	SectorStartEpoch abi.ChainEpoch // -1 if not yet included in proven sector
	LastUpdatedEpoch abi.ChainEpoch // -1 if deal state never updated
	SlashEpoch       abi.ChainEpoch // -1 if deal never slashed
}
