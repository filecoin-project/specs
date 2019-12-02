package crypto

type DomainSeparationTag int

const (
	DomainSeparationTag_TicketProduction DomainSeparationTag = 1 + iota
	DomainSeparationTag_ElectionPoSt
	// ...
)
