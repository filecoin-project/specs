combos = makeQuery([{
  windows: 256,
  window_size_gib: 4,
  nodes_in_sequence: 8,
  post_window_challenges: 2,
  expander_degree: 384,
  butterfly_degree: 16,
  expander_layers: 8,
  butterfly_layers: 7,
  porep_lambda: 10,
  // porep_challenges: 2080,
  node_size: 32,
  snark_partition: 100000000
}])
  .add({
    mtree_hash_name: 'poseidon',
    mtree_hash_time: 8.3e-7, // ((8/7)*(2^27/8 -1))*32*8, // GPU 4s per GiB // CPU 5.803e-5,
    mtree_hash_blocks: 8,
    mtree_hash_constraints: 508 + 56,
    kdf_constraints: 25849/2
  })
  .add({
    commd_hash_name: 'sha',
    commd_hash_constraints: 25840,
    commd_hash_time: 130e-9,
  })
  .add({
    rig_memaccess_throughput_tb_s:  3,
    rig_hashing_throughput_tb_s: 0.016 * 32,
    rig_lifetime_years: 2,
    rig_cost: 2000,
    rig_storage_lifetime_years: 2,
    rig_cost_storage_tb: 15
  })
  .add({
    snark_constraint_time: 0.00000317488,
    snark_size: 192
  })
  .add({
    spacegap: 0.15,
    delta: 0.05
  })
  .add({
    proving_period_hours: 24,
  })
  .add({
    network_size_eib: 10,
    block_time: 30,
    tipset_size: 1,
  })
  .compile()

solved_many

report_from_result(solved_many[0], {})

md`---`
md`## Dev`

md`### Vars`

class Query {
  constructor(query = []) {
    this.query = query
  }
  add (query) {
    this.query = extend_query(this.query, [query])
    return this
  }
  extend(query) {
    this.query = extend_query(this.query, query)
    return this
  }
  compile() {
    return this.query
  }
}

makeQuery = (query = []) => {
  return new Query(query)
}

range = d3.range


createJsonDownloadButton(combos)

solved_many_pre = (await solve_many_chunk(combos))
  .map(d => {
    d.construction = `${d.graph_name}_${d.proof_name}`
    d.proof_name = `${d.spacegap}_${d.sdr_delta}`
    return d
  })


solved_many = solved_many_pre
  .filter(d => d !== null)
 // solved_manys = (await solve_manys(combos)).flat()

createJsonDownloadButton(solved_many)

md`### Orient`

function dump_vars() {
  return fetch(orientServer + '/dump-vars')
    .then(response => response.json())
    .then(json => {

      const map = {}
      json.forEach(d => {
        map[d.name] = d
      })

      return map
    })
}

vars = dump_vars()

function solve_multiple(json) {
  return fetch(orientServer + '/solve', {
    body: JSON.stringify(json),
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    method: 'POST'
  })
    .then(res => {
      return res.json()
    })
    .then(res => {
      return res.map(d => {
        const results = {}
        Object.keys(res[0])
          .filter(d => !d.includes('%'))
          .map(d => {
            results[d] = res[0][d]
          })
        return results
      })
    })

}

chunk = (json, parts) => {
  return json.reduce((acc, curr) => {
    const index = acc.length - 1

    if (index < 0) {
      acc.push([curr])
      return acc
    }

    if (acc[index].length === parts) {
      acc.push([curr])
    } else {
      acc[index].push(curr)
    }

    return acc
  }, [])
}

async function solve_many_chunk(json) {
  const promised = await Promise.all(chunk(json, 10).map(chunk_json => {
    return fetch(orientServer + '/solve-many', {
      body: JSON.stringify(chunk_json),
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      method: 'POST'
    }).then(res => {
      return res.json()
    }).then(res => {
      return res
        .filter(d => d !== null)
        .map(d => d.flat())
        .flat()
        .filter(d => d !== null)
        .map(d => {
          const results = {}
          Object.keys(d)
            .filter(key => !key.includes('%'))
            .map(key => {
              results[key] = d[key]
            })

          return results
        })
    })
  }))
  return promised.flat()
}

