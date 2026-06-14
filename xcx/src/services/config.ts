import type { MiniShellConfig } from '@/types/config'
import { getStoredConfig, setStoredConfig } from '@/services/storage'

export const DEFAULT_SHELL_CONFIG: MiniShellConfig = {
  apiBaseUrl: 'http://127.0.0.1:8081',
  h5EntryUrl: 'http://localhost:3000/#/pages/store/home',
  merchantId: '1',
  source: 'xcx_shell',
  autoEnter: true,
  debugToken: '',
  debugOpenid: '',
}

export function getShellConfig(): MiniShellConfig {
  const storedConfig = getStoredConfig()
  if (!storedConfig) {
    return DEFAULT_SHELL_CONFIG
  }

  return {
    ...DEFAULT_SHELL_CONFIG,
    ...storedConfig,
  }
}

export function saveShellConfig(config: MiniShellConfig) {
  setStoredConfig(config)
}
