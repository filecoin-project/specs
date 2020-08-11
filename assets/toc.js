
export function initToc ({tocSelector, contentSelector}) {
  const model = buildTocModel(contentSelector)
  console.log(model)
  const toc = buildTocDom(model)
  document.querySelector(tocSelector).appendChild(toc)
  console.log('toc rendered')
}

  // [
  //   { text: "Foo", id: "foo", tagName: 'h1', children: [
  //     { text: "Xyz", id: "xyz", tagName: 'h2', children: [
  //       { text: "Bar", id: "bar", tagName: 'h3, children:[] }
  //       { text: "Bar", id: "bar", tagName: 'h3, children:[] }
  //     ]}
  //   ]},
  //   { text: "Baz", id: "baz", tag: 'h1', children: [] }
  // ]
function buildTocModel (contentSelector) {
  const model = []
  const headingList = document.querySelector(contentSelector).querySelectorAll('h1,h2,h3,h4,h5,h6')
  let parents = [{tagName: 'H0', children: model}]
  let prevSibling = null
  for (let el of headingList) {
    let node = {
      id: el.id,
      tagName: el.tagName,
      text: cleanHeadingText(el),
      children: []
    }
    if (!prevSibling || node.tagName[1] === prevSibling.tagName[1])  {
      parents[parents.length - 1].children.push(node)
      prevSibling = node
      
      // quick and dirty check for is h3 > h2 ?
    } else if (node.tagName[1] > prevSibling.tagName[1]) {
      parents.push(prevSibling)
      prevSibling.children.push(node)
      prevSibling = node
    } else {
      // h2 or h1 after an h3... gotta find out how far to unwind, parents may not be contiguous in a bad doc, so we walk.
      let prevParent = parents.pop()
      while (node.tagName[1] <= prevParent.tagName[1]) {
        prevParent = parents.pop()
      }
      prevParent.children.push(node)
      parents.push(prevParent)
      prevSibling = node
    }
  }
  return model
}


function buildTocDom (model) {
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

function cleanHeadingText (el) {
  // in the current dom, the first child of the h{1-6} el is the text we want
  console.log(el.textContent, el)
  return el.textContent.trim()
}
