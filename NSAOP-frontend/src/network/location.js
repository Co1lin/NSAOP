import service from "@/network/request";
import {firstTry} from "@/network/retry";

export function getLocationByUser(data) {
  return firstTry(data => service({
    url: '/location',
    method: 'get',
    params: data,
  }), data).then(res => res).catch(err => {
    // console.log(err.response)
    return Promise.reject(err)
  })
}

export function createLocation(data) {
  return firstTry(data => service({
    url: '/location',
    method: 'post',
    data
  }), data).then(res => res).catch(err => {
    // console.log(err.response)
    return Promise.reject(err)
  })
}

export function deleteLocationById(data) {
  return firstTry(() => service({
    url: '/location/' + data.id,
    method: 'delete',
  })).then(res => res).catch(err => {
    if(err.response.data.msg === "occupied"){
      return Promise.reject("该地址已有订单，不可删除")
    } else {
      return Promise.reject("删除地址失败，请重试")
    }
  })
}
