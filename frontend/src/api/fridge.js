import api from './index.js'

export const fridgeApi = {
  getList(params = {}) {
    return api.get('/fridge', { params })
  },

  getSummary() {
    return api.get('/fridge/summary')
  },

  getById(id) {
    return api.get(`/fridge/${id}`)
  },

  create(data) {
    return api.post('/fridge', data)
  },

  update(id, data) {
    return api.patch(`/fridge/${id}`, data)
  },

  deleteOne(id) {
    return api.delete(`/fridge/${id}`)
  },

  deleteBulk(ids) {
    return api.delete('/fridge', { data: { ids } })
  },
}
