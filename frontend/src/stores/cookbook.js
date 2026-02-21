import { defineStore } from 'pinia'
import { ref } from 'vue'
import { cookbookApi } from '@/api/cookbook.js'

export const useCookbookStore = defineStore('cookbook', () => {
  const savedRecipes = ref([])
  const labels = ref([])
  const loading = ref(false)
  const error = ref(null)

  // 저장된 레시피 목록 조회
  async function fetchSavedRecipes(params = {}) {
    loading.value = true
    error.value = null
    try {
      const { data } = await cookbookApi.getList(params)
      savedRecipes.value = Array.isArray(data) ? data : []
      return data
    } catch (e) {
      // Cookbook API가 아직 구현되지 않았을 수 있으므로 조용히 실패
      error.value = e
      savedRecipes.value = []
      return []
    } finally {
      loading.value = false
    }
  }

  // 레시피 저장
  async function saveRecipe(recipeData) {
    loading.value = true
    error.value = null
    try {
      const { data } = await cookbookApi.save(recipeData)
      await fetchSavedRecipes() // 목록 갱신
      return data
    } catch (e) {
      // Cookbook API가 아직 구현되지 않았을 수 있음
      error.value = e
      throw e
    } finally {
      loading.value = false
    }
  }

  // 레시피 삭제
  async function deleteRecipe(id) {
    loading.value = true
    error.value = null
    try {
      await cookbookApi.delete(id)
      savedRecipes.value = savedRecipes.value.filter((r) => r.id !== id)
    } catch (e) {
      error.value = e
      throw e
    } finally {
      loading.value = false
    }
  }

  // 라벨 목록 조회
  async function fetchLabels() {
    try {
      const { data } = await cookbookApi.getLabels()
      labels.value = Array.isArray(data) ? data : []
      return data
    } catch (e) {
      error.value = e
      throw e
    }
  }

  // 레시피가 이미 저장되어 있는지 확인
  function isSaved(recipeId) {
    return savedRecipes.value.some((r) => r.recipeId === recipeId || r.recipeSnapshot?.title)
  }

  return {
    savedRecipes,
    labels,
    loading,
    error,
    fetchSavedRecipes,
    saveRecipe,
    deleteRecipe,
    fetchLabels,
    isSaved,
  }
})
