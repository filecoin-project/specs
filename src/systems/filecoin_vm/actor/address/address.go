package address

import (
	"errors"

	util "github.com/filecoin-project/specs/util"
)

var Assert = util.Assert

type Int = util.Int

// Addresses for singleton system actors.
var (
	// Distinguished AccountActor that is the source of system implicit messages.
	SystemActorAddr        = Address_Make_ID(Address_NetworkID_Testnet, 0)
	InitActorAddr          = Address_Make_ID(Address_NetworkID_Testnet, 1)
	RewardActorAddr        = Address_Make_ID(Address_NetworkID_Testnet, 2)
	CronActorAddr          = Address_Make_ID(Address_NetworkID_Testnet, 3)
	StoragePowerActorAddr  = Address_Make_ID(Address_NetworkID_Testnet, 4)
	StorageMarketActorAddr = Address_Make_ID(Address_NetworkID_Testnet, 5)
	// Distinguished AccountActor that is the destination of all burnt funds.
	BurntFundsActorAddr = Address_Make_ID(Address_NetworkID_Testnet, 99)
)

const FirstNonSingletonActorId = 100

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

func (a *Address_I) String() string {
	return string(Serialize_Address_Compact(a))
}

func Serialize_Address_Compact(Address) util.Serialization {
	// TODO: custom encoding as in
	// https://github.com/filecoin-project/lotus/blob/master/chain/address/address.go
	panic("TODO")
}

func Deserialize_Address_Compact(util.Serialization) (Address, error) {
	// TODO: custom encoding as in
	// https://github.com/filecoin-project/lotus/blob/master/chain/address/address.go
	panic("TODO")
}

func Deserialize_Address_Compact_Assert(x util.Serialization) Address {
	ret, err := Deserialize_Address_Compact(x)
	Assert(err == nil)
	return ret
}

func (a *Address_I) IsIDType() bool {
	panic("TODO")
}

func (a *Address_I) IsKeyType() bool {
	
	panic("TODO")
}

func (a *Address_I) GetID() (ActorID, error) {
	if !a.IsIDType() {
		return ActorID(0), errors.New("not an ID address")
	}
	return a.Data_.As_ID(), nil
}

func (a *Address_I) GetKey() (KeyHash, error) {
	if !a.IsKeyType() {
		return KeyHash(0), errors.New("not a key address")
	}
	if a.Data_.Which() == Address_Data_Case_PublicKey_BLS {
		return a.Data_.As_PublicKey_BLS(), nil
	} else if a.Data_.Which() == Address_Data_Case_PublicKey_Secp256k1 {
		return a.Data_.As_PublicKey_Secp256k1(), nil
	}
	else {
		return KeyHash(0), errors.New("not a recognized key type")
	}
}

func Address_Make_ID(net Address_NetworkID, x ActorID) Address {
	return &Address_I{
		NetworkID_: net,
		Data_:      Address_Data_Make_ID(x),
	}
}

func Address_Make_Key(net Address_NetworkID, x KeyHash) (Address, error) {
	if util.IsBLS(x) {
		d := Address_Data_Make_PublicKey_BLS(x)
	}
	else if util.IsSECP(x) {
		d := Address_Data_Make_PublicKey_Secp256k1(x)
	} else {
		return nil, errors.New("Not a recognized key type")
	}
	return &Address_I{
		NetworkID_: net,
		Data_:      d,
	}, nil
}

func Address_Make_ActorExec(net Address_NetworkID, hash ActorExecHash) Address {
	return &Address_I{
		NetworkID_: net,
		Data_:      Address_Data_Make_ActorExec(hash),
	}
}

type Address_Ptr = *Address

func (a *Address_I) Ref() Address_Ptr {
	var ret Address = a
	return &ret
}
