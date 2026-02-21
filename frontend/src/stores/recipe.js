import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { recommendationApi } from '@/api/recommendation.js'

export const useRecipeStore = defineStore('recipe', () => {
  const recommendations = ref([])
  const todayRecommendations = ref([])
  const searchResults = ref([])
  const currentRecipe = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // 추천 레시피 조회
  async function fetchRecommendations(params = {}) {
    loading.value = true
    error.value = null
    try {
      const { data } = await recommendationApi.getRecommendations(params)
      recommendations.value = data.items ?? []
      // fridge_ingredient_count는 필요시 사용 가능
      return {
        items: data.items ?? [],
        total: data.total ?? 0,
        fridgeIngredientCount: data.fridge_ingredient_count ?? data.fridgeIngredientCount ?? 0,
        message: data.message,
      }
    } catch (e) {
      error.value = e
      throw e
    } finally {
      loading.value = false
    }
  }

  // 오늘의 추천 레시피 조회
  async function fetchTodayRecommendations() {
    loading.value = true
    error.value = null
    try {
      const { data } = await recommendationApi.getTodayRecommendations()
      todayRecommendations.value = data.items ?? []
      return {
        items: data.items ?? [],
        total: data.total ?? 0,
        urgentIngredients: data.urgent_ingredients ?? data.urgentIngredients ?? [],
      }
    } catch (e) {
      error.value = e
      throw e
    } finally {
      loading.value = false
    }
  }

  // 레시피 검색
  async function searchRecipes(params = {}) {
    loading.value = true
    error.value = null
    try {
      const { data } = await recommendationApi.searchRecipes(params)
      searchResults.value = data.items ?? []
      return data
    } catch (e) {
      error.value = e
      throw e
    } finally {
      loading.value = false
    }
  }

  // 레시피 상세 조회
  async function fetchRecipeById(id) {
    loading.value = true
    error.value = null
    try {
      const { data } = await recommendationApi.getRecipeById(id)
      currentRecipe.value = data
      return data
    } catch (e) {
      error.value = e
      throw e
    } finally {
      loading.value = false
    }
  }

  // RecommendationResult를 화면에서 사용하기 쉬운 형태로 변환
  function formatRecommendationResult(result) {
    const recipe = result.recipe || result
    // 백엔드 응답은 camelCase (requiredIngredients, matchScore 등)
    // API 스펙 문서는 snake_case지만 실제 구현은 camelCase
    const requiredIngredients = recipe.requiredIngredients || []
    const optionalIngredients = recipe.optionalIngredients || []
    
    return {
      id: recipe.id,
      title: recipe.title,
      description: recipe.description,
      ingredients: [...requiredIngredients, ...optionalIngredients], // 전체 재료 목록
      requiredIngredients,
      optionalIngredients,
      mainIngredient: recipe.mainIngredient,
      category: recipe.category,
      tags: recipe.tags || [],
      cookingTimeMin: recipe.cookingTimeMin,
      difficulty: recipe.difficulty,
      sourceType: recipe.sourceType,
      sourceUrl: recipe.sourceUrl,
      thumbnailUrl: recipe.thumbnailUrl,
      // RecommendationResult 전용 필드 (camelCase)
      tier: result.tier,
      matchScore: result.matchScore ?? result.match_score ?? 0,
      matchRate: Math.round((result.matchScore ?? result.match_score ?? 0) * 100),
      matchedIngredients: result.matchedIngredients ?? result.matched_ingredients ?? [],
      missingIngredients: result.missingIngredients ?? result.missing_ingredients ?? [],
      urgencyBonus: result.urgencyBonus ?? result.urgency_bonus ?? false,
    }
  }

  return {
    recommendations,
    todayRecommendations,
    searchResults,
    currentRecipe,
    loading,
    error,
    fetchRecommendations,
    fetchTodayRecommendations,
    searchRecipes,
    fetchRecipeById,
    formatRecommendationResult,
  }
})
