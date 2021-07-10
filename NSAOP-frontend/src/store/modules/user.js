import '@/common/store'
import {login, logout} from "@/network/user"
import {setToken, getToken, revokeToken} from "@/common/auth";
import {SET_TOKEN, SET_USERNAME, SET_COMPANY, SET_TEL, SET_EMAIL, SET_ROLE} from "@/common/store";

const state = () => ({
  token: getToken(),
  username: '',
  company: '',
  tel: '',
  email: '',
  role: '',
})

const mutations = {
  SET_TOKEN(state, token) {
    state.token = token
  },
  SET_USERNAME(state, username) {
    state.username = username
  },
  SET_COMPANY(state, company) {
    state.company = company
  },
  SET_TEL(state, tel) {
    state.tel = tel
  },
  SET_EMAIL(state, email) {
    state.email = email
  },
  SET_ROLE(state, role) {
    state.role = role
  },
}

const actions = {
  login({commit}, data) {
    return new Promise((resolve, reject) => {
      login(data).then(token => {
        setToken(token)
        commit(SET_TOKEN, token)
        resolve()
      }).catch(err => {
        reject(err)
      })
    })
  },
  logout({ commit }) {
    return new Promise((resolve, reject) => {
      logout().then(() => {
        revokeToken()
        commit(SET_TOKEN, '')
        commit(SET_USERNAME, '')
        commit(SET_COMPANY, '')
        commit(SET_TEL, '')
        commit(SET_EMAIL, '')
        commit(SET_ROLE, '')
        resolve()
      }).catch(() => reject())
    })
  }
}

export default {state, mutations, actions}
