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

    <!-- Tier 필터 탭 -->
    <div class="flex gap-2 mb-3 overflow-x-auto pb-2 scrollbar-hide">
      <button
        v-for="tab in TABS"
        :key="tab.value"
        :class="[
          'flex-shrink-0 px-3 py-1.5 rounded-full text-xs font-medium transition-colors',
          activeTab === tab.value
            ? 'bg-blue-600 text-white'
            : 'bg-white text-gray-600 border border-gray-200 hover:border-blue-300',
        ]"
        @click="activeTab = tab.value"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 카테고리 필터 -->
    <div class="flex gap-2 mb-3 overflow-x-auto pb-2 scrollbar-hide">
      <button
        v-for="cat in CATEGORIES"
        :key="cat.value"
        :class="[
          'flex-shrink-0 px-3 py-1.5 rounded-full text-xs font-medium transition-colors',
          selectedCategory === cat.value
            ? 'bg-gray-800 text-white'
            : 'bg-white text-gray-600 border border-gray-200 hover:border-gray-400',
        ]"
        @click="selectedCategory = cat.value"
      >
        {{ cat.label }}
      </button>
    </div>

    <!-- 최대 부족 재료 수 필터 -->
    <div class="flex items-center gap-2 mb-4">
      <label class="text-xs text-gray-600 whitespace-nowrap">부족 재료 최대:</label>
      <select
        v-model="maxMissing"
        class="flex-1 max-w-32 px-3 py-1.5 border border-gray-200 rounded-lg text-xs focus:outline-none focus:border-blue-500 bg-white"
      >
        <option :value="null">제한 없음</option>
        <option :value="0">0개 (완벽 매칭)</option>
        <option :value="1">1개 이하</option>
        <option :value="2">2개 이하</option>
        <option :value="3">3개 이하</option>
      </select>
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
    <div v-if="recipeStore.loading || ingredientStore.loading" class="text-sm text-gray-400 text-center py-12">
      불러오는 중...
    </div>

    <!-- 빈 결과 -->
    <div v-else-if="filteredRecipes.length === 0" class="text-center py-12">
      <p class="text-gray-400 text-sm">
        {{ searchQuery ? '검색 결과가 없어요.' : '추천 레시피가 없어요. 재료를 추가해보세요!' }}
      </p>
    </div>

    <!-- 레시피 목록 -->
    <ul v-else class="space-y-3">
      <li
        v-for="recipe in filteredRecipes"
        :key="recipe.id"
        class="bg-white rounded-xl shadow-sm overflow-hidden cursor-pointer hover:shadow-md transition-shadow"
        @click="$router.push(`/recipes/${recipe.id}`)"
      >
        <!-- 썸네일 이미지 -->
        <div v-if="recipe.thumbnailUrl" class="w-full h-40 bg-gray-100 overflow-hidden">
          <img
            :src="recipe.thumbnailUrl"
            :alt="recipe.title"
            class="w-full h-full object-cover"
          />
        </div>
        <div class="p-4">
          <div class="flex items-start justify-between gap-3 mb-2">
            <div class="flex-1">
              <h3 class="text-base font-bold text-gray-900 leading-snug">{{ recipe.title }}</h3>
              <div class="flex items-center gap-2 mt-1">
                <span
                  :class="matchBadgeClass(recipe.matchRate)"
                  class="text-xs font-bold px-2 py-0.5 rounded-full"
                >
                  Tier {{ recipe.tier }} · {{ recipe.matchRate }}%
                </span>
                <span v-if="recipe.urgencyBonus" class="text-xs text-red-600 font-medium">
                  🔥 긴급
                </span>
              </div>
            </div>
          </div>
          <p class="text-sm text-gray-500 mb-3 line-clamp-2">{{ recipe.description }}</p>

          <!-- 재료 태그 -->
          <div class="flex flex-wrap gap-1 mb-3">
            <span
              v-for="ing in recipe.ingredients.slice(0, 5)"
              :key="ing"
              :class="
                recipe.matchedIngredients?.includes(ing)
                  ? 'bg-green-50 text-green-700 border-green-200'
                  : 'bg-gray-50 text-gray-500 border-gray-200'
              "
              class="text-xs px-2 py-0.5 rounded-full border"
            >
              {{ recipe.matchedIngredients?.includes(ing) ? '✓ ' : '' }}{{ ing }}
            </span>
            <span v-if="recipe.ingredients.length > 5" class="text-xs text-gray-400 px-1 py-0.5">
              +{{ recipe.ingredients.length - 5 }}개
            </span>
            <span v-if="recipe.missingIngredients && recipe.missingIngredients.length > 0" class="text-xs text-orange-500 px-1 py-0.5">
              (부족: {{ recipe.missingIngredients.length }}개)
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
import { ref, computed, onMounted, watch } from 'vue'
import { Search } from 'lucide-vue-next'
import { useIngredientStore } from '@/stores/ingredient.js'
import { useRecipeStore } from '@/stores/recipe.js'

