import api from './index.js'

export const cookbookApi = {
  // GET /api/v1/cookbook
  getList(params = {}) {
    return api.get('/cookbook', { params })
  },

  // POST /api/v1/cookbook
  save(data) {
    return api.post('/cookbook', data)
  },

  // GET /api/v1/cookbook/:id
  getById(id) {
    return api.get(`/cookbook/${id}`)
  },

  // PATCH /api/v1/cookbook/:id
  update(id, data) {
    return api.patch(`/cookbook/${id}`, data)
  },

  // DELETE /api/v1/cookbook/:id
  delete(id) {
    return api.delete(`/cookbook/${id}`)
  },

  // GET /api/v1/cookbook/labels
  getLabels() {
    return api.get('/cookbook/labels')
  },
}
