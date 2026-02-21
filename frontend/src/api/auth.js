import api from './index.js'

export const authApi = {
  signup({ email, password, nickname }) {
    return api.post('/auth/signup', { email, password, nickname })
  },

  login({ email, password }) {
    return api.post('/auth/login', { email, password })
  },

  refresh(refreshToken) {
    return api.post('/auth/refresh', { refresh_token: refreshToken })
  },

  logout(refreshToken) {
    return api.post('/auth/logout', { refresh_token: refreshToken })
  },

  getMe() {
    return api.get('/users/me')
  },
}
