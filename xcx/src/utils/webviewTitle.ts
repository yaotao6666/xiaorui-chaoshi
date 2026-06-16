const DEFAULT_WEBVIEW_TITLE = '商家助手'

const ROUTE_TITLE_MAP: Record<string, string> = {
  '/pages/store/home': '店铺首页',
  '/pages/store/product': '商品详情',
  '/pages/store/confirm': '确认订单',
  '/pages/store/cart': '购物车',
  '/pages/store/my-orders': '我的订单',
  '/pages/store/order-detail': '订单详情',
  '/pages/store/address-list': '收货地址',
  '/pages/store/address-edit': '编辑地址',
}

function normalizeRoutePath(pathValue: string): string {
  const trimmed = String(pathValue || '').trim()
  if (!trimmed) {
    return ''
  }
  const withoutHash = trimmed.replace(/^.*#/, '')
  const [pathname] = withoutHash.split('?')
  if (!pathname) {
    return ''
  }
  return pathname.startsWith('/') ? pathname : `/${pathname}`
}

export function getDefaultWebviewTitle() {
  return DEFAULT_WEBVIEW_TITLE
}

export function getWebviewTitleByPath(pathValue: string) {
  const normalizedPath = normalizeRoutePath(pathValue)
  return ROUTE_TITLE_MAP[normalizedPath] || DEFAULT_WEBVIEW_TITLE
}

export function getWebviewTitleFromTargetUrl(targetUrl: string) {
  const trimmed = String(targetUrl || '').trim()
  if (!trimmed) {
    return DEFAULT_WEBVIEW_TITLE
  }

  try {
    const parsed = new URL(trimmed)
    const hashPath = parsed.hash.startsWith('#') ? parsed.hash.slice(1) : parsed.hash
    return getWebviewTitleByPath(hashPath)
  } catch {
    return getWebviewTitleByPath(trimmed)
  }
}

