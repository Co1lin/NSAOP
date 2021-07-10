import {includeChinese} from "@/common/regex";

export function lengthValidator(min = 1, max = 100) {
  return function (rule, value, callback) {
    if (value.length < min) {
      if (min === 1) {
        callback(new Error("该字段不能为空"))
      } else {
        callback(new Error("长度至少为" + min + "个字符"))
      }
    } else if (value.length > max) {
      callback(new Error("不能超过" + max + "个字符"))
    } else {
      callback()
    }
  }
}

export function shortMessageValidator(rule, value, callback) {
  if(value === "") {
    callback(new Error("该字段不能为空"))
  } else if (value.length > 10) {
    callback(new Error("不能超过10个字符"))
  } else {
    callback()
  }
}

export function validateTel(rule, value, callback) {
  const patten = /^0?(13[0-9]|14[5-9]|15[012356789]|166|17[0-8]|18[0-9]|19[8-9])[0-9]{8}$/
  if (!(patten.test(value))) {
    callback(new Error("请输入有效的电话号码"))
  } else {
    callback()
  }
}

export function validateLocation(rule, value, callback) {
  if (value.length > 10) {
    callback(new Error("不能超过10个字符"))
  } else {
    callback()
  }
}

export function validateEmail(rule, value, callback) {
  const patten = /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/

  if (value === "") {
    callback(new Error("邮箱名不能为空"))
  } else if (value.length > 40) {
    callback(new Error("不能超过40个字符"))
  } else if(!(patten.test(value))) {
    callback(new Error("请输入有效的邮箱"))
  } else if (includeChinese(value)) {
    callback(new Error("不能包含汉字"))
  } else {
    callback()
  }
}
