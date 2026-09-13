<template>
    <div>
        <h2>Application Keys</h2>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { post } from '../../../lib/request.js'

const router = useRouter()
const adminInfo = ref({})

// Permission check: root is always allowed; regular admins read the permission field; returns false while data is not loaded
const canOperateAppKey = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_operate_app_key)
})

onMounted(async () => {
    try {
        const res = await post('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 is already handled by the request wrapper, which redirects to the login page; skip further checks on network errors
        return
    }
    if (!canOperateAppKey.value) {
        router.push('/main/home')
        return
    }
})
</script>