function solve_manys(json) {
  return fetch(orientServer + '/solve-many', {
    body: JSON.stringify(json),
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    method: 'POST'
  }).then(res => {
    return res.json()
  }).then(res => {
    return res
      .filter(d => d !== null)
      .map(d => d.flat())
      .flat()
      .filter(d => d !== null)
      .map(d => {
        const results = {}
        Object.keys(d)
          .filter(key => !key.includes('%'))
          .map(key => {
            results[key] = d[key]
          })

        return results
      })
  })
}

function ev (func, data) {
  let res

  try {
    res = (1, eval)(`(${func})`)(data)
  } catch(err) {
    throw err
  }

  return res
}

function solve_many(json) {
  return Promise.all(json.map(j => solve(j))).then(json => json.map(d => d[0]))
}

function solve(json) {
  return fetch(orientServer + '/solve', {
    body: JSON.stringify(json),
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    method: 'POST'
  }).then(res => {
    return res.json()
  }).then(res => {
    const results = {}
    Object.keys(res[0])
      .filter(d => !d.includes('%'))
      .map(d => {
        results[d] = res[0][d]
      })
    return results
  })
}

orientServer = `http://${window.location.hostname}:8000`

md`### Orientable`

function report_from_result(result, starting_assignments, simplify_terms) {
  const html = md`

| name | val | type | desc |
| ---- | --- | ---- | ---- |
${Object.keys(result)
// .filter(d => vars && vars[d] && vars[d].description)
.sort()
  .map(d => `| ${!starting_assignments[d] ? `**${d}**` : d} | ${result[d]} | ${vars[d] && vars[d].type ? vars[d].type : ''} | ${vars[d] && vars[d].description ? vars[d].description : ''} |\n`)}
`
  html.value = result
  return html
}

bar_chart = (raw, title, variables, group_by, opts) => {
  let data = raw
      .map(d => {
        return variables.map(key => ({
          construction: group_by.map(g => `${d[g]}`).join(', '),
          type: key,
          value: d[key],
          title: d[title]
        }))
      })
      .flat()
      .reduce((acc, curr) => {
        let exists = acc.count[`${curr.construction}_${curr.type}`]
        if (exists) {
          return acc
        }
        acc.count[`${curr.construction}_${curr.type}`] = 1
        acc.res.push(curr)
        return acc
      }, {res: [], count: {}})
      .res

  let discarded_data = []
  let organized_data = []

  if (opts && opts.filter) {
    organized_data = data.filter(d => opts.filter(d['title']))
    discarded_data = data.filter(d => !opts.filter(d['title']))
  }
  let graph = {
    "$schema": "https://vega.github.io/schema/vega-lite/v4.json",
    "title": `Composition of ${title} (${vars[title].type || ''})`,
    vconcat: [
      {
        layer: [{
          "width": 800,
          "mark": "bar",
          "data": { values: organized_data },
          "encoding": {
            "x": {"aggregate": "sum", "field": "value", "type": "quantitative"},
            "y": {"field": "construction", "type": "nominal", "sort": {"op": "sum", "field": "title"}},
            "color": {"field": "type", "type": "nominal"}
          }
        }]
      },
    ]
  }

  if (discarded_data.length > 0) {
    graph.vconcat.push({
      "mark": "bar",
      "width": 800,
      "title": "Data filtered out",
      "data": { values: discarded_data },
      "encoding": {
        "x": {"aggregate": "sum", "field": "value", "type": "quantitative"},
        "y": {"field": "construction", "type": "nominal", "sort": {"op": "sum", "field": "title"}},
        "color": {"field": "type", "type": "nominal"}
      }
    })
  }

  if (opts && opts.yrule) {
    const rule = [{}]
    rule[0][title] = opts.yrule
    graph.vconcat[0].layer.push({
      "data": { "values": rule },
      "layer": [{
        "mark": "rule",
        "encoding": {
          "x": {"field": title, "type": "quantitative"},
          "color": {"value": "red"},
          "size": {"value": 3}
        }
      }]
    })
  }
  return vl(graph)
}

