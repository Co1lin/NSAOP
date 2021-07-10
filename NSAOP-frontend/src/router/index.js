import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const Home = () => import('@/views/home/Home')
const Login = () => import('@/views/login/Login')
const Register = () => import('@/views/login/Register')
const ForgetPassword = () => import('@/views/login/ForgetPassword')
const ResetPassword = () => import('@/views/login/ResetPassword')
const PersonalCenter = () => import('@/views/personalCenter/PersonalCenter')
const OrderList = () => import('@/views/orderList/OrderList')
const LocationList = () => import('@/views/locationList/LocationList')
const Detail = () => import('@/views/order/Detail')
const Deploy = () => import('@/views/order/Deploy')
const Page404 = () => import('@/views/404')

const routes = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/home',
    component: Home,
  },
  {
    path: '/login',
    component: Login,
  },
  {
    path: '/register/:role',
    component: Register,
  },
  {
    path: '/login/forget',
    component: ForgetPassword,
  },
  {
    path: '/reset',
    component: ResetPassword,
  },
  {
    path: '/profile',
    component: PersonalCenter
  },
  {
    path: '/locations',
    component: LocationList
  },
  {
    path: '/orders',
    component: OrderList
  },
  {
    name: 'detail',
    path: '/orders/:id',
    component: Detail,
  },
  {
    name: 'deploy',
    path: '/orders/deploy/:id',
    component: Deploy,
  },
  {
    path: '/404',
    component: Page404,
  },
  {
    path: '*',
    redirect: '/404',
  }
]

const router = new VueRouter({
  mode: 'history',
  routes,
})

export default router

