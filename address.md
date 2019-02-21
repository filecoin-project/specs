# Address

A Filecoin address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. This address encodes information about the network to which an actor belongs, the specific type of address encoding, the address payload itself, and a checksum. The goal of this format is to provide a robust address format that is both easy to use and resistant to errors.

#### Design criteria

1. **Identifiable**: The address must be easily identifiable as a Filecoin address.
2. **Reliable**: Addresses must provide a mechanism for error detection when they might be transmitted outside the network.
3. **Upgradable**: Addresses must be versioned to permit the introduction of new address formats.
4. **Compact**: Given the above constraints, addresses must be as short as possible.

## Specification

A Filecoin address is a byte string consisting of:
* The **letter 'F'** to identify the address as belonging to a Filecoin network.
* A **network indicator** that identifies which network this address belongs to.
* A **protocol indicator** that identifies the type and version of this address.
* The **payload** used to uniquely identify the actor according to the protocol.

```
|---|---------|----------|---------|
| f | network | protocol | payload |
|---|---------|----------|---------|
```

An example of a Filecoin address in golang:

```go
type Address struct {
    // 0: main-net
    // 1: test-net
    network leb128
    
    // 0: SECP256K1 Blake2b160-Hash of Public Key
    // 1: Actor Blake2b160-Hash of random data
    // 2: BLS Public Key
    // 3: ID
    protocol leb128
    
    // raw bytes containing the data associated with protocol
    payload []byte
}
```

#### Network Indicator

The **network** is a [LEB128](https://en.wikipedia.org/wiki/LEB128) variable-length integer that identifies the network:

- `0` :  Mainnet
- `1` : Testnet

An example description in golang:

```go
// String Prefix
const (
    MAINNET_PREFIX = "FC"
    TESTNET_PREFIX = "FT"
)

// Network varint
leb128 Mainnet = 0
leb128 Testnet = 1
```

#### Protocol Indicator

The **protocol indicator** is a [LEB128](https://en.wikipedia.org/wiki/LEB128) variable-length integer that describes how a method should interpret the information in the `payload` field of an address. Any deviation for the algorithms and data types specified by the protocol must be assigned a new protocol number. In this way, protocols also act as versions.

- `0` : SECP256K1
- `1` : Actor
- `2` : BLS Public Key
- `3` : ID

An example description in golang:

```go
// Protocol varint
const (
	SECP256K1 = 0
    Actor     = 1
    BLS       = 2
    ID        = 3 
)
```

###### Protocol 0: libsecpk1 Elliptic Curve Public Keys

**Protocol 0** addresses represent secp256k1 public encryption keys. The payload field contains a hash and a checksum. The hash is a [Blake2b 160](https://blake2.net/) hash of the public key. The checksum is a six character double-sha256 checksum. 

###### Protocol 1: Actor

**Protocol 1** addresses represent a builtin Actor. The payload field contains a hash and a checksum. The hash is a [Blake2b 160](https://blake2.net/) hash of random date (**TODO**: better definition). The checksum is a six character double-sha256 checksum.

###### Protocol 2: BLS Public Keys

**Protocol 2** addresses represent BLS public encryptions keys. The payload field contains the public key as bytes and a checksum.  The checksum is a six character double sha256 checksum. 

###### Protocol 3: IDs

**Protocol 3** addresses are simple ids. These addresses are expected to remain within the network. All actors have a numeric ID even if they don't have public keys. The key part of an ID address is the decimal representation of the id only characters '0'-'9'. IDs are not hashed and do not have a checksum.

#### Payload

**TODO** Encoding inside the byte array. Add serialization for different protocols.


#### Expected Methods

All implementations in Filecoin must have methods for creating, encoding, and decoding addresses in addition to checksum creation and validation. The follwing is a golang version of the Address Interface:

```go
type Address interface {
    New(n varint, p varint, payload []byte) Address
    Encode(a Adress) string
    Decode(s string) Address
    Checksum(a Address) []byte
    Validate(a Address) bool
}
```
##### New()

New returns an Address for the specified network and protocol encapsulating corresponding payload. New fails for unknown network or protocol.

```go
func New(network byte, protocol byte, payload []byte) Address {
	if network != Mainnet && network != Testnet {
        Fatal(ErrUnknownNetwork)
	}
	if protocol < SECP256K1 || protocol > BLS {
        Fatal(ErrUnknownType)
	}
    return Address{
        Network:  network,
        protocol: protocol,
        Payload:  payload,
    }
}
```

##### Encode()

Software encoding a Filecoin address must:

- produce an address encoded to a known network
- produce an address encoded to a known protocol
- produce an address with a valid checksum

Encodes an Address as a string, converting the network value to the corresponding prefix and encoding the protocol and payload to base32.

```go
func Encode(a Address) string {
	var prefix string
    switch a.network {
    case Mainnet:
        prefix = MAINNET_PREFIX // "FC"
    case Testnet:
        prefix = TESTNET_PREFIX // "FT"
    default:
        Fatal("invalid address network")
    }
    
    switch a.protocol {
        case SECP256K1:
        case Actor:
        case BLS:
        	suffix := leb128.Encode(a.protocol) + a.payload
            return prefix + base32.Encode(suffix + Checksum(suffix))
    	case ID:
        	suffix := leb128.Encode(a.protocol) + a.payload
        	return prefix + base32.Encode(suffix)
        default:
        	Fatal("invalid address protocol")
    }
}
```

##### Decode()

Software decoding a Filecoin address must:
* verify the network is a known network.
* verify the protocol is a number of a known protocol.

* verify the checksum is valid

Decoded an Address from a string, converting the network prefix to the corresponding LEB128 value, decoding the protocol and payload, converting the protocol to the corresponding LEB128 value, and validating the checksum.

```go
func DecodeString(a string) Address {
    var ntwk leb128
    switch a[:2] {
    case MAINNET_PREFIX:
        ntwk = Mainnet
    case TESTNET_PREFIX:
    	ntwk = Testnet
    default:
        Fatal(ErrUnknownNetwork)
    }
    
    raw := base32.Decode(a[2:])
    if raw[0] != ID {
        return Address{
            network: ntwk,
            protocol: leb128.Decode(raw),
            payload: raw[len(ntwk)+len(protocol):]
        }
    }
    
    if !Validate(raw) {
        Fatal(ErrInvalidChecksum)
    }
    
    return Address{
        network: ntwk,
        protocol: leb128.Decode(raw),
        payload: raw[len(protocol): len(raw)- CKSM_LEN]
    }
}
```

##### Checksum()

Checksum produces a byte array of length CKSM_LEN derived from the double-sha256 of an addresses protocol and payload.

```go
const CKSM_LEN = 6

func Checksum(a Address) [CKSM_LEN]byte {
    digest = doubleSha256([]byte{a.protocol, a.payload})
    return digest[len(digest)-CKSM_LEN:]
}
```

##### Validate()

Validate returns true if the Checksum of data matches the expected checksum.

```go
func Validate(data, expected []byte) bool {
    digest := Checksum(data)
    return digest == expected
}
```

### Checksums

> TODO a better definition of checksums

Filecoin checksums are a double sha256 checksum