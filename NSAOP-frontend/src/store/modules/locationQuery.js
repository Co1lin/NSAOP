import {SET_LOCATION_QUERY} from "@/common/store"

const state = () => ({
  target: '',
})

const mutations = {
  SET_LOCATION_QUERY(state, target) {
    state.target = target
  }
}

const actions = {
  setLocationTarget({commit}, {data}) {
    return new Promise((resolve) => {
      commit(SET_LOCATION_QUERY, data)
      resolve()
    })
  }
}

export default {state, mutations, actions}
