import md5 from 'js-md5'

const salt1 = "NSAOP"
const salt2 = "litangdingzhen"

export function encrypt(username, password) {
  return md5(md5(md5(password) + username + salt1) + md5(password) + salt2)
}
