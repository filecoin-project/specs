#!/usr/bin/env node

const chokidar = require ('chokidar')
const jsdom = require('jsdom')
const path = require('path')
const fs = require('fs')

const { buildTocModel } = require('./build-model')

async function processHtml (pathToFile) {
  const dom = await jsdom.JSDOM.fromFile(pathToFile)
  const model = buildTocModel(dom.window.document.querySelector('.markdown'))
  const dest = path.join(__dirname, '../../data/toc.json')
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

// function run () {
//   const args = process.argv.slice(2)
//   if (args[0] === '--watch') {
//     chokidar.watch()
//   }
// }

processHtml(path.join(__dirname, '../../public/index.html'))
