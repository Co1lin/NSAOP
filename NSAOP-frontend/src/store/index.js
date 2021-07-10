import Vue from 'vue'
import Vuex from 'vuex'

import getters from "./getters";
import userModule from './modules/user'
import settings from './modules/layout'
import orderQueryModule from './modules/orderQuery'
import locationQueryModule from "@/store/modules/locationQuery";

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    user: userModule,
    settings: settings,
    orderQuery: orderQueryModule,
    locationQuery: locationQueryModule
  },
  state: {
    price: 1
  },
  getters
})

export default store
