<template>
  <div class="home">
    <el-card>
      <template #header>
        <h2>Welcome to GoAccountHub Admin</h2>
      </template>
      <p>Logged in as: {{ adminInfo.username || adminName || 'Unknown' }}</p>
    </el-card>

    <el-row :gutter="20" class="card-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <h2>System Info</h2>
          </template>
          <p>Admin Count: {{ appInfo.admin_count || 'Unknown' }}</p>
          <p>User Count: {{ appInfo.user_count || 'Unknown' }}</p>
          <p>Character Count: {{ appInfo.character_count || 'Unknown' }}</p>
          <p>Key Count: {{ appInfo.key_count ?? 'Unknown' }}</p>
          <p>Total Tokens: {{ appInfo.total_token_count || 'Unknown' }}</p>
          <p>Admin Tokens: {{ appInfo.admin_token_count || 'Unknown' }}</p>
          <p>User Tokens: {{ appInfo.user_token_count || 'Unknown' }}</p>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <h2>Permissions</h2>
          </template>
          <p>Add Admin: <el-icon :class="hasPerm('can_add_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_add_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <p>Delete Admin: <el-icon :class="hasPerm('can_delete_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_delete_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <p>Edit Admin Info: <el-icon :class="hasPerm('can_edit_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_edit_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <p>View Admin Info: <el-icon :class="hasPerm('can_get_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_get_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <p>Operate App Key: <el-icon :class="hasPerm('can_operate_app_key') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_operate_app_key')" /><CircleCloseFilled v-else /></el-icon></p>
          <br><br>      <!--Align card content-->
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { CircleCheckFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import { post } from '../../../lib/request'

// Parse admin_name from the document.cookie string (fallback display before the API returns)
const adminName = ref('')
const match = document.cookie.match(/(?:^|;\s*)admin_name=([^;]*)/)
if (match) {
  adminName.value = decodeURIComponent(match[1])
}

const appInfo = ref({})
const adminInfo = ref({})

// Permission check: root is always allowed; regular admins read the permission field; returns false while data is not loaded
const hasPerm = (key) => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.[key])
}

onMounted(() => {
  post('/info')
    .then(res => {
      appInfo.value = res.data
    })
  post('/admin/info')
    .then(res => {
      adminInfo.value = res.data
    })
})
</script>

<style scoped>
.card-row {
  margin-top: 20px;
}
.home h2 {
  margin: 0 0 12px;
  color: #303133;
}
.home p {
  margin: 0 0 8px;
  color: #606266;
}
.perm-yes {
  color: #67c23a;
}
.perm-no {
  color: #F56C6C;
}
</style>
