md`# Proofs Tradeoff Report`

viewof config = {
  const form = formToObject(html`
<form>
  <div><input type=range name=post_lambda min=1 max=128 step=1 value=10> <i>post_lambda</i></div>
</form>`)
  return form
}

constants = Object.assign({}, base, constraints, filecoin, bench, rig)

combos = [wrapperVariant, wrapper, stackedReplicas]
  .map(d => [
    Object.assign({}, constants, d, stackedChungParams, config),
    Object.assign({}, constants, d, stackedSDRParams, config)
   ]).flat()


solved_many = (await solve_many(combos)).map(d => d[0])
// solved_manys = (await solve_manys(combos)).flat()

md`#### Vars that matter`
table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'decoding_time_parallel',
  'porep_time_parallel',
  'porep_proof_size',
  'epost_time_parallel',
], [])

md`#### Graphs`
table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'porep_lambda',
  'porep_challenges',
  'post_lambda',
  'post_challenges',
  'stacked_layers',
  'expander_parents',
  'drg_parents',
  'windows',
  'window_size_mib',
  'sector_size_gib',
], [])

md`#### Encoding time`
table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'encoding_time',
  'encoding_time_parallel',
], [])

md`#### PoRep`
table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'porep_proof_size',
  'porep_snark_constraints',
  'porep_time'
], [])

md`#### PoSt`
table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'post_proof_size',
  'post_snark_constraints',
  'post_snark_time',
  'post_snark_time_parallel',
  'post_time',
  'post_time_parallel',
  'post_inclusions_time',
  'post_inclusions_time_parallel',
  'post_data_access',
  'post_data_access_parallel'
], [])

md`#### EPoSt`
table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'epost_time',
  'epost_time_parallel',
  'epost_inclusions_time',
  'epost_inclusions_time_parallel',
  'epost_data_access',
  'epost_data_access_parallel'
], [])
// report_from_result(solved_many[1], combos[1])

md`---`

md`### Parameters`

base = ({
  "porep_lambda": 10,
  "post_lambda": 10,
  "sector_size_gib": 32,
  "window_size_mib": 64,
  "wrapper_parents": 10
})

md`### Constants`

md`#### Graph`

stackedChungParams = ({
  "graph_name": "Chung",
  "!StackedChungParameters": true,
  "!StackedSDRParameters": false,
  "chung_delta": 0.01,
  "expander_parents": 16
})

stackedSDRParams = ({
  "graph_name": "SDR",
  "!StackedChungParameters": false,
  "!StackedSDRParameters": true,
  "sdr_delta": 0.01
})

md`#### Proofs`

wrapper = ({
  "proof_name": "wrapping",
  "!ElectionWithFallbackPoSt": true,
  "!SectorEncoding": true,
  "!VectorR": true,
  "!Wrapping": true,
})

wrapperVariant = ({
  "proof_name": "wrappingVariant",
  "!ElectionWithFallbackPoSt": true,
  "!SectorEncoding": true,
  "!VectorR": true,
  "!WrappingVariant": true,
})

stackedReplicas = ({
  "proof_name": "stackedReplicas",
  "!ElectionWithFallbackPoSt": true,
  "!SectorEncoding": true,
  "!VectorR": true,
  "!StackedReplicas": true
})

md`#### Protocol`

filecoin = ({
  "ec_e": 5,
  "fallback_period_days": 1,
  "fallback_ratio": 0.05,
  "filecoin_reseals_per_year": 1,
  "filecoin_storage_capacity_eib": 10,
  "node_size": 32,
  "polling_time": 15,
  "post_mtree_layers_cached": 25,
  "cost_amax": 1,
  "hashing_amax": 2,
  "spacegap": 0.2,
  "proofs_block_fraction": 0.3,
  "epost_challenged_sectors_fraction": 0.04,
})

md`### Miner`

md`#### Hardware Config`

rig = ({
  "rig_cores": 16,
  "rig_snark_parallelization": 2,
  "rig_malicious_cost_per_year": 2.5,
  "rig_ram_gib": 32,
  "rig_storage_latency": 0.003,
  "rig_storage_min_tib": 100,
  "rig_storage_parallelization": 30,
  "rig_storage_read_mbs": 80,
  "cost_gb_per_month": 0.005,
  "extra_storage_time": 0,
  "hash_gb_per_second": 5,
})

md`#### Benchmarks`

bench = ({
  "column_leaf_hash_time": 2.56e-7,
  "kdf_time": 1.28e-8,
  "merkle_tree_datahash_time": 1.28e-8,
  "merkle_tree_hash_time": 2.56e-7,
  "snark_constraint_time": 0.000004642,
  "ticket_hash": 2.56e-7,
})

md`### SNARKs`

constraints = ({
  "merkle_tree_hash_constraints": 1376,
  "ticket_constraints": 1376,
  "merkle_tree_datahash_constraints": 56000,
  "kdf_constraints": 25000,
  "column_leaf_hash_constraints": 1376,
  "snark_size": 192,
  "porep_snark_partition_constraints": 100000000,
  "post_snark_partition_constraints": 3000000,
})

md`---`

md`## Dev`

md`### Orient`

function solve_multiple(json) {
  return fetch('http://localhost:8888/solve', {
    body: JSON.stringify(json),
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    method: 'POST'
  }).then(res => {
    return res.json()
  }).then(res => {
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

function solve_manys(json) {
  return fetch('http://localhost:8888/solve-many', {
    body: JSON.stringify(json),
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    method: 'POST'
  }).then(res => {
    return res.json()
  }).then(res => {
    return res.flat().map(d => {
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

function solve_many(json) {
  return Promise.all(json.map(j => solve(j)))
}

function solve(json) {
  return fetch('http://localhost:8888/solve', {
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

orientServer = 'http://127.0.0.1:8888'

md`### Orientable`

function report_from_result(result, starting_assignments, simplify_terms) {
  const html = md`

| name | val |
| ---- | --- |
${Object.keys(result).sort()
  .map(d => `| ${!starting_assignments[d] ? `**${d}**` : d} | ${result[d]} |\n`)}
`
  html.value = result
  return html
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
    results = results.sort((a, b) => +a[sort_by] > +b[sort_by])
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

md`### Utils`

_f = (d) => typeof d == 'number' || !Number.isNaN(+d) ? d3.format('0.3~f')(d) : d

function formToObject (form) {
  // Array.from(form.children).forEach(el => {
  //   el.append(html`<span>hey</span>`)
  // })
  Array.from(form.querySelectorAll('input')).forEach(el => {
    el.parentNode.append(html`<output name=output_${el.name} style="font: 14px Menlo, Consolas, monospace; margin-left: 0.5em;"></output>`)
  })
  

  form.oninput = () => {
    form.value = Array.from(form.elements)
      .reduce(function(map, _, i) {
        if (form.elements[i].name.substr(0,6) !== 'output') {
          map[form.elements[i].name] = form.elements[i].valueAsNumber
        }
        return map;
      }, {});
    
    Object.keys(form.value).forEach(k => { form[`output_${k}`].value = form[k].value })
  }

  form.oninput()
  
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

d3 = require('d3')
vl = require('@observablehq/vega-lite')
