<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboard } from '@/api/sp'
import type { DashboardData } from '@/types/sp'
import { formatAmount } from '@/utils/format'

const router = useRouter()
const loading = ref(false)
const dashboard = ref<DashboardData>({
  total_merchants: 0,
  today_orders: 0,
  today_revenue: 0,
  distribution: [],
  trend: []
})

async function loadDashboard() {
  loading.value = true
  try {
    dashboard.value = await getDashboard()
  } finally {
    loading.value = false
  }
}

function goTo(path: string) {
  router.push(path)
}

function getTrendWidth(orders: number) {
  const max = Math.max(...(dashboard.value.trend || []).map((item) => Number(item.orders || 0)), 1)
  return `${Math.max((Number(orders || 0) / max) * 100, 8)}%`
}

onMounted(loadDashboard)
</script>

<template>
  <div class="page-shell">
    <div class="page-header">
      <div class="page-title-wrap">
        <h1 class="page-title">总部工作台</h1>
        <p class="page-subtitle">统一查看门店经营总览，并快速进入门店、订单和数据分析能力。</p>
      </div>
      <el-button type="primary" @click="goTo('/merchants/new')">新增门店</el-button>
    </div>

    <el-skeleton :rows="8" animated :loading="loading">
      <div class="metric-grid">
        <div class="metric-card">
          <div class="metric-label">门店总数</div>
          <div class="metric-value">{{ dashboard.total_merchants }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-label">今日订单</div>
          <div class="metric-value">{{ dashboard.today_orders }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-label">今日成交</div>
          <div class="metric-value">¥{{ formatAmount(dashboard.today_revenue) }}</div>
        </div>
        <div class="metric-card">
          <div class="metric-label">营业门店</div>
          <div class="metric-value">{{ dashboard.total_merchants }}</div>
        </div>
      </div>

      <div class="section-grid">
        <el-card class="page-card" shadow="never">
          <template #header>
            <span>快捷入口</span>
          </template>
          <div class="simple-list">
            <div class="simple-list-item">
              <div>
                <div>门店列表</div>
                <small>统一查看资料与经营概览</small>
              </div>
              <el-button text type="primary" @click="goTo('/merchants')">进入</el-button>
            </div>
            <div class="simple-list-item">
              <div>
                <div>数据分析</div>
                <small>查看订单趋势、转化排行和金额走势</small>
              </div>
              <el-button text type="primary" @click="goTo('/analytics')">进入</el-button>
            </div>
            <div class="simple-list-item">
              <div>
                <div>后台设置</div>
                <small>维护管理员信息并修改登录密码</small>
              </div>
              <el-button text type="primary" @click="goTo('/settings')">进入</el-button>
            </div>
          </div>
        </el-card>

        <el-card class="page-card" shadow="never">
          <template #header>
            <span>门店行业分布</span>
          </template>
          <div v-if="dashboard.distribution?.length" class="simple-list">
            <div v-for="item in dashboard.distribution" :key="item.category || 'unknown'" class="simple-list-item">
              <span>{{ item.category || '未分类' }}</span>
              <strong>{{ item.count }}</strong>
            </div>
          </div>
          <el-empty v-else description="暂无行业分布数据" />
        </el-card>
      </div>

      <el-card class="page-card" shadow="never" style="margin-top: 20px;">
        <template #header>
          <span>近 7 日订单趋势</span>
        </template>
        <div v-if="dashboard.trend?.length" class="simple-list">
          <div v-for="item in dashboard.trend" :key="item.date" class="simple-list-item">
            <div style="min-width: 92px;">{{ item.date }}</div>
            <div class="progress-row" style="flex: 1;">
              <div class="progress-track">
                <div class="progress-bar" :style="{ width: getTrendWidth(item.orders) }"></div>
              </div>
              <strong>{{ item.orders }}</strong>
            </div>
          </div>
        </div>
        <el-empty v-else description="暂无订单趋势数据" />
      </el-card>
    </el-skeleton>
  </div>
</template>
