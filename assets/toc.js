import {html, render} from 'lit-html';


function buildDashboard(model) {

  const tableData = (model, depth, output=[]) => {
    const index = depth.length-1
    for (const node of model) {
      depth[index] += 1
      output.push({number: depth.join("."), ...node})
      if (node.children) {
        tableData(node.children, [...depth, 0], output)
      }
    }
    return output
  }

  const data = tableData(model, [0])
  console.log("buildDashboard -> data", data)
  const tpl = html`
<table id="dashboard" class="sort Dashboard tablesort">
  <thead>
    <tr>
        <th>Section</th>
        <th>Weight</th>
        <th>State</th>
        <th>Theory Audit</th>
    </tr>
  </thead>
  <tbody>
    ${data.map((i)=> i.page ? html`
    <tr>
      <td class="Dashboard-section">${i.number} <a href="#${i.id}">${i.text}</a></td>
      <td>${i.dashboardWeight}</td>
      <td class="text-black bg-na bg-${i.dashboardState}">${i.dashboardState}</td>
      <td class="text-transparent ${i.dashboardAudit > 0 ? 'bg-stable' : 'bg-incorrect'}">${i.dashboardAudit}</td>
    </tr>
    `: '')}
  </tbody>
</table>  
  `
  render(tpl, document.querySelector('#test-dash'))
}
export function initToc ({tocSelector, contentSelector}) {
  const model = buildTocModel(contentSelector)
  console.log(model)
  const toc = buildTocDom(model)
  document.querySelector(tocSelector).appendChild(toc)
  console.log('toc rendered')
  buildDashboard(model)
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
      page: Boolean(el.dataset.page),
      dashboardWeight: el.dataset.dashboardWeight,
      dashboardAudit: el.dataset.dashboardAudit,
      dashboardState: el.dataset.dashboardState,
      children: []
    }
    if (!prevSibling || headingNum(node) === headingNum(prevSibling))  {
      parents[parents.length - 1].children.push(node)
      prevSibling = node
      
      // is h3 > h2 ?
    } else if (headingNum(node) > headingNum(prevSibling)) {
      parents.push(prevSibling)
      prevSibling.children.push(node)
      prevSibling = node
    } else {
      // h2 or h1 after an h3... gotta find out how far to unwind, parents may not be contiguous in a bad doc, so we walk.
      let prevParent = parents.pop()
      while (headingNum(node) <= headingNum(prevParent)) {
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
  return el.textContent.trim()
}

function headingNum (el) {
  return Number(el.tagName[1])
}
