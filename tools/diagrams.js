#!/usr/bin/env node

const globby = require('globby');
const execa = require('execa')
const path = require('path')
const fs = require('fs')

const runMmd = (p) => {
    const outDir = path.dirname(p).replace('content/', 'static/_gen/diagrams/')
    const outFile = path.basename(p).replace('.mmd', '.svg')

    fs.mkdirSync(outDir, { recursive: true })

    return execa('mmdc', [
        '-p', 'tools/pptr.config',
        '-i', p,
        '-o', path.join(outDir, outFile)
    ], { preferLocal: true, stdio: "inherit" });
}

const runMmdAll = async () => {
    const paths = await globby(['content/**/*.mmd']);
    await Promise.all(paths.map(runMmd))
}

const runDot = (p) => {
    const outDir =path.dirname(p).replace('content/', 'static/_gen/diagrams/')
    const outFile = path.basename(p).replace('.dot', '.svg')
    fs.mkdirSync(outDir, { recursive: true })
    return execa('graphviz', [
        '-Tsvg',
        `-o${path.join(outDir, outFile)}`,
        p
    ], { preferLocal: true });
}

const runDotAll = async () => {
    const paths = await globby(['content/**/*.dot']);
    await Promise.all(paths.map(runDot))
}

const run = async () => {
    const args = process.argv.slice(2);
    console.log('Processing *.{mmd,dot}');
    console.time('Processed *.{mmd,dot}')
    await Promise.all([
        runDotAll(),
        runMmdAll()
    ])
    console.timeEnd('Processed *.{mmd,dot}')
}

module.exports.configureWatcher = (watcher) => {
  watcher.on('all', async (_, p) => {
    const ext = path.extname(p)
    switch (ext) {
        case ".dot":
            await runDot(p)
            console.log('done ', p)
            break;
        case ".mmd":
            await runMmd(p)
            console.log('done ', p)
            break
        default:
            break;
    }
  })
  watcher.add('content/**/*.{mmd,dot}')
}

// run as script, so do the thing
if (require.main === module) {
    run()
}
