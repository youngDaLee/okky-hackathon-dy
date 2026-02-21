import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import HomeView from '@/views/HomeView.vue'
import LoginView from '@/views/LoginView.vue'
import SignupView from '@/views/SignupView.vue'
import DashboardView from '@/views/DashboardView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // 인증 전 (공개)
    { path: '/home',   name: 'home',   component: HomeView },
    { path: '/login',  name: 'login',  component: LoginView },
    { path: '/signup', name: 'signup', component: SignupView },

    // 인증 후 (AppLayout 하위)
    {
      path: '/',
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        { path: '',               name: 'dashboard',      component: DashboardView },
        { path: 'recipes',        name: 'recipes',        component: () => import('@/views/RecipesView.vue') },
        { path: 'recipes/:id',    name: 'recipe-detail',  component: () => import('@/views/RecipeDetailView.vue') },
        { path: 'add-ingredient', name: 'add-ingredient', component: () => import('@/views/AddIngredientView.vue') },
        { path: 'fridge',         name: 'fridge',         component: () => import('@/views/FridgeView.vue') },
      ],
    },

    // 기본 리다이렉트
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

router.beforeEach((to) => {
  const token = localStorage.getItem('access_token')
  if (to.meta.requiresAuth && !token) {
    return { name: 'login' }
  }
  if (!to.meta.requiresAuth && token && (to.name === 'login' || to.name === 'signup')) {
    return { name: 'dashboard' }
  }
})

export default router
