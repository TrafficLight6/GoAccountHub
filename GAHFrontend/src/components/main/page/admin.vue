<template>
    <el-card>
        <template #header>
            操作
        </template>
        <p><el-input placeholder="管理员用户名"></el-input></p>
        <p><el-input placeholder="管理员UUHash"></el-input></p>
        <p>权限筛选：</p>
        <el-checkbox label="是否可以添加管理员"></el-checkbox>
        <el-checkbox label="是否可以删除管理员"></el-checkbox>
        <el-checkbox label="是否可以修改管理员"></el-checkbox>
        <el-checkbox label="是否可以获取管理员列表"></el-checkbox>
        <el-checkbox label="是否可以操作用户"></el-checkbox>
        <el-checkbox label="是否可以操作角色"></el-checkbox>
        <template #footer>
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button type="success" @click="handleAdd">添加管理员</el-button>
        </template>
    </el-card>
    <br>
    <el-card>
        <template #header>
            管理员列表
        </template>
        <!-- data here -->

    </el-card>
</template>
<script setup>
import { get } from '.././../../lib/request.js'
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const adminInfo = ref({})

const canGetAdmin = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_get_admin)
})

onMounted(async () => {
    try {
        const res = await get('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 已由 request 封装自动跳转登录页，网络错误时不再继续鉴权
        return
    }
    if (!canGetAdmin.value) {
        router.push('/')
    }
})
</script>
