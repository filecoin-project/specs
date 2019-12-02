package address

import (
	util "github.com/filecoin-project/specs/util"
)

type Int = util.Int

// Addresses for singleton system actors
var (
	InitActorAddr           = &Address_I{} // TODO
	CronActorAddr           = &Address_I{} // TODO
	StoragePowerActorAddr   = &Address_I{} // TODO
	StorageMarketActorAddr  = &Address_I{} // TODO
	PaymentChannelActorAddr = &Address_I{} // TODO
	BurntFundsActorAddr     = &Address_I{} // TODO
)

func (a *Address_I) VerifySyntax() bool {
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

func (a *Address_I) Equals(Address) bool {
	panic("TODO")
}

func (a *Address_I) String() AddressString {
	return Serialize_Address_Compact(a)
}

func Serialize_Address_Compact(Address) AddressString {
	// TODO: custom encoding as in
	// https://github.com/filecoin-project/lotus/blob/master/chain/address/address.go
	panic("TODO")
}

func Deserialize_Address_Compact(AddressString) (Address, error) {
	// TODO: custom encoding as in
	// https://github.com/filecoin-project/lotus/blob/master/chain/address/address.go
	panic("TODO")
}

func (a *Address_I) IsKeyType() bool {
	panic("TODO")
}

func Address_Make_ID(net Address_NetworkID, x ActorID) Address {
	return &Address_I{
		NetworkID_: net,
		Data_:      Address_Data_Make_ID(x),
	}
}

func Address_Make_ActorExec(net Address_NetworkID, hash ActorExecHash) Address {
	return &Address_I{
		NetworkID_: net,
		Data_:      Address_Data_Make_ActorExec(hash),
	}
}
