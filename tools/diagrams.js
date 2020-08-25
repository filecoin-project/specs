#!/usr/bin/env node

const globby = require('globby');
const execa = require('execa')
const path = require('path')
const fs = require('fs')
const chokidar = require('chokidar');

const runMmd = (p) => {
    const outDir =path.dirname(p).replace('content/', 'static/_gen/diagrams/')
    const outFile = path.basename(p).replace('.mmd', '.svg')

    fs.mkdirSync(outDir, { recursive: true })
    
    return execa('mmdc', [
        '-i', p,
        '-o', path.join(outDir, outFile)
    ], { preferLocal: true });
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

    if(args[0] === '--all') {
        console.log('Processing *.{mmd,dot}');
        await Promise.all([
            runDotAll(),
            runMmdAll()
        ])
        console.log('Done');
    }

    if(args[0] === '--watch') {
        chokidar
            .watch('content/**/*.{dot,mmd}', {
                awaitWriteFinish: {
                    stabilityThreshold: 1000,
                    pollInterval: 100
                },
                ignoreInitial: true
            })
            .on('all', async (event, p) => {
                console.log(event, p);
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
            .on('ready', () => {
                console.log(`Watching 'content/**/*.{dot,mmd}'`);
            })
            .on('error', err => console.error('error watching: ', err));
    }
    
}

run()
