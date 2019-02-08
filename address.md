# Address

An address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. An address encodes information about the network it belongs to, the type of data it contains, the data itself, and depending on the type, a checksum.

```go
type Address struct {
    // 0: main-net
    // 1: test-net
    network byte
    
    // 0: SECP256K1 Public Key
    // 1: ActorID
    // 2: Builtin Actor
    // 3: BLS Public Key
    typ byte
    
    // raw bytes associated with typ
    data []byte
}

// checksum size
const CKSM_LEN = 4

// Netowork Byte
const (
    Mainnet = 0
	Testnet = 1
)

// Type Byte
const (
    SECP256K1 = 0
    ID = 1
    Actor = 2
    BLS = 3
)

// String Prefix
const (
    MAINNET_PREFIX = "fc"
    TESTNET_PREFIX = "tf"
)

// New returns an Address for the specified network and type containing data and fails otherwise.
func New(n byte, t byte, d []byte) Address {
	if n != Mainnet && n != Testnet {
        Fatal(ErrUnknownNetwork)
	}
	if t != SECP256K1 && t != ID && t != Actor && t != BLS {
        Fatal(ErrUnknownType)
	}
    return Address{
        network: n,
        typ: t,
        data: d,
    }
}

// Marshal an Address to bytes
func Marshal(a Address) []byte {
    if a.typ == ID {
        return []byte{a.network, a.typ, a.data}
    }
    return []byte{a.network, a.typ, a.data, Checksum(a)}
}

// Unmarshal bytes to an Address
func Unnarshal(a []byte) Address {
    if a[1] == ID {
        return New(a[0], a[1], a[2:])
    }
    
    raw := []byte{a[1], a[2:]}
    data := raw[1 : len(raw) - CKSM_LEN]
    cksm := raw[len(raw) - CKSM_LEN : len(raw)]
    
    if !validChecksum(data, cksm) {
        Fatal(ErrInvalidCksm)
    }
    
    return New(a[0], a[1], data)
}

// EncodeString encodes an Address to a string
func EncodeString(a Address) string {
	var prefix string
    switch a.network {
    case Mainnet:
        prefix = MAINNET_PREFIX
    case Testnet:
        prefix = TESTNET_PREFIX
    default:
        return "<INVALID ADDRESS>"
    }
    
    suffix := a.typ + a.data
    if a.typ == ID {
        return prefix + base58.Encode(suffix)
    }
    return prefix + base58.Encode(suffix + Checksum(suffix))
}

// DecodeString decoes a string to an Address
func DecodeString(a string) Address {
    var ntwk byte
    switch a[:2] {
    case MAINNET_PREFIX:
        ntwk = Mainnet
    case TESTNET_PREFIX:
    	ntwk = Testnet
    default:
        Fatal(ErrUnknownNetwork)
    }
    
    raw := base58.Decode(a[2:])
    if raw[0] != ID {
        return Address{
            network: ntwk,
            typ: raw[0],
            data: raw[1:]
        }
    }
    
    if !validChecksum(raw) {
        Fatal(ErrInvalidChecksum)
    }
    
    return Address{
        network: ntwk,
        typ: raw[0],
        data: raw[1: len(raw)- CKSM_LEN]
    }
}

// Checksum is the last CKSM_LEN bytes the sha256 of Address Type and its data.
func Checksum(a Address) [CKSM_LEN]byte {
    digest = sha256([]byte{a.typ, a.data})
    return digest[len(digest)-CKSM_LEN:]
}

// validChecksum returns true if cksm is valid for data and false otherwise.
func validChecksum(data, cksm []byte) bool {
    digest := Checksum(data)
    return digest == cksm
}
```

