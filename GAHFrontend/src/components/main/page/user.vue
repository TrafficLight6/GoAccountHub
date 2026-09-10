<template>
    <h2>
        用户管理
    </h2>
</template>
<script setup>
import { get } from '.././../../lib/request.js'
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const adminInfo = ref({})

const canGetUser = computed(() => {
  if (adminInfo.value.is_root) return true
  return Boolean(adminInfo.value.permission?.can_operate_user)
})

onMounted(async () => {
  try {
    const res = await get('/admin/info')
    adminInfo.value = res.data
  } catch {
    // 401 已由 request 封装自动跳转登录页，网络错误时不再继续鉴权
    return
  }
  if (!canGetUser.value) {
    router.push('/')
  }
})
</script>