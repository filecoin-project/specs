import { buildTocModel } from './content-model'
import { buildToc } from './toc.js'
import { buildDashboard } from './dashboard-spec'

const model = buildTocModel('.markdown')
buildToc({tocSelector:'.toc', model })
buildDashboard('#dashboard-container', model)