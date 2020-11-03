'use strict'

const merge = require('merge-options')
const nanoid = require('nanoid/non-secure')

const cacheBust = nanoid.nanoid()

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
        const data = await transform(await response.json())

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

module.exports = get
