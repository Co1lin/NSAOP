import {load} from 'recaptcha-v3'

const siteKey = '6Lf4R6UaAAAAADivRqmYp60lOU46qm23ciftBsxH'

export async function reCAPTCHA(action) {
  const reCAPTCHA = await load(siteKey, {useRecaptchaNet: true, autoHideBadge: true})
  const token = await reCAPTCHA.execute(action)
  return token
}
