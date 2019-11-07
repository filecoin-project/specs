package crypto

type DomainSeparationTag byte

const (
	TicketTag 			DomainSeparationTag = '0'
	ElectionTag		 	DomainSeparationTag = '1'
	BlockTag		   	DomainSeparationTag = '2'
)
