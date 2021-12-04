import {
  createRouter,
  createWebHashHistory
} from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior: () => ({ y: 0 }),
  routes: [
    {
      path: '/',
      component: () => import('@/views/index.vue'),
      redirect: '/home',
      name: 'index',
      hidden: true,
      children: [
        { name: 'home', path: '/home', component: () => import('@/views/home.vue') },
        { name: 'imp_xlsx', path: '/imp_xlsx', component: () => import('@/views/imp_xlsx.vue') },
        { name: 'imp_docx', path: '/imp_docx', component: () => import('@/views/imp_docx.vue') },
        { name: 'exp_xlsx', path: '/exp_xlsx', component: () => import('@/views/exp_xlsx.vue') }
      ]
    }
  ]
})

export default router
