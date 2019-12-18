md`# Proofs Tradeoff Report

1. Calculator may contain errors and definitely TODOs
2. Adding new graphs is easy
3. Fixing errors or change construction is easy
`

md`## Sliders`

md`### Parameters`

viewof base = jsonToSliders(
  {
    "porep_lambda": {value: 10, min: 0, max: 128, step: 1},
    "post_mtree_layers_cached": {value: 25, min: 0, max: 40, step: 1},
    "post_lambda": {value: 10, min: 0, max: 128, step: 1},
  },
  { "!StackedReplicaUnaligned": true}
)

md`### Constants`

md`#### Graph`

md`##### Chung params`
viewof stackedChungParams = jsonToSliders(
  {
    "chung_delta": {value: 0.06, min: 0.01, max: 0.08, step: 0.01},
    "expander_parents": {value: 16, min: 0, max: 128, step: 1},
    "rig_malicious_cost_per_year": {value: 2000, min: 0, max: 10000, step: 0.1},
    "hash_gb_per_second": {value: 640000, min: 0, max: 100000000, step: 1},
    "sector_size_gib": {value: 32, min: 1, max: 1024, step: 1},
    "window_size_mib": {value: 1024, min: 4, max: 32 * 1024, step: 1},
    "wrapper_parents": {value: 100, min: 0, max: 1000, step: 1},
  },
  {
    "graph_name": "Chung",
    "!StackedChungParameters": true,
  })

md`##### SDR params`
viewof stackedSDRParams = jsonToSliders(
  {
    "sdr_delta": {value: 0.01, min: 0.01, max: 0.08, step: 0.01},
    "rig_malicious_cost_per_year": {value: 3, min: 0, max: 10000, step: 0.1},
    "hash_gb_per_second": {value: 5, min: 0, max: 10000, step: 1},
  },
  {
    "graph_name": "SDR",
    "!StackedSDRParameters": true,
    "!TimingAssumption": true,
    "sdr_delta": 0.01,
    "windows": 1,
    "time_amax": 2
  })


md`#### Proofs`

wrapper = ({
  "proof_name": "wrapping",
  "!ElectionWithFallbackPoSt": true,
  "!VectorR": true,
  "!Wrapping": true,
})

