package serde

import autil "github.com/filecoin-project/specs/actors/util"

// Serializes a structure or value to CBOR.
func Serialize(o interface{}) ([]byte, error) {
	autil.TODO("CBOR-serialization")
	return nil, nil
}

func MustSerialize(o interface{}) []byte {
	s, err := Serialize(o)
	autil.AssertMsg(err == nil, "serialization failed")
	return s
}

// Serializes an array of method invocation params.
func MustSerializeParams(o ...interface{}) []byte {
	return MustSerialize(o)
}

// Deserializes a structure or value from CBOR.
func Deserialize(b []byte, out interface{}) error {
	autil.TODO("CBOR-deseriaization")
	return nil
}
