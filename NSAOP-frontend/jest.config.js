let testURL = 'http://localhost:8081/'
if (process.env.NODE_ENV === 'production')
  testURL = 'https://nsaop.enjoycolin.top/'
module.exports = {
  preset: '@vue/cli-plugin-unit-jest',
  coverageDirectory: "./tests/coverage",
  testURL,
}