wrapperVariant = ({
  "proof_name": "wrappingVariant",
  "!ElectionWithFallbackPoSt": true,
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

viewof filecoin = jsonToSliders({
  "ec_e": {value: 5, min: 1, max: 20, step: 1},
  "fallback_period_days": {value: 1, min: 1, max: 10, step: 1},
  "fallback_ratio": {value: 0.05, min: 0.01, max: 1, step: 0.01},
  "filecoin_reseals_per_year": {value: 1, min: 0, max: 365, step: 1},
  "filecoin_storage_capacity_eib": {value: 1, min: 0.5, max: 20, step: 0.5},
  "block_time": {value: 15, min: 15, max: 60, step: 1},
  "cost_amax": {value: 1, min: 1, max: 10, step: 1},
  "hashing_amax": {value: 2, min: 1, max: 10, step: 1},
  "spacegap": {value: 0.2, min: 0.01, max: 0.2, step: 0.01},
  "proofs_block_fraction": {value: 0.3, min: 0.01, max: 1, step: 0.01},
  "epost_challenged_sectors_fraction": {value: 0.04, min: 0.01, max: 1, step: 0.01},
}, {
  node_size: 32
})

md`### Miner`

md`#### Hardware Config`

viewof rig = jsonToSliders({
  "rig_cpu_cost": {value: 1200, min: 0, max: 10000, step: 100},
  "rig_gpu_cost": {value: 1200, min: 0, max: 10000, step: 100},
  "rig_cpu_lifetime_years": {value: 2, min: 1, max: 10, step: 1},
  "rig_gpu_lifetime_years": {value: 2, min: 1, max: 10, step: 1},
  "rig_cores": {value: 16, min: 1, max: 512, step: 1},
  "rig_snark_parallelization": {value: 2, min: 1, max: 64, step: 1},
  "rig_ram_gib": {value: 32, min: 1, max: 128, step: 1},
  "rig_storage_latency": {value: 0.003, min: 0.0003, max: 0.01, step: 0.0001},
  "rig_storage_min_tib": {value: 100, min: 0.5, max: 1024, step: 0.5},
  "rig_storage_parallelization": {value: 4, min: 1, max: 128, step: 1},
  "rig_storage_read_mbs": {value: 80, min: 80, max: 1000, step: 1},
  "cost_gb_per_month": {value: 0.0025, min: 0.0001, max: 0.1, step: 0.0001},
  "extra_storage_time": {value: 0, min: 0, max: 10, step: 1 },
})

md`#### Benchmarks`

bench = ({
  "column_leaf_hash_time": 1.7028e-5/2,
  "kdf_time": 5.4e-7/2,
  "merkle_tree_datahash_time": 1.28e-8/2,
  "merkle_tree_hash_time": 1.7028e-5/2,
  "snark_constraint_time": 3.012e-5/2,
  "ticket_hash": 1.7028e-5/2,
  "kdf_time_fastcpu": 0.0000000128/2
})

md`### SNARKs`

poseidon = ({
  "mtree_hash_name": "poseidon",
  "merkle_tree_hash_constraints": 1376/8/2,
  "ticket_constraints": 1376/8/2,
  "column_leaf_hash_constraints": 1376/8/2,
})

pedersen = ({
  "mtree_hash_name": "pedersen"
})


constraints = ({
  "merkle_tree_hash_constraints": 1376/2,
  "ticket_constraints": 1376/2,
  "merkle_tree_datahash_constraints": 25000/2,
  "kdf_constraints": 14066/2,
  "column_leaf_hash_constraints": 1376/2,
  "snark_size": 192,
  "porep_snark_partition_constraints": 100000000,
  "post_snark_partition_constraints": 3000000,
})

md`## Filters`

// Window size MiB
// window_size_mib_choices = [4, 64, 128, 1024, 16384, 32768]
// viewof window_size_mib_config = checkbox({
//   title: "Window Sizes",
//   options: solved_many_pre.map(d => d['window_size_mib']).map(d => ({value: d, label: d})),
//   value: solved_many_pre.map(d => d['window_size_mib']),
// })

md`## Utility`

html`Link to <a href="?utility_raw=${encodeURIComponent(utility_raw)}&utility_cols=${encodeURIComponent(utility_cols.join(','))}">current setting</a>`

viewof utility_raw = codeView({
  localStorageKey: 'utility',
  value: qs('utility_raw') || `function (d) {
  return 0.5 * d['porep_time_parallel']
}`,
  mode: 'javascript',
  height: 300
})

// viewof reset_button = button({value: "Reset" })

viewof utility_cols = checkbox({
  title: "Vars to show in utility table",
  options: Object.keys(vars).map(d => ({value: d, label: d})),
  value: qs('utility_cols') || ['decoding_time_parallel', 'proofs_per_block_kib', 'epost_time_parallel', 'porep_cost', 'porep_decoding_cost', 'epost_cost']
})

table_constraints(
  solved_many,
  ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name', 'utility'].concat(utility_cols),
  [],
  'utility'
)

md`## Graphs`
md`### On-chain footprint

This graphshows the average proofs per block (assuming a network size of ${filecoin.filecoin_storage_capacity_eib}EiB)
`

viewof proofs_per_block_kib_ruler = chooser(solved_many, 'proofs_per_block_kib', 1000)
bar_chart(solved_many, 'proofs_per_block_kib', [
  'seals_size_per_block_kib',
  'posts_size_per_block_kib',
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, proofs_per_block_kib_ruler),
  yrule: Math.pow(10, proofs_per_block_kib_ruler)
})

md`### Encoding time (estimated from benchmarks)`

viewof encoding_time_ruler = chooser(solved_many, 'encoding_time', 3*60*60)

bar_chart(solved_many, 'encoding_time', [
  'encoding_time',
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, encoding_time_ruler),
  yrule: Math.pow(10, encoding_time_ruler)
})

md`### Encoding time (estimated from fastest CPU)`

viewof encoding_time_fastcpu_ruler = chooser(solved_many, 'encoding_time_fastcpu', 3*60*60)

bar_chart(solved_many, 'encoding_time_fastcpu', [
  'encoding_time_fastcpu',
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, encoding_time_fastcpu_ruler),
  yrule: Math.pow(10, encoding_time_fastcpu_ruler)
})



md`### Retrieval`

viewof decoding_time_parallel_ruler = chooser(solved_many, 'decoding_time_parallel', 2)

bar_chart(solved_many, 'decoding_time_parallel', [
  'encoding_window_time_parallel',
  'window_read_time_parallel',
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, decoding_time_parallel_ruler),
  yrule: Math.pow(10, decoding_time_parallel_ruler)
})

viewof decoding_time_ruler = chooser(solved_many, 'decoding_time', 16)

bar_chart(solved_many, 'decoding_time', [
  'encoding_window_time',
  'window_read_time',
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, decoding_time_ruler),
  yrule: Math.pow(10, decoding_time_ruler)
})


// table_constraints(solved_many, [
//   'proof_name',
//   'graph_name',
//   'window_size_mib',
//   'decoding_time_parallel',
//   'encoding_window_time_parallel',
//   'window_read_time_parallel'
// ], [])

md`### PoRep time`

viewof porep_time_parallel_ruler = chooser(solved_many, 'porep_time_parallel', 12 * 60 * 60)

