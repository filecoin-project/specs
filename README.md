# Filecoin Specification

![CI](https://github.com/filecoin-project/specs/workflows/CI/badge.svg)

This is the [Filecoin Specification](https://github.com/filecoin-project/specs), a repository that contains documents, code, models, and diagrams that constitute the specification of the [Filecoin Protocol](https://filecoin.io). This repository is the singular source of truth for the Filecoin Protocol. All implementations of the Filecoin Protocol should match and comply with the descriptions, interfaces, code, and models defined in this specification.

<https://spec.filecoin.io> is the user-friendly website rendering, which we recommend for reading this repository. The website is updated automatically with every merge to `master`.

## Table of Contents

-   [Install](#install)
-   [Writing the spec](#writing-the-spec)
-   [Check your markdown](#check-your-markdown)
-   [Page Template](#page-template)
-   [Code](#code)
-   [Images](#images)
-   [Links](#links)
-   [Shortcodes](#shortcodes)
    -   [`embed`](#embed)
    -   [`listing`](#listing)
    -   [`mermaid`](#mermaid)
    -   [`hint`](#hint)
-   [Math mode](#math-mode)
    -   [Wrap `def`, `gdef`, etc.](#wrap-def-gdef-etc)
    -   [Wrap inline math text with code blocks](#wrap-inline-math-text-with-code-blocks)
    -   [Wrap math blocks with code fences](#wrap-math-blocks-with-code-fences)
-   [Front-matter](#front-matter)
-   [External modules](#external-modules)
-   [Solving Common problems](#solving-common-problems)
-   [References](#references)

## Install

To build the spec website you need

-   [`node` & `npm`](https://nodejs.org/en/download)
-   [`go`](https://golang.org/doc/install)
-   `bzr` (required to build lotus)

On macOS you can get go and bzr from Homebrew

```bash
brew install go bzr
```

Clone the repo, and use `npm install` to fetch the dependencies

```sh
git clone https://github.com/filecoin-project/specs.git
npm install
```

To run it locally with live-reload in the browser when you change a file run

```sh
npm start
```

Then open http&#x3A;//localhost:1313/ in the browser

## Writing the spec

The spec is written in markdown. Each section is markdown document in the `content` directory. The first level of the directory structure denotes the top level sections of the spec; (Intrduction, Systems, etc.) The `_index.html` file in each folder is used as the starting point for each section. For example the introduction starts in `content/intro/_index.md`. 

Sections can be split out into multiple markdown documents. The build process combines them into a single html page. The sections are ordered by the `weight` front-matter property. The introduction appears at the start of the html page because `content/intro/_index.md` has `weight: 1`, while `content/systems/_index.html` has `weight: 2` so it appears as the second section. 

You can split out sub-sections by adding additional pages to a section directory. The `content/intro/concepts.md` defines the Key Concepts sub-section of the the Introduction. The order of sub-sections within a section is again controlled by setting the `weight` property. This pattern repeats for sub sub folders which represent sub sub sections.

The markdown documents should all be well formed, with a single h1, and headings should increment by a single level.

## Check your markdown

Use `npm test` to run a markdown linter set up to check for common errors. It runs in CI and you can run it locally with:

```bash
npm test
content/algorithms/crypto/randomness.md
  15:39-15:46  warning  Found reference to undefined definition  no-undefined-references  remark-lint
  54:24-54:31  warning  Found reference to undefined definition  no-undefined-references  remark-lint

⚠ 2 warnings
```

## Page Template

A spec document should start with a YAML front-matter section and contain at least a single h1, as below.

```md
---
title: Important thing
weight: 1
dashboardState: wip
dashboardAudit: missing
---

# Important thing
```

## Code

Wrap code blocks in _code fences_. Code fences should **always** have a lang. It is used to provide syntax heighlighting. Use `text` as the language flag for pseudocode for no highlighting.

````text
```text
Your algorithm here
```
````

You can embed source code from other repos. Mount the repo as a hugo modules as descibed in [External Modules](#external-modules) then use the [`embed shorcode`](#embed) to link to a specific symbol.

```go
{{<embed src="/externals/go-data-transfer/types.go"  lang="go" symbol="Channel">}}
```

## Images

Use normal markdown syntax to include images.

For `dot` and `mermaid` diagrams you link to the source file and the pipelines will handle converting that to `svg`.

```md
# relative to the markdown file
![Alt text](picture.jpg)

# relative to the content folder
![Alt text](/content/intro/diagram1.mmd)

![Alt text](graph.dot "Graph title")
```

> the alt text as title is used as the title where it is not provided.

## Links

Use markdown syntax `[text](markdown-document-name)`. 

These links use "portable links" just like `relref`. Just give it the name of the file and it will fetch the correct relative path and title automatically. You can override the title by passing a second `string` in the link definition.

> **Note**: When using anchors the title can't be fetched automatically.

```md
[](storage_power_consensus)

# Renders to
<a href="/systems/filecoin_blockchain/storage_power_consensus" title="Storage Power Consensus">Storage Power Consensus</a>


[Storage Power](storage_power_consensus "Title to override the page original title")

# Renders to
<a href="/systems/filecoin_blockchain/storage_power_consensus" title="Title to override the page original title">Storage Power</a>


[Tickets](storage_power_consensus#the-ticket-chain-and-drawing-randomness "The Ticket chain and drawing randomness")

# Renders to
<a href="/systems/filecoin_blockchain/storage_power_consensus#the-ticket-chain-and-drawing-randomness" title="The Ticket chain and drawing randomness">Tickets</a>

```

## Shortcodes

hugo shortcodes you can add to your markdown.

### `embed`

```md
# src relative to the page
{{<embed src="piece_store.id" lang="go">}}

# src relative to content folder
{{<embed src="/systems/piece_store.id" lang="go">}}

# can just embed a markdown file
{{<embed src="section.md" markdown="true">}}

# can embed symbols from Go files
# extracts comments and symbol body
{{<embed src="/externals/go-data-transfer/types.go"  lang="go" symbol="Channel">}}
```

### `listing`

The listing shortcode creates tables from externals sources, supports Go `struct`.

```md
# src relative to the page
{{<listing src="piece_store.go" symbol="Channel">}}

# src relative to content folder
{{<listing src="/systems/piece_store.go" symbol="Channel">}}

# src can also be from the externals repos
{{<listing src="/externals/go-data-transfer/types.go"  symbol="Channel">}}
```

### `mermaid`

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

## Math mode

For short snippets of math text you can just use the `{{<katex>}}` shortcode, but if you need to write lots of math in a page you can just use `math-mode` and avoid writting the katex shortcode everywhere.

Parses math typesetting with [KaTeX](https://katex.org/docs/api.html)   

Check this example [example](https://deploy-preview-969--fil-spec-staging.netlify.app/math-mode/)

> Some syntax like `\_` can't go through HUGO markdown parser and for that reason we need to wrap math text with code blocks, code fendes or the shortcode `{{<plain>}}`. See examples below.
>
> ### Add `math-mode` prop to the Frontmatter
>
> ```md
> ---
> title: Math Mode
> math-mode: true
> ---
> ```

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

````md
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
````

## Front-matter

Description for all the available frontmatter properties

```yaml
# Page Title to be used in the navigation
title: Libraries 
# Small description for html metadata, if not present the first couple of paragraphs will be used instead -->
description: Libraries used from Filecoin
# This will be used to order the ToC, navigation and any other listings of pages -->
weight: 3
# This will make a page section collapse in the navigation -->
bookCollapseSection: true
# This will hidden the page from the navigation -->
bookhidden: true
# This is used in the dashboard to describe the importance of the page content -->
dashboardWeight: 2
# This is used in the dashboard to describe the state of the page content options are "missing", "incorrect", "wip", "reliable", "stable" or "n/a" -->
dashboardState: stable
# This is used in the dashboard to describe if the theory of the page has been audited, options are "missing", "wip", "done" or "n/a" -->
dashboardAudit: wip
# When dashboardAudit is stable we should have a report url -->
dashboardAuditURL: https://url.to.the.report
# The date that the report at dashboardAuditURL was completed -->
dashboardAuditDate: "2020-08-01"
# This is used in the dashboard to describe if the page content has compliance tests, options are 0 or numbers of tests -->
dashboardTests: 0
```

## External modules

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

-   `path`: Repository's URL without protocol.
-   `source`: Folder from the repository referenced in the `path` to be mounted into the local Hugo filesystem.
-   `target`: Folder where `source` will be mounted locally, this should follow this structure `content/modules/<target value>`.

Example: if you want to link/embed to the file `xyz.go` that lives in `https://github.com/filecoin-project/specs-actors/actors/xyz-folder/xyz.go`, from within a markdown file then with the above configuration the `src` for shortcodes or markdown image syntax would be:

    {{<embed src="/externals/specs-actors/actors/xyz-folder/xyz.go"  lang="go">}}

> The first foward slash is important it means the path is relative to the content folder.

These modules can be updated with 

```sh
hugo mod get -u
```

or to a specific version with

```sh
hugo mod get github.com/filecoin-project/specs-actors@v0.7.2
```

## Solving Common problems

**Problem** - Site fails to build with an error that states it failed to download modules on macos

    Error: failed to download modules: go command failed ...

**Solution** - run `npm run clean` - the cache dir hugo uses can get corrupted, and this resets it. See [#1048](https://github.com/filecoin-project/specs/issues/1048)

## References

-   [hugo theme book](https://themes.gohugo.io//theme/hugo-book/docs/shortcodes/columns/)
-   [Katex](https://katex.org/)
-   [Mermaid](https://mermaid-js.github.io/mermaid/#/)
    -   [config](https://github.com/mermaid-js/mermaid/blob/master/docs/mermaidAPI.md#mermaidapi-configuration-defaults)
    -   [editor](https://mermaid-js.github.io/mermaid-live-editor)
-   [Pan/Zoom for SVG](https://github.com/anvaka/panzoom)
-   [Icons](https://css.gg/)
-   [Working with submodules](https://github.blog/2016-02-01-working-with-submodules/)
