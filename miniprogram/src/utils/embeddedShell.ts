import { isMiniProgramWebview, syncMiniProgramTitle } from './miniProgramBridge'

const SHELL_SOURCE_KEY = 'store_h5_entry_source'

const STORE_PAGE_TITLE_MAP: Record<string, string> = {
  '/pages/store/home': '店铺首页',
  '/pages/store/product': '商品详情',
  '/pages/store/confirm': '确认订单',
  '/pages/store/cart': '购物车',
  '/pages/store/my-orders': '我的订单',
  '/pages/store/order-detail': '订单详情',
  '/pages/store/address-list': '收货地址',
  '/pages/store/address-edit': '编辑地址',
}

function getHashQueryParams() {
  if (typeof window === 'undefined') {
    return new URLSearchParams()
  }

  const hash = window.location.hash || ''
  const queryString = hash.includes('?') ? hash.slice(hash.indexOf('?') + 1) : ''
  return new URLSearchParams(queryString)
}

function getCurrentHashPath() {
  if (typeof window === 'undefined') {
    return ''
  }

  const hash = window.location.hash || ''
  const routePart = hash.startsWith('#') ? hash.slice(1) : hash
  const [path] = routePart.split('?')
  if (!path) {
    return ''
  }
  return path.startsWith('/') ? path : `/${path}`
}

export function cacheEntrySourceFromUrl() {
  const source = getHashQueryParams().get('source') || ''
  if (source) {
    uni.setStorageSync(SHELL_SOURCE_KEY, source)
  }
  return source
}

export function isXcxShellSource() {
  const source = cacheEntrySourceFromUrl() || uni.getStorageSync(SHELL_SOURCE_KEY) || ''
  return String(source).trim() === 'xcx_shell'
}

export function shouldUseEmbeddedShellTitle() {
  return isMiniProgramWebview() && isXcxShellSource()
}

export function getStorePageTitleByRoute(pathValue?: string) {
  const currentPath = (pathValue || getCurrentHashPath()).trim()
  return STORE_PAGE_TITLE_MAP[currentPath] || '商家助手'
}

export async function syncCurrentPageTitle(pathValue?: string) {
  if (typeof window !== 'undefined') {
    document.title = getStorePageTitleByRoute(pathValue)
  }

  if (!shouldUseEmbeddedShellTitle()) {
    return
  }

  await syncMiniProgramTitle(getStorePageTitleByRoute(pathValue))
}

