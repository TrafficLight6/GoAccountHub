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
          <template #title>Home</template>
        </el-menu-item>
        <el-menu-item index="/main/admin" v-if="canGetAdmin">
          <el-icon>
            <SetUp />
          </el-icon>
          <template #title>Admin Management</template>
        </el-menu-item>
        <el-menu-item index="/main/user" v-if="canGetUser">
          <el-icon>
            <User />
          </el-icon>
          <template #title>User Management</template>
        </el-menu-item>
        <el-menu-item index="/main/character" v-if="canGetCharacter">
          <el-icon>
            <Grid />
          </el-icon>
          <template #title>Character Management</template>
        </el-menu-item>
        <el-menu-item index="/main/key" v-if="canOperateKey">
          <el-icon>
            <Key />
          </el-icon>
          <template #title>Key Management</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <!-- Top current page bar -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="collapsed = !collapsed">
            <Fold v-if="!collapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/main/home' }">Home</el-breadcrumb-item>
            <el-breadcrumb-item v-if="pageTitle !== 'Home'">{{ pageTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="handleLogout">Logout</el-button>
        </div>
      </el-header>

      <!-- Content area -->
      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Odometer, Fold, Expand, User, SetUp, Grid, Key } from '@element-plus/icons-vue'
import { post, put, del } from '../../lib/request'
import VueCookies from 'vue-cookies'

const route = useRoute()
const router = useRouter()
const collapsed = ref(false)

// Currently active menu: from the route path
const activeMenu = computed(() => route.path)
// Current page title: prefer route meta.title
const pageTitle = computed(() => route.meta.title || 'Home')

const handleLogout = async () => {
  try {
    await del('/admin/logout', {}, { showError: false })
  } catch {
    // Continue cleanup and redirect even if the backend errors
  }
  // Remove local cookies
  VueCookies.remove('admin_token')
  VueCookies.remove('admin_name')
  router.push('/')
}

const adminInfo = ref({})

// Permission check: allow everything for root; read the permission field for normal admins; return false before data loads
const canGetAdmin = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_get_admin)
})

// Permission check: allow everything for root; read the permission field for normal admins; return false before data loads
const canGetUser = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_get_user)
})

// Permission check: allow everything for root; read the permission field for normal admins; return false before data loads
const canGetCharacter = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_get_character)
})

// Permission check: allow everything for root; read the permission field for normal admins; return false before data loads
const canOperateKey = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_operate_app_key)
})

onMounted(()=>{
  // Check Admin Token
  post('/admin/check_token').then(respond=>{
    if(respond.code === 200){
    }else{
    }
  })
  // Get Admin Info
  post('/admin/info').then(respond=>{
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