bar_chart(solved_many, 'porep_time_parallel', [
  'porep_snark_time_parallel',
  'porep_commit_time_parallel',
  'encoding_time_parallel'
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, porep_time_parallel_ruler),
  yrule: Math.pow(10, porep_time_parallel_ruler)
})

md`### PoRep cost`

viewof porep_cost_ruler = chooser(solved_many, 'porep_cost', 2)

bar_chart(solved_many, 'porep_cost', [
  'porep_commit_cost',
  'porep_encoding_cost',
  'porep_snark_cost'
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, porep_cost_ruler),
  yrule: Math.pow(10, porep_cost_ruler)
})

md`### PoRep constraints`

viewof porep_snark_constraints_ruler = chooser(solved_many, 'porep_snark_constraints', 1000 * 1000 * 1000)

bar_chart(solved_many, 'porep_snark_constraints', [
  'porep_commc_leaves_constraints',
  'porep_commc_inclusions_constraints',
  'porep_commr_inclusions_constraints',
  'porep_commd_inclusions_constraints',
  'porep_labeling_proofs_constraints'
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, porep_snark_constraints_ruler),
  yrule: Math.pow(10, porep_snark_constraints_ruler)
})

table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'window_size_mib',
  'mtree_hash_name',
  'porep_commc_leaves_constraints',
  'porep_commc_inclusions_constraints',
  'porep_commr_inclusions_constraints',
  'porep_commd_inclusions_constraints',
  'porep_labeling_proofs'
], [])



// bar_chart(solved_many, 'porep_commit_time', [
//   'commr_time',
//   'commq_time',
//   'commc_time'
// ], ['proof_name', 'graph_name', 'window_size_mib'])

// bar_chart(solved_many, 'commc_time', [
//   'commc_tree_time',
//   'commc_leaves_time',
// ], ['proof_name', 'graph_name', 'window_size_mib'])

md`### EPoSt`

viewof epost_time_parallel_ruler = chooser(solved_many, 'epost_time_parallel', 30)

bar_chart(solved_many, 'epost_time_parallel', [
  'epost_leaves_read_parallel',
  'epost_mtree_read_parallel',
  // 'epost_data_access_parallel',
  'post_ticket_gen',
  'epost_inclusions_time_parallel',
  'post_snark_time_parallel'
], ['proof_name', 'graph_name', 'window_size_mib', 'mtree_hash_name'], {
  filter: d => d < Math.pow(10, epost_time_parallel_ruler),
  yrule: Math.pow(10, epost_time_parallel_ruler)
})

table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'window_size_mib',
  'mtree_hash_name',
  'epost_time_parallel',
  'epost_data_access_parallel',
  'epost_mtree_read_parallel',
  'epost_leaves_read_parallel',
  'epost_leaves_read',
  'epost_challenged_sectors',
  'post_leaf_read',
  'post_ticket_gen',
  'epost_inclusions_time_parallel',
  'post_snark_time_parallel',
  'post_challenges',
  'post_challenge_read',
  'windows'
], [])

md`### Merkle tree caching`


// graph_constraints(mtree_solved, 'post_mtree_layers_cached', 'epost_time_parallel', ['proof_name'], { height: 100, yrule: 15 })
// graph_constraints(mtree_solved, 'post_mtree_layers_cached', 'post_inclusion_time', ['proof_name'], { height: 100 })

md`### Impact of \`chung_delta\` in StackedChung`
queries = [...Array(8)].map((_, i) => {
  const query = Object.assign(
    {},
    constants,
    stackedChungParams,
    { chung_delta: 0.01 * (i+1) },
    { window_size_mib: 128 }
  )

  return [
    Object.assign({}, query, stackedReplicas),
    Object.assign({}, query, wrapper),
    Object.assign({}, query, wrapperVariant),
  ]
}).flat()

// delta_solved = (await solve_many(queries)).map(d => d[0])

// // graph_constraints(delta_solved, 'chung_delta', 'stacked_layers', [], { height: 100 })
// // graph_constraints(delta_solved, 'chung_delta', 'porep_challenges', [], { height: 100 })
// // graph_constraints(delta_solved, 'chung_delta', 'post_challenges', [], { height: 100 })
// graph_constraints(delta_solved, 'chung_delta', 'decoding_time_parallel', ['proof_name'], {yrule: 0.5, height: 100})
// graph_constraints(delta_solved, 'chung_delta', 'epost_time_parallel', ['proof_name'], {yrule: 10, height: 100})
// graph_constraints(delta_solved, 'chung_delta', 'porep_proof_size_kib', ['proof_name'], { height: 100 })
// graph_constraints(delta_solved, 'chung_delta', 'block_size_kib', ['proof_name'], { height: 100 })
// graph_constraints(delta_solved, 'chung_delta', 'onboard_tib_time_days', ['proof_name'], { height: 100 })
// graph_constraints(delta_solved, 'chung_delta', 'porep_time_parallel', ['proof_name'], { height: 100 })
// plot3d(delta_solved, 'chung_delta', 'epost_time_parallel', 'onboard_tib_time_days')

