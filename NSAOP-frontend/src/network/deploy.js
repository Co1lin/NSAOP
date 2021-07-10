import service from './request'
import {firstTry} from './retry'
import {filterNCEError, NCE_ERROR} from "@/common/utils";

export function createDevices(data) {
  return firstTry(data => service({
    url: '/service/' + data.service_id + '/device',
    method: 'post',
    data,
  }), data).then(res => res).catch(err => {
    // console.log(err.response);
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    } else {
      return Promise.reject(err)
    }
  })
}

export function createSSID(data) {
  return firstTry(data => service({
    url: '/service/' + data.service_id + '/ssid',
    method: 'post',
    data,
  }), data).then(res => res).catch(err => {
    // console.log(err.response);
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    } else {
      return Promise.reject(err)
    }
  })
}

export async function createSSIDs(service_id, ssids) {
  let errors = []
  for (let ssid of ssids) {
    await createSSID({service_id, ssid,}).then(() => {}).catch(err => errors.push(err))
  }
  if (errors.length !== 0) {
    return filterNCEError(errors)
  } else {
    return Promise.resolve("success")
  }
}

export function getDevices(data) {
  return firstTry(() => service({
    url: '/service/' + data + '/device',
    method: 'get',
  })).then(res => {
    return res.data
  }).catch(err => {
    // console.log(err.response);
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    } else {
      return Promise.reject(err)
    }
  })
}

export function getSSIDs(data) {
  return firstTry(() => service({
    url: '/service/' + data + '/ssid',
    method: 'get',
  })).then(res => {
    return res.data
  }).catch(err => {
    // console.log(err.response);
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    } else {
      return Promise.reject(err)
    }
  })
}

export function deleteDevices(data) {
  return firstTry(() => service({
    url: '/service/' + data.service_id + '/device',
    method: 'delete',
    data,
  })).then(res => res).catch(err => {
    // console.log(err.response);
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    } else {
      return Promise.reject(err)
    }
  })
}

export function deleteSSIDs(data) {
  return firstTry(data => service({
    url: '/service/' + data.service_id + '/ssid',
    method: 'delete',
    data,
  }), data).then(res => res).catch(err => {
    // console.log(err.response);
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    } else {
      return Promise.reject(err)
    }
  })
}
