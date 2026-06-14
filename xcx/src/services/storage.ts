import Taro from '@tarojs/taro'
import type { MiniShellConfig } from '@/types/config'
import type { UserSession } from '@/types/auth'

const SHELL_CONFIG_KEY = 'xcx_shell_config'
const USER_SESSION_KEY = 'xcx_user_session'

export function getStoredConfig(): MiniShellConfig | null {
  try {
    const rawValue = Taro.getStorageSync(SHELL_CONFIG_KEY)
    if (!rawValue) {
      return null
    }
    return JSON.parse(rawValue) as MiniShellConfig
  } catch (error) {
    console.error('[XCX][Storage] 读取壳配置失败', error)
    return null
  }
}

export function setStoredConfig(config: MiniShellConfig) {
  Taro.setStorageSync(SHELL_CONFIG_KEY, JSON.stringify(config))
}

export function getStoredUserSession(): UserSession | null {
  try {
    const rawValue = Taro.getStorageSync(USER_SESSION_KEY)
    if (!rawValue) {
      return null
    }
    return JSON.parse(rawValue) as UserSession
  } catch (error) {
    console.error('[XCX][Storage] 读取用户会话失败', error)
    return null
  }
}

export function setStoredUserSession(session: UserSession) {
  Taro.setStorageSync(USER_SESSION_KEY, JSON.stringify(session))
}

export function clearStoredUserSession() {
  Taro.removeStorageSync(USER_SESSION_KEY)
}
