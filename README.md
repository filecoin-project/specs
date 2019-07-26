# The Filecoin Spec

This repo contains the documents that comprise the Filecoin spec.

Every document in the top level of the repo is part of the official spec, and
is canon. Documents in the 'drafts' folder are work-in-progress draft documents
that aren't yet accepted as part of the spec, but exist here for discussion.
Documents in the notes repo are various notes from different meetings and
discussions.

### Viewing the spec

*Recommended:* You can view the spec [here](https://filecoin-project.github.io/specs).

You can also view it locally by using [hugo](https://gohugo.io/).

```
> git submodule update --init --recursive
> hugo serve
```


If you're just browsing on GitHub, start with [INTRO.md](INTRO.md). But really, we recommend using
the rendered output. It is much easier to read and use.

## Updates process for specs

For info on how this spec changes, please see [the process doc](process.md).

## Questions on the spec?

Issues are a great way to ask these questions. In general, your issue is much more likely to get a response
if you tag an interested party in your question. Some folks you may consider tagging (based on subject):
- [@whyrusleeping](https://github.com/whyrusleeping) - node behavior, storage market, networking behavior, protocol stewardship (upgrading, versioning, governance, etc)
- [@dignifiedquire](https://github.com/dignifiedquire) - PoSTs, proofs, data structures
- [@nicola](https://github.com/nicola) - PoSTs, proofs
- [@pooja](https://github.com/pooja) - protocol stewardship, project status
- [@henri](https://github.com/sternhenri) - chain sync, consensus

## Owners/ Points of Contact

- [@whyrusleeping](https://github.com/whyrusleeping)
- [@dignifiedquire](https://github.com/dignifiedquire)
