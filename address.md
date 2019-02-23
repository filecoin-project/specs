# Address

A Filecoin address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. This address encodes information about the network to which an actor belongs, the specific type of address encoding, the address payload itself, and a checksum. The goal of this format is to provide a robust address format that is both easy to use and resistant to errors.

#### Design criteria

1. **Identifiable**: The address must be easily identifiable as a Filecoin address.
2. **Reliable**: Addresses must provide a mechanism for error detection when they might be transmitted outside the network.
3. **Upgradable**: Addresses must be versioned to permit the introduction of new address formats.
4. **Compact**: Given the above constraints, addresses must be as short as possible.

## Specification

A Filecoin address is a byte string consisting of:
* A **network indicator** that identifies which network this address belongs to.
* A **protocol indicator** that identifies the type and version of this address.
* The **payload** used to uniquely identify the actor according to the protocol.
* A **checksum** verifies the address is valid.

```
|---------|----------|---------|----------|
| network | protocol | payload | checksum |
|---------|----------|---------|----------|
| 1 byte  |  1 byte  | n bytes |  4 bytes |
```

An example of a Filecoin address in golang:

```go
type Address struct {
    // 0: main-net
    // 1: test-net
    network byte
    
    // 0: ID
    // 1: Blake2b160-Hash of secp256k1 Public Key
    // 2: Blake2b160-Hash of Actor creation data
    // 3: BLS Public Key
    protocol byte
    
    // raw bytes containing the data associated with protocol
    payload []byte
}
```

#### Network Indicator

The **network** byte identifies the network the address belongs to:

- `0` : Mainnet
- `1` : Testnet

An example description in golang:

```go
// String Prefix
Mainnet_Prefix = 'f'
Testnet_Prefix = 't'

type Protocol = byte
const (
    Mainnet Network = iota
    Testnet
)
```

#### Protocol Indicator

The **protocol indicator** byte describes how a method should interpret the information in the payload field of an address. Any deviation for the algorithms and data types specified by the protocol must be assigned a new protocol number. In this way, protocols also act as versions.

- `0`: ID
- `1` : SECP256K1 Public Key
- `2` : Actor
- `3` : BLS Public Key

An example description in golang:

```go
// Protocol byte
type Protocol = byte
const (
    ID Protocol = iota
    SECP256K1
    Actor
    BLS
)
```

###### Protocol 0: IDs

**Protocol 0** addresses are simple ids. These addresses are expected to remain within the network. All actors have a numeric ID even if they don't have public keys. The key part of an ID address is the decimal representation of the id only characters '0'-'9'. IDs are not hashed and do not have a checksum.

###### Protocol 1: libsecpk1 Elliptic Curve Public Keys

**Protocol 1** addresses represent secp256k1 public encryption keys. The payload field contains the [Blake2b 160](https://blake2.net/) hash of the public key. 

###### Protocol 2: Actor

**Protocol 2** addresses representing an Actor. The payload field contains the [Blake2b 160](https://blake2.net/) hash of meaningful data produced as a result of creating the actor. 

###### Protocol 3: BLS

**Protocol 3** addresses represent BLS public encryptions keys. The payload filed contains the BLS public key.

#### Payload

The payload represents the data specified by the protocol. All payloads excecpt the payload of the ID protocol are base32 encoded when seralized to their human readable format.

#### Checksum

The checksum is calculated using the first 4 bytes of the blake2b-160 hash of the addresses network, protocol, and payload bytes. Checksums are only used when seralizing and deseralizing an address to its human readable format.


#### Expected Methods

All implementations in Filecoin must have methods for creating, encoding, and decoding addresses in addition to checksum creation and validation. The follwing is a golang version of the Address Interface:

```go
type Address interface {
    New(n byte, p byte, payload []byte) Address
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

Encodes an Address as a string, converting the network byte to the corresponding prefix and encoding the payload to [base32](https://tools.ietf.org/html/rfc4648).

```go
func Encode(a Address) string {
	var prefix string
    switch a.network {
    case Mainnet:
        prefix = Mainnet_Prefix // "f"
    case Testnet:
        prefix = Testnet_Prefix // "t"
    default:
        Fatal("invalid address network")
    }
    
    switch a.protocol {
        case SECP256K1:
        case Actor:
        case BLS:
        	cksm := Checksum(a)
            return prefix + a.protocol + base32.Encode(a.payload + cksm)
        case ID:
        	return prefix + a.protocol + base32.Encode(a.payload)
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

Decode an Address from a string, converting the network prefix and protocol prefix to their corresponding byte values, decoding the payload from base32, and validating the checksum.

```go
func DecodeString(a string) Address {
    if len(a) < 4 {
        Fatal(ErrInvalidLength)
    }
    
    var ntwk byte
    switch a[:2] {
    case Mainnet_Prefix:
        ntwk = Mainnet
    case Testnet_Prefix:
    	ntwk = Testnet
    default:
        Fatal(ErrUnknownNetwork)
    }
    
    protocol := a[3]
    payload := a[3:]
    
    if protocol == ID {
        return Address{
            network: ntwk,
            protocol: protocol
            payload: payload
        }
    }
    
    cksm := payload[len(payload)-CksmLen:len(payload)]
    
    if !Validate(a, cksm) {
        Fatal(ErrInvalidChecksum)
    }
    
    return Address{
        network: ntwk,
        protocol: protocol,
        payload: payload,
    }
}
```

##### Checksum()

Checksum produces a byte array of the blake2b-16 checksum of an address, returning the first 4 bytes of the digest.

```go
const CksmLen = 4

func Checksum(a Address) [CKSM_LEN]byte {
    digest = blake2b(a)
    return digest[:4]
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

Filecoin checksums are the first 4 bytes of the blake2b-160 hash of an Address. Checksums are only added to an address when encoding the address to its human readable format. Addresses using Protocol 0 do not have a checksum.