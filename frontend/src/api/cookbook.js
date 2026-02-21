import api from './index.js'

export const cookbookApi = {
  // GET /api/v1/cookbook
  getList(params = {}) {
    return api.get('/cookbook', { params })
      .catch((err) => {
        // Cookbook API가 아직 구현되지 않았을 수 있으므로 404는 조용히 처리
        if (err.response?.status === 404) {
          return { data: [] }
        }
        throw err
      })
  },

  // POST /api/v1/cookbook
  save(data) {
    return api.post('/cookbook', data)
      .catch((err) => {
        // Cookbook API가 아직 구현되지 않았을 수 있음
        if (err.response?.status === 404) {
          throw new Error('요리책 기능이 아직 준비되지 않았어요.')
        }
        throw err
      })
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
      .catch((err) => {
        // Cookbook API가 아직 구현되지 않았을 수 있으므로 404는 조용히 처리
        if (err.response?.status === 404) {
          return { data: [] }
        }
        throw err
      })
  },
}
