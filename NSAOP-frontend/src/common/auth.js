import {setCookie, getCookie, delCookie} from "@/common/cookie"
const tokenName = "token"

export function setToken(token) {
  setCookie(tokenName, token)
}

export function getToken() {
  return getCookie(tokenName)
}

export function revokeToken() {
  delCookie(tokenName)
}
