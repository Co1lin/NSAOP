/**
 * @description 导出默认通用配置
 */
const setting = {
  //标题 （包括初次加载雪花屏的标题 页面的标题 浏览器的标题）
  title: 'NSAOP',
  //简写
  abbreviation: 'NSAOP',
  //开发环境端口号
  devPort: '80',
  //版本号
  version: process.env.VUE_APP_VERSION,
  //是否显示页面底部自定义版权信息
  footerCopyright: true,
  //是否显示顶部进度条
  progressBar: true,
  //不经过token校验的路由
  routesWhiteList: ['/login', '/register', '/404', '/401'],
  //加载时显示文字
  loadingText: '正在加载中...',
  //intelligence和all两种方式，前者后端权限只控制permissions不控制view文件的import（前后端配合，减轻后端工作量），all方式完全交给后端前端只负责加载
  authentication: 'intelligence',
  //需要自动注入并加载的模块
  providePlugin: { maptalks: 'maptalks', 'window.maptalks': 'maptalks' },
}
module.exports = setting
