<template>
    <div class="not-found">
        <h1>404</h1>
        <p>The page you are looking for does not exist.</p>
        <el-button type="primary" @click="goHome">Back to Home</el-button>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
    router.push('/main/home')
}
</script>

<style scoped>
.not-found {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background-color: #f5f7fa;
}

.not-found h1 {
    margin: 0;
    font-size: 72px;
    color: #303133;
}

.not-found p {
    margin: 0 0 16px;
    color: #606266;
}
</style>
