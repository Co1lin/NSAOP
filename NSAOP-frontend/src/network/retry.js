import store from "@/store";
import {setToken} from "@/common/auth"
import {refreshToken} from "@/network/user";
import {CONNECTION_LOST} from "@/common/utils";

export function firstTry(callback, data) {
  return callback(data).catch(err => refreshTry(callback, data, err))
}

function refreshTry(callback, data, err) {
  if (err === CONNECTION_LOST) {
    return Promise.reject(err)
  }
  if (err.response.status !== 401) {
    return Promise.reject(err)
  }
  return refreshToken().then(async res => {
    const token = res.data.token
    await store.commit("SET_TOKEN", token)
    await setToken(token)
    return secondTry(callback, data)
  })
}

function secondTry(callback, data) {
  return callback(data)
}
