---
title: "Verifiable Random Function"
---

{{<label vrf>}}

Filecoin uses the notion of a [Verifiable Random
Function](https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Pseudo%20Randomness/Verifiable_Random_Functions.pdf)
(VRF). A VRF uses a private key to produce a digest of
an arbitrary message such that the output is unique per signer and per message.
Any third party in possession of the corresponding public key, the message, and
the VRF output, can verify if the digest has been computed correctly and by the
correct signer. Using a VRF in the ticket generation process allows anyone to
verify if a block comes from an eligible block producer (see {{<sref tickets
"Ticket Generation" >}} for more details).

BLS signature can be used as the basis to construct a VRF. Filecoin transforms
the BLS signature scheme it uses (see {{<sref signatures Signatures>}} into a
VRF, Filecoin uses the random oracle model and deterministically hashes the
signature to produce the final digest. Filecoin uses
SHA256 as the hash function. The algorithm is the following:
```
VRFOutput = SHA256(DST || index || Serialization(BLSSignature(message)))
```

where `DST` is a domain separation tag in order to treat the hash
function as an independent random oracle in the VRF output (see {{<sref randomness>}}).


{{< readfile file="vrf.id" code="true" lang="go" >}}
{{< readfile file="vrf.go" code="true" lang="go" >}}