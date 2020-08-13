import renderMathInElement from 'katex/dist/contrib/auto-render.mjs'

function initKatex () {
  console.log('init katex')
  document.querySelectorAll('.math-mode').forEach(function (el) {
    renderMathInElement(el, {
      ignoredTags: ["script", "noscript", "style", "textarea"],
      throwOnError: false,
      delimiters: [
          {left: "$$", right: "$$", display: true},
          {left: "$", right: "$", display: false},
          {left: "\\(", right: "\\)", display: false},
          {left: "\\[", right: "\\]", display: true}
          ]
    })
  })
}

window.addEventListener('load', () => {
  setTimeout(initKatex, 2000)
});
