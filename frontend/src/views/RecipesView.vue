<template>
  <div class="max-w-4xl mx-auto px-4 py-4">
    <!-- 헤더 -->
    <header class="mb-4">
      <h1 class="text-xl font-bold text-gray-900 mb-1">🍲 레시피 찾기</h1>
      <p class="text-sm text-gray-600">가진 재료로 만들 수 있는 요리</p>
    </header>

    <!-- 검색 -->
    <div class="relative mb-3">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-gray-400" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="레시피 검색..."
        class="w-full pl-9 pr-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:border-blue-500 bg-white"
      />
    </div>

    <!-- 필터 탭 -->
    <div class="flex gap-2 mb-4">
      <button
        v-for="tab in TABS"
        :key="tab.value"
        :class="[
          'px-3 py-1.5 rounded-full text-xs font-medium transition-colors',
          activeTab === tab.value
            ? 'bg-blue-600 text-white'
            : 'bg-white text-gray-600 border border-gray-200 hover:border-blue-300',
        ]"
        @click="activeTab = tab.value"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 냉장고 재료 요약 -->
    <div v-if="ingredientStore.ingredients.length > 0" class="bg-blue-50 rounded-xl px-4 py-2.5 mb-4 flex items-center gap-2">
      <span class="text-sm text-blue-700">
        🧊 냉장고 재료 <strong>{{ ingredientStore.ingredients.length }}개</strong> 기준으로 매칭 중
      </span>
    </div>
    <div v-else class="bg-gray-50 rounded-xl px-4 py-2.5 mb-4">
      <span class="text-sm text-gray-500">재료를 추가하면 냉장고 재료로 매칭율을 계산해 드려요 🥕</span>
    </div>

    <!-- 로딩 -->
    <div v-if="ingredientStore.loading" class="text-sm text-gray-400 text-center py-12">
      불러오는 중...
    </div>

    <!-- 빈 결과 -->
    <div v-else-if="filteredRecipes.length === 0" class="text-center py-12">
      <p class="text-gray-400 text-sm">검색 결과가 없어요.</p>
    </div>

    <!-- 레시피 목록 -->
    <ul v-else class="space-y-3">
      <li
        v-for="recipe in filteredRecipes"
        :key="recipe.id"
        class="bg-white rounded-xl shadow-sm overflow-hidden cursor-pointer hover:shadow-md transition-shadow"
        @click="$router.push(`/recipes/${recipe.id}`)"
      >
        <div class="p-4">
          <div class="flex items-start justify-between gap-3 mb-2">
            <h3 class="text-base font-bold text-gray-900 leading-snug">{{ recipe.title }}</h3>
            <span
              :class="matchBadgeClass(recipe.matchRate)"
              class="flex-shrink-0 text-xs font-bold px-2 py-0.5 rounded-full"
            >
              {{ recipe.matchRate }}%
            </span>
          </div>
          <p class="text-sm text-gray-500 mb-3 line-clamp-2">{{ recipe.description }}</p>

          <!-- 재료 태그 -->
          <div class="flex flex-wrap gap-1 mb-3">
            <span
              v-for="ing in recipe.ingredients.slice(0, 5)"
              :key="ing"
              :class="
                hasIngredient(ing)
                  ? 'bg-green-50 text-green-700 border-green-200'
                  : 'bg-gray-50 text-gray-500 border-gray-200'
              "
              class="text-xs px-2 py-0.5 rounded-full border"
            >
              {{ hasIngredient(ing) ? '✓ ' : '' }}{{ ing }}
            </span>
            <span v-if="recipe.ingredients.length > 5" class="text-xs text-gray-400 px-1 py-0.5">
              +{{ recipe.ingredients.length - 5 }}개
            </span>
          </div>

          <!-- 요리 정보 -->
          <div class="flex gap-4 text-xs text-gray-400">
            <span v-if="recipe.cookTime">⏱ {{ recipe.cookTime }}</span>
            <span v-if="recipe.servings">👥 {{ recipe.servings }}</span>
            <span v-if="recipe.difficulty">{{ recipe.difficulty }}</span>
          </div>
        </div>

        <!-- 매칭율 바 -->
        <div class="h-1 bg-gray-100">
          <div
            :class="matchBarClass(recipe.matchRate)"
            class="h-full transition-all"
            :style="{ width: `${recipe.matchRate}%` }"
          />
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Search } from 'lucide-vue-next'
import { useIngredientStore } from '@/stores/ingredient.js'
import { DUMMY_RECIPES, calcMatchRate } from '@/data/recipes.js'

const ingredientStore = useIngredientStore()

const searchQuery = ref('')
const activeTab = ref('all')

const TABS = [
  { label: '전체', value: 'all' },
  { label: '100% 매칭', value: 'perfect' },
  { label: '75% 이상', value: 'high' },
]

const fridgeNames = computed(() =>
  ingredientStore.ingredients.map((i) => i.name.toLowerCase()),
)

function hasIngredient(name) {
  return fridgeNames.value.some((n) => n.includes(name.toLowerCase()) || name.toLowerCase().includes(n))
}

const recipes = computed(() =>
  DUMMY_RECIPES.map((r) => ({ ...r, matchRate: calcMatchRate(r, fridgeNames.value) }))
    .sort((a, b) => b.matchRate - a.matchRate),
)

const filteredRecipes = computed(() => {
  let list = recipes.value

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(
      (r) =>
        r.title.toLowerCase().includes(q) ||
        r.ingredients.some((ing) => ing.toLowerCase().includes(q)),
    )
  }

  if (activeTab.value === 'perfect') {
    list = list.filter((r) => r.matchRate === 100)
  } else if (activeTab.value === 'high') {
    list = list.filter((r) => r.matchRate >= 75)
  }

  return list
})

function matchBadgeClass(rate) {
  if (rate === 100) return 'bg-green-100 text-green-700'
  if (rate >= 75) return 'bg-blue-100 text-blue-700'
  if (rate >= 50) return 'bg-yellow-100 text-yellow-700'
  return 'bg-gray-100 text-gray-500'
}

function matchBarClass(rate) {
  if (rate === 100) return 'bg-green-500'
  if (rate >= 75) return 'bg-blue-500'
  if (rate >= 50) return 'bg-yellow-400'
  return 'bg-gray-300'
}

onMounted(async () => {
  if (ingredientStore.ingredients.length === 0) {
    await ingredientStore.fetchIngredients()
  }
})
</script>
