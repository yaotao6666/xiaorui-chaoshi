<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { useAuthStore } from '@/stores/auth'
import { cacheEntrySourceFromUrl, shouldUseEmbeddedShellTitle } from '@/utils/embeddedShell'

function applyEmbeddedShellClass() {
  // #ifdef H5
  if (typeof document === 'undefined') {
    return
  }

  cacheEntrySourceFromUrl()
  const className = 'xcx-shell-h5'
  if (shouldUseEmbeddedShellTitle()) {
    document.documentElement.classList.add(className)
    document.body?.classList.add(className)
    return
  }

  document.documentElement.classList.remove(className)
  document.body?.classList.remove(className)
  // #endif
}

onLaunch(() => {
  console.log('App Launch')
  const authStore = useAuthStore()

  // 统一收口商家登录失效事件，避免 401 跳转时残留 WebSocket 连接。
  uni.$off('merchant-session-expired')
  uni.$on('merchant-session-expired', () => {
    authStore.handleSessionExpired()
  })

  // 检查登录状态
  authStore.checkLogin()
  applyEmbeddedShellClass()
})

onShow(() => {
  console.log('App Show')
  applyEmbeddedShellClass()
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style lang="scss">
@import '@/uni.scss';

/* 全局样式 */

/* CSS变量 */
page {
  --primary-color: #007aff;
  --success-color: #07c160;
  --warning-color: #ff9500;
  --danger-color: #ff4d4f;
  --text-color: #333333;
  --text-color-secondary: #666666;
  --text-color-placeholder: #999999;
  --border-color: #eeeeee;
  --bg-color: #f5f5f5;
  --white: #ffffff;
}

.xcx-shell-h5 uni-page-head,
.xcx-shell-h5 .uni-page-head,
.xcx-shell-h5 uni-page-head-fixed,
.xcx-shell-h5 .uni-page-head-fixed {
  display: none !important;
}

.xcx-shell-h5 uni-page-wrapper,
.xcx-shell-h5 .uni-page-wrapper,
.xcx-shell-h5 uni-page-body,
.xcx-shell-h5 .uni-page-body {
  margin-top: 0 !important;
  padding-top: 0 !important;
  top: 0 !important;
}
</style>
