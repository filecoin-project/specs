#!/usr/bin/env node
import { run as mmdc } from '@mermaid-js/mermaid-cli'
import globby from 'globby'
import path from 'path'
import fs from 'fs'
import { renderGraphFromSource } from 'graphviz-cli'
import * as url from 'node:url';

const runMmd = async (p) => {
  const outDir = path.dirname(p).replace('content/', 'static/_gen/diagrams/')
  const outFile = path.basename(p).replace('.mmd', '.svg')
  fs.mkdirSync(outDir, { recursive: true })
  const config = process.env.CI ? { puppeteerConfig: 'tools/pptr.config' } : {}
  return await mmdc(p, path.join(outDir, outFile), {
    puppeteerConfig: 'tool/ppt.config',
  })
}

const runMmdAll = async () => {
  const paths = await globby(['content/**/*.mmd'])
  await Promise.all(paths.map(runMmd))
}

const runDot = async (p) => {
  const outDir = path.dirname(p).replace('content/', 'static/_gen/diagrams/')
  const outFile = path.basename(p).replace('.dot', '.svg')
  fs.mkdirSync(outDir, { recursive: true })

  return await renderGraphFromSource(
    { name: p },
    { format: 'svg', name: path.join(outDir, outFile) }
  )
}

const runDotAll = async () => {
  const paths = await globby(['content/**/*.dot'])
  await Promise.all(paths.map(runDot))
}

const run = async () => {
  const args = process.argv.slice(2)
  console.log('Processing *.{mmd,dot}')
  console.time('Processed *.{mmd,dot}')
  await Promise.all([runDotAll(), runMmdAll()])
  console.timeEnd('Processed *.{mmd,dot}')
}

export const configureWatcher = (watcher) => {
  watcher.on('all', async (_, p) => {
    const ext = path.extname(p)
    switch (ext) {
      case '.dot':
        await runDot(p)
        console.log('done ', p)
        break
      case '.mmd':
        await runMmd(p)
        console.log('done ', p)
        break
      default:
        break
    }
  })
  watcher.add('content/**/*.{mmd,dot}')
}

// run as script, so do the thing
if (import.meta.url.startsWith('file:')) { // (A)
  const modulePath = url.fileURLToPath(import.meta.url);
  if (process.argv[1] === modulePath) { // (B)
    run()
  }
}
