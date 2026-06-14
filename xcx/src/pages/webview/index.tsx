import React, { useState } from 'react'
import { Text, View, WebView } from '@tarojs/components'
import { useLoad } from '@tarojs/taro'
import styles from './index.module.scss'

const WebviewPage: React.FC = () => {
  const [targetUrl, setTargetUrl] = useState('')

  useLoad((options) => {
    const target = typeof options?.target === 'string' ? decodeURIComponent(options.target) : ''
    setTargetUrl(target)
  })

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

  return <WebView src={targetUrl} />
}

export default WebviewPage
