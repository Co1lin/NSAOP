import service from "@/network/request"
import {firstTry} from "@/network/retry"

export function login(data) {
  return service({
    url: '/user/login',
    method: 'post',
    data
  }).then(res => {
    return res.data.token
  }).catch(err => {
    const msg = err.response.data.msg
    let alertMsg = err.response.data.msg
    if (err.response.status === 400) {
      alertMsg = "400 format error" // 这种情形不应当出现，因为前端对格式应当进行检查。
    } else if (msg === "user not found") {
      alertMsg = "用户名不存在"
    } else if (msg === "wrong password") {
      alertMsg = "密码错误"
    }
    return Promise.reject(alertMsg)
  })
}

export function register(data) {
  return service({
    url: '/user/signup',
    method: 'post',
    data
  }).then(res => res).catch(err => {
    const msg = err.response.data.msg
    let alertMsg = ""
    if (msg === "username already exist") {
      alertMsg = "该用户名已被占用"
    } else if (err.response.status === 400) {
      if(data.role === 'customer')
        alertMsg = "400 format error" // 这种情形不应当出现，因为前端对格式应当进行检查。
      else
        alertMsg = "邀请码错误"
    } else {
      // console.log(err.response);
    }
    return Promise.reject(alertMsg)
  })
}

export function usernameCheck(data) {
  // console.log(data);
  return service({
    url: '/user/check/username',
    method: 'post',
    data
  }).then(() => Promise.resolve()).catch(err => {
    return Promise.reject(err)
  })
}

export function checkPassword(data) {
  return service({
    url: 'user/check/password',
    method: 'post',
    data
  }).then(() => Promise.resolve()).catch(() => {
    return Promise.reject()
  })
}

export function logout(data) {
  return service({
    url: '/user/logout',
    method: 'post',
    data
  }).then(() => Promise.resolve()).catch(() => {

  })
}

export function authTest() {
  return firstTry(() => service({
    url: '/user/detail',
    method: 'get',
  })).then(res => {
    return res.data
  }).catch(err => {
    // console.log(err.response);
    return Promise.reject(err)
  })
}

export function refreshToken() {
  return service({
    url: '/user/refresh',
    method: 'post',
  })
}

export function changeInfo(data) {
  return firstTry(data => service({
    url: '/user/detail',
    method: 'put',
    data,
  }), data).catch(err => {
    // console.log(err.response)
    if (err.response.status === 400) {
      return Promise.reject("修改信息格式错误")
    }
  })
}

export function sendEmail() {
  return firstTry(() => service({
    url: '/user/email',
    method: 'get',
  })).then(res => {
    return res.data
  }).catch(err => {
    return Promise.reject(err)
  })
}

export function requestResetPassword(data) {
  return firstTry(() => service({
    url: '/user/resetpasswd',
    method: 'get',
    params: data,
  })).then(res => res).catch(err => Promise.reject(err))
}

export function resetPassword(data) {
  return firstTry(data => service({
    url: '/user/resetpasswd',
    method: 'post',
    data,
  }), data).then(res => res).catch(err => Promise.reject(err))
}
