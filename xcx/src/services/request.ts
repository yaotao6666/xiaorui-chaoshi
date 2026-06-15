import Taro from '@tarojs/taro'
import type { ApiEnvelope } from '@/types/auth'
import { getStoredUserSession } from '@/services/storage'

export async function postJson<T>(
  url: string,
  data: Record<string, unknown>,
  options: { withAuth?: boolean } = {}
) {
  const { withAuth = false } = options
  const session = withAuth ? getStoredUserSession() : null
  const response = await Taro.request<ApiEnvelope<T> | T>({
    url,
    method: 'POST',
    data,
    header: {
      'Content-Type': 'application/json',
      ...(session?.token ? { Authorization: `Bearer ${session.token}` } : {}),
    },
  })

  const responseData = response.data as ApiEnvelope<T> | T

  if (typeof (responseData as ApiEnvelope<T>)?.code === 'number') {
    const envelope = responseData as ApiEnvelope<T>
    if (envelope.code !== 0) {
      throw new Error(envelope.message || '请求失败')
    }
    return envelope.data as T
  }

  return responseData as T
}
