# Address

A Filecoin address is an identifier that refers to an actor in the Filecoin state. All actors (miner actors, the storage market actor, account actors) have an address. This address encodes information about the network to which an actor belongs, the specific type of address encoding, the address data itself, and a checksum. The goal of this format is to provide a robust address format that is both easy to use and resistant to errors. Much of this specification was extended from lessons learned in both Bitcoin and Ethereum.

#### Design criteria

1. **Identifiable**: The address must be easily identifiable as a Filecoin address.
2. **Reliable**: Addresses must provide a mechanism for error detection when they might be transmitted outside the network.
3. **Upgradable**: Addresses must be versioned to permit the introduction of new address formats.
4. **Compact**: Given the above constraints, addresses must be as short as possible.

#### Design Considerations

Learning from the Bitcoin:

* Base58 needs a lot of space in QR codes, as it cannot use the ''alphanumeric mode''.
* The mixed case in base58 makes it inconvenient to reliably write down, type on mobile keyboards, or read out loud.
* Most of the research on error-detecting codes only applies to character-set sizes that are a [prime power](https://en.wikipedia.org/wiki/Prime_power), which 58 is not.
* Base58 decoding is complicated and relatively slow.

Learning from the Ethereum:

* Hex encoding is inefficient with respect to storage. Representing addresses as Hex increases the address size by 100%.
* Checksums require mixed cases, resulting in the same issues as above.

## Specification

A Filecoin address is a byte string consisting of:
* The **letter 'f'** to identify the address as belonging to a Filecoin network.
* A **network indicator** that identifies which network this address belongs to.
* A **protocol indicator** that identifies the type and version of this address.
* The **data** used to uniquely identify the actor according to the protocol.

```
|---|---------|----------|------|
| f | network | protocol | data |
|---|---------|----------|------|
```

An example of a Filecoin address in golang,

```go
type Address struct {
    // 0: main-net
    // 1: test-net
    network varint
    
    // 0: Blake2b160-Hash of SECP256K1 Public Key
    // 1: ActorID
    // 2: BLS Public Key -> in bytes?
    protocol varint
    
    // raw bytes containing the data associated with typ
    data []byte
}
```

#### Network Indicator

The **network** is a `var32encoded` integer that identifies the network:

- `0` :  Mainnet
- `1` : Testnet 1
- `2` : Testnet 2
- `3` : etc.

The size of this integer is `4` bytes. An example description in golang,

> Note(bvo): Double check that this is the integer 4 bytes?

```go
// String Prefix
const (
    MAINNET_PREFIX = "fc"
    TESTNET_PREFIX = "tf"
)

// Network Byte
const (
	Mainnet = 0
	Testnet = 1
)
```

#### Protocol Indicator

The **protocol indicator** is a *var32encoded* integer that describes how a method should interpret the information in the `data` field of an address. This indicator  specifies the hashing algorithms, error correcting protocols, and the origin of the hashed data (e.g. a public key) when needed. Any deviation for the algorithms and data types specified by the protocol must be assigned a new protocol number. In this way, protocols also act as versions.

```go
// Type varint
const (
	SECP256K1 = 0
	ID = 1
	Actor = 2
	BLS = 3
)
```

###### Protocol 0: libsecpk1 Elliptic Curve Public Keys

**Protocol 0** addresses represent public encryption keys. These addresses represent Filecoin user accounts and are likely to be transmitted outside the network over noisy channels (e.g. hand written on a napkin). As such, they require robust error detection. The multi-base encoded key contains a hash and a checksum. The hash is a [Blake2b 160](https://blake2.net/) hash of the public key. The checksum is a six character BCH checksum. The algorithm details can be found in Appendix - 9.2. Algorithms - Checksums.

###### Protocol 1: IDs

**Protocol 1** addresses are simple ids. These addresses are expected to remain within the network. All actors have a numeric ID even if they don't have public keys. The key part of an ID address is the decimal representation of the id only characters '0'-'9'.

###### Protocol 2: BLS Public Keys

**TODO**

#### Data

**TODO** Encoding inside the byte array. Add data serialization for different protocols

- `varint` is an LEB128 encoded varint


#### Expected Methods

All implementations in Filecoin must have methdos for creating, serializaing, and encoding addresses. The follwing is a golang version of the Address Interface:

```go
type IAddress interface {
    New(n varint, p varint, data []byte) IAddress
    Marshal(a Addressi) []byte
    Unmarshal(data []byte) IAddress
    EncodeString(a IAddress) string
    DecodeString(s string) IAddress
    Checksum(a IAddress) []byte
    validateChecksum(data, cksm []byte) bool
}
```
> TODO(@frrist): add exposition about these methods

##### New()

```go
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
```

##### Marshal()

```go
// Marshal an Address to bytes
func Marshal(a Address) []byte {
    if a.typ == ID {
        return []byte{a.network, a.typ, a.data}
    }
    return []byte{a.network, a.typ, a.data, Checksum(a)}
}
```

##### Unmarshal()

```go
// Unmarshal bytes to an Address
func Unmarshal(a []byte) Address {
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
```

##### EncodeString()

```go
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
```

##### DecodeString()

Software interpreting a Filecoin address:
* MUST verify that the address begins with `f`.
* MUST verify that the network is a known network.
* MUST verify that the protocol is a number of a known protocol.

For public key addresses, software must:
* MUST verify the address is multibase encoded with a known encoding.
* MUST verify that the hash is 160 bits.
* MUST verify the checksum is valid (see section 9.2).

```go
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
```

##### Checksum()

```go
// checksum size
const CKSM_LEN = 4

// Checksum is the last CKSM_LEN bytes the sha256 of Address Type and its data.
func Checksum(a Address) [CKSM_LEN]byte {
    digest = sha256([]byte{a.typ, a.data})
    return digest[len(digest)-CKSM_LEN:]
}

```

**Uppercase/lowercase**

The lowercase form is used when determining a character's value for checksum purposes.

Encoders MUST always output an all lowercase Base32Checked string.
If an uppercase version of the encoding result is desired, (e.g.- for presentation purposes, or QR code use),
then an uppercasing procedure can be performed external to the encoding process.

Decoders MUST NOT accept strings where some characters are uppercase and some are lowercase (such strings are referred to as mixed case strings).

For presentation, lowercase is usually preferable, but inside QR codes uppercase SHOULD be used, as those permit the use of
[`alphanumeric mode`](http://www.thonky.com/qr-code-tutorial/alphanumeric-mode-encoding), which is 45% more compact than the normal
[`byte mode`](http://www.thonky.com/qr-code-tutorial/byte-mode-encoding).



##### validateChecksum()

```go
// validChecksum returns true if cksm is valid for data and false otherwise.
func validChecksum(data, cksm []byte) bool {
    digest := Checksum(data)
    return digest == cksm
}
```

### Checksums

> TODO(@frrist) - let's make sure we look into this. I think the output of a BCH operation is encoded data without a checksum but I could be wrong; it's been a while. Double check?

Filecoin checksums are a [BCH code](https://en.wikipedia.org/wiki/BCH_code) that
guarantees detection of **any error affecting at most 4 characters**
and has less than a 1 in 10<sup>9</sup> chance of failing to detect more
errors. More details about the properties can be found in the
Checksum Design appendix. The human-readable part is processed by first
feeding the higher bits of each character's US-ASCII value into the
checksum calculation followed by a zero and then the lower bits of each.

Valid strings MUST pass the criteria for validity specified by the Python3 code snippet below. The function,

`base32checked_verify_checksum` must return true when its arguments are:
* `hrp`: the human-readable part as a string
* `protocol`: the protocol as an integer
* `key`: the decoded bits of the encoded part (minus the checksum) divided into 5 bit characters.

### Error Correction Codes

One of the properties of these BCH codes is that they can be used for
error correction. An unfortunate side effect of error correction is that it erodes error detection: correction changes invalid inputs into valid inputs, but if more than a few errors were made then the valid input may not be the correct input. Use of an incorrect but valid input can cause funds to be lost irrecoverably. Because of this, implementations SHOULD NOT implement correction beyond potentially suggesting to the user where in the string an error might be found, without suggesting the correction to make.

#### Error Correction Design Choices
BCH codes can be constructed over any prime-power alphabet and can be chosen to have a good trade-off between size and error-detection capabilities. While most work around BCH codes uses a binary alphabet, that is not a requirement.
This makes them more appropriate for our use case than [CRC codes](https://en.wikipedia.org/wiki/Cyclic_redundancy_check). Unlike
[Reed-Solomon codes](https://en.wikipedia.org/wiki/Reed%E2%80%93Solomon_error_correction), they are not restricted in length to one less than the alphabet size. While they also support efficient error correction, the implementation of just error detection is very simple.
We pick 6 checksum characters as a trade-off between length of the addresses and the error-detection capabilities, as 6 characters is the lowest number sufficient for a random failure chance below 1 per billion. For the length of data we're interested in protecting (32 bytes), BCH codes can be constructed that guarantee detecting up to 4 errors.
**TODO: verify that we can't drop the character count for 32bytes.**


#### Selected Properties
Many of these codes perform badly when dealing with more errors than they are designed to detect, but not all.
For that reason, we consider codes that are designed to detect only 3 errors as well as 4 errors,
and analyse how well they perform in practice.
The specific code chosen here is the result
of:
* Starting with an exhaustive list of 159605 BCH codes designed to detect 3 or 4 errors up to length 93, 151, 165, 341, 1023, and 1057.
* From those, requiring the detection of 4 errors up to length 71, resulting in 28825 remaining codes.
* From those, choosing the codes with the best worst-case window for 5-character errors, resulting in 310 remaining codes.
* From those, picking the code with the lowest chance for not detecting small numbers of ''bit'' errors.
As a naive search would require over 6.5 * 10<sup>19</sup> checksum evaluations, a collision-search approach was used for
analysis. The code can be found [here](https://github.com/sipa/ezbase32/).

### Var32encoded integers

Var32encoded integers map integers 0-9 to '0'-'9', and permit much larger numbers. These numbers are encoded using the [base32hex alphabet](https://tools.ietf.org/html/rfc4648#section-6). The rationale and mechanism for encoding and decoding these numbers are described in the algorithms section. The biggest advantage of this scheme is that numbers less than 10 encode to a single digit character.

### Var32encoded Integer Encoding
Address network and protocol need a space efficient, human readable, and future proof numbering scheme that is safe for transmission. Specifically the scheme must:
1. permit arbitrarily large numbers (up to some astronomic bound) to future proof the format
2. be self delimiting, so that there will be no confusion about where the numbers end
3. be compact, so that small numbers can be stored in a single char and no number needs more than one extra char than its minimal representation in the character encoding.
4. encode in a safe character encoding that is easily transmissible.
5. encode small numbers in a human readable fashion. Ideally this is '0', '1', ... '9'.
To satisfy these constraints, we introduce var32encoded integers. Base32Hex encoding lets us map 0-9 to '0'-'9', so small, common numbers are very readable. The encoding is standard and safe. The numeric encoding is inspired by [SQLite's Varints](https://sqlite.org/src4/doc/trunk/www/varint.wiki), but modified from bytes to 5 bit chars to fit our character encoding directly.
Specifically for 5 bit characters A0,A1... will be decoded as follows:
``` python
A0 < 10: A0
A0 < 19: 10 + 32*(A0-10) + A1
A0 < 20: 298 + 32*A1 + A2
# A0 specifies number of chars. We decode A1..AN as a base32hex encoded number
else   : sum(A[i+1]*32**(a0-18-i) for i in range(A0-17))
```
The numbers can be encoded as follows:
``` python
n <   10: base32hex(n)
n <  298: base32hex((n-10)//32 + 10) + base32hex((n-10)%32)
n < 1322: 'k' + base32hex((n-298)//32) + base32hex((n-298)%32)
else    : base32hex(ceil(n.bit_length()/5)+17) + base32hex(n)
```
The encoded numbers look like this:
```
0: 0
1: 1
...
9: 9
10: a0
11: a1
12: a2
...
42: b0
43: b1
...
296: iy
297: iz
298: j00
299: j01
...
1321: jzz
1322: k19a
1323: k10b
32767: kzzz
32768: l1000
1048575: lzzzz
1048576: m10000
...
2^70-1: vvvvvvvvvvvvvvv
```

## Acknowledgements:

- The address format is inspired by [BIP173](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki).

--- 


