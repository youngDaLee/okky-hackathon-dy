import api from './index.js'

export const recommendationApi = {
  // GET /api/v1/recommendations
  getRecommendations(params = {}) {
    return api.get('/recommendations', { params })
  },

  // GET /api/v1/recommendations/today
  getTodayRecommendations() {
    return api.get('/recommendations/today')
  },

  // GET /api/v1/recipes
  searchRecipes(params = {}) {
    return api.get('/recipes', { params })
  },

  // GET /api/v1/recipes/:id
  getRecipeById(id) {
    return api.get(`/recipes/${id}`)
  },
}
