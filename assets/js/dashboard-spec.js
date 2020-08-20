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
      <td class="text-black bg-na bg-${i.dashboardState}">${i.dashboardState}</td>
      <td class="text-transparent ${i.dashboardAudit > 0 ? 'bg-stable' : 'bg-incorrect'}">${i.dashboardAudit}</td>
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