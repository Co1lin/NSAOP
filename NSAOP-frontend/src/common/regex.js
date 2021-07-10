
export function consistsOfOnlyEnglish(str) {
  const pattern = /^[a-z|A-Z]{1,64}$/
  return pattern.test(str)
}

export function ssidNameValid(str) {
  const pattern = /^[^"][^&=?$%+]{1,32}$/
  return pattern.test(str)
}

export function includeChinese(str) {
  const pattern = new RegExp('[\\u4E00-\\u9FFF]', 'g')
  return pattern.test(str)
}
