# Address

A Filecoin address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. This address encodes information about the network to which an actor belongs, the specific type of address encoding, the address payload itself, and a checksum. The goal of this format is to provide a robust address format that is both easy to use and resistant to errors.

#### Design criteria

1. **Identifiable**: The address must be easily identifiable as a Filecoin address.
2. **Reliable**: Addresses must provide a mechanism for error detection when they might be transmitted outside the network.
3. **Upgradable**: Addresses must be versioned to permit the introduction of new address formats.
4. **Compact**: Given the above constraints, addresses must be as short as possible.

## Specification

There are 2 ways a filecoin address can be represented. An address appearing on chain will always be formatted as raw bytes. An address may also be encoded to a string, this encoding includes a checksum and network prefix. An address encoded as a string will never appear on chain, this format is used for sharing among humans.

##### Bytes

When represented as bytes a filecoin address contains the following:

* A **protocol indicator** byte that identifies the type and version of this address.
* The **payload** used to uniquely identify the actor according to the protocol.
```
|----------|---------|
| protocol | payload |
|----------|---------|
|  1 byte  | n bytes |
```

##### String

When encoded to a string a filecoin address contains the following:

* A **network prefix** character that identifies the network the address belongs to.
* A **protocol indicator** byte that identifies the type and version of this address.
* A **payload** used to uniquely identify the actor according to the protocol.
* A **checksum** used to validate the address.

```
|------------|----------|---------|----------|
|  network   | protocol | payload | checksum |
|------------|----------|---------|----------|
| 'f' or 't' |  1 byte  | n bytes | 4 bytes  |
```

An example of a Filecoin address in golang:

```go
type Address struct {

	// 0: Mainnet
	// 1: Testnet
	Network byte

	// 0: ID
	// 1: Blake2b160-Hash of secp256k1 Public Key
	// 2: Blake2b160-Hash of Actor creation data
	// 3: BLS Public Key
	Protocol byte

	// raw bytes containing the data associated with protocol
	Payload []byte
}
```

#### Network Prefix

The **network prefix** is prepended to an address when encoding to a string. The network prefix indicates which network an address belongs in. The network prefix may either be `f` for filecoin mainnet or `t` for filecoin testnet. Is it worth noting that a network prefix will never appear on chain and is only used when encoding an address to a human readable format.

#### Protocol Indicator

The **protocol indicator** byte describes how a method should interpret the information in the payload field of an address. Any deviation for the algorithms and data types specified by the protocol must be assigned a new protocol number. In this way, protocols also act as versions.

- `0 ` : ID
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

**Protocol 0** addresses are simple IDs.  All actors have a numeric ID even if they don't have public keys. The payload of an ID address is base10 encoded. IDs are not hashed and do not have a checksum.

**Bytes**

```
|----------|---------------|
| protocol |    payload    |
|----------|---------------|
|    0     | leb128-varint |
```

**String**

```
|------------|----------|---------------|
|  network   | protocol |    payload    |
|------------|----------|---------------|
| 'f' or 't' |    '0'   | leb128-varint |
                  base10[...............]
```

###### Protocol 1: libsecpk1 Elliptic Curve Public Keys

