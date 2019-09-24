
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
        $ul.classList.remove('hidden')
      } else {
        $ul.previousElementSibling.classList.remove('highlight')
        $ul.classList.add('hidden')
      }
    }
  }
}

// lead this at the beginning to set the slider correctly
window.addEventListener('DOMContentLoaded', onMenuDetailChange, false)
