import renderMathInElement from 'katex/dist/contrib/auto-render.mjs'

function renderKatex(target) {
  renderMathInElement(target, {
    ignoredTags: ['script', 'noscript', 'style', 'textarea'],
    throwOnError: false,
    delimiters: [
      { left: '$$', right: '$$', display: true },
      { left: '$', right: '$', display: false },
      { left: '\\(', right: '\\)', display: false },
      { left: '\\[', right: '\\]', display: true },
    ],
  })
}

export { renderKatex }
