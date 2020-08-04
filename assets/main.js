import '@pwabuilder/pwaupdate'
import toc from './toc/index.js'
import panzoom from 'panzoom'
import tablesort from 'tablesort'

// Note: the tablesort lib is not ESM friendly, and the sorts expect `Tablesort` to be available on the global
window.Tablesort = tablesort
require('tablesort/dist/sorts/tablesort.number.min.js')

function initToc () {
  console.log('init toc')
    toc.init({
        tocSelector: '.toc',
        contentSelector: '.markdown',
        headingSelector: 'h1, h2, h3, h4, h5, h6',
        hasInnerContainers: false,
        orderedList: true,
        smoothScroll: false,
        collapseDepth: 2,
        headingLabelCallback: (label) => {
          return label.replace('#', '')
        },
        headingsOffset: 1
    });
}

function initPanZoom () {
  console.log('init panzoom')
  var elements = document.querySelectorAll(".zoomable")
  elements.forEach(function (el) {
    panzoom(el.querySelector('*:first-child'), {
      maxZoom: 10,
      minZoom: 0.5
    })
  })
}

function initTableSort () {
  console.log('init tablesort')
  var elements = document.querySelectorAll(".tablesort")
  elements.forEach(function (el) {
    tablesort(el);
  })
}

initToc()
initPanZoom()
initTableSort()
