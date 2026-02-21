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

    <!-- 로딩 -->
    <div v-if="loading" class="text-center py-16">
      <p class="text-gray-400 text-sm">레시피를 불러오는 중...</p>
    </div>

    <!-- 에러 -->
    <div v-else-if="error" class="text-center py-16">
      <p class="text-red-500 text-sm">레시피를 불러오는데 실패했어요.</p>
      <button
        @click="router.back()"
        class="mt-3 text-blue-600 text-sm hover:underline"
      >
        돌아가기
      </button>
    </div>

    <!-- 레시피를 찾지 못한 경우 -->
    <div v-else-if="!recipe" class="text-center py-16">
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
        
        <!-- 필수 재료 -->
        <div v-if="recipe.requiredIngredients && recipe.requiredIngredients.length > 0" class="mb-4">
          <p class="text-xs text-gray-500 mb-2 font-medium">필수 재료</p>
          <ul class="space-y-2">
            <li
              v-for="ing in recipe.requiredIngredients"
              :key="ing"
              class="flex items-center gap-2 text-sm"
            >
              <span
                :class="hasIngredient(ing) ? 'text-green-500' : 'text-red-500'"
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
              <span v-else class="text-xs text-red-600 bg-red-50 px-1.5 py-0.5 rounded-full">
                부족
              </span>
            </li>
          </ul>
        </div>

        <!-- 선택 재료 -->
        <div v-if="recipe.optionalIngredients && recipe.optionalIngredients.length > 0">
          <p class="text-xs text-gray-500 mb-2 font-medium">선택 재료</p>
          <ul class="space-y-2">
            <li
              v-for="ing in recipe.optionalIngredients"
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
        </div>

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
      <div v-if="recipe.steps && recipe.steps.length > 0" class="bg-white rounded-xl shadow-sm p-4 mb-6">
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

      <!-- 외부 링크 -->
      <div v-if="recipe.sourceUrl" class="bg-white rounded-xl shadow-sm p-4 mb-6">
        <h3 class="text-base font-bold text-gray-900 mb-3">레시피 링크</h3>
        <a
          :href="recipe.sourceUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="text-blue-600 text-sm hover:underline flex items-center gap-2"
        >
          {{ recipe.sourceType === 'YOUTUBE' ? '📺 YouTube' : '🔗 외부 링크' }}에서 보기
        </a>
      </div>

      <!-- 요리책 저장 버튼 -->
      <div class="bg-white rounded-xl shadow-sm p-4 mb-6">
        <button
          v-if="!isSaved"
          @click="handleSaveRecipe"
          :disabled="saving"
          class="w-full py-3 bg-blue-600 text-white rounded-xl text-sm font-medium hover:bg-blue-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
        >
          <BookmarkPlus class="size-4" />
          {{ saving ? '저장 중...' : '📖 요리책에 저장' }}
        </button>
        <div v-else class="text-center py-2">
          <p class="text-sm text-green-600 font-medium">✓ 요리책에 저장됨</p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, CheckCircle, Circle, BookmarkPlus } from 'lucide-vue-next'
import { useIngredientStore } from '@/stores/ingredient.js'
import { useRecipeStore } from '@/stores/recipe.js'
import { useCookbookStore } from '@/stores/cookbook.js'

const router = useRouter()
const route = useRoute()
const ingredientStore = useIngredientStore()
const recipeStore = useRecipeStore()
const cookbookStore = useCookbookStore()

const loading = ref(false)
const error = ref(null)
const saving = ref(false)

const fridgeNames = computed(() =>
  ingredientStore.ingredients.map((i) => i.name.toLowerCase()),
)

function hasIngredient(name) {
  return fridgeNames.value.some((n) => n.includes(name.toLowerCase()) || name.toLowerCase().includes(n))
}

// API에서 받은 레시피를 화면용으로 변환
const recipe = computed(() => {
  if (!recipeStore.currentRecipe) return null

  const r = recipeStore.currentRecipe
  // 백엔드는 camelCase로 응답 (requiredIngredients, optionalIngredients)
  // API 스펙 문서는 snake_case지만 실제 구현은 camelCase
  const requiredIngredients = r.requiredIngredients || r.required_ingredients || []
  const optionalIngredients = r.optionalIngredients || r.optional_ingredients || []
  const allIngredients = [...requiredIngredients, ...optionalIngredients]

  // 매칭 계산
  const matched = allIngredients.filter((ing) => hasIngredient(ing))
  const matchRate = allIngredients.length > 0
    ? Math.round((matched.length / allIngredients.length) * 100)
    : 0

  return {
    id: r.id,
    title: r.title,
    description: r.description,
    ingredients: allIngredients,
    requiredIngredients,
    optionalIngredients,
    mainIngredient: r.mainIngredient || r.main_ingredient,
    category: r.category,
    tags: r.tags || [],
    cookingTimeMin: r.cookingTimeMin || r.cooking_time_min,
    difficulty: r.difficulty,
    sourceType: r.sourceType || r.source_type,
    sourceUrl: r.sourceUrl || r.source_url,
    thumbnailUrl: r.thumbnailUrl || r.thumbnail_url,
    matchRate,
    cookTime: (r.cookingTimeMin || r.cooking_time_min) ? `${r.cookingTimeMin || r.cooking_time_min}분` : null,
    servings: null,
    steps: [], // API에 조리 순서가 없으면 빈 배열
  }
})

const matchedCount = computed(() =>
  recipe.value ? recipe.value.ingredients.filter((ing) => hasIngredient(ing)).length : 0,
)

const isSaved = computed(() => {
  if (!recipe.value) return false
  return cookbookStore.savedRecipes.some(
    (saved) => saved.recipeId === recipe.value.id || saved.recipeSnapshot?.title === recipe.value.title
  )
})

async function handleSaveRecipe() {
  if (!recipe.value) return
  
  saving.value = true
  try {
    await cookbookStore.saveRecipe({
      recipeId: recipe.value.id || null,
      recipeSnapshot: {
        title: recipe.value.title,
        sourceUrl: recipe.value.sourceUrl || '',
        thumbnailUrl: recipe.value.thumbnailUrl || '',
        sourceType: recipe.value.sourceType || 'INTERNAL',
        mainIngredient: recipe.value.mainIngredient || '',
        category: recipe.value.category || '',
      },
      label: '',
      note: '',
    })
    // 저장 후 목록 갱신
    await cookbookStore.fetchSavedRecipes()
  } catch (e) {
    console.error('Failed to save recipe:', e)
    const message = e?.message || '레시피 저장에 실패했어요. 다시 시도해주세요.'
    alert(message)
  } finally {
    saving.value = false
  }
}

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
  loading.value = true
  error.value = null

  try {
    // 재료 목록 로드
    if (ingredientStore.ingredients.length === 0) {
      await ingredientStore.fetchIngredients()
    }
    // 레시피 상세 조회
    await recipeStore.fetchRecipeById(route.params.id)
    // 저장된 레시피 목록 로드 (저장 여부 확인용) - 실패해도 계속 진행
    try {
      await cookbookStore.fetchSavedRecipes()
    } catch (cookbookErr) {
      console.warn('Failed to load saved recipes (non-critical):', cookbookErr)
    }
  } catch (e) {
    error.value = e
    console.error('Failed to load recipe:', e)
  } finally {
    loading.value = false
  }
})
</script>
