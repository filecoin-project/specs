const globby = require('globby');
const execa = require('execa')
const path = require('path')
const fs = require('fs')

const runMmd = async () => {
	const paths = await globby(['content/**/*.mmd']);
    await Promise.all(paths.map(p => {
        const outDir =path.dirname(p).replace('content/', 'static/_gen/diagrams/')
        const outFile = path.basename(p).replace('.mmd', '.svg')
        fs.mkdirSync(outDir, { recursive: true })
        return execa('mmdc', [
            '-i', p,
            '-o', path.join(outDir, outFile)
        ], { preferLocal: true });
    }))
}


const runDot = async () => {
    const paths = await globby(['content/**/*.dot']);
    await Promise.all(paths.map(p => {
        const outDir =path.dirname(p).replace('content/', 'static/_gen/diagrams/')
        const outFile = path.basename(p).replace('.dot', '.svg')
        fs.mkdirSync(outDir, { recursive: true })
        return execa('graphviz', [
            '-Tsvg',
            `-o${path.join(outDir, outFile)}`,
            p
        ], { preferLocal: true });
    }))
}


const run = async () => {
    console.log('Processing *.{mmd,dot}');
    await Promise.all([
        runDot(),
        runMmd()
    ])
    console.log('Done');
}

run()