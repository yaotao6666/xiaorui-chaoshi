import React, { useCallback, useState } from 'react'
import { Text, View, WebView } from '@tarojs/components'
import Taro, { useLoad } from '@tarojs/taro'
import { getDefaultWebviewTitle, getWebviewTitleFromTargetUrl } from '@/utils/webviewTitle'
import styles from './index.module.scss'

const WebviewPage: React.FC = () => {
  const [targetUrl, setTargetUrl] = useState('')

  const updateNavigationTitle = useCallback((title: string) => {
    const resolvedTitle = title.trim() || getDefaultWebviewTitle()
    Taro.setNavigationBarTitle({
      title: resolvedTitle,
    })
  }, [])

  useLoad((options) => {
    const target = typeof options?.target === 'string' ? decodeURIComponent(options.target) : ''
    setTargetUrl(target)
    updateNavigationTitle(getWebviewTitleFromTargetUrl(target))
  })

  const handleMessage = useCallback((event: Record<string, any>) => {
    const payloadList = Array.isArray(event?.detail?.data) ? event.detail.data : []
    const payload = payloadList.find((item) => item && typeof item === 'object' && item.type === 'page_title_sync')
    if (!payload || typeof payload.title !== 'string') {
      return
    }
    updateNavigationTitle(payload.title)
  }, [updateNavigationTitle])

  if (!targetUrl) {
    return (
      <View className={styles.emptyState}>
        <View className={styles.card}>
          <Text className={styles.title}>缺少 H5 地址</Text>
          <Text className={styles.desc}>请先从启动页完成登录，再进入当前容器页。</Text>
        </View>
      </View>
    )
  }

  return <WebView src={targetUrl} onMessage={handleMessage} />
}

export default WebviewPage
