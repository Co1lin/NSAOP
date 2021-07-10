// For authoring Nightwatch tests, see
// https://nightwatchjs.org/guide

module.exports = {
  'default e2e tests': browser => {
    browser
      .init()
      .waitForElementVisible('.container')
      .assert.elementPresent('.container')
      .assert.containsText('h3', '登录')
      .end()
  },
  'test login button with no input': browser => {
    browser
      .init()
      .waitForElementVisible('.container')
      .assert.elementPresent('.el-button')
      .assert.elementPresent('#login-button')
      .click('#login-button', res => {
        console.log('click', res)
      })
      .assert.elementPresent('.el-form-item__error')
      .assert.containsText('.container', "请输入密码")
      .assert.containsText('.container', "请输入用户名")

  }
}
