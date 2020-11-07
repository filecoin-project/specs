#!/usr/bin/env node
const vfile = require('to-vfile')
const remark = require('remark')
const toc = require('remark-toc')

const readme = vfile.readSync('README.md')

// inject toc into readme
remark()
  .use(toc, { tight: true })
  .process(readme, function (err) {
    if (err) throw err
    vfile.writeSync(readme)
    console.log('Updated README.md')
  })
