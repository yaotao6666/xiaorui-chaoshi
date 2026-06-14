<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  batchUpdateMerchantProductStatus,
  deleteMerchantProduct,
  getMerchantCategories,
  getMerchantDetail,
  getMerchantProducts,
  merchantProductOffSale,
  merchantProductOnSale,
} from '@/api/sp'
import type { MerchantCategory, MerchantDetail, MerchantProduct } from '@/types/sp'
import { formatAmount, formatDateTime } from '@/utils/format'

const route = useRoute()
const router = useRouter()

const merchantId = computed(() => Number(route.params.id || 0))
const merchant = ref<MerchantDetail | null>(null)
const categories = ref<MerchantCategory[]>([])
const products = ref<MerchantProduct[]>([])
const loading = ref(false)
const selectedIds = ref<number[]>([])

const filters = reactive({
  category_id: undefined as number | undefined,
  status: undefined as number | undefined,
  keyword: '',
  page: 1,
  page_size: 10,
})

const pagination = reactive({
  total: 0,
  page: 1,
  page_size: 10,
})

async function loadBaseData() {
  if (!merchantId.value) return
  const [merchantDetail, categoryList] = await Promise.all([
    getMerchantDetail(merchantId.value),
    getMerchantCategories(merchantId.value),
  ])
  merchant.value = merchantDetail
  categories.value = categoryList
}

async function loadProducts() {
  if (!merchantId.value) return
  loading.value = true
  try {
    const data = await getMerchantProducts(merchantId.value, {
      page: filters.page,
      page_size: filters.page_size,
      category_id: filters.category_id,
      status: filters.status,
      keyword: filters.keyword.trim() || undefined,
    })
    products.value = data.list
    pagination.total = data.pagination.total
    pagination.page = data.pagination.page
    pagination.page_size = data.pagination.page_size
  } finally {
    loading.value = false
  }
}

async function initializePage() {
  loading.value = true
  try {
    await loadBaseData()
    await loadProducts()
  } finally {
    loading.value = false
  }
}

function handleSelectionChange(rows: MerchantProduct[]) {
  selectedIds.value = rows.map((item) => item.id)
}

async function handleSearch() {
  filters.page = 1
  await loadProducts()
}

async function handleReset() {
  filters.category_id = undefined
  filters.status = undefined
  filters.keyword = ''
  filters.page = 1
  await loadProducts()
}

async function handlePageChange(page: number) {
  filters.page = page
  await loadProducts()
}

async function handlePageSizeChange(pageSize: number) {
  filters.page_size = pageSize
  filters.page = 1
  await loadProducts()
}

async function toggleStatus(product: MerchantProduct) {
  if (!merchantId.value) return
  if (product.status === 1) {
    await merchantProductOffSale(merchantId.value, product.id)
    ElMessage.success('商品已下架')
  } else {
    await merchantProductOnSale(merchantId.value, product.id)
    ElMessage.success('商品已上架')
  }
  await loadProducts()
}

async function handleDelete(product: MerchantProduct) {
  if (!merchantId.value) return
  try {
    await ElMessageBox.confirm(`确认删除商品「${product.name}」？`, '删除商品', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  await deleteMerchantProduct(merchantId.value, product.id)
  ElMessage.success('商品删除成功')
  await loadProducts()
}

async function batchUpdateStatus(status: number) {
  if (!merchantId.value || selectedIds.value.length === 0) {
    return ElMessage.warning('请先选择商品')
  }
  await batchUpdateMerchantProductStatus(merchantId.value, selectedIds.value, status)
  ElMessage.success(status === 1 ? '批量上架成功' : '批量下架成功')
  selectedIds.value = []
  await loadProducts()
}

onMounted(initializePage)
</script>

<template>
  <div class="page-shell">
    <div class="page-header">
      <div class="page-title-wrap">
        <h1 class="page-title">商品管理</h1>
        <p class="page-subtitle">
          {{ merchant ? `当前门店：${merchant.name}` : '维护当前门店商品信息' }}
        </p>
      </div>
      <el-space>
        <el-button @click="router.push(`/merchants/${merchantId}`)">返回门店</el-button>
        <el-button @click="router.push(`/merchants/${merchantId}/categories`)">分类管理</el-button>
        <el-button type="primary" @click="router.push(`/merchants/${merchantId}/products/new`)">添加商品</el-button>
      </el-space>
    </div>

    <el-card class="page-card" shadow="never">
      <div class="filter-bar">
        <el-select v-model="filters.category_id" clearable placeholder="全部分类" style="width: 180px;">
          <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
        <el-select v-model="filters.status" clearable placeholder="全部状态" style="width: 160px;">
          <el-option label="上架" :value="1" />
          <el-option label="下架" :value="2" />
        </el-select>
        <el-input
          v-model="filters.keyword"
          placeholder="搜索商品名称"
          clearable
          style="width: 240px;"
          @keyup.enter="handleSearch"
        />
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
        <div class="filter-spacer" />
        <el-button :disabled="selectedIds.length === 0" @click="batchUpdateStatus(1)">批量上架</el-button>
        <el-button :disabled="selectedIds.length === 0" @click="batchUpdateStatus(2)">批量下架</el-button>
      </div>

      <el-table :data="products" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column label="图片" width="92">
          <template #default="{ row }">
            <img v-if="row.images?.[0]" :src="row.images[0]" class="product-image" alt="商品图" />
            <div v-else class="empty-image">无图</div>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" min-width="200" />
        <el-table-column prop="category_name" label="分类" min-width="140" />
        <el-table-column label="价格" width="140">
          <template #default="{ row }">
            <div class="price-block">
              <span class="price-now">¥{{ formatAmount(row.price) }}</span>
              <span v-if="Number(row.original_price || 0) > 0" class="price-origin">
                ¥{{ formatAmount(row.original_price) }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="100" />
        <el-table-column prop="sales" label="销量" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column label="更新时间" min-width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-space>
              <el-button link type="primary" @click="router.push(`/merchants/${merchantId}/products/${row.id}/edit`)">
                编辑
              </el-button>
              <el-button link type="primary" @click="toggleStatus(row)">
                {{ row.status === 1 ? '下架' : '上架' }}
              </el-button>
              <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && products.length === 0" description="暂无商品，请先添加" />

      <div class="pagination-wrap">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next"
          :total="pagination.total"
          :current-page="filters.page"
          :page-size="filters.page_size"
          :page-sizes="[10, 20, 50]"
          @current-change="handlePageChange"
          @size-change="handlePageSizeChange"
        />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.page-shell {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.page-title-wrap {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.page-title {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  color: #111827;
}

.page-subtitle {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.page-card {
  border-radius: 20px;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.filter-spacer {
  flex: 1;
}

.product-image,
.empty-image {
  width: 52px;
  height: 52px;
  border-radius: 10px;
}

.product-image {
  object-fit: cover;
  display: block;
}

.empty-image {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  color: #9ca3af;
  font-size: 12px;
}

.price-block {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.price-now {
  color: #ef4444;
  font-weight: 600;
}

.price-origin {
  color: #9ca3af;
  text-decoration: line-through;
  font-size: 12px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
