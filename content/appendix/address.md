---
title: 'Address'
weight: 2
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Filecoin Address

A Filecoin address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. This address encodes information about the network to which an actor belongs, the specific type of address encoding, the address payload itself, and a checksum. The goal of this format is to provide a robust address format that is both easy to use and resistant to errors.

Note that each `ActorAddress` in the protocol contains a unique `ActorID` given to it by the `InitActor`. Throughout the protocol, actors are referenced by their ID-addresses. ID-addresses are computable from IDs (using `makeIdAddress(id)`), and vice versa (using the `AddressMap` to go from addr to ID).

Most actors have an alternative address which is not their key in the state tree, but is resolved to their ID-address during message processing.

Accounts have a public key-based address (e.g. to receive funds), and non-singleton actors have a temporary reorg-proof address.

An account actor's crypto-address (for signature verification) is found by looking up its actor state, keyed by the canonical ID-address. There is no map from ID-address to pubkey address.

The reference implementation of the Filecoin Address can be found in the [`go-address` Github repository](https://github.com/filecoin-project/go-address).

## Design criteria

1. **Identifiable**: The address must be easily identifiable as a Filecoin address.
2. **Reliable**: Addresses must provide a mechanism for error detection when they might be transmitted outside the network.
3. **Upgradable**: Addresses must be versioned to permit the introduction of new address formats.
4. **Compact**: Given the above constraints, addresses must be as short as possible.

## Specification

There are 2 ways a filecoin address can be represented. An address appearing on chain will always be formatted as raw bytes. An address may also be encoded to a string, this encoding includes a checksum and network prefix. An address encoded as a string will never appear on chain, this format is used for sharing among humans.

### Bytes

When represented as bytes a filecoin address contains the following:

- A **protocol indicator** byte that identifies the type and version of this address.
- The **payload** used to uniquely identify the actor according to the protocol.

```text
|----------|---------|
| protocol | payload |
|----------|---------|
|  1 byte  | n bytes |
```

### String

When encoded to a string a filecoin address contains the following:

- A **network prefix** character that identifies the network the address belongs to.
- A **protocol indicator** byte that identifies the type and version of this address.
- A **payload** used to uniquely identify the actor according to the protocol.
- A **checksum** used to validate the address.

```text
|------------|----------|---------|----------|
|  network   | protocol | payload | checksum |
|------------|----------|---------|----------|
| 'f' or 't' |  1 byte  | n bytes | 4 bytes  |
```

### Network Prefix

The **network prefix** is prepended to an address when encoding to a string. The network prefix indicates which network an address belongs to. The network prefix may either be `f` for filecoin mainnet or `t` for filecoin testnet. It is worth noting that a network prefix will never appear on chain and is only used when encoding an address to a human readable format.

### Protocol Indicator

The **protocol indicator** byte describes how a method should interpret the information in the payload field of an address. Any deviation for the algorithms and data types specified by the protocol must be assigned a new protocol number. In this way, protocols also act as versions.

- `0` : ID
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

#### Protocol 0: IDs

**Protocol 0** addresses are simple IDs. All actors have a numeric ID even if they don't have public keys. The payload of an ID address is base10 encoded. IDs are not hashed and do not have a checksum.

**Bytes**

```text
|----------|---------------|
| protocol |    payload    |
|----------|---------------|
|    0     | leb128-varint |
```

**String**

```text
|------------|----------|---------------|
|  network   | protocol |    payload    |
|------------|----------|---------------|
| 'f' or 't' |    '0'   | leb128-varint |
                  base10[...............]
```

#### Protocol 1: libsecpk1 Elliptic Curve Public Keys

**Protocol 1** addresses represent secp256k1 public encryption keys. The payload field contains the [Blake2b 160](https://blake2.net/) hash of the **uncompressed** public key (65 bytes).

**Bytes**

```text
|----------|----------------------------------|
| protocol |               payload            |
|----------|----------------------------------|
|    1     | blake2b-160( PubKey [65 bytes] ) |
```

**String**

```text
|------------|----------|--------------------------------|----------|
|  network   | protocol |      payload                   | checksum |
|------------|----------|--------------------------------|----------|
| 'f' or 't' |    '1'   | blake2b-160(PubKey [65 bytes]) |  4 bytes |
                  base32[...........................................]
```

#### Protocol 2: Actor

**Protocol 2** addresses representing an Actor. The payload field contains the SHA256 hash of meaningful data produced as a result of creating the actor.

**Bytes**

```text
|----------|---------------------|
| protocol |        payload      |
|----------|---------------------|
|    2     | blake2b-160(Random) |
```

**String**

```text
|------------|----------|-----------------------|----------|
|  network   | protocol |         payload       | checksum |
|------------|----------|-----------------------|----------|
| 'f' or 't' |    '2'   |  blake2b-160(Random)  |  4 bytes |
                  base32[..................................]
```

#### Protocol 3: BLS

**Protocol 3** addresses represent BLS public encryption keys. The payload field contains the BLS public key.

**Bytes**

```text
|----------|---------------------|
| protocol |        payload      |
|----------|---------------------|
|    3     | 48 byte BLS PubKey  |
```

**String**

```text
|------------|----------|---------------------|----------|
|  network   | protocol |      payload        | checksum |
|------------|----------|---------------------|----------|
| 'f' or 't' |    '3'   |  48 byte BLS PubKey |  4 bytes |
                  base32[................................]
```

### Payload

The payload represents the data specified by the protocol. All payloads except the payload of the ID protocol are [base32](https://tools.ietf.org/html/rfc4648#section-6) encoded using the lowercase alphabet when seralized to their human readable format.

### Checksum

Filecoin checksums are calculated over the address protocol and payload using blake2b-4. Checksums are base32 encoded and only added to an address when encoding to a string. Addresses following the ID Protocol do not have a checksum.

### Expected Methods

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

#### New()

`New()` returns an Address for the specified protocol encapsulating corresponding payload. New fails for unknown protocols.

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

#### Encode()

Software encoding a Filecoin address must:

- produce an address encoded to a known network
- produce an address encoded to a known protocol
- produce an address with a valid checksum (if applicable)

Encodes an Address as a string, prepending the network prefix, calculating the checksum, and encoding the payload and checksum to [base32](https://tools.ietf.org/html/rfc4648).

```go
func Encode(network string, a Address) string {
	if network != "f" && network != "t" {
		Fatal("Invalid Network")
	}

	switch a.Protocol {
	case SECP256K1, Actor, BLS:
		cksm := Checksum(a)
		return network + string(a.Protocol) + base32.Encode(a.Payload+cksm)
	case ID:
		return network + string(a.Protocol) + base10.Encode(leb128.Decode(a.Payload))
	default:
		Fatal("invalid address protocol")
	}
}
```

#### Decode()

Software decoding a Filecoin address must:

- verify the network is a known network.
- verify the protocol is a number of a known protocol.
- verify the checksum is valid

Decode an Address from a string by removing the network prefix, validating the address is of a known protocol, decoding the payload and checksum, and validating the checksum.

```go
func Decode(a string) Address {
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
			Payload:  leb128.Encode(base10.Decode(raw)),
		}
	}

	raw = base32.Decode(raw)
	payload = raw[:len(raw)-CksmLen]
	if protocol == SECP256K1 || protocol == Actor {
		if len(payload) != 20 {
			Fatal(ErrInvalidBytes)
		}
	}

	cksm := payload[len(payload)-CksmLen:]
	if !ValidateChecksum(a, cksm) {
		Fatal(ErrInvalidChecksum)
	}

	return Address{
		Protocol: protocol,
		Payload:  payload,
	}
}
```

#### Checksum()

`Checksum` produces a byte array by taking the blake2b-4 hash of an address protocol and payload.

```go

func Checksum(a Address) [4]byte {
	blake2b4(a.Protocol + a.Payload)
}
```

#### ValidateChecksum()

`ValidateChecksum` returns true if the `Checksum` of data matches the expected checksum.

```go
func ValidateChecksum(data, expected []byte) bool {
	digest := Checksum(data)
	return digest == expected
}
```

## Test Vectors

These are a set of test vectors that can be used to test an implementation of
this address spec. Test vectors are presented as newline-delimited address/hex
fields. The 'address' field, when parsed, should produce raw bytes that match
the corresponding item in the 'hex' field. For example:

```text
address1
hex1

address2
hex2
```

### ID Type Addresses

```text
f00
0000

f0150
009601

f01024
008008

f01729
00c10d

f018446744073709551615
00ffffffffffffffffff01
```

### Secp256k1 Type Addresses

```text
f17uoq6tp427uzv7fztkbsnn64iwotfrristwpryy
01fd1d0f4dfcd7e99afcb99a8326b7dc459d32c628

f1xcbgdhkgkwht3hrrnui3jdopeejsoatkzmoltqy
01b882619d46558f3d9e316d11b48dcf211327026a

f1xtwapqc6nh4si2hcwpr3656iotzmlwumogqbuaa
01bcec07c05e69f92468e2b3e3bf77c874f2c5da8c

f1wbxhu3ypkuo6eyp6hjx6davuelxaxrvwb2kuwva
01b06e7a6f0f551de261fe3a6fe182b422ee0bc6b6

f12fiakbhe2gwd5cnmrenekasyn6v5tnaxaqizq6a
01d1500504e4d1ac3e89ac891a4502586fabd9b417
```

### Actor Type Addresses

```text
f24vg6ut43yw2h2jqydgbg2xq7x6f4kub3bg6as6i
02e54dea4f9bc5b47d261819826d5e1fbf8bc5503b

f25nml2cfbljvn4goqtclhifepvfnicv6g7mfmmvq
02eb58bd08a15a6ade19d0989674148fa95a8157c6

f2nuqrg7vuysaue2pistjjnt3fadsdzvyuatqtfei
026d21137eb4c4814269e894d296cf6500e43cd714

f24dd4ox4c2vpf5vk5wkadgyyn6qtuvgcpxxon64a
02e0c7c75f82d55e5ed55db28033630df4274a984f

f2gfvuyh7v2sx3patm5k23wdzmhyhtmqctasbr23y
02316b4c1ff5d4afb7826ceab5bb0f2c3e0f364053
```

### BLS Type Addresses

To aid in readability, these addresses are line-wrapped. Address and hex pairs
are separated by `---`.

```text
f3vvmn62lofvhjd2ugzca6sof2j2ubwok6cj4xxbfzz
4yuxfkgobpihhd2thlanmsh3w2ptld2gqkn2jvlss4a
---
03ad58df696e2d4e91ea86c881e938ba4ea81b395e12
797b84b9cf314b9546705e839c7a99d606b247ddb4f9
ac7a3414dd

f3wmuu6crofhqmm3v4enos73okk2l366ck6yc4owxwb
dtkmpk42ohkqxfitcpa57pjdcftql4tojda2poeruwa
---
03b3294f0a2e29e0c66ebc235d2fedca5697bf784af
605c75af608e6a63d5cd38ea85ca8989e0efde9188b
382f9372460d

f3s2q2hzhkpiknjgmf4zq3ejab2rh62qbndueslmsdz
ervrhapxr7dftie4kpnpdiv2n6tvkr743ndhrsw6d3a
---
0396a1a3e4ea7a14d49985e661b22401d44fed402d1
d0925b243c923589c0fbc7e32cd04e29ed78d15d37d
3aaa3fe6da33

f3q22fijmmlckhl56rn5nkyamkph3mcfu5ed6dheq53
c244hfmnq2i7efdma3cj5voxenwiummf2ajlsbxc65a
---
0386b454258c589475f7d16f5aac018a79f6c1169d2
0fc33921dd8b5ce1cac6c348f90a3603624f6aeb91b
64518c2e8095

f3u5zgwa4ael3vuocgc5mfgygo4yuqocrntuuhcklf4
xzg5tcaqwbyfabxetwtj4tsam3pbhnwghyhijr5mixa
---
03a7726b038022f75a384617585360cee629070a2d9
d28712965e5f26ecc40858382803724ed34f2720336
f09db631f074
```
