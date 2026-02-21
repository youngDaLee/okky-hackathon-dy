<template>
  <div class="max-w-4xl mx-auto px-4 py-4 space-y-4">
    <header class="mb-2 flex items-start justify-between">
      <div>
        <h1 class="text-xl font-bold text-gray-900">🍳 나의 냉장고</h1>
        <p class="text-sm text-gray-500">똑똑한 재료 관리와 맞춤 레시피 추천</p>
      </div>
      <button
        @click="handleLogout"
        class="text-xs text-gray-400 hover:text-red-500 transition-colors mt-1 flex items-center gap-1"
      >
        <LogOut class="size-3.5" />
        로그아웃
      </button>
    </header>

    <AlertBar :count="store.urgentCount" />

    <!-- 오늘의 추천 -->
    <div v-if="todayRecipes.length > 0" class="bg-white rounded-xl shadow-sm p-4">
      <div class="mb-3">
        <h2 class="text-base font-bold text-gray-900">🔥 오늘의 추천</h2>
        <p class="text-xs text-gray-500">유통기한 임박 재료로 만들 수 있는 추천 메뉴</p>
      </div>
      <RecipeCarousel :recipes="todayRecipes" />
    </div>

    <MiniInventory :items="store.ingredients.slice(0, 5)" :loading="store.loading" />

    <Cookbook :bookmarks="[]" />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { LogOut } from 'lucide-vue-next'
import { useIngredientStore } from '@/stores/ingredient.js'
import { useAuthStore } from '@/stores/auth.js'
import { useRecipeStore } from '@/stores/recipe.js'
import AlertBar from '@/components/AlertBar.vue'
import RecipeCarousel from '@/components/RecipeCarousel.vue'
import MiniInventory from '@/components/MiniInventory.vue'
import Cookbook from '@/components/Cookbook.vue'

const router = useRouter()
const store = useIngredientStore()
const auth = useAuthStore()
const recipeStore = useRecipeStore()

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}

// 오늘의 추천 레시피 (URGENT 재료 기반)
const todayRecipes = computed(() => {
  return recipeStore.todayRecommendations
    .map((result) => {
      const formatted = recipeStore.formatRecommendationResult(result)
      return {
        ...formatted,
        cookTime: formatted.cookingTimeMin ? `${formatted.cookingTimeMin}분` : null,
        servings: null,
        difficulty: formatted.difficulty || null,
      }
    })
    .slice(0, 5)
})

onMounted(async () => {
  await Promise.all([
    store.fetchIngredients(),
    store.fetchSummary(),
    recipeStore.fetchTodayRecommendations(),
  ])
})
</script>
