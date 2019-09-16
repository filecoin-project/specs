package codeGen

type NetworkID int

const (
	NetworkID_Testnet NetworkID = 0
	NetworkID_Mainnet NetworkID = 1
)

func (id NetworkID) ToByte() byte {
	switch id {
	case NetworkID_Testnet:
		return 't'
	case NetworkID_Mainnet:
		return 'f'
	default:
		panic("Unknown NetworkID")
	}
}

type AddressProtocol int

const (
	AddressProtocol_ID         AddressProtocol = 0
	AddressProtocol_Secp256k1  AddressProtocol = 1
	AddressProtocol_Actor      AddressProtocol = 2
	AddressProtocol_BLS        AddressProtocol = 3
)

func (p AddressProtocol) ToByte() byte {
	Assert(p >= 0)
	Assert(p <= 9)
	return '0' + byte(p)
}

type Address interface {
	Network() NetworkID
}
