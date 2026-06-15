export default defineAppConfig({
  pages: [
    'pages/index/index',
    'pages/settings/index',
    'pages/debug/index',
    'pages/payment/index',
    'pages/webview/index',
  ],
  window: {
    backgroundTextStyle: 'light',
    backgroundColor: '#f6f8fb',
    navigationBarBackgroundColor: '#f6f8fb',
    navigationBarTitleText: '超市 H5 验证壳',
    navigationBarTextStyle: 'black',
  },
})
