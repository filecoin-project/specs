#!/usr/bin/env node

const chokidar = require('chokidar')
const jsdom = require('jsdom')
const path = require('path')
const fs = require('fs')
const { buildTocModel } = require('./toc/build-model')

const src = path.join(__dirname, '../public/index.html')
const dest = path.join(__dirname, '../data/toc.json')

run(src, dest)

async function run (src, dest) {
  const args = process.argv.slice(2)
  if (args[0] === '--watch') {
    chokidar.watch(src, {
      awaitWriteFinish: {
        stabilityThreshold: 1000,
        pollInterval: 100
      },
      ignoreInitial: true
    })
      .on('all', async (event, p) => {
        console.log(event, p)
        await processHtml(src, dest)
      })
      .on('ready', () => {
        console.log(`Watching ${src}`)
      })
      .on('error', err => console.error('error watching: ', err))
  } else {
    console.time('Building toc.json')
    await processHtml(src, dest)
    console.timeEnd('Building toc.json')
  }
}

async function processHtml (src, dest) {
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
    } catch (err) {
      return console.error(err)
    }
    console.log('Updated toc.json')
  }
}
