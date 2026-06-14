/// <reference types="vite/client" />

/**
 * 小程序运行环境配置
 * 统一收口 HTTP 地址，避免业务代码分散维护。
 */

const DEFAULT_API_BASE_URL = 'http://127.0.0.1:8080'
const DEFAULT_APP_NAME = '商家助手'

function trimTrailingSlash(url: string): string {
  return url.replace(/\/+$/, '')
}

function getEnvValue(
  key: 'VITE_API_BASE_URL' | 'VITE_APP_NAME' | 'VITE_H5_ENTRY_URL',
  fallback: string
): string {
  const value = import.meta.env?.[key]
  return typeof value === 'string' && value.trim() ? value.trim() : fallback
}

export const API_BASE_URL = trimTrailingSlash(
  getEnvValue('VITE_API_BASE_URL', DEFAULT_API_BASE_URL)
)

export const H5_ENTRY_URL = trimTrailingSlash(
  getEnvValue('VITE_H5_ENTRY_URL', API_BASE_URL)
)

export const APP_NAME = getEnvValue('VITE_APP_NAME', DEFAULT_APP_NAME)
