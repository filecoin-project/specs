#!/usr/bin/env node
const chokidar = require('chokidar')
const diagrams = require('./diagrams')
const toc = require('./toc')

watch()

function watch () {
  const watcher = chokidar.watch([], {
    awaitWriteFinish: {
      stabilityThreshold: 1000,
      pollInterval: 100
    },
    ignoreInitial: true
  })

  watcher
    .on('error', err => console.error('error watching: ', err))
    .on('all', (event, p) => console.log(event, p))

  toc.configureWatcher(watcher)
  diagrams.configureWatcher(watcher)

  console.log('Watching', watcher.getWatched())
}