## Install

```sh
git clone https://github.com/filecoin-project/specs.git
cd next # until we move this top level
yarn install
```

## Develop

### Update submodules
```sh
git submodule update --init
```

### Serve with Live Reload
```sh
yarn serve
# open http://localhost:1313/ in the browser
```
# Shortcodes
### `Mermaid`
```html
<!-- Relative path -->
{{< mermaid file="full-deals-on-chain.mmd" />}}

<!-- From hugo content folder -->
{{< mermaid file="/intro/full-deals-on-chain.mmd" />}}

<!-- Inline -->
{{< mermaid >}}
graph TD
  A[Christmas] -->|Get money| B(Go shopping)
  B --> C{Let me think}
  C -->|One| D[Laptop]
  C -->|Two| E[iPhone]
  C -->|Three| F[fa:fa-car Car]
		
{{</ mermaid >}}
```
## References
- [hugo theme book](https://themes.gohugo.io//theme/hugo-book/docs/shortcodes/columns/)
- [Katex](https://katex.org/)
- [Mermaid](https://mermaid-js.github.io/mermaid/#/)
  - [config](https://github.com/mermaid-js/mermaid/blob/master/docs/mermaidAPI.md#mermaidapi-configuration-defaults)
  - [editor](https://mermaid-js.github.io/mermaid-live-editor)
- [Pan/Zoom for SVG](https://github.com/anvaka/panzoom)
- [Icons](https://css.gg/)
- [Working with submodules](https://github.blog/2016-02-01-working-with-submodules/)