import axios from 'axios'

const state = {
  profile: null,
  files: [],
  following: [],
  followers: []
}

const mutations = {
  SET_PROFILE(state, profile) {
    state.profile = profile
  },
  SET_FILES(state, files) {
    state.files = files
  },
  ADD_FILE(state, file) {
    state.files.unshift(file)
  },
  REMOVE_FILE(state, fileId) {
    state.files = state.files.filter(file => file.id !== fileId)
  },
  SET_FOLLOWING(state, following) {
    state.following = following
  },
  SET_FOLLOWERS(state, followers) {
    state.followers = followers
  },
  ADD_FOLLOWING(state, follow) {
    state.following.push(follow)
  },
  REMOVE_FOLLOWING(state, userId) {
    state.following = state.following.filter(follow => follow.followed_id !== userId)
  }
}

const actions = {
  async fetchProfile({ commit }) {
    try {
      const response = await axios.get('/profile')
      commit('SET_PROFILE', response.data.user)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '获取用户信息失败' 
      }
    }
  },

  async fetchUserFiles({ commit }, userId = 'me') {
    try {
      const response = await axios.get(`/files/user/${userId}`)
      commit('SET_FILES', response.data.files)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '获取文件列表失败' 
      }
    }
  },

  async uploadFile({ commit }, formData) {
    try {
      const response = await axios.post('/files', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      commit('ADD_FILE', response.data.file)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '文件上传失败' 
      }
    }
  },

  async deleteFile({ commit }, fileId) {
    try {
      await axios.delete(`/files/${fileId}`)
      commit('REMOVE_FILE', fileId)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '删除文件失败' 
      }
    }
  },

  async fetchFollowing({ commit }) {
    try {
      const response = await axios.get('/following')
      commit('SET_FOLLOWING', response.data.following)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '获取关注列表失败' 
      }
    }
  },

  async fetchFollowers({ commit }) {
    try {
      const response = await axios.get('/followers')
      commit('SET_FOLLOWERS', response.data.followers)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '获取粉丝列表失败' 
      }
    }
  },

  async followUser({ commit }, userId) {
    try {
      const response = await axios.post(`/follow/${userId}`)
      commit('ADD_FOLLOWING', response.data.follow)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '关注失败' 
      }
    }
  },

  async unfollowUser({ commit }, userId) {
    try {
      await axios.delete(`/follow/${userId}`)
      commit('REMOVE_FOLLOWING', userId)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '取消关注失败' 
      }
    }
  }
}

const getters = {
  userProfile: state => state.profile,
  userFiles: state => state.files,
  following: state => state.following,
  followers: state => state.followers,
  isFollowing: state => userId => {
    return state.following.some(follow => follow.followed_id === userId)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}
