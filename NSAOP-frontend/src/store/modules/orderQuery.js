import {SET_ORDER_QUERY} from "@/common/store"

const state = () => ({
  target: '',
})

const mutations = {
  SET_ORDER_QUERY(state, target) {
    state.target = target
  }
}

const actions = {
  setTarget({commit}, {data}) {
    return new Promise((resolve) => {
      commit(SET_ORDER_QUERY, data)
      resolve()
    })
  }
}

export default {state, mutations, actions}
