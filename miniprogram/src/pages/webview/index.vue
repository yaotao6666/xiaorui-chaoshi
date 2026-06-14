<template>
  <web-view v-if="targetUrl" :src="targetUrl"></web-view>
  <view v-else class="empty-state">
    <text>页面地址无效</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const targetUrl = ref('')

onLoad((options) => {
  const target = typeof options?.target === 'string' ? decodeURIComponent(options.target) : ''
  targetUrl.value = target
})
</script>

<style scoped>
.empty-state {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666666;
  font-size: 28rpx;
}
</style>
