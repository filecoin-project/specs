import '@pwabuilder/pwaupdate'
import { initToc } from './toc.js'
import panzoom from 'panzoom'
import tablesort from 'tablesort'
import Gumshoe from 'gumshoejs'

// Note: the tablesort lib is not ESM friendly, and the sorts expect `Tablesort` to be available on the global
window.Tablesort = tablesort
require('tablesort/dist/sorts/tablesort.number.min.js')

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

function initTocDepthSlider () {
  var slider = document.getElementById('toc-depth-slider')
  var toc = document.querySelector('.toc')
  
  slider.addEventListener('change', (event) => {
    handleSliderChange(Number(event.target.value))
  })

  function handleSliderChange (depth) {
    console.log('handleSliderChange', depth)
    for (let i = 0; i < 6; i++) {
      toc.querySelectorAll(`.depth-${i}`).forEach(el => {
        if (i < depth) {
          el.classList.remove('maybe-hide')
        } else {
          el.classList.add('maybe-hide')
        }
      })
    }
  }
  // init to the current value
  handleSliderChange(slider.value)
}

function initTocScrollSpy () {
  console.log('initTocScrollSpy')
  var spy = new Gumshoe('.toc a', {
    nested: true,
    nestedClass: 'active-parent'
  })
}

window.addEventListener('DOMContentLoaded', () => {
  initToc({tocSelector:'.toc', contentSelector: '.markdown'})
  initTocDepthSlider()
  initTocScrollSpy()
  initPanZoom()
  initTableSort()
});
