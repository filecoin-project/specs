package address

// import (
// 	util "github.com/filecoin-project/specs/util"
// )

// Addresses for singleton system actors
var (
	InitActorAddr           = &Address_I{} // TODO
	CronActorAddr           = &Address_I{} // TODO
	StoragePowerActorAddr   = &Address_I{} // TODO
	StorageMarketActorAddr  = &Address_I{} // TODO
	PaymentChannelActorAddr = &Address_I{} // TODO
	BurntFundsActorAddr     = &Address_I{} // TODO
)

func (a *Address_I) VerifySyntax(addrType Address_Type) bool {
	panic("TODO")
	// switch aType := addrType; aType {
	// case Address_Protocol.Secp256k1():
	// 	// 80 Bytes
	// 	return len(self)
	// case Address_Protocol.ID():
	// 	// ?
	// case Address_Protocol.Actor():
	// 	// Blake2b - 64 Bytes
	// case Address_Protocol.BLS():
	// 	// BLS-12_381 - 48 Byte PK
	// }
}

func (a *Address_I) String() AddressString {
	return AddressString("") // TODo
}

func (a *Address_I) IsKeyType() bool {
	panic("TODO")
}

func MakeAddress(net Address_NetworkID, t Address_Type) Address {
	return &Address_I{
		NetworkID_: net,
		Type_:      t,
	}
}
