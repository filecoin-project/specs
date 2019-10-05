
# Account Actor

- **Code Cid**: `<codec:raw><mhType:identity><"account">`

The Account actor is the actor used for normal keypair backed accounts on the filecoin network.

```sh
type AccountActorState struct {
    address Address
}
```

## Methods

| Name | Method ID |
|--------|-------------|
| `AccountConstructor` | 1 |
| `GetAddress` | 2 |

```
type AccountConstructor struct {
}
```

## `GetAddress`

**Parameters**

```sh
type GetAddress struct {
} representation tuple
```

**Algorithm**

```go
func GetAddress() Address {
  return self.address
}
```
