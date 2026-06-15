declare global {
  interface Window {
    wx?: {
      miniProgram?: {
        navigateTo: (options: { url: string; success?: () => void; fail?: (error: unknown) => void }) => void
      }
      ready?: (callback: () => void) => void
    }
    __wxjs_environment?: string
  }
}

const WECHAT_JSSDK_URL = 'https://res.wx.qq.com/open/js/jweixin-1.6.0.js'

let sdkLoadingPromise: Promise<void> | null = null

function isH5Runtime() {
  return typeof window !== 'undefined' && typeof document !== 'undefined'
}

export function isWechatBrowser() {
  if (!isH5Runtime()) {
    return false
  }
  return /micromessenger/i.test(window.navigator.userAgent || '')
}

export function isMiniProgramWebview() {
  if (!isWechatBrowser()) {
    return false
  }
  const userAgent = window.navigator.userAgent || ''
  return window.__wxjs_environment === 'miniprogram' || /miniProgram/i.test(userAgent)
}

function loadWechatJSSDK() {
  if (!isH5Runtime()) {
    return Promise.reject(new Error('当前环境不支持小程序桥接'))
  }
  if (window.wx?.miniProgram?.navigateTo) {
    return Promise.resolve()
  }
  if (sdkLoadingPromise) {
    return sdkLoadingPromise
  }

  sdkLoadingPromise = new Promise<void>((resolve, reject) => {
    const existingScript = document.querySelector<HTMLScriptElement>(`script[src="${WECHAT_JSSDK_URL}"]`)
    if (existingScript) {
      existingScript.addEventListener('load', () => resolve(), { once: true })
      existingScript.addEventListener('error', () => reject(new Error('微信 JSSDK 加载失败')), { once: true })
      return
    }

    const script = document.createElement('script')
    script.src = WECHAT_JSSDK_URL
    script.async = true
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('微信 JSSDK 加载失败'))
    document.head.appendChild(script)
  })

  return sdkLoadingPromise
}

export async function openXcxPaymentPage(params: {
  orderId: number
  merchantId: number
  returnTarget: string
}) {
  if (!isMiniProgramWebview()) {
    throw new Error('当前不是小程序壳环境，请使用小程序完成支付')
  }

  await loadWechatJSSDK()
  if (!window.wx?.miniProgram?.navigateTo) {
    throw new Error('当前环境缺少小程序跳转能力')
  }

  const targetUrl = `/pages/payment/index?orderId=${encodeURIComponent(String(params.orderId))}&merchantId=${encodeURIComponent(String(params.merchantId))}&returnTarget=${encodeURIComponent(params.returnTarget)}`
  return new Promise<void>((resolve, reject) => {
    window.wx?.miniProgram?.navigateTo({
      url: targetUrl,
      success: () => resolve(),
      fail: (error) => reject(error instanceof Error ? error : new Error('跳转小程序支付页失败'))
    })
  })
}
