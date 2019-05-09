# Style

## Wording
Any content that is written with `code ticks` has a specific definition to Filecoin and is defined in [the glossary](definitions.md).

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.

## Code Blocks

Many sections of the spec use go type notation to describe the functionality of certain components. This is entirely a style preference by the authors and does not imply in any way that one must use go to implement Filecoin. 

## Hints

Throughout this document you will find hints as follows. These _are not part_ of spec itself, but are merely helpers and hints for the authors, implementers or readers, with the meanings as follows. If an implementation doesn't adhere to what is said in these hints it may still be fully implementing the spec.

{% hint style="info" %}
This is an informational note, highlighting things you might wouldn't have noticed otherwise and potentially give you a hint regarding implementation following that.
{% endhint %}


{% hint style="tip" %}
This is a tip for implementation: We recommend you implement it this way.
{% endhint %}

{% hint style="working" %}
We are still working on this part, giving hints what is missing and how you could help us fix this.
{% endhint %}


{% hint style="danger" %}
This is a very important (and rarely used) alert, with information you **still must be aware of**, like potential attack vectors to keep in mind when implementing or security issues that surfaced after this version of the spec was finalised.
{% endhint %}