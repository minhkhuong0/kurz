import api from './api.ts'

export default {
  createCode(data: String) {
    return api().post('/api/url', data)
  }
}
