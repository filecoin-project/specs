## Install

```sh
git clone https://github.com/filecoin-project/specs.git
yarn install
```

## Develop
You need to have Go installed https://golang.org/doc/install

```bash
brew install go
```

### Serve with Live Reload
```sh
yarn serve
# open http://localhost:1313/ in the browser
```
# Shortcodes
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

### `svg`
This shortcode includes zoom and pad features.
```html
<!-- Relative path -->
{{< svg src="pull-flow.mmd.svg" title="Data Transfer - Pull Flow" >}}

<!-- From hugo content folder -->
{{< svg src="/systems/pull-flow.mmd.svg" title="Data Transfer - Pull Flow" >}}
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
### `figure`
```md
{{< figure src="diagrams/pieces.png" title="Pieces, Proving Trees, and Piece Data Structures" zoom="true">}}
```

### `embed`
```md
# src relative to the page
{{<embed src="piece_store.id" lang="go">}}

# src relative to content folder
{{<embed src="/systems/piece_store.id" lang="go">}}
```


# Frontmatter
```md
title: Libraries
description: Libraries used from Filecoin
weight: 3
bookCollapseSection: true
bookhidden: true
dashboardAudit: 1
dashboardState: wip
dashboardInterface: stable
```

# Code fences

They should **always** have a lang, if you don't know or don't care just use `text`

```text

```text
Random plain text context ...
``

```
# Document header
The first heading should be # Head with `---` like below and should refer to the overall title of the document. The right nav **only** starts on the second level of headings. 

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


# References
## Markdown links **(Recommended)**
These links use "portable links" just like `relref` so you can just give it the name of the file and it will fetch the correct relative link and title for the `<a href="/relative/path" title="page title">` automatically.
You can override the `<a>` title by passing a second `string` in the link definition.

**Note**: When using anchors the title can't be fetched automatically.
```md
[Storage Power](storage_power_consensus)
# <a href="/systems/filecoin_blockchain/storage_power_consensus" title="Storage Power Consensus">Storage Power</a>

[Storage Power](storage_power_consensus "Title to override the page original title")
# <a href="/systems/filecoin_blockchain/storage_power_consensus" title="Title to override the page original title">Storage Power</a>

[Tickets](storage_power_consensus#the-ticket-chain-and-drawing-randomness "The Ticket chain and drawing randomness")
# <a href="/systems/filecoin_blockchain/storage_power_consensus#the-ticket-chain-and-drawing-randomness" title="The Ticket chain and drawing randomness">Tickets</a>

```

## Hugo Cross Refs
Check Hugo's documentation [here](https://gohugo.io/content-management/shortcodes/#ref-and-relref)
```md
[Random]({{<relref "randomness">}})
[Pledge Collateral]({{<relref "storage_power_actor#pledge-collateral">}})
```
## Link shortcode
Theres also `link` shortcode which will fetch the title of the page automatically and use it for the `<a>` text and title, but **DOES NOT** work with anchors (`#anchor-id`)
```md
{{<link storage_power_consensus>}}
# <a href="/systems/filecoin_blockchain/storage_power_consensus" title="Storage Power Consensus">Storage Power Consensus</a>
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