graph(solutions, 'spacegap', 'block_size_kib', ['sector_size_gib'])

md`## Block size`

table(solutions, ['expected_winning_miners', 'block_size_kib', 'chain_size_year_gib', 'spacegap'], ['sector_size_gib'], 'block_size_kib')

md`## PoSt sizes`

table(solutions, ['post_snark_circuit', 'post_snark_proof_size', 'post_snark_proof_partitions'], ['spacegap'])

md`## PoSt in blockchain`

table(solutions, ['sector_size_gib', 'post_size_per_block', 'avg_posts_messages_per_block'], ['spacegap'])

graph(solutions, 'spacegap', 'post_size_per_block', ['sector_size_gib'])

md`## PoRep in blockchain`

graph(solutions, 'spacegap', 'seal_size_per_block', ['sector_size_gib', 'lambda'])

graph(solutions, 'spacegap', 'block_size_kib', ['sector_size_gib', 'lambda'])

md`## Interactivity`

table(solutions, ['block_size_kib', 'sector_size_gib', 'spacegap'], ['lambda'])

solutions = JSON.parse(await fetch_model('http://localhost:8000/solved-parameters.json'))

md`---`

clean = (solutions, group_by, filter) => {
  return solutions.map(s => {
    const solution = {}

    filter.forEach(d => {
      solution[d] = s[d]
    })
    solution.name = group_by.map(g => `${g}=${s[g]}`).join(', ')

    return solution
  })
}

table = (solutions, filter, group_by, sort_by) => {
  let results = clean(solutions, group_by, filter)
  if (sort_by) {
    results = results.sort((a, b) => +a[sort_by] > +b[sort_by])
  }
  const header = `
  Sorted by: ${sort_by}

  | name | ${filter.join(' | ')} |`
  const divider = `| --- | ${filter.map(f => '--:').join(' | ')} |`
  const rows = results.map(r => {
    return `| ${r.name} | ${filter.map(f => r[f]).join(' | ')}`
  })
  const table = [header, divider, rows.join('\n')].join('\n')

  return md`${table}`
}

graph = (solutions, x, y, group_by) => {
  const results = clean(solutions, group_by, [x, y])
  return vl({
    "title": `Plotting:  ${x} vs ${y}`,
    "width": 600,
    "data": {"values": results},
    "mark": {"type": "line"},
    "encoding": {
      "x": {
        "field": x,
        "type": "quantitative",
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
  })
}


fetch_model = (model_url) => {
  return fetch(model_url).then(response => {
      return (response.ok) ? response.text() : false;
    });
    // form.dispatchEvent(new CustomEvent("input"));
}

vl = require('@observablehq/vega-lite')

html`<img src="http://localhost:8000/filecoin.svg"/>`


chart = {
  const svg = d3.create("svg")

  const g = svg.append("image")
    .attr("xlink:href","http://localhost:8000/filecoin.svg")
      .attr("viewBox", [0, 0, width, 400]);

  svg.call(d3.zoom()
      .extent([[0, 0], [width, 400]])
      .scaleExtent([1, 8])
      .on("zoom", zoomed));

  function zoomed() {
    g.attr("transform", d3.event.transform);
  }

  return svg.node();
}

style = html`<link rel="stylesheet" type="text/css" href="https://static.observablehq.com/style.ff32b4006c53039c311bef34594462b0ef9290929816e4cdaf29cfb8cb6dc4d1.css"/>`

d3 = require("d3@5")
