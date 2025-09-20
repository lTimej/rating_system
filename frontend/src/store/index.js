import Vue from 'vue'
import Vuex from 'vuex'
import auth from './modules/auth'
import user from './modules/user'
import rating from './modules/rating'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    auth,
    user,
    rating
  }
})