add_query = (query, ext) => {
  return query.map(d => Object.assign({}, d, ext))
}

extend_query = (array, ...exts) => {
  let query = array

  const extend_one = (arr, ext) => arr.map(d => ext.map((_, i) => Object.assign({}, d, ext[i])))

  exts.forEach(ext => {
    query = extend_one(query, ext).flat()
  })

  return query
}

multiple_solutions = (solutions, group_by, filter) => {
  return solutions.map(s => {
    const solution = {}

    filter.forEach(d => {
      solution[d] = s[d]
    })
    solution.name = group_by.map(g => `${g}=${s[g]}`).join(', ')

    return solution
  })
}

table_constraints = (solutions, filter, group_by, sort_by) => {
  let results = multiple_solutions(solutions, group_by, filter)
  if (sort_by) {
    results = results.sort((a, b) => +a[sort_by] - +b[sort_by])
  }
  const header = `
  ${sort_by ? `Sorted by: ${sort_by}` : ''}

  ${group_by.length ? `| name ` : ``}| ${filter.join(' | ')} |`
  const divider = `${group_by.length ? `| --- ` : ``}| ${filter.map(f => '--:').join(' | ')} |`
  const rows = results.map(r => {
    return `${group_by.length ? `| ${r.name} ` : ``}| ${filter.map(f => `\`${_f(r[f])}\``).join(' | ')}`
  })
  const table = [header, divider, rows.join('\n')].join('\n')

  return md`${table}`
}

chooser = (data, field, base) => {
  const log_base = base ? Math.log10(base) : false
  const maximum = Math.log10(Math.max(...solved_many.map(d => d[field])))+0.5
  const minimum = Math.log10(Math.min(...solved_many.map(d => d[field])))
  const format = v => `${_f(Math.pow(10, v))} ${vars[field].type || ''}`

  return slider({
    title: field,
    description: vars[field].description,
    min: minimum,
    max: maximum,
    value: log_base || maximum,
    step: 0.01,
    format: format,
  })
}

md`### Utils`

_f = (d) => typeof d == 'number' || !Number.isNaN(+d) ? d3.format('0.10~f')(d) : d

jsonToSliders = (obj, assigned) => {
  const inputs = Object.keys(obj).map(d => `
<div style="padding-bottom: 10px; padding-top: 10px;">
  <div style="font: 700 0.9rem sans-serif;">${d}</div>
  <div class="input">
    <input type=range name=${d} min=${obj[d].min} max=${obj[d].max} step=${obj[d].step} value=${obj[d].value}>
  </div>
<div style="font-size: 0.85rem; font-style: italic;">${vars && vars[d] && vars[d].description ? vars[d].description : ''}</div>
</div>`)
  const form = formToObject(html`
<form>
  ${inputs.join('\n')}
</form>`)

  if (assigned) {
    form.value = Object.assign({}, form.value, assigned)
  }

  return form
}

function formToObject (form) {
  // Array.from(form.children).forEach(el => {
  //   el.append(html`<span>hey</span>`)
  // })
  Array.from(form.querySelectorAll('input')).forEach(el => {
    el.parentNode.append(html`<output name=output_${el.name} style="font: 14px Menlo, Consolas, monospace; margin-left: 0.5em;"></output>`)
  })


  Array.from(form.querySelectorAll('input')).forEach(el => {
    el.oninput = (e) => {
      form[`output_${el.name}`].value = `${el.value} ${vars[el.name].type || ''}`
      e.stopPropagation()
    }
  })

  form.oninput = (e) => {
    const value = Array.from(form.elements)
      .reduce(function(map, _, i) {
        if (form.elements[i].name.substr(0,6) !== 'output') {
          map[form.elements[i].name] = form.elements[i].valueAsNumber
        }
        return map;
      }, {});

    Object.keys(value).forEach(k => {
      form[`output_${k}`].value = `${form[k].value} ${vars && vars[k] ? vars[k].type || '' : ''}`
    })

    e.stopPropagation()
  }

  form.onmouseup = form.onkeyup = form.ontouchend = e => {
    form.value = Array.from(form.elements)
      .reduce(function(map, _, i) {
        if (form.elements[i].name.substr(0,6) !== 'output') {
          map[form.elements[i].name] = form.elements[i].valueAsNumber
        }
        return map;
      }, {});

    Object.keys(form.value).forEach(k => {
      form[`output_${k}`].value = `${form[k].value} ${vars && vars[k] ? vars[k].type || '' : ''}`
    })

    form.dispatchEvent(new CustomEvent('input'));
  };

  form.onmouseup()
  return form
}

function flatten(items) {
  const flat = [];

  items.forEach(item => {
    if (Array.isArray(item)) {
      flat.push(...flatten(item));
    } else {
      flat.push(item);
    }
  });

  return flat;
}

md`### Imports`

import {slider, checkbox, number, button} from "@jashkenas/inputs"
d3 = require('d3')
vl = require('@observablehq/vega-lite')
import { createJsonDownloadButton } from "@trebor/download-json"
import {localStorage} from "@mbostock/safe-local-storage"

md`### Styles`

html`<style>
.markdown-body table td, .markdown-body table t {
  padding: 4px !important;
}
table {
  font-size: 12px
}
th {
  font-size: 10px;
}
label {
  font-size: 10px;
}
</style>`

graph_constraints = (solutions, x, y, group_by, opts) => {
  const results = multiple_solutions(solutions, group_by, [x, y])
  const graph = {
    "title": `Plotting:  ${x} vs ${y}`,
    "width": 600,
    layer: [{
      "data": {"values": results},
      "mark": {"type": "line"},
      "encoding": {
        "x": {
          "field": x,
          "type": "quantitative",
          "axis": {
            "labelLimit": 400,
            "labelPadding": 30
          }
        },
        "y": {
          "field": y,
          "type": "quantitative",
        },
        "color": {
          "field": "name",
          "type": "nominal",
          "scale": {"scheme": "category10"}
        },
      },
    }]
  }
  if (opts && opts.height) {
    graph.height = opts.height
  }
  if (opts && opts.yrule) {
    const rule = [{}]
    rule[0][y] = opts.yrule

    graph.layer.push({
      "data": { "values": rule },
      "layer": [{
        "mark": "rule",
        "encoding": {
          "y": {"field": y, "type": "quantitative"},
          "color": {"value": "red"},
          "size": {"value": 3}
        }
      }]
    })
  }
  return vl(graph)

}

Plotly = require("https://cdn.plot.ly/plotly-latest.min.js")


plot3d = (rows, x, y, z) =>  {
  var zData = rows.map(d => {
    return [d[x], d[y], d[z]]
  });

  var data = [{
    z: zData,
    type: 'surface'
  }];

  var data2 = [{
    x: rows.map(d => d[x]),
    y: rows.map(d => d[y]),
    z: rows.map(d => d[z]),
    type: 'scatter3d'
  }]

  var layout = {
    title: `${x} vs ${y} vs ${z}`,
    autosize: false,
    width: width * 0.7,
    height: width * 0.7,
    scene: {
      xaxis: {
        title:{
          text: x,
          font: {
            family: 'Courier New, monospace',
            size: 18,
            color: '#7f7f7f'
          }
        }
      },
      yaxis: {
        title:{
          text: y,
          font: {
            family: 'Courier New, monospace',
            size: 18,
            color: '#7f7f7f'
          }
        }
      },
      zaxis: {
        title:{
          text: z,
          font: {
            family: 'Courier New, monospace',
            size: 18,
            color: '#7f7f7f'
          }
        }
      }
    },
    margin: {
      l: 65,
      r: 50,
      b: 65,
      t: 90,
    }
  };

  const div = DOM.element('div');
  Plotly.newPlot(div, data2, layout);
  return div
}

// plotMultiLine(delta_solved, 'chung_delta', ['decoding_time_parallel', 'epost_time_parallel', 'onboard_tib_time_days'])

plotMultiLine = (solutions, x, names) => {

  const traces = names.map(y => {
    return {
      x: solutions.map(d => d[x]),
      y: solutions.map(d => d[y]),
      name: `${y} data`,
      yaxis: y,
      type: 'scatter'
    }
  })

  var layout = {
    title: 'multiple y-axes example',
    // width: 800,
    autosize: false,
    xaxis: {domain: [0.01, 0.20]},
    yaxis: {
      title: 'yaxis title',
      titlefont: {color: '#1f77b4'},
      tickfont: {color: '#1f77b4'}
    },
    yaxis2: {
      title: 'yaxis2 title',
      titlefont: {color: '#ff7f0e'},
      tickfont: {color: '#ff7f0e'},
      anchor: 'free',
      overlaying: 'y',
      side: 'left',
      position: 0.15
    },
    yaxis3: {
      title: 'yaxis4 title',
      titlefont: {color: '#d62728'},
      tickfont: {color: '#d62728'},
      anchor: 'right',
      overlaying: 'y',
      side: 'left'
    },
  };

  const div = DOM.element('div');
  Plotly.newPlot(div, traces, layout);
  return div
}

codeView =({value, mode, height, localStorageKey}) => {

  value = localStorage.getItem(localStorageKey) || value

  const fn = ({CodeMirror} = {}) => {
    return ({id, value, mode}) => {
      const cm = CodeMirror(document.body, {
        value,
        mode,
        lineNumbers: true
      })
      CodeMirror.modeURL = 'https://codemirror.net/mode/%N/%N.js'
      CodeMirror.autoLoadMode(cm, mode)

      cm.on('keypress', (cm, event) => {
        if (event.key === 'Enter' && event.shiftKey) {
          event.preventDefault()
          window.parent.postMessage({
            id,
            value: cm.getValue(),
            height: document.body.offsetHeight
          }, document.origin)
        }
      })
      setInterval(() => {
        window.parent.postMessage({
          id,
          height: document.body.offsetHeight
        }, document.origin)
      }, 100)
    }
  }
  const randomId = `el${Math.floor(Math.random() * 1000000)}`
  const frameSrc = `
    <link rel="stylesheet" href="https://unpkg.com/codemirror@5.39.2/lib/codemirror.css" />
    <script src="https://unpkg.com/codemirror@5.39.2/lib/codemirror.js"></script>
    <script src="https://codemirror.net/addon/mode/loadmode.js"></script>
    <script src="https://codemirror.net/mode/meta.js"></script>
    <style type="text/css">
      body, html {
        margin: 0;
        padding: 0;
        overflow-y: hidden;
      }
      .CodeMirror {
        border: 1px solid #eee;
        height: auto;
      }
      .CodeMirror-scroll {
        height: ${height}px;
        overflow-y: hidden;
        overflow-x: auto;
      }
    </style>
    <script type="text/javascript">
      document.addEventListener('DOMContentLoaded', () => {
        (${fn().toString().trim()})(${JSON.stringify({id: randomId, mode, value})})
      })
    </script>
  `
  const frameStyle = `width: 100%; height: 300px; border: 0; overflow-y: hidden;`
  const frame = html`<iframe style="${frameStyle}"></iframe>`
  const messageListener = event => {
    if (document.contains(frame)) {
      if (event.data.id === randomId) {
        if (event.data.value !== undefined) {
          frame.value = event.data.value
          localStorage.setItem(localStorageKey, frame.value)
          frame.dispatchEvent(new CustomEvent("input"))
        }
        frame.style.height = `${event.data.height}px`
      }
    } else {
      window.removeEventListener('message', messageListener)
    }
  }
  window.addEventListener('message', messageListener, false)
  frame.srcdoc = frameSrc
  frame.value = value
  return frame
}

// {
//   reset_button;
//   if (this) {
//     localStorage.removeItem('utility')
//   }

//   return !this
// }

function qs(variable) {
  var query = window.location.search.substring(1);
  var vars = query.split('&');
  for (var i = 0; i < vars.length; i++) {
    var pair = vars[i].split('=');
    if (decodeURIComponent(pair[0]) == variable) {
      return decodeURIComponent(pair[1]);
    }
  }
  return false
}
