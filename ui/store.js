// eslint-disable-next-line camelcase
import { EosProjects_GetProjects } from './client/eosprojects'
import axios from 'axios'

const state = {
  config: null,
  projects: ''
}

const getters = {
  config: state => state.config,
  projects: state => state.projects,
  getServerForJsClient: (state, getters, rootState, rootGetters) => rootGetters.configuration.server.replace(/\/$/, '')
}

const actions = {
  // Used by ocis-web.
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  getProjects ({ commit, dispatch, getters, rootGetters }) {
    injectAuthToken(rootGetters)
    EosProjects_GetProjects({
      $domain: getters.getServerForJsClient,
      body: { }
    })
      .then(response => {
        console.log(response)

        if (response.status === 200 || response.status === 201) {
          commit('SET_PROJECTS', response.data.projects)
        } else {
          dispatch('showMessage', {
            title: 'Response failed',
            desc: response.statusText,
            status: 'danger'
          }, { root: true })
        }
      })
      .catch(error => {
        console.error(error)

        dispatch('showMessage', {
          title: 'Saving your name failed',
          desc: error.message,
          status: 'danger'
        }, { root: true })
      })
  }
}

const mutations = {
  SET_PROJECTS (state, payload) {
    state.projects = payload
  },

  LOAD_CONFIG (state, config) {
    state.config = config
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}

function injectAuthToken (rootGetters) {
  axios.interceptors.request.use(config => {
    if (typeof config.headers.Authorization === 'undefined') {
      const token = rootGetters.user.token
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    return config
  })
}
