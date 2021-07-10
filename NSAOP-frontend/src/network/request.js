import Axios from "axios"
import store from "@/store";
import qs from 'qs'
import {getToken} from "@/common/auth";
import {CONNECTION_LOST} from "@/common/utils";

let baseURL = ''

if (process.env.NODE_ENV === 'production') {
  baseURL = 'https://nsaop.enjoycolin.top/api/v2/'
} else if (process.env.NODE_ENV === 'development') {
  // baseURL = 'https://dev-api.nsaop.enjoycolin.top/v2/'
  baseURL = 'https://nsaop.enjoycolin.top/api/v2/'
}

const instance = Axios.create({
  baseURL,
  timeout: 10000,
  paramsSerializer: params => {
    return qs.stringify(params, { indices: false })
  },
})


instance.interceptors.request.use(config => {
  if (store.getters.token !== undefined && store.getters.token !== "") {
    config.headers['Authorization'] = "Bearer " + getToken()
    config.headers['Content-type'] = "application/json"
  }

  return config
},err => {
  // console.log(err);
  return Promise.reject(err)
})

instance.interceptors.response.use(res => {
  return res.data
},err => {
  // console.log(err);
  // console.log(err.request);
  // console.log(err.message);
  if (err.message === 'Network Error' || !navigator.onLine) {
    return Promise.reject(CONNECTION_LOST)
  } else {
    return Promise.reject(err)
  }
})

export default instance
