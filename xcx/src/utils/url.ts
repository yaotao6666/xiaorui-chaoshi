export interface EmbeddedH5Params {
  merchantId: string
  token: string
  openid: string
  source: string
  appId?: string
}

function normalizeUrl(url: string) {
  return url.trim().replace(/\s+/g, '')
}

export function maskToken(token: string) {
  if (!token) {
    return '未获取'
  }
  if (token.length <= 12) {
    return `${token.slice(0, 4)}...${token.slice(-2)}`
  }
  return `${token.slice(0, 6)}...${token.slice(-6)}`
}

export function buildEmbeddedH5Url(entryUrl: string, params: EmbeddedH5Params) {
  const safeUrl = normalizeUrl(entryUrl)
  const query = new URLSearchParams()

  query.set('merchant_id', params.merchantId)
  query.set('token', params.token)
  query.set('openid', params.openid)
  query.set('source', params.source)

  if (params.appId) {
    query.set('app_id', params.appId)
  }

  if (safeUrl.includes('#')) {
    const [originPart, hashPart] = safeUrl.split('#')
    const routePart = hashPart || '/pages/store/home'
    const [hashPath, hashQuery = ''] = routePart.split('?')
    const mergedQuery = new URLSearchParams(hashQuery)

    query.forEach((value, key) => {
      mergedQuery.set(key, value)
    })

    return `${originPart}#${hashPath}${mergedQuery.toString() ? `?${mergedQuery.toString()}` : ''}`
  }

  const [pathPart, currentQuery = ''] = safeUrl.split('?')
  const mergedQuery = new URLSearchParams(currentQuery)
  query.forEach((value, key) => {
    mergedQuery.set(key, value)
  })

  return `${pathPart}${mergedQuery.toString() ? `?${mergedQuery.toString()}` : ''}`
}
