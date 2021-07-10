import { mount } from '@vue/test-utils'
import store from "@/store";
import router from "@/router";
import Login from "@/views/login/Login";
import localVue from "./index";

describe('tests for login view', () => {
  let wrapper = mount(Login,{
    localVue,
    router,
    store,
  })

  it('is correct component', () => {
    expect(wrapper.is(Login)).toBe(true)
    expect(wrapper.find('.el-form').exists()).toBe(true)
  })

  it('background img attributes correct', () => {
    expect(wrapper.find('a').exists()).toBe(true)
    expect(wrapper.find('a').attributes('href'))
  })

  it('form render correct', () => {
    let items = wrapper.findAll('.el-form-item')
    expect(items.length).toEqual(3)
  })
})

