module.exports = function updateTocScroll(options, toc) {
  if (toc && toc.scrollHeight > toc.clientHeight) {
    var activeItem = toc.querySelector('.' + options.activeListItemClass)
    if (activeItem ) {
        toc.scrollTop = activeItem.offsetTop/2
    }
  }
}