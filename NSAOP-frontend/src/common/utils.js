
export const TIMEOUT = {
  msg: "timeout",
  status: 504,
}

export const CONNECTION_LOST = {
  msg: "unable to connect",
  status: 503,
}

export const NCE_ERROR = {
  msg: "NCE Error, please contact admin",
  status: 505,
}

export function deepcopy(item) {
  return JSON.parse(JSON.stringify(item))
}

export function debounce(callback, wait = 300) {
  let timeout
  return (...args) => {
    let context = this
    if (timeout) {
      clearTimeout(timeout)
    }
    timeout = setTimeout(() => {
      callback.apply(context, args)
    }, wait)
  }
}

export function filterNCEError(errors) {
  for (let err of errors) {
    if (err.response.status === 500 || err.response.status === 408) {
      return Promise.reject(NCE_ERROR)
    }
  }
  return Promise.reject(errors[0])
}
