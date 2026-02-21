<template>
  <div class="max-w-4xl mx-auto px-4 py-4 space-y-4">
    <header class="mb-2">
      <h1 class="text-xl font-bold text-gray-900">🍳 나의 냉장고</h1>
      <p class="text-sm text-gray-500">똑똑한 재료 관리와 맞춤 레시피 추천</p>
    </header>

    <AlertBar :count="store.urgentCount" />

    <RecipeCarousel :recipes="topRecipes" />

    <MiniInventory :items="store.ingredients.slice(0, 5)" :loading="store.loading" />

    <Cookbook :bookmarks="[]" />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useIngredientStore } from '@/stores/ingredient.js'
import { DUMMY_RECIPES, calcMatchRate } from '@/data/recipes.js'
import AlertBar from '@/components/AlertBar.vue'
import RecipeCarousel from '@/components/RecipeCarousel.vue'
import MiniInventory from '@/components/MiniInventory.vue'
import Cookbook from '@/components/Cookbook.vue'

const store = useIngredientStore()

const fridgeNames = computed(() =>
  store.ingredients.map((i) => i.name.toLowerCase()),
)

const topRecipes = computed(() =>
  DUMMY_RECIPES
    .map((r) => ({ ...r, matchRate: calcMatchRate(r, fridgeNames.value) }))
    .filter((r) => r.matchRate === 100)
    .slice(0, 5),
)

onMounted(async () => {
  await Promise.all([store.fetchIngredients(), store.fetchSummary()])
})
</script>
