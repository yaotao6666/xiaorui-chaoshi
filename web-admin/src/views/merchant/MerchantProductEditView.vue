<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  createMerchantProduct,
  getMerchantCategories,
  getMerchantDetail,
  getMerchantProduct,
  updateMerchantProduct,
} from '@/api/sp'
import type {
  MerchantCategory,
  MerchantDetail,
  MerchantProduct,
  MerchantProductPayload,
  MerchantProductSpec,
} from '@/types/sp'
import { uploadSpImage } from '@/utils/qiniu'

type ProductSpecFormOption = {
  name: string
  price?: number
}

type ProductSpecForm = {
  name: string
  options: ProductSpecFormOption[]
}

const route = useRoute()
const router = useRouter()

const merchantId = computed(() => Number(route.params.id || 0))
const productId = computed(() => Number(route.params.productId || 0))
const isEditMode = computed(() => productId.value > 0)

const merchant = ref<MerchantDetail | null>(null)
const categories = ref<MerchantCategory[]>([])
const loading = ref(false)
const submitting = ref(false)
const uploading = ref(false)
const imageInputRef = ref<HTMLInputElement | null>(null)

const form = reactive({
  category_id: undefined as number | undefined,
  name: '',
  description: '',
  price: undefined as number | undefined,
  original_price: undefined as number | undefined,
  stock: 0,
  unit: '份',
  images: [] as string[],
  specs: [] as ProductSpecForm[],
  sort: 0,
})

function normalizeSpecs(specs: MerchantProductSpec[] = []): ProductSpecForm[] {
  return specs.map((spec) => ({
    name: spec.name || '',
    options: Array.isArray(spec.options)
      ? spec.options.map((option) => ({
          name: option.name || '',
          price: Number(option.price || 0) || undefined,
        }))
      : [],
  }))
}

function fillForm(product: MerchantProduct) {
  form.category_id = product.category_id || undefined
  form.name = product.name || ''
  form.description = product.description || ''
  form.price = Number(product.price || 0) || undefined
  form.original_price = Number(product.original_price || 0) || undefined
  form.stock = Number(product.stock || 0)
  form.unit = product.unit || '份'
  form.images = Array.isArray(product.images) ? [...product.images] : []
  form.specs = normalizeSpecs(product.specs)
  form.sort = Number(product.sort || 0)
}

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

    if (isEditMode.value) {
      const product = await getMerchantProduct(merchantId.value, productId.value)
      fillForm(product)
    }
  } finally {
    loading.value = false
  }
}

function openImagePicker() {
  if (uploading.value) return
  imageInputRef.value?.click()
}

async function handleImageChange(event: Event) {
  const target = event.target as HTMLInputElement
  const files = Array.from(target.files || [])
  target.value = ''
  if (files.length === 0) return

  uploading.value = true
  try {
    for (const file of files) {
      if (form.images.length >= 9) break
      const uploadResult = await uploadSpImage(file)
      form.images.push(uploadResult.url)
    }
    ElMessage.success('图片上传成功')
  } finally {
    uploading.value = false
  }
}

function removeImage(index: number) {
  form.images.splice(index, 1)
}

function addSpec() {
  form.specs.push({
    name: '',
    options: [{ name: '', price: undefined }],
  })
}

function removeSpec(index: number) {
  form.specs.splice(index, 1)
}

function addOption(specIndex: number) {
  form.specs[specIndex].options.push({ name: '', price: undefined })
}

function removeOption(specIndex: number, optionIndex: number) {
  form.specs[specIndex].options.splice(optionIndex, 1)
}

function validateForm() {
  if (!form.name.trim()) {
    ElMessage.warning('请输入商品名称')
    return false
  }
  if (!form.category_id) {
    ElMessage.warning('请选择商品分类')
    return false
  }
  if (!form.price || form.price <= 0) {
    ElMessage.warning('请输入正确的售价')
    return false
  }
  if (form.stock < 0) {
    ElMessage.warning('库存不能小于 0')
    return false
  }

  for (const spec of form.specs) {
    if (!spec.name.trim()) {
      ElMessage.warning('规格名称不能为空')
      return false
    }
    if (spec.options.length === 0) {
      ElMessage.warning('请至少添加一个规格选项')
      return false
    }
    for (const option of spec.options) {
      if (!option.name.trim()) {
        ElMessage.warning('规格选项名称不能为空')
        return false
      }
      if (option.price !== undefined && option.price < 0) {
        ElMessage.warning('规格加价不能小于 0')
        return false
      }
    }
  }

  return true
}

function buildPayload(): MerchantProductPayload {
  const specs: MerchantProductSpec[] = form.specs.map((spec) => ({
    name: spec.name.trim(),
    options: spec.options.map((option) => ({
      name: option.name.trim(),
      price: Number(option.price || 0),
    })),
  }))

  return {
    category_id: form.category_id,
    name: form.name.trim(),
    description: form.description.trim(),
    images: form.images,
    price: Number(form.price || 0),
    original_price: Number(form.original_price || 0),
    stock: Number(form.stock || 0),
    unit: form.unit.trim() || '份',
    sort: Number(form.sort || 0),
    specs,
  }
}