md`---`

md`#### Other important vars`

table_constraints(solved_many, [
  'proof_name',
  'graph_name',
  'window_size_mib',
  'mtree_hash_name',
  'decoding_time_parallel',
  'porep_time_parallel',
  'porep_proof_size_kib',
  'block_size_kib',
  'epost_time_parallel',
], [])

// md`#### Graphs`
// table_constraints(solved_many, [
//   'proof_name',
//   'graph_name',
//   'window_size_mib',
//   'porep_lambda',
//   'porep_challenges',
//   'post_lambda',
//   'post_challenges',
//   'stacked_layers',
//   'expander_parents',
//   'drg_parents',
//   'windows',
//   'window_size_mib',
//   'sector_size_gib',
// ], [])

// md`#### PoRep`
// table_constraints(solved_many, [
//   'proof_name',
//   'graph_name',
//   'window_size_mib',
//   'encoding_time',
//   'encoding_time_parallel',
//   'porep_commit_time',
//   'porep_commit_time_parallel',
//   'porep_snark_time',
//   'porep_snark_time_parallel',
//   'porep_proof_size',
//   'porep_snark_constraints',
//   'porep_time'
// ], [])

// md`#### PoSt`
// table_constraints(solved_many, [
//   'proof_name',
//   'graph_name',
//   'window_size_mib',
//   'post_proof_size',
//   'post_snark_constraints',
//   'post_snark_time',
//   'post_snark_time_parallel',
//   'post_time',
//   'post_time_parallel',
//   'post_inclusions_time',
//   'post_inclusions_time_parallel',
//   'post_data_access',
//   'post_data_access_parallel'
// ], [])

// md`#### EPoSt`
// table_constraints(solved_many, [
//   'proof_name',
//   'graph_name',
//   'window_size_mib',
//   'epost_time',
//   'epost_time_parallel',
//   'epost_inclusions_time',
//   'epost_inclusions_time_parallel',
//   'epost_data_access',
//   'epost_data_access_parallel'
// ], [])

// md`## Debug`
// report_from_result(solved_many[0], combos[0])
// report_from_result(solved_many[1], combos[1])
// report_from_result(solved_many[2], combos[2])
// report_from_result(solved_many[3], combos[3])
// report_from_result(solved_many[4], combos[4])
// report_from_result(solved_many[5], combos[5])

md`---`
md`## Dev`

md`### Vars`
constants = Object.assign({}, base, constraints, filecoin, bench, rig)

combos = {
  let query = [constants]
  query = extend_query(query, [stackedReplicas])
  query = extend_query(query, [stackedSDRParams])
  query = extend_query(query, [poseidon, pedersen])
  return query
}

// combos = {
//   let start = [constants]
//   let proofs = extend_query(start, [wrapperVariant, wrapper, stackedReplicas])
//   let graphs = extend_query(proofs, [stackedChungParams, stackedSDRParams])
//   let windows = extend_query(graphs, window_size_mib_choices.map(d => ({window_size_mib: +d})))
//   let query = extend_query(windows, [poseidon, pedersen])

//   return query
// }

createJsonDownloadButton(combos)

solved_many_pre = (await solve_manys(combos))
  .map(d => {
    d.construction = `${d.graph_name}_${d.proof_name}`
    return d
  })

utility_fun = (data) => ev(utility_raw, data)

solved_many = solved_many_pre
  // .filter(d => window_size_mib_config.some(c => d['window_size_mib'] === +c))
  .map(d => {
    const utility = utility_fun(d)
    return Object.assign({}, d, {utility: utility})
  })
// solved_manys = (await solve_manys(combos)).flat()

createJsonDownloadButton(solved_many)

// mtree_query = {
//   let query = [constants]
//   const proofs = [wrapper, wrapperVariant, stackedReplicas]
//   const post_mtree_layers_cached = [...Array(10)].map((_, i) => ({post_mtree_layers_cached: i+20}))

//   query = extend_query(query, proofs, post_mtree_layers_cached, [stackedChungParams])

//   return query
// }

// mtree_solved = (await solve_many(mtree_query)).map(d => d[0])


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

| name | val |
| ---- | --- |
${Object.keys(result).sort()
  .map(d => `| ${!starting_assignments[d] ? `**${d}**` : d} | ${result[d]} |\n`)}
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
            "y": {"field": "construction", "type": "nominal"},
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
        "y": {"field": "construction", "type": "nominal"},
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

_f = (d) => typeof d == 'number' || !Number.isNaN(+d) ? d3.format('0.6~f')(d) : d

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
