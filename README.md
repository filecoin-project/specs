
# Filecoin Specification
![CI](https://github.com/filecoin-project/specs/workflows/CI/badge.svg)

This is the [Filecoin Specification](https://github.com/filecoin-project/specs), a repository that contains documents, code, models, and diagrams that constitute the specification of the [Filecoin Protocol](https://filecoin.io). This repository is the singular source of truth for the Filecoin Protocol. All implementations of the Filecoin Protocol should match and comply with the descriptions, interfaces, code, and models defined in this specification.

Note that the `beta` branch of the specs moves quickly. We work to merge PRs as fast as possible into master, which means changes or reversals are possible here. Accordingly, we periodically compile swaths of spec along with a high-level difflog into the `release` branch. As the spec stabilizes, this practice will change.

## Website

https://beta.spec.filecoin.io is the user-friendly website rendering, which we recommend for reading this repository. The website is updated automatically with every merge to `beta`.

## Previous version
You can find the previous version of the Filecoin Spec in the master branch [here](https://github.com/filecoin-project/specs/tree/master).

## Install

```sh
git clone https://github.com/filecoin-project/specs.git
yarn install
```

## Develop
You need to have [Go](https://golang.org/doc/install) and `bzr` installed (required to build lotus)

```bash
brew install go bzr
```

### Serve with Live Reload
```sh
yarn start
# open http://localhost:1313/ in the browser
```

### Solving Common problems

**Problem** - Site fails to build with an error that states it faled to download modules on macos

```
Error: failed to download modules: go command failed ...
```

**Solution** - run `npm run clean-modules` - the cache dir hugo uses can get corrupted, and this resets it. See [#1048](https://github.com/filecoin-project/specs/issues/1048)


### External modules
External modules should be added as [Hugo Modules](https://gohugo.io/hugo-modules/)
You can find examples in the `config.toml`

```toml
[module]
  [[module.imports]]
    path = "github.com/filecoin-project/specs-actors"
    [[module.imports.mounts]]
    source = "."
    target = "content/externals/specs-actors"
```
> `target` should **ALWAYS** use the folder `content/externals`

This makes files from external repos available for Hugo rendering and allows for linking to up-to-date files that are directly pulled from other repositories.

The configuration above gives the following information:

- `path`: Repository's URL without protocol.
- `source`: Folder from the repository referenced in the `path` to be mounted into the local Hugo filesystem.
- `target`: Folder where `source` will be mounted locally, this should follow this structure `content/modules/<target value>`.

Example: if you want to link/embed to the file `xyz.go` that lives in `https://github.com/filecoin-project/specs-actors/actors/xyz-folder/xyz.go`, from within a markdown file then with the above configuration the `src` for shortcodes or markdown image syntax would be:

```
{{<embed src="/externals/specs-actors/actors/xyz-folder/xyz.go"  lang="go">}}
```
> The first foward slash is important it means the path is relative to the content folder.

These modules can be updated with 

```sh
hugo mod get -u
```
or to a specific version with

```sh
hugo mod get github.com/filecoin-project/specs-actors@v0.7.2
```

## Page Header
The first heading should be `# Page Title` with `---` like below and should refer to the overall title of the document.

```md
---
title: Storage Power Actor
---

# Storage Power Actor
---

## Header for a section in this document
Some text

### Sub header for the a nested section

## Another top level header
```

## Frontmatter

Description for all the available frontmatter properties

```md
<!-- Page Title to be used in the navigation -->
title: Libraries 
<!-- Small description for html metadata, if not present the first couple of paragraphs will be used instead -->
description: Libraries used from Filecoin
<!-- This will be used to order the ToC, navigation and any other listings of pages -->
weight: 3
<!-- This will make a page section collapse in the navigation -->
bookCollapseSection: true
<!-- This will hidden the page from the navigation -->
bookhidden: true
<!-- This is used in the dashboard to describe the importance of the page content -->
dashboardWeight: 2
<!-- This is used in the dashboard to describe the state of the page content options are "missing", "incorrect", "wip", "reliable", "stable" or "n/a" -->
dashboardState: stable
<!-- This is used in the dashboard to describe if the theory of the page has been audited, options are "missing", "wip", "done" or "n/a" -->
dashboardAudit: wip
<!-- When dashboardAudit is stable we should have a report url -->
dashboardAuditURL: https://url.to.the.report
<!-- The date that the report at dashboardAuditURL was completed -->
dashboardAuditDate: "2020-08-01"
<!-- This is used in the dashboard to describe if the page content has compliance tests, options are 0 or numbers of tests -->
dashboardTests: 0
```

## Code fences

Code fences should **always** have a lang, if you don't know or don't care just use `text`

```text

```text
Random plain text context ...
``

```

## Images, diagrams - Markdown images
To add an image or diagram you just need to use normal markdown syntax. 
For diagrams you can use the source files and the pipelines will handle converting that to `svg`, we support mermaid and dot source files.

```md
# relative to the markdown file
![Alt text](picture.jpg)

# relative to the content folder
![Alt text](/content/intro/diagram1.mmd)

![Alt text](graph.dot "Graph title")
```
> When there's no title we use the alt text as title   


## References - Markdown links
These links use "portable links" just like `relref` so you can just give it the name of the file and it will fetch the correct relative link and title for the `<a href="/relative/path" title="page title">` automatically.
You can override the `<a>` title by passing a second `string` in the link definition.

> **Note**: When using anchors the title can't be fetched automatically.
```md
[](storage_power_consensus)
# <a href="/systems/filecoin_blockchain/storage_power_consensus" title="Storage Power Consensus">Storage Power Consensus</a>

[Storage Power](storage_power_consensus "Title to override the page original title")
# <a href="/systems/filecoin_blockchain/storage_power_consensus" title="Title to override the page original title">Storage Power</a>

[Tickets](storage_power_consensus#the-ticket-chain-and-drawing-randomness "The Ticket chain and drawing randomness")
# <a href="/systems/filecoin_blockchain/storage_power_consensus#the-ticket-chain-and-drawing-randomness" title="The Ticket chain and drawing randomness">Tickets</a>

```


## Shortcodes
### `Mermaid` 
Inline mermaid syntax rendering
```html
{{< mermaid >}}
graph TD
  A[Christmas] -->|Get money| B(Go shopping)
  B --> C{Let me think}
  C -->|One| D[Laptop]
  C -->|Two| E[iPhone]
  C -->|Three| F[fa:fa-car Car]
		
{{</ mermaid >}}
```

### `hint`
```md
<!-- info|warning|danger -->
{{< hint info >}}
**Markdown content**  
Lorem markdownum insigne. Olympo signis Delphis! Retexi Nereius nova develat
stringit, frustra Saturnius uteroque inter! Oculis non ritibus Telethusa
{{< /hint >}}
```
### `embed`
```md
# src relative to the page
{{<embed src="piece_store.id" lang="go">}}

# src relative to content folder
{{<embed src="/systems/piece_store.id" lang="go">}}
```

## Math mode
For short snippets of math text you can just use the `{{<katex>}}` shortcode, but if you need to write lots of math in a page you can just use `math-mode` and avoid writting the katex shortcode everywhere.

Parses math typesetting with [KaTeX](https://katex.org/docs/api.html)   

Check this example [example](https://deploy-preview-969--fil-spec-staging.netlify.app/math-mode/)

> Some syntax like `\_` can't go through HUGO markdown parser and for that reason we need to wrap math text with code blocks, code fendes or the shortcode `{{<plain>}}`. See examples below.
### Add `math-mode` prop to the Frontmatter
```md
---
title: Math Mode
math-mode: true
---
```

### Wrap `def`, `gdef`, etc.
Math text needs to be wrapped to avoid Hugo's Markdown parser. When wrapping defs or any math block that doesn't need to be rendered the recommended option is to use the shortcode `{{<plain hidden}}` with the hidden argument.

```md
{{<plain hidden>}}
$$
\gdef\createporepbatch{\textsf{create_porep_batch}}
\gdef\GrothProof{\textsf{Groth16Proof}}
\gdef\Groth{\textsf{Groth16}}
\gdef\GrothEvaluationKey{\textsf{Groth16EvaluationKey}}
\gdef\GrothVerificationKey{\textsf{Groth16VerificationKey}}
{{</plain>}}
```

### Wrap inline math text with code blocks
```md
The index of a node in a `$\BinTree$` layer `$l$`. The leftmost node in a tree has `$\index_l = 0$`.
```

### Wrap math blocks with code fences
~~~md
```text
$\overline{\underline{\Function \BinTree\dot\createproof(c: \NodeIndex) \rightarrow \BinTreeProof_c}}$
$\line{1}{\bi}{\leaf: \Safe = \BinTree\dot\leaves[c]}$
$\line{2}{\bi}{\root: \Safe = \BinTree\dot\root}$

$\line{3}{\bi}{\path: \BinPathElement^{[\BinTreeDepth]}= [\ ]}$
$\line{4}{\bi}{\for l \in [\BinTreeDepth]:}$
$\line{5}{\bi}{\quad \index_l: [\len(\BinTree\dot\layer_l)] = c \gg l}$
$\line{6}{\bi}{\quad \missing: \Bit = \index_l \AND 1}$
$\line{7}{\bi}{\quad \sibling: \Safe = \if \missing = 0:}$
$\quad\quad\quad \BinTree\dot\layer_l[\index_l + 1]$
$\quad\quad\thin \else:$
$\quad\quad\quad \BinTree\dot\layer_l[\index_l - 1]$
$\line{8}{\bi}{\quad \path\dot\push(\BinPathElement \thin \{\ \sibling, \thin \missing\ \} \thin )}$

$\line{9}{\bi}{\return \BinTreeProof_c \thin \{\ \leaf, \thin \root, \thin \path\ \}}$
```
~~~

## References
- [hugo theme book](https://themes.gohugo.io//theme/hugo-book/docs/shortcodes/columns/)
- [Katex](https://katex.org/)
- [Mermaid](https://mermaid-js.github.io/mermaid/#/)
  - [config](https://github.com/mermaid-js/mermaid/blob/master/docs/mermaidAPI.md#mermaidapi-configuration-defaults)
  - [editor](https://mermaid-js.github.io/mermaid-live-editor)
- [Pan/Zoom for SVG](https://github.com/anvaka/panzoom)
- [Icons](https://css.gg/)
- [Working with submodules](https://github.blog/2016-02-01-working-with-submodules/)
