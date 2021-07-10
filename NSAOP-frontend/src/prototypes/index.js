import { title } from '@/config'
import store from '@/store'

const device = store.getters['settings/device']

// eslint-disable-next-line no-unused-vars
const install = (Vue, opts = {}) => {
  // 全局标题
  Vue.prototype.$baseTitle = (() => {
    return title
  })()
  //全局表格高度
  Vue.prototype.$baseTableHeight = (formType) => {
    let height = window.innerHeight
    let paddingHeight = 250
    const formHeight = 50

    if (window.innerWidth <= 661) {
      paddingHeight = 300
    }
    if (device === "mobile") {
      paddingHeight = 400
    }

    if ('number' == typeof formType) {
      height = height - paddingHeight - formHeight * formType
    } else {
      height = height - paddingHeight
    }
    return height
  }
}


if (typeof window !== 'undefined' && window.Vue) {
  install(window.Vue)
}

export default install
