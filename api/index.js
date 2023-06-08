const Router = require('./src/router')
const dlv = require('dlv')
const get = require('./src/get')
const map = require('p-map')

addEventListener('fetch', (event) => {
  event.respondWith(handleRequest(event))
})

async function handleRequest(event) {
  const r = new Router()
  r.get('.*/cov', () => cov(event))
  r.get('.*/github', () => github(event))
  r.get('.*/releases', () => releases(event))
  r.get('/', () => new Response('Hello worker!')) // return a default message for the root route

  try {
    return await r.route(event.request)
  } catch (err) {
    console.log('handleRequest -> err', err.stack)
    return new Response(err.message, {
      status: 500,
      statusText: 'internal server error',
      headers: {
        'content-type': 'text/plain',
      },
    })
  }
}

async function cov(event) {
  const url = new URL(event.request.url)
  // https://github.com/filecoin-project/lotus
  const [owner, repo] = url.searchParams.get('repo').split('/').slice(3)
  const headers = {
    'User-Agent': 'ianconsolata',
    Accept: 'application/json',
    Authorization: `Bearer ${CODECOV_TOKEN}`,
  }
  const repo_resp = await get(event, {
    url: `https://api.codecov.io/api/v2/github/${owner}/repos/${repo}/`,
    transform: (data) => {
      const out = {
        branch: dlv(data, 'branch', 'master'),
        lang: dlv(data, 'language', 'N/A'),
      }
      return out
    },
    headers,
  })
  const repo_data = await repo_resp.json()

  const cov_data = await get(event, {
    url: `https://api.codecov.io/api/v2/github/${owner}/repos/${repo}/branches/${repo_data.branch}/`,
    transform: (data) => {
      const out = {
        cov: dlv(data, 'head_commit.totals.coverage', 0),
        ci: dlv(data, 'head_commit.ci_passed', false),
        repo: repo,
        org: owner,
        lang: repo_data.lang,
        branch: repo_data.branch,
      }
      return out
    },
    headers,
  })
  return cov_data
}

async function github(event) {
  const url = new URL(event.request.url)
  const file = url.searchParams.get('file').split('/')
  // https://github.com/filecoin-project/lotus/blob/master/paychmgr/paych.go
  const repo = file.slice(3, 5).join('/')
  const path = file.slice(7).join('/')
  const ref = file[6]
  const headers = {
    'User-Agent': 'ianconsolata',
    Authorization: `token ${GITHUB_TOKEN}`,
  }

  const treeUrlRsp = await get(event, {
    url: `https://api.github.com/repos/${repo}/commits?sha=${ref}&per_page=1&page=1`,
    headers,
  })
  const treeUrl = await treeUrlRsp.json()

  const data = await get(event, {
    url: `https://api.github.com/repos/${repo}/contents/${path}?ref=${ref}`,
    transform: (data) => {
      return {
        content: data.content,
        size: data.size,
        url: `https://github.com/${repo}/tree/${treeUrl[0].sha}/${path}`,
      }
    },
    headers,
  })
  return data
}

async function releases(event) {
  const headers = {
    'User-Agent': 'ianconsolata',
    Authorization: `token ${GITHUB_TOKEN}`,
  }
  const rsp = await get(event, {
    url: `https://api.github.com/repos/filecoin-project/specs/releases?per_page=100&page=1`,
    headers,
    force: true,
    transform: async (releases) => {
      return (
        await map(
          releases,
          async (r) => {
            const status = await get(event, {
              url: `https://api.github.com/repos/filecoin-project/specs/commits/${r.tag_name}/status`,
              headers,
            })
            const statusData = await status.json()
            const preview = dlv(statusData, 'statuses').find(
              (d) => d.description === 'Preview ready'
            )

            if (preview) {
              return {
                state: dlv(statusData, 'state'),
                preview: preview.target_url,
                tag_name: r.tag_name,
                name: r.name,
                author: r.author,
                created_at: r.created_at,
                published_at: r.published_at,
                body: r.body,
              }
            }
          },
          { concurrency: 3 }
        )
      ).filter(Boolean)
    },
  })

  return rsp
}
