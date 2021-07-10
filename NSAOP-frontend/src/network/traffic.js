import {firstTry} from "@/network/retry";
import service from  './request'

export function getTraffic(data) {
  return firstTry(() => service({
    url: '/service/' + data + '/traffic',
    method: 'get',
  })).then(res => res.data).catch(err => {
    // console.log(err.response);
    return Promise.reject(err)
  })
}
