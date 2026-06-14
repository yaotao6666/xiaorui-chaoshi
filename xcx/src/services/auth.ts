import Taro from '@tarojs/taro'
import { getShellConfig } from '@/services/config'
import { postJson } from '@/services/request'
import { clearStoredUserSession, getStoredUserSession, setStoredUserSession } from '@/services/storage'
import type { UserSession, WechatLoginPayload } from '@/types/auth'

function mapPayloadToSession(payload: WechatLoginPayload): UserSession {
  return {
    token: payload.token,
    openid: payload.user.openid,
    appId: payload.app_id || '',
    nickname: payload.user.nickname || '微信用户',
    obtainedAt: Date.now(),
  }
}

function createDebugSession() {
  const config = getShellConfig()
  if (!config.debugToken || !config.debugOpenid) {
    throw new Error('当前不是微信小程序环境，请先在配置页填写调试 token 和调试 openid')
  }

  return {
    token: config.debugToken,
    openid: config.debugOpenid,
    appId: 'debug-preview',
    nickname: '预览调试用户',
    obtainedAt: Date.now(),
  } satisfies UserSession
}

export async function ensureUserSession(forceRefresh = false): Promise<UserSession> {
  if (!forceRefresh) {
    const storedSession = getStoredUserSession()
    if (storedSession?.token && storedSession.openid) {
      return storedSession
    }
  }

  if (process.env.TARO_ENV !== 'weapp') {
    const debugSession = createDebugSession()
    setStoredUserSession(debugSession)
    return debugSession
  }

  console.info('[XCX][Auth] 开始调用 Taro.login')
  const loginResult = await Taro.login()
  if (!loginResult.code) {
    throw new Error('未获取到微信登录 code')
  }

  const { apiBaseUrl } = getShellConfig()
  const payload = await postJson<WechatLoginPayload>(`${apiBaseUrl}/api/v1/auth/user/wechat-login`, {
    code: loginResult.code,
  })
  const session = mapPayloadToSession(payload)
  setStoredUserSession(session)
  console.info('[XCX][Auth] 微信登录成功，已写入本地会话')
  return session
}

export function clearUserSession() {
  clearStoredUserSession()
}
