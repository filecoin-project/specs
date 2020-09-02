#!/usr/bin/env node

const jsdom = require('jsdom')
const path = require('path')
const fs = require('fs')
const { buildTocModel } = require('./toc/build-model')

const src = 'public/index.html'
const dest = 'data/toc.json'

// run as script, so do the thing
if (require.main === module) {
  run(src, dest)
}

async function run (src, dest) {
  console.log('Building toc.json')
  console.time('Built toc.json')
  await buildToc(src, dest)
  console.timeEnd('Built toc.json')
}

async function buildToc (src, dest) {
  if (!fs.existsSync(path.dirname(dest))) {
    fs.mkdirSync(path.dirname(dest))
  }
  const dom = await jsdom.JSDOM.fromFile(src)
  const model = buildTocModel(dom.window.document.querySelector('.markdown'))
  const json = JSON.stringify(model, null, 2)
  let prev = Buffer.from('')
  try {
    prev = fs.readFileSync(dest)
  } catch {
    // ok, no previous data.
  }
  if (!Buffer.from(json).equals(prev)) {
    try {
      fs.writeFileSync(dest, json)
      console.log(`Updated toc.json`)
    } catch (err) {
      return console.error(err)
    }
  }
}

module.exports.configureWatcher = (watcher) => {
  watcher.on('all', async (_, p) => {
    if (p === src) {
      await buildToc(src, dest)
    }
  })

  watcher.add(src)
}
