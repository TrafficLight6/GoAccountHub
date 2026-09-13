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
        meta: { title: 'Home' },
      },
      {
        path: 'admin',
        name: 'Admin',
        component: Admin,
        meta: { title: 'Admin Management' },
      },
      {
        path: 'user',
        name: 'User',
        component: User,
        meta: { title: 'User Management' },
      },
      {
        path: 'character',
        name: 'Character',
        component: Character,
        meta: { title: 'Character Management' },
      },
      {
        path: 'key',
        name: 'Key',
        component: () => import('../components/main/page/key.vue'),
        meta: { title: 'Key Management' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
