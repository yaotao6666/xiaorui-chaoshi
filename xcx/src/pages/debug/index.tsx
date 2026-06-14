import React, { useMemo, useState } from 'react'
import { Button, Text, View } from '@tarojs/components'
import Taro, { useDidShow } from '@tarojs/taro'
import { clearUserSession } from '@/services/auth'
import { getShellConfig } from '@/services/config'
import { getStoredUserSession } from '@/services/storage'
import { buildEmbeddedH5Url, maskToken } from '@/utils/url'
import styles from './index.module.scss'

const DebugPage: React.FC = () => {
  const [version, setVersion] = useState(0)

  useDidShow(() => {
    setVersion((value) => value + 1)
  })

  const config = useMemo(() => getShellConfig(), [version])
  const session = useMemo(() => getStoredUserSession(), [version])
  const previewUrl = useMemo(() => {
    if (!session?.token || !session.openid) {
      return '尚未生成'
    }

    return buildEmbeddedH5Url(config.h5EntryUrl, {
      merchantId: config.merchantId,
      token: session.token,
      openid: session.openid,
      source: config.source,
      appId: session.appId,
    })
  }, [config, session])

  const handleCopy = async () => {
    if (previewUrl === '尚未生成') {
      Taro.showToast({ title: '当前还没有可复制地址', icon: 'none' })
      return
    }

    await Taro.setClipboardData({ data: previewUrl })
  }

  const handleClear = () => {
    clearUserSession()
    setVersion((value) => value + 1)
    Taro.showToast({ title: '会话已清空', icon: 'success' })
  }

  return (
    <View className={styles.container}>
      <View className={styles.card}>
        <Text className={styles.title}>当前调试快照</Text>
        <Text className={styles.desc}>
          这里展示本地缓存的登录结果和最终拼好的 H5 地址，便于核对 token 是否已经从小程序壳带入 H5。
        </Text>

        <View className={styles.section}>
          <Text className={styles.sectionTitle}>环境配置</Text>
          <View className={styles.item}>
            <Text className={styles.label}>当前环境</Text>
            <Text className={styles.value}>{process.env.TARO_ENV || 'unknown'}</Text>
          </View>
          <View className={styles.item}>
            <Text className={styles.label}>API 基址</Text>
            <Text className={styles.value}>{config.apiBaseUrl}</Text>
          </View>
          <View className={styles.item}>
            <Text className={styles.label}>H5 入口</Text>
            <Text className={styles.value}>{config.h5EntryUrl}</Text>
          </View>
          <View className={styles.item}>
            <Text className={styles.label}>门店 ID</Text>
            <Text className={styles.value}>{config.merchantId}</Text>
          </View>
        </View>

        <View className={styles.section}>
          <Text className={styles.sectionTitle}>登录态</Text>
          <View className={styles.item}>
            <Text className={styles.label}>Token</Text>
            <Text className={styles.value}>{maskToken(session?.token || '')}</Text>
          </View>
          <View className={styles.item}>
            <Text className={styles.label}>OpenID</Text>
            <Text className={styles.value}>{session?.openid || '未获取'}</Text>
          </View>
          <View className={styles.item}>
            <Text className={styles.label}>AppID</Text>
            <Text className={styles.value}>{session?.appId || '未获取'}</Text>
          </View>
        </View>

        <View className={styles.section}>
          <Text className={styles.sectionTitle}>最终 H5 地址</Text>
          <View className={styles.item}>
            <Text className={styles.value}>{previewUrl}</Text>
          </View>
        </View>

        <View className={styles.actions}>
          <Button className={styles.ghostButton} onClick={handleCopy}>
            <Text className={styles.ghostText}>复制最终地址</Text>
          </Button>
          <Button className={styles.ghostButton} onClick={handleClear}>
            <Text className={styles.ghostText}>清空会话缓存</Text>
          </Button>
          <Button
            className={styles.ghostButton}
            onClick={() => Taro.navigateTo({ url: '/pages/settings/index' })}
          >
            <Text className={styles.ghostText}>返回配置页</Text>
          </Button>
        </View>
      </View>
    </View>
  )
}

export default DebugPage
