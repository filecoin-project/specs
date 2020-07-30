import toc from './toc/index.js'

const run = function() {
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
        headingsOffset: 50,
    });
}

run()