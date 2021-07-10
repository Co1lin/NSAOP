import service from "@/network/request";
import {firstTry} from "@/network/retry";
import {NCE_ERROR} from "@/common/utils";

export function getOrderInfoById(data) {
  return firstTry(data => service({
    url: '/service/' + data.id,
    method: 'get',
  }), data).then(res => {
    return res.data
  }).catch(err => {
    // console.log(err.response);
    return Promise.reject(err)
  })
}

export function getOrderByUser(data) {
  return firstTry(data => service({
    url: '/service',
    method: 'get',
    params: data
  }), data).then(res => res).catch(err => {
    // console.log(err.response);
    return Promise.reject(err)
  })
}

export function changeOrder(data) {
  return firstTry(data => service({
    url: '/service/' + data.id,
    method: 'put',
    data,
  }), data).then(res => res).catch(err => {
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    } else {
      return Promise.reject(err.response.data.msg)
    }
  })
}


export function submitOrder(data) {
  return firstTry(data => service({
    url: '/service',
    method: 'post',
    data
  }), data).catch(err => {
    // console.log(err.response)
    if (err.response.status === 400) {
      return Promise.reject("提交失败：格式错误")
    }
    return Promise.reject(err)
  })
}
