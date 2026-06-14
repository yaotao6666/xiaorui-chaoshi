export interface ApiEnvelope<T> {
  code?: number
  message?: string
  data?: T
}

export interface WechatUserInfo {
  id: number
  openid: string
  nickname: string
}

export interface WechatLoginPayload {
  token: string
  app_id: string
  user: WechatUserInfo
}

export interface UserSession {
  token: string
  openid: string
  appId: string
  nickname: string
  obtainedAt: number
}
