<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createMerchantCategory,
  deleteMerchantCategory,
  getMerchantCategories,
  getMerchantDetail,
  updateMerchantCategory,
} from '@/api/sp'
import type { MerchantCategory, MerchantCategoryPayload, MerchantDetail } from '@/types/sp'

const route = useRoute()
const router = useRouter()

const merchantId = computed(() => Number(route.params.id || 0))
const merchant = ref<MerchantDetail | null>(null)
const categories = ref<MerchantCategory[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogSubmitting = ref(false)
const editingCategory = ref<MerchantCategory | null>(null)
const categoryName = ref('')

async function loadData() {
  if (!merchantId.value) return
  loading.value = true
  try {
    const [merchantDetail, categoryList] = await Promise.all([
      getMerchantDetail(merchantId.value),
      getMerchantCategories(merchantId.value),
    ])
    merchant.value = merchantDetail
    categories.value = categoryList
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingCategory.value = null
  categoryName.value = ''
  dialogVisible.value = true
}

function openEdit(category: MerchantCategory) {
  editingCategory.value = category
  categoryName.value = category.name
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
  editingCategory.value = null
  categoryName.value = ''
}

async function submitDialog() {
  if (!merchantId.value || dialogSubmitting.value) return
  const name = categoryName.value.trim()
  if (!name) {
    return ElMessage.warning('请输入分类名称')
  }

  const payload: MerchantCategoryPayload = { name }
  dialogSubmitting.value = true
  try {
    if (editingCategory.value) {
      await updateMerchantCategory(merchantId.value, editingCategory.value.id, payload)
      ElMessage.success('分类更新成功')
    } else {
      await createMerchantCategory(merchantId.value, payload)
      ElMessage.success('分类创建成功')
    }
    closeDialog()
    await loadData()
  } finally {
    dialogSubmitting.value = false
  }
}

async function handleDelete(category: MerchantCategory) {
  if (!merchantId.value) return
  try {
    await ElMessageBox.confirm(`确认删除分类「${category.name}」？`, '删除分类', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  await deleteMerchantCategory(merchantId.value, category.id)
  ElMessage.success('分类删除成功')
  await loadData()
}

onMounted(loadData)
</script>

<template>
  <div class="page-shell">
    <div class="page-header">
      <div class="page-title-wrap">
        <h1 class="page-title">分类管理</h1>
        <p class="page-subtitle">
          {{ merchant ? `当前门店：${merchant.name}` : '维护当前门店商品分类' }}
        </p>
      </div>
      <el-space>
        <el-button @click="router.push(`/merchants/${merchantId}`)">返回门店</el-button>
        <el-button @click="router.push(`/merchants/${merchantId}/products`)">商品管理</el-button>
        <el-button type="primary" @click="openCreate">添加分类</el-button>
      </el-space>
    </div>

    <el-card class="page-card" shadow="never">
      <el-table :data="categories" v-loading="loading">
        <el-table-column prop="name" label="分类名称" min-width="220" />
        <el-table-column label="商品数" width="120">
          <template #default="{ row }">
            {{ row.product_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="120" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" min-width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-space>
              <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
              <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && categories.length === 0" description="暂无分类，请先添加" />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingCategory ? '编辑分类' : '添加分类'" width="420px" @closed="closeDialog">
      <el-input v-model="categoryName" placeholder="请输入分类名称" maxlength="32" />
      <template #footer>
        <el-space>
          <el-button :disabled="dialogSubmitting" @click="closeDialog">取消</el-button>
          <el-button type="primary" :loading="dialogSubmitting" @click="submitDialog">保存</el-button>
        </el-space>
      </template>
    </el-dialog>
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
</style>
