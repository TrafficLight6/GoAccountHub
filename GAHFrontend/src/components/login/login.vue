<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">
          <h2>Login</h2>
        </div>
      </template>

      <el-form ref="loginFormRef" :model="loginForm" :rules="rules" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="Username" prop="username">
          <el-input v-model="loginForm.username" placeholder="Please enter username" :prefix-icon="User" clearable />
        </el-form-item>

        <el-form-item label="Password" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="Please enter password" :prefix-icon="Lock" show-password
            @keyup.enter="handleLogin" />
        </el-form-item>

        <el-form-item>
          <el-checkbox v-model="loginForm.keepLoggedIn">Keep me logged in</el-checkbox>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" class="login-btn" @click="handleLogin">
            Login
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
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
    { required: true, message: 'Please enter username', trigger: 'blur' },
    { min: 0, max: 32, message: 'Username must be at most 32 characters', trigger: 'blur' },
  ],
  password: [
    { required: true, message: 'Please enter password', trigger: 'blur' },
    { min: 0, max: 64, message: 'Password must be at most 64 characters', trigger: 'blur' },
  ],
}

// Redirect to home if already logged in
onMounted(async () => {
  try {
    await post('/admin/check_token', {}, { showError: false, autoRedirect401: false })
    router.push('/main/home')
  } catch {
    // Not logged in, stay on the login page
  }
})

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
    ElMessage.success('Logged in')
    router.push('/main/home')
  } catch (e) {
    if (e.code === 400) {
      // Invalid credentials (wrong username/password), unified message
      ElMessage.error('Incorrect username or password')
    } else if (e.code) {
      // Other backend business errors
      ElMessage.error(e.message || 'Login failed, please try again')
    } else {
      // Network layer error (offline / backend not running)
      ElMessage.error('Network request failed, please check your network or backend service')
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
