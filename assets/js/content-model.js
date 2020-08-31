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
        dashboardAuditURL: el.dataset.dashboardAuditUrl,
        dashboardAuditDate: el.dataset.dashboardAuditDate,
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

  function cleanHeadingText (el) {
    return el.textContent.trim()
  }
  
  function headingNum (el) {
    return Number(el.tagName[1])
  }

  export {
    buildTocModel
  }