**Protocol 1** addresses represent secp256k1 public encryption keys. The payload field contains the [Blake2b 160](https://blake2.net/) hash of the public key.

**Bytes**

```
|----------|---------------------|
| protocol |        payload      |
|----------|---------------------|
|    1     | blake2b-160(PubKey) |
```

**String**

```
|------------|----------|---------------------|----------|
|  network   | protocol |      payload        | checksum |
|------------|----------|---------------------|----------|
| 'f' or 't' |    '1'   | blake2b-160(PubKey) |  4 bytes |
                  base32[................................]
```

###### Protocol 2: Actor

**Protocol 2** addresses representing an Actor. The payload field contains the [Blake2b 160](https://blake2.net/) hash of meaningful data produced as a result of creating the actor.

**Bytes**

```
|----------|---------------------|
| protocol |        payload      |
|----------|---------------------|
|    1     | blake2b-160(Random) |
```

**String**

```
|------------|----------|-----------------------|----------|
|  network   | protocol |         payload       | checksum |
|------------|----------|-----------------------|----------|
| 'f' or 't' |    '2'   |  blake2b-160(Random)  |  4 bytes |
                  base32[..................................]
```

###### Protocol 3: BLS

**Protocol 3** addresses represent BLS public encryption keys. The payload field contains the BLS public key.

**Bytes**

```
|----------|---------------------|
| protocol |        payload      |
|----------|---------------------|
|    1     | 48 byte BLS PubKey  |
```

**String**

```
|------------|----------|---------------------|----------|
|  network   | protocol |      payload        | checksum |
|------------|----------|---------------------|----------|
| 'f' or 't' |    '3'   |  48 byte BLS PubKey |  4 bytes |
                  base32[................................]
```

#### Payload

The payload represents the data specified by the protocol. All payloads except the payload of the ID protocol are [base32](https://tools.ietf.org/html/rfc4648) encoded using the lowercase alphabet when seralized to their human readable format.

#### Checksum

Filecoin checksums are calculated over the address protocol and payload using blake2b-4. Checksums are base32 encoded and only added to an address when encoding to a string. Addresses following the ID Protocol do not have a checksum.


#### Expected Methods

All implementations in Filecoin must have methods for creating, encoding, and decoding addresses in addition to checksum creation and validation. The follwing is a golang version of the Address Interface:

```go
func New(protocol byte, payload []byte) Address

type Address interface {
	Encode(network Network, a Adress) string
	Decode(s string) Address
	Checksum(a Address) []byte
	ValidateChecksum(a Address) bool
}
```
##### New()

New returns an Address for the specified protocol encapsulating corresponding payload. New fails for unknown protocols.

```go
func New(protocol byte, payload []byte) Address {
	if protocol < SECP256K1 || protocol > BLS {
		Fatal(ErrUnknownType)
	}
	return Address{
		Protocol: protocol,
		Payload:  payload,
	}
}
```

##### Encode()

Software encoding a Filecoin address must:

- produce an address encoded to a known network
- produce an address encoded to a known protocol
- produce an address with a valid checksum

Encodes an Address as a string, prepending the network prefix, calculating the checksum, and encoding the payload and checksum to [base32](https://tools.ietf.org/html/rfc4648).

```go
func Encode(network string, a Address) string {
	if network != "f" && network != "t" {
		Fatal("Invalid Network")
	}

	switch a.Protocol {
		case SECP256K1, Actor, BLS:
			cksm := Checksum(a)
			return network + a.Protocol + base32.Encode(a.Payload + cksm)
		case ID:
			return network + a.Protocol + base10.Encode(a.Payload)
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

Decode an Address from a string by removing the network prefix, validating the address is of a know protocol, decoding the payload and checksum, and validating the checksum.

```go
func DecodeString(a string) Address {
	if len(a) < 3 {
		Fatal(ErrInvalidLength)
	}

	if a[0] != "f" && a[0] != "t" {
		Fatal(ErrUnknownNetwork)
	}

	protocol := a[1]
	raw := a[2:]
	if protocol == ID {
		return Address{
			Protocol: protocol,
			Payload:  base10.Decode(raw),
		}
	}

	payload = raw[:len(raw)-CksmLen]
	if protocol == SECP256K1 || protocol == Actor {
		if len(payload) != 20 {
			Fatal(ErrInvalidBytes)
		}
	}

	cksm := base32.Decode(payload[len(payload)-CksmLen:])
	if !ValidateChecksum(a, cksm) {
		Fatal(ErrInvalidChecksum)
	}

	return Address{
		Protocol: protocol,
		Payload:  base32.Decode(payload),
	}
}
```

##### Checksum()

Checksum produces a byte array by taking the blake2b-4 hash of an address protocol and payload.

```go

func Checksum(a Address) [4]byte {
	blake2b4(a.Protocol + a.Payload)
}
```

##### ValidateChecksum()

ValidateChecksum returns true if the Checksum of data matches the expected checksum.

```go
func ValidateChecksum(data, expected []byte) bool {
	digest := Checksum(data)
	return digest == expected
}
```
