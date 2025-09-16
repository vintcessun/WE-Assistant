import type { Plugin } from 'vite'
import fs from 'fs'

interface ApiPluginOptions {
  templatePath?: string
}

export default function apiPlugin(options: ApiPluginOptions = {}): Plugin {
  const { templatePath = 'api-template.html' } = options
  
  return {
    name: 'vite-plugin-api',
    configureServer(server) {
      // 处理 /render 的 POST 请求
      server.middlewares.use('/render', async (req, res) => {
        if (req.method !== 'POST') {
          res.statusCode = 405
          res.end('Method Not Allowed')
          return
        }
        
        try {
          // 读取请求体
          let body = ''
          req.on('data', chunk => {
            body += chunk.toString()
          })
          
          req.on('end', async () => {
            try {
              const data = JSON.parse(body)
              
              // 读取模板文件
              const template = fs.readFileSync(templatePath, 'utf-8')
              
              // 替换模板中的占位符
              const dataJson = JSON.stringify(data).replace(/"/g, '"')
              const html = template
                .replace('{{TITLE}}', 'FakeQQ-UI 消息渲染')
                .replace(/{{DATA}}/g, dataJson)
              
              res.setHeader('Content-Type', 'text/html')
              res.end(html)
            } catch (error) {
              res.statusCode = 400
              res.end(`Bad Request: ${error instanceof Error ? error.message : 'Unknown error'}`)
            }
          })
        } catch (error) {
          res.statusCode = 500
          res.end(`Internal Server Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
        }
      })
    }
  }
}

