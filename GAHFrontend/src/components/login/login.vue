<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">
          <h2>登录</h2>
        </div>
      </template>

      <el-form ref="loginFormRef" :model="loginForm" :rules="rules" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名" :prefix-icon="User" clearable />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" :prefix-icon="Lock" show-password
            @keyup.enter="handleLogin" />
        </el-form-item>

        <el-form-item>
          <el-checkbox v-model="loginForm.keepLoggedIn">保持登录</el-checkbox>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" class="login-btn" @click="handleLogin">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { post } from '../../lib/request'

const router = useRouter()
const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
  keepLoggedIn: false,
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 0, max: 32, message: '用户名长度不能超过 32 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 0, max: 64, message: '密码长度不能超过 64 个字符', trigger: 'blur' },
  ],
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  try {
    await loginFormRef.value.validate()
  } catch {
    return
  }
  loading.value = true
  try {
    await post(
      '/admin/login',
      {
        username: loginForm.username,
        password: loginForm.password,
        is_remember: loginForm.keepLoggedIn,
      },
      { showError: false }
    )
    ElMessage.success('登录成功')
    router.push('/main/home')
  } catch (e) {
    if (e.code === 400) {
      // 凭证错误（用户名/密码错误），统一中文提示
      ElMessage.error('用户名或密码错误')
    } else if (e.code) {
      // 其他后端业务错误
      ElMessage.error(e.message || '登录失败，请重试')
    } else {
      // 网络层异常（断网/后端未启动）
      ElMessage.error('网络请求失败，请检查网络或后端服务')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 380px;
  max-width: 90vw;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
}

.login-header h2 {
  margin: 0;
  color: #303133;
}

.login-btn {
  width: 100%;
}
</style>
