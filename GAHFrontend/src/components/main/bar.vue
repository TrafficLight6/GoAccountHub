<template>
  <el-container class="layout">
    <!-- Left Navigation -->
    <el-aside :width="collapsed ? '64px' : '220px'" class="layout-aside">
      <div class="logo">
        <span v-if="!collapsed">GoAccountHub</span>
        <span v-else>GAH</span>
      </div>
      <el-menu :default-active="activeMenu" :collapse="collapsed" :collapse-transition="false" router
        class="layout-menu">
        <el-menu-item index="/main/home">
          <el-icon>
            <Odometer />
          </el-icon>
          <template #title>首页</template>
        </el-menu-item>
        <el-menu-item index="/main/admin" v-if="canGetAdmin">
          <el-icon>
            <SetUp />
          </el-icon>
          <template #title>管理员管理</template>
        </el-menu-item>
        <el-menu-item index="/main/user" v-if="canGetUser">
          <el-icon>
            <User />
          </el-icon>
          <template #title>用户管理</template>
        </el-menu-item>
        <el-menu-item index="/main/character" v-if="canGetCharacter">
          <el-icon>
            <Grid />
          </el-icon>
          <template #title>角色管理</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <!-- 顶部当前页栏 -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="collapsed = !collapsed">
            <Fold v-if="!collapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/main/home' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="pageTitle !== '首页'">{{ pageTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Odometer, Fold, Expand, User, SetUp, Grid } from '@element-plus/icons-vue'
import { get, post, put, del } from '../../lib/request'

const route = useRoute()
const router = useRouter()
const collapsed = ref(false)

// 当前激活菜单：取路由路径
const activeMenu = computed(() => route.path)
// 当前页面名称：优先取路由 meta.title
const pageTitle = computed(() => route.meta.title || '首页')

const handleLogout = async () => {
  try {
    await del('/admin/logout', {}, { showError: false })
  } catch {
    // 即使后端出错也继续清理并跳转
  }
  // 删除本地 cookie
  document.cookie = 'admin_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;'
  document.cookie = 'admin_name=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;'
  router.push('/')
}

const adminInfo = ref({})

// 权限判断：root 全部放行；普通管理员读 permission 字段；数据未加载时返回 false
const canGetAdmin = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_get_admin)
})

// 权限判断：root 全部放行；普通管理员读 permission 字段；数据未加载时返回 false
const canGetUser = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_get_user)
})

// 权限判断：root 全部放行；普通管理员读 permission 字段；数据未加载时返回 false
const canGetCharacter = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_get_character)
})

onMounted(()=>{
  // Check Admin Token
  get('/admin/check_token').then(respond=>{
    if(respond.code === 200){
    }else{
    }
  })
  // Get Admin Info
  get('/admin/info').then(respond=>{
    if(respond.code === 200){
      adminInfo.value = respond.data
    }
  })  
})
</script>

<style scoped>
.layout {
  height: 100vh;
}

.layout-aside {
  background-color: #ffffff;
  border-right: 1px solid #e4e7ed;
  transition: width 0.28s;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #409eff;
  border-bottom: 1px solid #e4e7ed;
  white-space: nowrap;
}

.layout-menu {
  border-right: none;
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #ffffff;
  border-bottom: 1px solid #e4e7ed;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
}

.collapse-btn:hover {
  color: #409eff;
}

.page-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.layout-main {
  background-color: #f5f7fa;
  padding: 20px;
}
</style>
