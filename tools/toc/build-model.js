// [
//   { text: "Foo", id: "foo", tagName: 'h1', children: [
//     { text: "Xyz", id: "xyz", tagName: 'h2', children: [
//       { text: "Bar", id: "bar", tagName: 'h3, children:[] }
//       { text: "Bar", id: "bar", tagName: 'h3, children:[] }
//     ]}
//   ]},
//   { text: "Baz", id: "baz", tag: 'h1', children: [] }
// ]

function buildTocModel (root) {
  const model = []
  const headingList = root.querySelectorAll('h1,h2,h3,h4,h5,h6')
  let parents = [{tagName: 'H0', children: model}]
  let prevSibling = null
  let sectionNumber = [0]

  function addSibling(node) {
    sectionNumber[sectionNumber.length - 1] = sectionNumber[sectionNumber.length - 1] + 1
    node.number = sectionNumber.join('.')
    parents[parents.length - 1].children.push(node)
    prevSibling = node
  }

  function addChild(node) {
    sectionNumber.push(1)
    node.number = sectionNumber.join('.')
    parents.push(prevSibling)
    prevSibling.children.push(node)
    prevSibling = node
  }

  for (let el of headingList) {
    let node = {
      id: el.id,
      number: '',
      tagName: el.tagName,
      text: cleanHeadingText(el),
      page: Boolean(el.dataset.page),
      dashboardWeight: el.dataset.dashboardWeight,
      dashboardAudit: el.dataset.dashboardAudit,
      dashboardAuditURL: el.dataset.dashboardAuditUrl,
      dashboardAuditDate: el.dataset.dashboardAuditDate,
      dashboardState: el.dataset.dashboardState,
      children: []
    }
    
    if (!prevSibling || headingNum(node) === headingNum(prevSibling))  {
      // sibling: h2 == h2
      addSibling(node)
      
    } else if (headingNum(node) > headingNum(prevSibling)) {
      // child: h3 > h2
      addChild(node)

    } else {
      // h2 or h1 after an h3... gotta find out how far to unwind. Parents may not be contiguous, so walk till we find a parent
      let target = headingNum(node)
      let rmCount = 0
      while (target <= headingNum(parents[parents.length - (rmCount + 1)])) {
        rmCount++
      }
      parents = parents.slice(0, parents.length - rmCount)
      sectionNumber = sectionNumber.slice(0, sectionNumber.length - rmCount)

      addSibling(node)
    }
  }

  return model
}

function cleanHeadingText (el) {
  return el.textContent.trim()
}

function headingNum (el) {
  return Number(el.tagName[1])
}

module.exports = {
  buildTocModel
}