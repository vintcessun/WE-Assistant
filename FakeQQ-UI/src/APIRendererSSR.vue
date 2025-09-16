<script setup lang="ts">
import { computed } from 'vue'
import './assets/reset.css'
import './package/styles/light.scss'
import './package/styles/dark.scss'

import { QTagColors, QMain, QReply, QText, QImage } from './package'

// 定义props类型
interface MessageData {
  text?: string
  url?: string
  head?: string
  message?: MessageUnit[]
}

interface MessageUnit {
  type: 'text' | 'image' | 'reply'
  data: MessageData
}

interface Message {
  type: 'self' | 'other'
  message: MessageUnit[]
}

interface RequestData {
  id: {
    other: string
    self: string
  }
  messages: Message[]
}

const props = defineProps<{
  data: RequestData
}>()

// 定义组件配置类型
interface ComponentConfig {
  component: unknown
  props: Record<string, unknown>
  slots?: Record<string, string>
}

// 生成组件配置
const messageComponents = computed(() => {
  if (!props.data) return []

  return props.data.messages
    .map(msg => {
      const components: ComponentConfig[] = []

      // 根据消息类型生成不同的avatar URL
      const avatarQQ = msg.type === 'self' ? props.data.id.self : props.data.id.other
      const avatarUrl = `http://q1.qlogo.cn/g?b=qq&nk=${avatarQQ}&s=640`

      // 检查第一个消息单元是否为回复类型
      const firstUnit = msg.message[0]
      let contentHtml = ''

      if (firstUnit && firstUnit.type === 'reply') {
        for (let i = 1; i < msg.message.length; i++) {
          const unit = msg.message[i]

          if (unit.type === 'text') {
            contentHtml += unit.data.text || ''
          } else if (unit.type === 'image') {
            contentHtml += `<img src="${unit.data.url || ''}" alt="图片" style="max-width: 200px; margin: 5px 0;">`
          }
        }
        const replyContent =
          typeof firstUnit.data.message === 'string'
            ? firstUnit.data.message
            : firstUnit.data.message?.[0]?.data.text || ''
        components.push({
          component: QReply,
          props: {
            name: '',
            avatar: avatarUrl,
            tagColor: msg.type === 'self' ? 'purple' : 'sage_green',
            target: firstUnit.data.head || '',
            replyText: replyContent,
            self: msg.type === 'self'
          },
          slots: {
            default: contentHtml
          }
        })
      } else {
        // 处理剩余的消息单元
        // 如果只有一个图片消息，使用QImage组件
        if (msg.message.length === 1 && msg.message[0].type === 'image') {
          const unit = msg.message[0]
          components.push({
            component: QImage,
            props: {
              name: '',
              avatar: avatarUrl,
              tagColor: msg.type === 'self' ? QTagColors.orange : QTagColors.sage_green,
              src: unit.data.url || '',
              self: msg.type === 'self'
            }
          })
        }
        // 如果有多个消息单元，使用QText组件组合内容
        else if (msg.message.length > 0) {
          // 构建组合内容HTML
          for (let i = 0; i < msg.message.length; i++) {
            const unit = msg.message[i]

            if (unit.type === 'text') {
              contentHtml += unit.data.text || ''
            } else if (unit.type === 'image') {
              contentHtml += `<img src="${unit.data.url || ''}" alt="图片" style="max-width: 200px; margin: 5px 0;">`
            }
          }

          components.push({
            component: QText,
            props: {
              name: '',
              avatar: avatarUrl,
              tagColor: msg.type === 'self' ? 'purple' : 'sage_green',
              self: msg.type === 'self'
            },
            slots: {
              default: contentHtml
            }
          })
        }
      }

      return components
    })
    .flat()
})
</script>

<template>
  <section class="qq-api-renderer">
    <q-main>
      <component
        :is="comp.component"
        v-for="(comp, index) in messageComponents"
        :key="index"
        v-bind="comp.props"
      >
        <template v-if="comp.slots">
          <span
            v-for="(slotContent, slotName) in comp.slots"
            :key="slotName"
            v-html="slotContent"
          ></span>
        </template>
      </component>
    </q-main>
  </section>
</template>

<style lang="scss" scoped>
.qq-api-renderer {
  width: 100%;
  background-color: var(--qq-background-03);
  border-radius: 10px;
  transition: background-color 0.2s;
  padding: 20px;
}
</style>
