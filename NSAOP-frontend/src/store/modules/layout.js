import defaultSettings from '@/theme/config'

const { tabsBar, logo, layout, header, themeBar } = defaultSettings
const state = () => ({
  tabsBar: tabsBar,
  logo,
  collapse: true,
  layout: layout,
  header: header,
  device: 'desktop',
  themeBar,
})
const getters = {
  collapse: (state) => state.collapse,
  device: (state) => state.device,
  header: (state) => state.header,
  layout: (state) => state.layout,
  logo: (state) => state.logo,
  tabsBar: (state) => state.tabsBar,
  themeBar: (state) => state.themeBar,
}
const mutations = {
  changeLayout: (state, layout) => {
    if (layout) state.layout = layout
  },
  changeHeader: (state, header) => {
    if (header) state.header = header
  },
  changeTabsBar: (state, tabsBar) => {
    if (tabsBar) state.tabsBar = tabsBar
  },
  changeCollapse: (state) => {
    state.collapse = !state.collapse
  },
  foldSideBar: (state) => {
    state.collapse = true
  },
  openSideBar: (state) => {
    state.collapse = false
  },
  toggleDevice: (state, device) => {
    state.device = device
  },
}
const actions = {
  changeLayout({ commit }, layout) {
    commit('changeLayout', layout)
  },
  changeHeader({ commit }, header) {
    commit('changeHeader', header)
  },
  changeTabsBar({ commit }, tabsBar) {
    commit('changeTabsBar', tabsBar)
  },
  changeCollapse({ commit }) {
    commit('changeCollapse')
  },
  foldSideBar({ commit }) {
    commit('foldSideBar')
  },
  openSideBar({ commit }) {
    commit('openSideBar')
  },
  toggleDevice({ commit }, device) {
    commit('toggleDevice', device)
  },
}
export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions }
