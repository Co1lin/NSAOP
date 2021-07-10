import {createLocalVue} from '@vue/test-utils'
import ElementUI from 'element-ui';
import VueRouter from 'vue-router'

const localVue = createLocalVue()
localVue.use(ElementUI)
localVue.use(VueRouter)

export default localVue
