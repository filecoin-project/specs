import {html, render} from 'lit-html';

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

const humanize = (input) => {
  if(input === 'wip') {
    return 'Draft/WIP'
  }

  if(input === 'n/a') {
    return 'N/A'
  }
  
  return input.charAt(0).toUpperCase() + input.substr(1);
}

const stateToNumber = (s) => {
  switch (s) {
    case 'done':
    case 'stable':
      return 1
    case 'reliable':
      return 2
    case 'wip':
      return 3
    case 'incorrect':
      return 4
    case 'missing':
    return 5
    default:
      return 6
  }
}

function buildDashboard(selector, model) {
  const data = tableData(model, [0])
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
      <td data-sort="${stateToNumber(i.dashboardState)}" class="text-black bg-na bg-${i.dashboardState}">${humanize(i.dashboardState)}</td>
      <td data-sort="${stateToNumber(i.dashboardAudit)}" class="text-black bg-na bg-${i.dashboardAudit}">
        ${i.dashboardAuditURL
          ? html`<a href="${i.dashboardAuditURL}" title="Read the audit report" target="_blank" rel="noopener noreferrer" class="text-black">${i.dashboardAuditDate}<i class="gg-external gg-s-half"></i></a>`
          : humanize(i.dashboardAudit) } 
      </td>
    </tr>
    `: '')}
  </tbody>
</table>  
  `
  render(tpl, document.querySelector(selector))
}

export {
    buildDashboard
}