<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { changeSpPassword } from '@/api/sp'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const passwordSaving = ref(false)
const info = computed(() => authStore.adminInfo)

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

async function savePassword() {
  if (!passwordForm.old_password.trim()) {
    ElMessage.warning('请输入旧密码')
    return
  }
  if (passwordForm.new_password.trim().length < 6) {
    ElMessage.warning('新密码至少 6 位')
    return
  }
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    ElMessage.warning('两次密码不一致')
    return
  }

  passwordSaving.value = true
  try {
    await changeSpPassword({
      old_password: passwordForm.old_password.trim(),
      new_password: passwordForm.new_password.trim()
    })
    ElMessage.success('密码修改成功')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } finally {
    passwordSaving.value = false
  }
}
</script>

<template>
  <div class="page-shell">
    <div class="page-header">
      <div class="page-title-wrap">
        <h1 class="page-title">后台设置</h1>
        <p class="page-subtitle">查看当前管理员信息并调整登录密码。</p>
      </div>
    </div>

    <div>
      <div class="section-grid">
        <el-card class="page-card" shadow="never">
          <template #header>
            <span>当前账号</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="后台名称">{{ info?.name || '未设置' }}</el-descriptions-item>
            <el-descriptions-item label="显示名称">{{ info?.display_name || '未设置' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>

      <el-card class="page-card" shadow="never" style="margin-top: 20px;">
        <template #header>
          <span>修改密码</span>
        </template>
        <div class="detail-grid">
          <el-form label-position="top">
            <el-form-item label="旧密码">
              <el-input v-model="passwordForm.old_password" show-password type="password" placeholder="请输入旧密码" />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="passwordForm.new_password" show-password type="password" placeholder="请输入新密码" />
            </el-form-item>
            <el-form-item label="确认新密码">
              <el-input v-model="passwordForm.confirm_password" show-password type="password" placeholder="请再次输入新密码" />
            </el-form-item>
            <el-button type="primary" :loading="passwordSaving" @click="savePassword">确认修改</el-button>
          </el-form>
        </div>
      </el-card>
    </div>
  </div>
</template>
