<template>
  <div class="max-w-4xl mx-auto px-4 py-4">
    <!-- 상단 내비 -->
    <div class="flex items-center gap-2 mb-4">
      <button
        @click="router.back()"
        class="p-2 hover:bg-gray-100 rounded-lg transition-colors -ml-2"
        aria-label="뒤로 가기"
      >
        <ArrowLeft class="size-5" />
      </button>
      <h1 class="font-semibold text-base">레시피 상세</h1>
    </div>

    <!-- 레시피를 찾지 못한 경우 -->
    <div v-if="!recipe" class="text-center py-16">
      <p class="text-gray-400 text-sm">레시피를 찾을 수 없어요.</p>
      <button
        @click="router.back()"
        class="mt-3 text-blue-600 text-sm hover:underline"
      >
        돌아가기
      </button>
    </div>

    <template v-else>
      <!-- 헤더 -->
      <div class="bg-gradient-to-br from-orange-50 to-yellow-50 rounded-2xl p-5 mb-4">
        <div class="flex items-start justify-between gap-3 mb-2">
          <h2 class="text-2xl font-bold text-gray-900">{{ recipe.title }}</h2>
          <span
            :class="matchBadgeClass(recipe.matchRate)"
            class="flex-shrink-0 text-sm font-bold px-3 py-1 rounded-full"
          >
            {{ recipe.matchRate }}% 매칭
          </span>
        </div>
        <p class="text-sm text-gray-600 mb-4">{{ recipe.description }}</p>
        <div class="flex gap-4 text-sm text-gray-500">
          <span v-if="recipe.cookTime">⏱ {{ recipe.cookTime }}</span>
          <span v-if="recipe.servings">👥 {{ recipe.servings }}</span>
          <span v-if="recipe.difficulty">{{ recipe.difficulty }}</span>
        </div>
      </div>

      <!-- 재료 -->
      <div class="bg-white rounded-xl shadow-sm p-4 mb-4">
        <h3 class="text-base font-bold text-gray-900 mb-3">재료</h3>
        <ul class="space-y-2">
          <li
            v-for="ing in recipe.ingredients"
            :key="ing"
            class="flex items-center gap-2 text-sm"
          >
            <span
              :class="hasIngredient(ing) ? 'text-green-500' : 'text-gray-300'"
            >
              <CheckCircle v-if="hasIngredient(ing)" class="size-4" />
              <Circle v-else class="size-4" />
            </span>
            <span :class="hasIngredient(ing) ? 'text-gray-900 font-medium' : 'text-gray-500'">
              {{ ing }}
            </span>
            <span v-if="hasIngredient(ing)" class="text-xs text-green-600 bg-green-50 px-1.5 py-0.5 rounded-full">
              냉장고 있음
            </span>
          </li>
        </ul>

        <!-- 매칭율 바 -->
        <div class="mt-4 pt-3 border-t border-gray-100">
          <div class="flex items-center justify-between text-xs text-gray-500 mb-1">
            <span>냉장고 재료 매칭</span>
            <span class="font-medium">{{ matchedCount }}/{{ recipe.ingredients.length }}개</span>
          </div>
          <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
            <div
              :class="matchBarClass(recipe.matchRate)"
              class="h-full rounded-full transition-all"
              :style="{ width: `${recipe.matchRate}%` }"
            />
          </div>
        </div>
      </div>

      <!-- 조리 순서 -->
      <div class="bg-white rounded-xl shadow-sm p-4 mb-6">
        <h3 class="text-base font-bold text-gray-900 mb-3">조리 순서</h3>
        <ol class="space-y-3">
          <li
            v-for="(step, idx) in recipe.steps"
            :key="idx"
            class="flex gap-3"
          >
            <span
              class="flex-shrink-0 w-6 h-6 bg-blue-600 text-white text-xs font-bold rounded-full flex items-center justify-center mt-0.5"
            >
              {{ idx + 1 }}
            </span>
            <p class="text-sm text-gray-700 leading-relaxed">{{ step }}</p>
          </li>
        </ol>
      </div>

      <!-- 안내 문구 (백엔드 미연동) -->
      <div class="bg-gray-50 rounded-xl px-4 py-3 text-center">
        <p class="text-xs text-gray-400">🚧 레시피 북마크 및 평점 기능은 준비 중이에요.</p>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, CheckCircle, Circle } from 'lucide-vue-next'
import { useIngredientStore } from '@/stores/ingredient.js'
import { DUMMY_RECIPES, calcMatchRate } from '@/data/recipes.js'

const router = useRouter()
const route = useRoute()
const ingredientStore = useIngredientStore()

const fridgeNames = computed(() =>
  ingredientStore.ingredients.map((i) => i.name.toLowerCase()),
)

function hasIngredient(name) {
  return fridgeNames.value.some((n) => n.includes(name.toLowerCase()) || name.toLowerCase().includes(n))
}

const recipe = computed(() => {
  const found = DUMMY_RECIPES.find((r) => r.id === route.params.id)
  if (!found) return null
  return { ...found, matchRate: calcMatchRate(found, fridgeNames.value) }
})

const matchedCount = computed(() =>
  recipe.value ? recipe.value.ingredients.filter((ing) => hasIngredient(ing)).length : 0,
)

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
