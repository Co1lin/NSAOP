import Cookies from 'js-cookie'

export function setCookie(name, value) {
  Cookies.set(name, value)
}

export function getCookie(name) {
  return Cookies.get(name)
}

export function delCookie(name) {
  Cookies.remove(name)
}
