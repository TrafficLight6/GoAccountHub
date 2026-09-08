import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/login/login.vue'
import Home from '../components/main/home.vue'

const routes = [
  {
    path: '/',
    name: 'Login',
    component: Login,
  },
  {
    path: '/home',
    name: 'Home',
    component: Home,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
