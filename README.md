# The Filecoin Spec

This repo contains the documents that comprise the Filecoin spec.

Every document in the top level of the repo is part of the official spec, and
is canon. Documents in the 'drafts' folder are work-in-progress draft documents
that aren't yet accepted as part of the spec, but exist here for discussion.
Documents in the notes repo are various notes from different meetings and
discussions.

### Viewing the spec

*Recommended:* You can view the spec [here](https://filecoin-1.gitbook.io/spec/).

You can also use gitbook tooling to view the spec locally.

First, install `gitbook-cli` via npm then run the following to install all necessary plugins:
```
gitbook install
```

Now you can have gitbook serve it to you (with live-reload) via:
```
gitbook serve
```

Alternatively, you can use gitbook to print a pdf (or epub, or mobi) using:
```
gitbook pdf
```

(Note: on macOS using the pdf printer may require some extra installation, epub and mobi might require callibre to be installed, too)

If you're just browsing on GitHub, start with [INTRO.md](INTRO.md) and the
table of contents in [SUMMARY.md](SUMMARY.md). But really, we recommend using 
the Gitbook output. It is much easier to read and use.

## Updates process for specs
For info on how this spec changes, please see [the process doc](process.md).

## Owners/ Points of Contact
- [@whyrusleeping](https://github.com/whyrusleeping)
- [@bvohaska](https://github.com/bvohaska)
