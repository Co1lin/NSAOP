import Vue from 'vue'
import App from './App.vue'
import router from "@/router";
import store from "@/store";
import './plugins'
import 'element-ui/lib/theme-chalk/index.css'
import '@/assets/theme/index.css'

const EventBus = new Vue();

Object.defineProperties(Vue.prototype, {
  $bus: {
    get: function () {
      return EventBus
    }
  }
})

new Vue({
  el: '#app',
  router,
  store,
  render: (h) => h(App),
})

