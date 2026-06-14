import React, { useState } from 'react'
import { Button, Input, Switch, Text, View } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { DEFAULT_SHELL_CONFIG, getShellConfig, saveShellConfig } from '@/services/config'
import type { MiniShellConfig } from '@/types/config'
import styles from './index.module.scss'

const SettingsPage: React.FC = () => {
  const [form, setForm] = useState<MiniShellConfig>(getShellConfig())

  const updateField = <K extends keyof MiniShellConfig>(key: K, value: MiniShellConfig[K]) => {
    setForm((current) => ({
      ...current,
      [key]: value,
    }))
  }

  const handleSave = () => {
    saveShellConfig(form)
    Taro.showToast({
      title: '配置已保存',
      icon: 'success',
    })
  }

  const handleReset = () => {
    setForm(DEFAULT_SHELL_CONFIG)
    saveShellConfig(DEFAULT_SHELL_CONFIG)
    Taro.showToast({
      title: '已恢复默认配置',
      icon: 'success',
    })
  }

  return (
    <View className={styles.container}>
      <View className={styles.card}>
        <Text className={styles.title}>XCX 嵌入参数</Text>
        <Text className={styles.desc}>
          这里保存“小程序登录接口地址、H5 入口地址、门店 ID、预览调试 token”等信息。保存后，启动页会自动读取这些配置。
        </Text>

        <View className={styles.field}>
          <Text className={styles.label}>API 基址</Text>
          <Input
            className={styles.input}
            type='text'
            value={form.apiBaseUrl}
            onInput={(event) => updateField('apiBaseUrl', event.detail.value)}
          />
        </View>

        <View className={styles.field}>
          <Text className={styles.label}>H5 入口地址</Text>
          <Input
            className={styles.input}
            type='text'
            value={form.h5EntryUrl}
            onInput={(event) => updateField('h5EntryUrl', event.detail.value)}
          />
        </View>

        <View className={styles.field}>
          <Text className={styles.label}>门店 ID</Text>
          <Input
            className={styles.input}
            type='number'
            value={form.merchantId}
            onInput={(event) => updateField('merchantId', event.detail.value)}
          />
        </View>

        <View className={styles.field}>
          <Text className={styles.label}>来源标识</Text>
          <Input
            className={styles.input}
            type='text'
            value={form.source}
            onInput={(event) => updateField('source', event.detail.value)}
          />
        </View>

        <View className={styles.field}>
          <Text className={styles.label}>H5 预览调试 token</Text>
          <Input
            className={styles.input}
            type='text'
            value={form.debugToken}
            onInput={(event) => updateField('debugToken', event.detail.value)}
          />
        </View>

        <View className={styles.field}>
          <Text className={styles.label}>H5 预览调试 openid</Text>
          <Input
            className={styles.input}
            type='text'
            value={form.debugOpenid}
            onInput={(event) => updateField('debugOpenid', event.detail.value)}
          />
        </View>

        <View className={styles.switchRow}>
          <Text className={styles.switchText}>进入启动页后自动尝试登录并跳转</Text>
          <Switch
            checked={form.autoEnter}
            color='#14b8a6'
            onChange={(event) => updateField('autoEnter', event.detail.value)}
          />
        </View>

        <View className={styles.actions}>
          <Button className={styles.primaryButton} onClick={handleSave}>
            <Text className={styles.primaryText}>保存配置</Text>
          </Button>
          <Button className={styles.ghostButton} onClick={handleReset}>
            <Text className={styles.ghostText}>恢复默认值</Text>
          </Button>
          <Button
            className={styles.ghostButton}
            onClick={() => Taro.navigateTo({ url: '/pages/debug/index' })}
          >
            <Text className={styles.ghostText}>查看调试页</Text>
          </Button>
        </View>
      </View>

      <View className={styles.tipBox}>
        <Text className={styles.tipTitle}>测试说明</Text>
        <Text className={styles.tipText}>
          真机小程序 `web-view` 不能直接访问 `localhost`。如果要在微信里验证真实嵌入，请把 H5 入口替换成可公网访问且已配置业务域名的 HTTPS 地址。
        </Text>
      </View>
    </View>
  )
}

export default SettingsPage