const ingredientStore = useIngredientStore()
const recipeStore = useRecipeStore()

const searchQuery = ref('')
const activeTab = ref('all')
const selectedCategory = ref('')
const maxMissing = ref(null)

const TABS = [
  { label: '전체', value: 'all' },
  { label: 'Tier 1', value: 'tier1' },
  { label: 'Tier 2', value: 'tier2' },
  { label: 'Tier 3', value: 'tier3' },
]

const CATEGORIES = [
  { label: '전체', value: '' },
  { label: '한식', value: '한식' },
  { label: '중식', value: '중식' },
  { label: '일식', value: '일식' },
  { label: '양식', value: '양식' },
  { label: '간식', value: '간식' },
  { label: '기타', value: '기타' },
]

const fridgeNames = computed(() =>
  ingredientStore.ingredients.map((i) => i.name.toLowerCase()),
)

function hasIngredient(name) {
  return fridgeNames.value.some((n) => n.includes(name.toLowerCase()) || name.toLowerCase().includes(n))
}

// API에서 받은 RecommendationResult를 화면용으로 변환
const recipes = computed(() => {
  return recipeStore.recommendations.map((result) => {
    const formatted = recipeStore.formatRecommendationResult(result)
    return {
      ...formatted,
      cookTime: formatted.cookingTimeMin ? `${formatted.cookingTimeMin}분` : null,
      servings: null, // API에 없음
      difficulty: formatted.difficulty || null,
    }
  })
})

const filteredRecipes = computed(() => {
  let list = recipes.value

  // 검색 필터
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(
      (r) =>
        r.title.toLowerCase().includes(q) ||
        r.ingredients.some((ing) => ing.toLowerCase().includes(q)),
    )
  }

  // 카테고리 필터
  if (selectedCategory.value) {
    list = list.filter((r) => r.category === selectedCategory.value)
  }

  // 최대 부족 재료 수 필터
  if (maxMissing.value !== null) {
    list = list.filter((r) => (r.missingIngredients?.length || 0) <= maxMissing.value)
  }

  // Tier 필터
  if (activeTab.value === 'tier1') {
    list = list.filter((r) => r.tier === 1)
  } else if (activeTab.value === 'tier2') {
    list = list.filter((r) => r.tier === 2)
  } else if (activeTab.value === 'tier3') {
    list = list.filter((r) => r.tier === 3)
  }

  // Tier와 urgencyBonus 기준으로 정렬
  list.sort((a, b) => {
    // urgencyBonus가 있으면 우선
    if (a.urgencyBonus && !b.urgencyBonus) return -1
    if (!a.urgencyBonus && b.urgencyBonus) return 1
    // Tier 우선
    if (a.tier !== b.tier) return a.tier - b.tier
    // matchRate 높은 순
    return b.matchRate - a.matchRate
  })

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

// 검색어, Tier, 카테고리, max_missing 변경 시 API 호출
watch([searchQuery, activeTab, selectedCategory, maxMissing], async () => {
  if (searchQuery.value) {
    await recipeStore.searchRecipes({
      keyword: searchQuery.value,
      category: selectedCategory.value || undefined,
      limit: 20,
    })
  } else {
    // Tier 필터에 따라 추천 API 호출
    const params = { limit: 20 }
    if (activeTab.value === 'tier1') params.tier = 1
    else if (activeTab.value === 'tier2') params.tier = 2
    else if (activeTab.value === 'tier3') params.tier = 3
    
    // 카테고리 필터 추가
    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }
    
    // 최대 부족 재료 수 필터 추가
    if (maxMissing.value !== null) {
      params.max_missing = maxMissing.value
    }

    await recipeStore.fetchRecommendations(params)
  }
})

onMounted(async () => {
  // 재료 목록 로드
  if (ingredientStore.ingredients.length === 0) {
    await ingredientStore.fetchIngredients()
  }
  // 추천 레시피 로드
  await recipeStore.fetchRecommendations({ limit: 20 })
})
</script>
