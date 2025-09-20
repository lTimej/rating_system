import axios from 'axios'

const state = {
  ratings: [],
  publicRatings: [],
  userRatings: {
    received: [],
    given: []
  },
  commentRatings: {
    received: [],
    given: []
  }
}

const mutations = {
  SET_RATINGS(state, ratings) {
    state.ratings = ratings
  },
  SET_PUBLIC_RATINGS(state, ratings) {
    state.publicRatings = ratings
  },
  SET_USER_RATINGS(state, { received, given }) {
    state.userRatings.received = received || []
    state.userRatings.given = given || []
  },
  SET_COMMENT_RATINGS(state, { received, given }) {
    state.commentRatings.received = received || []
    state.commentRatings.given = given || []
  },
  ADD_RATING(state, rating) {
    state.ratings.unshift(rating)
    state.userRatings.given.unshift(rating)
  },
  ADD_FEEDBACK(state, { ratingId, feedback }) {
    const rating = state.ratings.find(r => r.id === ratingId)
    if (rating) {
      if (!rating.feedbacks) rating.feedbacks = []
      rating.feedbacks.push(feedback)
    }
  }
}

const actions = {
  async fetchPublicRatings({ commit }) {
    try {
      const response = await axios.get('/ratings/public')
      commit('SET_PUBLIC_RATINGS', response.data.ratings)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '获取公开评价失败' 
      }
    }
  },

  async fetchUserRatings({ commit }, userId = 'me') {
    try {
      const response = await axios.get(`/ratings/user/${userId}`)
      commit('SET_USER_RATINGS', {
        received: response.data.received_ratings,
        given: response.data.given_ratings
      })
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '获取用户评价失败' 
      }
    }
  },

  async createRating({ commit }, ratingData) {
    try {
      const response = await axios.post('/ratings', ratingData)
      commit('ADD_RATING', response.data.rating)
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '创建评价失败' 
      }
    }
  },

  async createFeedback({ commit }, { ratingId, isHelpful }) {
    try {
      const response = await axios.post(`/ratings/${ratingId}/feedback`, {
        is_helpful: isHelpful
      })
      commit('ADD_FEEDBACK', {
        ratingId,
        feedback: response.data.feedback
      })
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '反馈失败' 
      }
    }
  },

  async fetchUserCommentRatings({ commit }, userId = 'me') {
    try {
      const response = await axios.get(`/users/${userId}/comment-ratings`)
      commit('SET_COMMENT_RATINGS', {
        received: response.data.received_ratings,
        given: response.data.given_ratings
      })
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '获取评论评价失败' 
      }
    }
  },

  async createCommentRating({ dispatch }, ratingData) {
    try {
      const response = await axios.post('/comment-ratings', ratingData)
      // 重新获取用户的评论评价数据
      await dispatch('fetchUserCommentRatings')
      return { success: true, rating: response.data.rating }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.error || '创建评论评价失败' 
      }
    }
  }
}

const getters = {
  allRatings: state => state.ratings,
  publicRatings: state => state.publicRatings,
  receivedRatings: state => state.userRatings.received,
  givenRatings: state => state.userRatings.given,
  receivedCommentRatings: state => state.commentRatings.received,
  givenCommentRatings: state => state.commentRatings.given,
  getRatingById: state => id => {
    return state.ratings.find(rating => rating.id === id)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}
