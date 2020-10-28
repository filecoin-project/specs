const Router = require('./router')
const dlv = require('dlv')
const merge = require('merge-options')
const nanoid = require('nanoid/non-secure')

const cacheBust = nanoid.nanoid()
/**
 * Example of how router can be used in an application
 *  */
addEventListener('fetch', (event) => {
  event.respondWith(handleRequest(event))
})

async function handleRequest(event) {
  const r = new Router()
  // Replace with the appropriate paths and handlers
  r.get('.*/cov', () => cov(event))
  r.get('.*/github', () => github(event))
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
  const repo = url.searchParams.get('repo').split('/').slice(3).join('/')
  const data = await get(event, {
    url: `https://codecov.io/api/gh/${repo}`,
    transform: (data) => {
      const out = {
        cov: dlv(data, 'commit.totals.c', 0),
        ci: dlv(data, 'commit.ci_passed', false),
        repo: dlv(data, 'repo.name', 'N/A'),
        org: dlv(data, 'owner.username', 'N/A'),
        lang: dlv(data, 'repo.language', 'N/A'),
      }
      return out
    },
  })
  return data
}

async function github(event) {
  const url = new URL(event.request.url)
  const file = url.searchParams.get('file').split('/')
  // https://github.com/filecoin-project/lotus/blob/master/paychmgr/paych.go
  const repo = file.slice(3, 5).join('/')
  const path = file.slice(7).join('/')
  const ref = file[6]
  const headers = {
    'User-Agent': 'hugomrdias',
    // Authorization: `token ${GITHUB_TOKEN}`
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

async function get(event, options) {
  const { url, transform, force, headers } = merge(
    {
      url: '',
      transform: (d) => d,
      force: false,
      headers: {},
    },
    options
  )

  const cache = caches.default
  const cacheKey = url + cacheBust
  const cacheTTL = 86400 * 2 // 2 days
  const cacheRevalidateTTL = 3600 * 2 // 2 hours
  const cachedResponse = await cache.match(cacheKey)

  if (force || !cachedResponse) {
    console.log('Cache miss for ', cacheKey)
    // if not in cache get from the origin
    const response = await fetch(url, {
      headers: {
        ...headers,
        'If-None-Match': cachedResponse
          ? cachedResponse.headers.get('ETag')
          : null,
      },
    })

    if (response.ok) {
      const { headers } = response
      const contentType = headers.get('content-type') || ''

      if (contentType.includes('application/json')) {
        // transform the data
        const data = transform(await response.json())

        // build new response with the transformed body
        const transformedResponse = new Response(JSON.stringify(data), {
          headers: {
            'Content-Type': 'application/json;charset=UTF-8',
            'Cache-Control': `max-age=${cacheTTL}`,
            'X-RateLimit-Limit': headers.get('X-RateLimit-Limit'),
            'X-RateLimit-Remaining': headers.get('X-RateLimit-Remaining'),
            'X-RateLimit-Reset': headers.get('X-RateLimit-Reset'),
            ETag: headers.get('ETag'),
          },
        })

        // save response to cache
        event.waitUntil(cache.put(cacheKey, transformedResponse.clone()))

        return transformedResponse
      } else {
        throw new Error(
          `Request error content type not supported. ${contentType}`
        )
      }
    } else if (response.status === 304) {
      // renew cache response
      event.waitUntil(cache.put(cacheKey, cachedResponse.clone()))
      return cachedResponse.clone()
    } else {
      return response
    }
  } else {
    console.log('Cache hit for ', cacheKey, cachedResponse.headers.get('age'))
    const cacheAge = cachedResponse.headers.get('age')

    if (cacheAge > cacheRevalidateTTL) {
      console.log('Cache is too old, revalidating...')
      event.waitUntil(get(event, { url, transform, force: true }))
    }
    return cachedResponse
  }
}
