import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { fridgeApi } from '@/api/fridge.js'

export const useIngredientStore = defineStore('ingredient', () => {
  const ingredients = ref([])
  const summary = ref({ total: 0, urgent: 0, soon: 0, normal: 0, no_expiry: 0 })
  const loading = ref(false)
  const error = ref(null)

  const expiringSoon = computed(() =>
    ingredients.value.filter(i => i.expiry_status === 'URGENT' || i.expiry_status === 'SOON')
  )

  const urgentCount = computed(() => summary.value.urgent ?? 0)

  async function fetchIngredients(params = {}) {
    loading.value = true
    error.value = null
    try {
      const { data } = await fridgeApi.getList(params)
      ingredients.value = data.items ?? []
    } catch (e) {
      error.value = e
    } finally {
      loading.value = false
    }
  }

  async function fetchSummary() {
    try {
      const { data } = await fridgeApi.getSummary()
      summary.value = data
    } catch (e) {
      // summary는 조용히 실패
    }
  }

  async function addIngredient(form) {
    const { data } = await fridgeApi.create(form)
    await fetchIngredients()
    await fetchSummary()
    return data
  }

  async function removeIngredient(id) {
    await fridgeApi.deleteOne(id)
    ingredients.value = ingredients.value.filter(i => i.id !== id)
    await fetchSummary()
  }

  return {
    ingredients,
    summary,
    loading,
    error,
    expiringSoon,
    urgentCount,
    fetchIngredients,
    fetchSummary,
    addIngredient,
    removeIngredient,
  }
})
