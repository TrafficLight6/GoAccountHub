import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/login/login.vue'
import Bar from '../components/main/bar.vue'
import Home from '../components/main/page/home.vue'
import Admin from '../components/main/page/admin.vue'
import User from '../components/main/page/user.vue'
import Character from '../components/main/page/character.vue'
import ComponentTest from '../components/component_test.vue'

const routes = [
  {
    path: '/',
    name: 'Login',
    component: Login,
  },
  {
    path: '/component',
    name: 'ComponentTest',
    component: () => import('../components/component_test.vue'),
  },
  {
    path: '/main',
    component: Bar,
    redirect: '/main/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: Home,
        meta: { title: '首页' },
      },
      {
        path: 'admin',
        name: 'Admin',
        component: Admin,
        meta: { title: '管理员管理' },
      },
      {
        path: 'user',
        name: 'User',
        component: User,
        meta: { title: '用户管理' },
      },
      {
        path: 'character',
        name: 'Character',
        component: Character,
        meta: { title: '角色管理' },
           },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
