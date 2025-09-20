import axios from 'axios'

const state = {
  token: localStorage.getItem('token') || '',
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  isAuthenticated: !!localStorage.getItem('token')
}

const mutations = {
  LOGIN_SUCCESS(state, { token, user }) {
    state.token = token
    state.user = user
    state.isAuthenticated = true
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
  },
  LOGOUT(state) {
    state.token = ''
    state.user = null
    state.isAuthenticated = false
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  },
  UPDATE_USER(state, user) {
    state.user = { ...state.user, ...user }
    localStorage.setItem('user', JSON.stringify(state.user))
  }
}

const actions = {
  async login({ commit }, credentials) {
    try {
      const response = await axios.post('/login', credentials)
      const { token, user } = response.data
      commit('LOGIN_SUCCESS', { token, user })
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '登录失败' 
      }
    }
  },

  async register({ commit }, userData) {
    try {
      const response = await axios.post('/register', userData)
      const { token, user } = response.data
      commit('LOGIN_SUCCESS', { token, user })
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '注册失败' 
      }
    }
  },

  logout({ commit }) {
    commit('LOGOUT')
  },

  async updateProfile({ commit }, profileData) {
    try {
      const response = await axios.put('/profile', profileData)
      commit('UPDATE_USER', response.data.user)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '更新失败' 
      }
    }
  }
}

const getters = {
  isAuthenticated: state => state.isAuthenticated,
  currentUser: state => state.user,
  userRole: state => state.user?.role || '',
  isStudent: state => state.user?.role === 'student',
  isExpert: state => state.user?.role === 'expert',
  isAdmin: state => state.user?.role === 'admin'
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}
