<template>
  <div class="home">
    <el-card>
      <template #header>
        <h2>欢迎使用 GoAccountHub 管理后台</h2>
      </template>
      <p>当前登录用户：{{ adminInfo.username || adminName || '未知' }}</p>
    </el-card>

    <el-row :gutter="20" class="card-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <h2>系统信息</h2>
          </template>
          <p>管理员数量：{{ appInfo.admin_count || '未知' }}</p>
          <p>用户数量：{{ appInfo.user_count || '未知' }}</p>
          <p>角色数量：{{ appInfo.character_count || '未知' }}</p>
          <p>总token数量：{{ appInfo.total_token_count || '未知' }}</p>
          <p>管理员token数量：{{ appInfo.admin_token_count || '未知' }}</p>
          <p>用户token数量：{{ appInfo.user_token_count || '未知' }}</p>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <h2>权限</h2>
          </template>
          <p>添加管理员：<el-icon :class="hasPerm('can_add_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_add_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <p>删除管理员：<el-icon :class="hasPerm('can_delete_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_delete_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <p>修改管理员信息：<el-icon :class="hasPerm('can_edit_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_edit_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <p>查看管理员信息：<el-icon :class="hasPerm('can_get_admin') ? 'perm-yes' : 'perm-no'"><CircleCheckFilled v-if="hasPerm('can_get_admin')" /><CircleCloseFilled v-else /></el-icon></p>
          <br><br><br>      <!--用于对齐card内容-->
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { CircleCheckFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import { post } from '../../../lib/request'

// 从 document.cookie 字符串中解析 admin_name（接口未返回前的兜底显示）
const adminName = ref('')
const match = document.cookie.match(/(?:^|;\s*)admin_name=([^;]*)/)
if (match) {
  adminName.value = decodeURIComponent(match[1])
}

const appInfo = ref({})
const adminInfo = ref({})

// 权限判断：root 全部放行；普通管理员读 permission 字段；数据未加载时返回 false
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
