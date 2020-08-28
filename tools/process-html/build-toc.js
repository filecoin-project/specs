function buildToc ({model}) {
  const parent = document.createDocumentFragment()
  buildList(parent, model, 0)
  return parent
}

function buildList (parent, children, depth) {
  let ol = createList(depth)
  parent.append(ol)
  for (node of children) {
    let li = document.createElement('li')
    let a = document.createElement('a')
    a.setAttribute('href', '#' + node.id)
    a.innerText = node.text
    li.appendChild(a)
    ol.append(li)
    if (node.children) {
      buildList(li, node.children, depth + 1)
    }
  }
}

function createList(depth) {
  const ol = document.createElement('ol')
  ol.className = 'depth-' + depth
  return ol
}

module.exports = {
  buildToc
}
