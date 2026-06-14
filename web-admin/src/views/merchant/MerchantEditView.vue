<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  createSpMerchant,
  getMerchantDetail,
  resetSpMerchantAdminPassword,
  updateSpMerchant,
} from '@/api/sp'
import type { MerchantDetail, SpMerchantFormData, UpdateSpMerchantFormData } from '@/types/sp'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const savingBasic = ref(false)
const savingAdminPassword = ref(false)

const merchantId = computed(() => Number(route.params.id || 0))
const isEditMode = computed(() => merchantId.value > 0)

const basicForm = reactive<SpMerchantFormData>({
  name: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  address: '',
  business_category: '',
  business_hours: '',
  announcement: '',
  username: '',
  password: '',
  staff_name: '',
  staff_phone: '',
})

const adminForm = reactive({
  new_password: '',
})

function fillForm(detail: MerchantDetail) {
  basicForm.name = detail.name || ''
  basicForm.contact_name = detail.contact_name || ''
  basicForm.contact_phone = detail.contact_phone || ''
  basicForm.contact_email = detail.contact_email || ''
  basicForm.address = detail.address || ''
  basicForm.business_category = detail.business_category || ''
  basicForm.business_hours = detail.business_hours || ''
  basicForm.announcement = detail.announcement || ''
  basicForm.username = detail.admin_username || ''
  basicForm.staff_name = detail.admin_name || ''
  basicForm.staff_phone = detail.admin_phone || ''
  basicForm.password = ''
}

async function loadDetail() {
  if (!isEditMode.value) return
  loading.value = true
  try {
    const detail = await getMerchantDetail(merchantId.value)
    fillForm(detail)
  } finally {
    loading.value = false
  }
}

function validateBasic() {
  if (!basicForm.name.trim()) {
    throw new Error('请输入商家名称')
  }
  if (!isEditMode.value) {
    if (!basicForm.username.trim()) {
      throw new Error('请输入管理员账号')
    }
    if (basicForm.password.trim().length < 6) {
      throw new Error('管理员密码至少 6 位')
    }
  }
}

function validateAdminPassword() {
  if (adminForm.new_password.trim().length < 6) {
    throw new Error('管理员新密码至少 6 位')
  }
}

async function submitBasic() {
  try {
    validateBasic()
  } catch (error: any) {
    ElMessage.warning(error.message)
    return
  }

  savingBasic.value = true
  try {
    if (isEditMode.value) {
      const payload: UpdateSpMerchantFormData = {
        name: basicForm.name.trim(),
        contact_name: basicForm.contact_name?.trim(),
        contact_phone: basicForm.contact_phone?.trim(),
        contact_email: basicForm.contact_email?.trim(),
        address: basicForm.address?.trim(),
        business_category: basicForm.business_category?.trim(),
        business_hours: basicForm.business_hours?.trim(),
        announcement: basicForm.announcement?.trim(),
      }
      await updateSpMerchant(merchantId.value, payload)
      ElMessage.success('基础信息已保存')
      return
    }

    const created = await createSpMerchant({
      ...basicForm,
      name: basicForm.name.trim(),
      contact_name: basicForm.contact_name?.trim(),
      contact_phone: basicForm.contact_phone?.trim(),
      contact_email: basicForm.contact_email?.trim(),
      address: basicForm.address?.trim(),
      business_category: basicForm.business_category?.trim(),
      business_hours: basicForm.business_hours?.trim(),
      announcement: basicForm.announcement?.trim(),
      username: basicForm.username.trim(),
      password: basicForm.password.trim(),
      staff_name: basicForm.staff_name?.trim(),
      staff_phone: basicForm.staff_phone?.trim(),
    })
    ElMessage.success('商家创建成功')
    await router.replace(`/merchants/${created.id}`)
  } finally {
    savingBasic.value = false
  }
}

async function submitAdminPassword() {
  if (!isEditMode.value) return

  try {
    validateAdminPassword()
  } catch (error: any) {
    ElMessage.warning(error.message)
    return
  }

  savingAdminPassword.value = true
  try {
    const result = await resetSpMerchantAdminPassword(merchantId.value, {
      new_password: adminForm.new_password.trim(),
    })
    adminForm.new_password = ''
    ElMessage.success(`${result.username} 的登录密码已重置`)
  } finally {
    savingAdminPassword.value = false
  }
}

