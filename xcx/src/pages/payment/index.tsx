import React, { useCallback, useMemo, useState } from 'react'
import { Button, Text, View } from '@tarojs/components'
import Taro, { useLoad } from '@tarojs/taro'
import { ensureUserSession } from '@/services/auth'
import { getShellConfig } from '@/services/config'
import { postJson } from '@/services/request'
import type { PrepareOrderPaymentPayload } from '@/types/payment'
import { buildEmbeddedH5Url } from '@/utils/url'
import styles from './index.module.scss'

type PaymentStatus = 'idle' | 'loading' | 'success' | 'cancelled' | 'error'

function isCancelError(error: unknown) {
  if (!(error instanceof Error)) {
    return false
  }
  return /cancel/i.test(error.message)
}

const PaymentPage: React.FC = () => {
  const shellConfig = useMemo(() => getShellConfig(), [])
  const [status, setStatus] = useState<PaymentStatus>('idle')
  const [errorMessage, setErrorMessage] = useState('')
  const [orderId, setOrderId] = useState(0)
  const [merchantId, setMerchantId] = useState(0)
  const [returnTarget, setReturnTarget] = useState('/pages/store/order-detail')
  const [orderNo, setOrderNo] = useState('')
  const [payAmount, setPayAmount] = useState(0)

  const openH5Target = useCallback(
    async (targetPath: string) => {
      const session = await ensureUserSession(false)
      const h5EntryUrl = shellConfig.h5EntryUrl.includes('#')
        ? `${shellConfig.h5EntryUrl.split('#')[0]}#${targetPath}`
        : `${shellConfig.h5EntryUrl}#${targetPath}`

      const targetUrl = buildEmbeddedH5Url(h5EntryUrl, {
        merchantId: String(merchantId || shellConfig.merchantId || '1'),
        token: session.token,
        openid: session.openid,
        source: shellConfig.source,
        appId: session.appId,
      })

      Taro.redirectTo({
        url: `/pages/webview/index?target=${encodeURIComponent(targetUrl)}`,
      })
    },
    [merchantId, shellConfig]
  )

  const handlePay = useCallback(async (nextOrderId = orderId, nextReturnTarget = returnTarget) => {
    if (!nextOrderId) {
      setStatus('error')
      setErrorMessage('缺少订单 ID')
      return
    }

    try {
      setStatus('loading')
      setErrorMessage('')
      await ensureUserSession(false)

      const apiBaseUrl = shellConfig.apiBaseUrl
      const payload = await postJson<PrepareOrderPaymentPayload>(
        `${apiBaseUrl}/api/v1/user/orders/${nextOrderId}/pay/prepare`,
        {
          return_path: nextReturnTarget,
          source: shellConfig.source,
        },
        { withAuth: true }
      )

      setOrderNo(payload.order_no)
      setPayAmount(payload.pay_amount)

      await Taro.requestPayment({
        timeStamp: payload.pay_params.timeStamp,
        nonceStr: payload.pay_params.nonceStr,
        package: payload.pay_params.package,
        signType: payload.pay_params.signType as 'RSA' | 'MD5' | 'HMAC-SHA256',
        paySign: payload.pay_params.paySign,
      })

      setStatus('success')
      await openH5Target(payload.return_target)
    } catch (error: unknown) {
      if (isCancelError(error)) {
        setStatus('cancelled')
        setErrorMessage('你已取消支付，可稍后继续发起支付')
        return
      }

      console.error('[XCX][Payment] 支付失败', error)
      setStatus('error')
      setErrorMessage(error instanceof Error ? error.message : '支付失败')
    }
  }, [openH5Target, orderId, returnTarget, shellConfig])

  useLoad((options) => {
    const nextOrderId = Number(options?.orderId || 0)
    const nextMerchantId = Number(options?.merchantId || 0)
    const nextReturnTarget = typeof options?.returnTarget === 'string'
      ? decodeURIComponent(options.returnTarget)
      : '/pages/store/order-detail'

    setOrderId(nextOrderId)
    setMerchantId(nextMerchantId)
    setReturnTarget(nextReturnTarget)

    if (nextOrderId > 0) {
      void handlePay(nextOrderId, nextReturnTarget)
    } else {
      setStatus('error')
      setErrorMessage('缺少订单 ID')
    }
  })

  return (
    <View className={styles.container}>
      <View className={styles.card}>
        <Text className={styles.eyebrow}>XCX 原生支付页</Text>
        <Text className={styles.title}>正在为当前订单准备小程序支付</Text>
        <Text className={styles.desc}>
          当前页面用于承接 H5 提交订单后的支付动作。支付成功后会自动回到 H5 订单详情页；取消或失败时可在此页重新发起。
        </Text>

        <View
          className={[
            styles.status,
            status === 'error' || status === 'cancelled' ? styles.statusError : '',
            status === 'success' ? styles.statusSuccess : '',
          ].join(' ')}
        >
          <Text>
            {{
              idle: '等待发起支付',
              loading: '正在请求支付参数并拉起微信支付',
              success: '支付成功，准备返回 H5',
              cancelled: '支付已取消，可重新发起',
              error: '支付失败，请查看错误信息后重试',
            }[status]}
          </Text>
        </View>

        {!!errorMessage && <Text className={styles.desc}>{errorMessage}</Text>}

        <View className={styles.kvList}>
          <View className={styles.kvItem}>
            <Text className={styles.kvLabel}>订单 ID</Text>
            <Text className={styles.kvValue}>{orderId || '-'}</Text>
          </View>
          <View className={styles.kvItem}>
            <Text className={styles.kvLabel}>订单号</Text>
            <Text className={styles.kvValue}>{orderNo || '-'}</Text>
          </View>
          <View className={styles.kvItem}>
            <Text className={styles.kvLabel}>商家 ID</Text>
            <Text className={styles.kvValue}>{merchantId || '-'}</Text>
          </View>
          <View className={styles.kvItem}>
            <Text className={styles.kvLabel}>支付金额</Text>
            <Text className={styles.kvValue}>{payAmount ? `¥${payAmount.toFixed(2)}` : '-'}</Text>
          </View>
          <View className={styles.kvItem}>
            <Text className={styles.kvLabel}>回跳页面</Text>
            <Text className={styles.kvValue}>{returnTarget}</Text>
          </View>
        </View>

        <View className={styles.actions}>
          <Button className={styles.button} onClick={() => void handlePay()}>
            重新发起支付
          </Button>
          <Button className={styles.ghostButton} onClick={() => void openH5Target(returnTarget)}>
            返回 H5 订单详情
          </Button>
        </View>
      </View>
    </View>
  )
}

export default PaymentPage
