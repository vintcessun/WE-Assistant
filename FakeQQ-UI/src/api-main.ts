import { createApp } from 'vue'
import APIRenderer from './APIRendererSSR.vue'
import { renderToString } from 'vue/server-renderer'

// 定义数据类型
interface RequestData {
  id: {
    other: string
    self: string
  }
  messages: Array<{
    type: 'self' | 'other'
    message: Array<{
      type: 'text' | 'image' | 'reply'
      data: {
        text?: string
        url?: string
        head?: string
        message?: string | Array<{ type: string; data: { text?: string } }>
      }
    }>
  }>
}

// 导出渲染函数
export async function renderMessageData(renderData: RequestData): Promise<string> {
  const app = createApp(APIRenderer, {
    data: renderData
  })
  const renderedContent = await renderToString(app)
  return renderedContent
}

// 从URL参数获取数据（保持原有功能）
const urlParams = new URLSearchParams(window.location.search)
const dataParam = urlParams.get('data')

let renderData = null

if (dataParam) {
  try {
    renderData = JSON.parse(decodeURIComponent(dataParam))

    // 调用渲染函数并输出结果
    renderMessageData(renderData).then(html => {
      document.getElementById('app')!.innerHTML = html
    })
  } catch (error) {
    console.error('解析数据失败:', error)
    // 显示错误信息
    document.getElementById('app')!.innerHTML = `
      <div style="padding: 20px; color: red;">
        <h3>数据解析失败</h3>
        <p>${error instanceof Error ? error.message : 'Unknown error'}</p>
      </div>
    `
    throw error
  }
}