onMounted(loadDetail)
</script>

<template>
  <div class="page-shell">
    <div class="page-header">
      <div class="page-title-wrap">
        <h1 class="page-title">{{ isEditMode ? '编辑门店配置' : '新增门店' }}</h1>
        <p class="page-subtitle">总部后台统一维护门店基础信息和管理员账号。</p>
      </div>
      <el-button @click="router.push(isEditMode ? `/merchants/${merchantId}` : '/merchants')">返回</el-button>
    </div>

    <el-skeleton :rows="8" animated :loading="loading">
      <div class="section-grid" style="align-items: start;">
        <el-card class="page-card" shadow="never">
          <template #header>
            <span>基础信息</span>
          </template>
          <el-form label-position="top">
            <el-form-item label="门店名称">
              <el-input v-model="basicForm.name" placeholder="请输入门店名称" />
            </el-form-item>
            <el-form-item label="联系人">
              <el-input v-model="basicForm.contact_name" placeholder="请输入联系人姓名" />
            </el-form-item>
            <el-form-item label="联系电话">
              <el-input v-model="basicForm.contact_phone" placeholder="请输入联系电话" />
            </el-form-item>
            <el-form-item label="联系邮箱">
              <el-input v-model="basicForm.contact_email" placeholder="请输入联系邮箱" />
            </el-form-item>
            <el-form-item label="经营分类">
              <el-input v-model="basicForm.business_category" placeholder="如：轻食简餐、茶饮甜品" />
            </el-form-item>
            <el-form-item label="营业时间">
              <el-input v-model="basicForm.business_hours" placeholder="如：09:00-21:00" />
            </el-form-item>
            <el-form-item label="门店地址">
              <el-input v-model="basicForm.address" type="textarea" :rows="3" placeholder="请输入门店地址" />
            </el-form-item>
            <el-form-item label="门店公告">
              <el-input v-model="basicForm.announcement" type="textarea" :rows="3" placeholder="请输入门店公告" />
            </el-form-item>
            <template v-if="!isEditMode">
              <el-divider>管理员账号</el-divider>
              <el-form-item label="登录账号">
                <el-input v-model="basicForm.username" placeholder="请输入登录账号" />
              </el-form-item>
              <el-form-item label="登录密码">
                <el-input v-model="basicForm.password" show-password type="password" placeholder="请输入 6 位以上密码" />
              </el-form-item>
              <el-form-item label="员工姓名">
                <el-input v-model="basicForm.staff_name" placeholder="默认同步联系人姓名" />
              </el-form-item>
              <el-form-item label="员工电话">
                <el-input v-model="basicForm.staff_phone" placeholder="默认同步联系人电话" />
              </el-form-item>
            </template>
            <el-button type="primary" :loading="savingBasic" style="margin-top: 12px;" @click="submitBasic">
              {{ isEditMode ? '保存基础信息' : '创建门店' }}
            </el-button>
          </el-form>
        </el-card>

        <el-card v-if="isEditMode" class="page-card" shadow="never">
          <template #header>
            <span>管理员账号</span>
          </template>
          <el-form label-position="top">
            <el-form-item label="登录账号">
              <el-input :model-value="basicForm.username || '未配置管理员账号'" disabled />
            </el-form-item>
            <el-form-item label="员工姓名">
              <el-input :model-value="basicForm.staff_name || '-'" disabled />
            </el-form-item>
            <el-form-item label="员工电话">
              <el-input :model-value="basicForm.staff_phone || '-'" disabled />
            </el-form-item>
            <el-form-item label="新登录密码">
              <el-input
                v-model="adminForm.new_password"
                show-password
                type="password"
                placeholder="请输入 6 位以上新密码"
              />
            </el-form-item>
            <el-alert
              title="重置后将直接覆盖当前商家管理员账号密码，仅影响负责人 owner 账号。"
              type="info"
              :closable="false"
              show-icon
            />
            <el-button
              type="warning"
              plain
              :loading="savingAdminPassword"
              :disabled="!basicForm.username"
              style="margin-top: 12px;"
              @click="submitAdminPassword"
            >
              重置管理员密码
            </el-button>
          </el-form>
        </el-card>
      </div>
    </el-skeleton>
  </div>
</template>
