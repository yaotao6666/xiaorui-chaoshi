import React, { useCallback, useMemo, useRef, useState } from 'react'
import { Button, Text, View } from '@tarojs/components'
import Taro, { useDidShow } from '@tarojs/taro'
import classNames from 'classnames'
import { clearUserSession, ensureUserSession } from '@/services/auth'
import { getShellConfig } from '@/services/config'
import { getStoredUserSession } from '@/services/storage'
import { buildEmbeddedH5Url, maskToken } from '@/utils/url'
import styles from './index.module.scss'

type LaunchStatus = 'idle' | 'loading' | 'success' | 'error'

function getStatusText(status: LaunchStatus) {
  switch (status) {
    case 'loading':
      return '正在登录并准备 H5 地址'
    case 'success':
      return '登录成功，可直接进入 H5'
    case 'error':
      return '准备失败，请查看说明后重试'
    default:
      return '等待启动'
  }
}

const IndexPage: React.FC = () => {
  const [status, setStatus] = useState<LaunchStatus>('idle')
  const [errorMessage, setErrorMessage] = useState('')
  const [launchCount, setLaunchCount] = useState(0)
  const autoLaunchDoneRef = useRef(false)

  const shellConfig = useMemo(() => getShellConfig(), [launchCount])
  const cachedSession = useMemo(() => getStoredUserSession(), [launchCount])

  const openWebview = useCallback(
    async (forceRefresh = false) => {
      try {
        setStatus('loading')
        setErrorMessage('')

        const session = await ensureUserSession(forceRefresh)
        const targetUrl = buildEmbeddedH5Url(shellConfig.h5EntryUrl, {
          merchantId: shellConfig.merchantId,
          token: session.token,
          openid: session.openid,
          source: shellConfig.source,
          appId: session.appId,
        })

        setStatus('success')
        console.info('[XCX][Launch] H5 地址已生成，准备进入容器页')
        Taro.navigateTo({
          url: `/pages/webview/index?target=${encodeURIComponent(targetUrl)}`,
        })
      } catch (error: unknown) {
        console.error('[XCX][Launch] 启动失败', error)
        setStatus('error')
        setErrorMessage(error instanceof Error ? error.message : '启动失败')
      }
    },
    [shellConfig]
  )

  const handleClearAndRetry = () => {
    clearUserSession()
    setLaunchCount((value) => value + 1)
    openWebview(true)
  }

  useDidShow(() => {
    if (shellConfig.autoEnter && !autoLaunchDoneRef.current) {
      autoLaunchDoneRef.current = true
      openWebview(false)
    }
  })

  return (
    <View className={styles.container}>
      <View className={styles.heroCard}>
        <Text className={styles.eyebrow}>XCX 验证壳</Text>
        <Text className={styles.title}>小程序先登录，再把 token 带进现有 H5</Text>
        <Text className={styles.desc}>
          当前页面用于验证“小程序无感登录 + WebView 打开现有 H5”主链路。首次进入会优先复用本地会话，不足时再调用微信登录接口。
        </Text>
      </View>

      <View className={styles.grid}>
        <View className={styles.panel}>
          <Text className={styles.panelTitle}>当前状态</Text>
          <View
            className={classNames(
              styles.statusTag,
              status === 'idle' && styles.statusTagWarning,
              status === 'error' && styles.statusTagError
            )}
          >
            <Text>{getStatusText(status)}</Text>
          </View>
          {!!errorMessage && <Text className={styles.tipText}>{errorMessage}</Text>}
        </View>

        <View className={styles.panel}>
          <Text className={styles.panelTitle}>启动参数</Text>
          <View className={styles.kvList}>
            <View className={styles.kvItem}>
              <Text className={styles.kvLabel}>API 基址</Text>
              <Text className={styles.kvValue}>{shellConfig.apiBaseUrl}</Text>
            </View>
            <View className={styles.kvItem}>
              <Text className={styles.kvLabel}>H5 入口</Text>
              <Text className={styles.kvValue}>{shellConfig.h5EntryUrl}</Text>
            </View>
            <View className={styles.kvItem}>
              <Text className={styles.kvLabel}>门店 ID</Text>
              <Text className={styles.kvValue}>{shellConfig.merchantId}</Text>
            </View>
            <View className={styles.kvItem}>
              <Text className={styles.kvLabel}>来源标识</Text>
              <Text className={styles.kvValue}>{shellConfig.source}</Text>
            </View>
            <View className={styles.kvItem}>
              <Text className={styles.kvLabel}>缓存 token</Text>
              <Text className={styles.kvValue}>{maskToken(cachedSession?.token || '')}</Text>
            </View>
          </View>
        </View>

        <View className={styles.panel}>
          <Text className={styles.panelTitle}>操作</Text>
          <View className={styles.actions}>
            <Button className={styles.button} onClick={() => openWebview(false)}>
              <Text className={styles.buttonText}>立即进入 H5</Text>
            </Button>
            <Button className={styles.ghostButton} onClick={handleClearAndRetry}>
              <Text className={styles.ghostButtonText}>清缓存并重新登录</Text>
            </Button>
            <Button
              className={styles.ghostButton}
              onClick={() => Taro.navigateTo({ url: '/pages/settings/index' })}
            >
              <Text className={styles.ghostButtonText}>打开配置页</Text>
            </Button>
            <Button
              className={styles.ghostButton}
              onClick={() => Taro.navigateTo({ url: '/pages/debug/index' })}
            >
              <Text className={styles.ghostButtonText}>查看调试信息</Text>
            </Button>
          </View>
        </View>
      </View>

      <View className={styles.tipCard}>
        <Text className={styles.tipTitle}>使用提醒</Text>
        <Text className={styles.tipText}>
          微信真机内会调用 `Taro.login -> /api/v1/auth/user/wechat-login`。如果当前是 H5 预览环境，请先去配置页填写调试 token 和调试 openid，再回到这里验证跳转。
        </Text>
      </View>
    </View>
  )
}

export default IndexPage