async function submitForm(successMessage: string) {
  if (!merchantId.value || submitting.value) return
  if (!validateForm()) return

  submitting.value = true
  try {
    const payload = buildPayload()
    if (isEditMode.value) {
      await updateMerchantProduct(merchantId.value, productId.value, payload)
    } else {
      await createMerchantProduct(merchantId.value, payload)
    }
    ElMessage.success(successMessage)
    await nextTick()
    await router.push(`/merchants/${merchantId.value}/products`)
  } finally {
    submitting.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <div class="page-shell" v-loading="loading">
    <div class="page-header">
      <div class="page-title-wrap">
        <h1 class="page-title">{{ isEditMode ? '编辑商品' : '新增商品' }}</h1>
        <p class="page-subtitle">
          {{ merchant ? `当前门店：${merchant.name}` : '维护当前门店商品资料' }}
        </p>
      </div>
      <el-space>
        <el-button @click="router.push(`/merchants/${merchantId}/products`)">返回商品列表</el-button>
        <el-button @click="router.push(`/merchants/${merchantId}/categories`)">分类管理</el-button>
      </el-space>
    </div>

    <el-card class="page-card" shadow="never">
      <template #header>
        <span>基本信息</span>
      </template>
      <el-form label-position="top">
        <el-form-item label="商品名称">
          <el-input v-model="form.name" placeholder="请输入商品名称" maxlength="64" />
        </el-form-item>
        <el-form-item label="商品分类">
          <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 100%;">
            <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品描述">
          <el-input v-model="form.description" type="textarea" :rows="4" maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="page-card" shadow="never">
      <template #header>
        <span>价格库存</span>
      </template>
      <div class="form-grid form-grid-2">
        <el-form-item label="售价">
          <el-input-number v-model="form.price" :min="0" :precision="2" :step="0.1" style="width: 100%;" />
        </el-form-item>
        <el-form-item label="原价">
          <el-input-number v-model="form.original_price" :min="0" :precision="2" :step="0.1" style="width: 100%;" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="form.stock" :min="0" :step="1" style="width: 100%;" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="form.unit" placeholder="份/个/盒" />
        </el-form-item>
      </div>
    </el-card>

    <el-card class="page-card" shadow="never">
      <template #header>
        <div class="card-header-between">
          <span>商品图片</span>
          <el-button type="primary" plain :loading="uploading" @click="openImagePicker">
            {{ uploading ? '上传中' : '上传图片' }}
          </el-button>
        </div>
      </template>
      <input
        ref="imageInputRef"
        hidden
        type="file"
        multiple
        accept="image/*"
        @change="handleImageChange"
      />
      <div class="image-grid">
        <div v-for="(image, index) in form.images" :key="`${image}-${index}`" class="image-item">
          <img :src="image" alt="商品图片" />
          <button class="image-delete" type="button" @click="removeImage(index)">×</button>
        </div>
        <div v-if="form.images.length === 0" class="image-empty">暂未上传图片</div>
      </div>
    </el-card>

    <el-card class="page-card" shadow="never">
      <template #header>
        <div class="card-header-between">
          <span>商品规格</span>
          <el-button type="primary" plain @click="addSpec">添加规格</el-button>
        </div>
      </template>
      <div v-if="form.specs.length === 0" class="spec-empty">暂未配置规格，可按需新增</div>
      <div v-for="(spec, specIndex) in form.specs" :key="specIndex" class="spec-card">
        <div class="spec-head">
          <el-input v-model="spec.name" placeholder="规格名称，例如：份量" />
          <el-button link type="danger" @click="removeSpec(specIndex)">删除规格</el-button>
        </div>
        <div v-for="(option, optionIndex) in spec.options" :key="optionIndex" class="spec-option-row">
          <el-input v-model="option.name" placeholder="选项名称" />
          <el-input-number
            v-model="option.price"
            :min="0"
            :precision="2"
            :step="0.1"
            placeholder="加价"
            style="width: 180px;"
          />
          <el-button link type="danger" @click="removeOption(specIndex, optionIndex)">删除</el-button>
        </div>
        <el-button type="primary" link @click="addOption(specIndex)">+ 添加选项</el-button>
      </div>
    </el-card>

    <el-card class="page-card" shadow="never">
      <template #header>
        <span>排序</span>
      </template>
      <el-form-item label="排序值">
        <el-input-number v-model="form.sort" :min="0" :step="1" style="width: 100%;" />
      </el-form-item>
    </el-card>

    <div class="bottom-actions">
      <el-button size="large" @click="submitForm('商品保存成功')">保存草稿</el-button>
      <el-button type="primary" size="large" :loading="submitting" @click="submitForm('商品发布成功')">
        {{ submitting ? '提交中...' : '发布商品' }}
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.page-shell {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-bottom: 88px;
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

.form-grid {
  display: grid;
  gap: 18px;
}

.form-grid-2 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.card-header-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(132px, 1fr));
  gap: 16px;
}

.image-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: 14px;
  overflow: hidden;
  background: #f3f4f6;
}

.image-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.image-delete {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 50%;
  background: rgba(17, 24, 39, 0.72);
  color: #fff;
  cursor: pointer;
}

.image-empty,
.spec-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 120px;
  border-radius: 14px;
  background: #f9fafb;
  color: #9ca3af;
}

.spec-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-radius: 16px;
  background: #f9fafb;
}

.spec-head,
.spec-option-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.spec-head > :deep(.el-input),
.spec-option-row > :deep(.el-input) {
  flex: 1;
}

.bottom-actions {
  position: fixed;
  right: 24px;
  bottom: 24px;
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.12);
}
</style>
