import React from 'react'
import { useDidHide, useDidShow } from '@tarojs/taro'
import './app.scss'

function App(props: React.PropsWithChildren) {
  useDidShow(() => {
    console.info('[XCX][App] 应用进入前台')
  })

  useDidHide(() => {
    console.info('[XCX][App] 应用进入后台')
  })

  return props.children
}

export default App
