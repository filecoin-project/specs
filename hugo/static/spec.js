
function onReady() {
  onMenuDetailChange()
  addHeaderLinks()
  scrollSpy()
}
// lead this at the beginning to set the slider correctly
window.addEventListener('DOMContentLoaded', onReady, false)


function onMenuDetailChange(event) {
  var $slider = document.querySelector('#menu-detail-slider')
  var $label = document.querySelector('#menu-detail-slider-label')

  var depth = Number($slider.value) // force number
  $label.innerText = ''+depth

  for (var i = 1; i < 6; i++) {
    $uls = document.querySelectorAll('.menu-item-section.depth-' + i)
    for (var j = 0; j < $uls.length; j++) {
      $ul = $uls[j]
      if (i < depth) {
        if ($ul.childElementCount > 0) {
          $ul.previousElementSibling.classList.add('highlight')
        }
        $ul.classList.remove('maybe-hide')
      } else {
        $ul.previousElementSibling.classList.remove('highlight')
        $ul.classList.add('maybe-hide')
      }
    }
  }
}

function addHeaderLinks() {
  var $elems = document.querySelectorAll("h1, h2, h3, h4, h5, h6")
  $elems.forEach(function(el) {
    el.onclick = function() {
      var id = el.getAttribute('id') || el.parentElement.getAttribute('id')
      window.location.hash = id
    }
  })
}

function scrollSpy() {
  new Gumshoe('#menu-toc a', {
    nested: true,
    nestedClass: 'active-parent'
  })
}
