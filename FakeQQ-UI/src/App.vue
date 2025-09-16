<script setup lang="ts">
import './assets/reset.css'
import './package/styles/light.scss'
import './package/styles/dark.scss'

import { QTagColors, QMain, QReply, QText, QImage } from './package'
</script>

<template>
  <section class="qq">
    <q-main>
      <q-reply
        name=""
        avatar="http://q1.qlogo.cn/g?b=qq&nk=2595201156&s=640"
        tag-color="sage_green"
        target="方糖"
        reply-text="111"
      >
        no
      </q-reply>
      <q-reply
        name=""
        avatar="http://q1.qlogo.cn/g?b=qq&nk=2595201156&s=640"
        tag-color="sage_green"
        target="方糖"
        reply-image-url="/85467403_p0.jpg"
      >
        yes
      </q-reply>
      <q-image
        name=""
        avatar="http://q1.qlogo.cn/g?b=qq&nk=2595201156&s=640"
        :tag-color="QTagColors.orange"
        src="/85467403_p0.jpg"
      >
      </q-image>
      <q-text name="" avatar="http://q1.qlogo.cn/g?b=qq&nk=2595201156&s=640" tag-color="purple">
        123
        <img src="/85467403_p0.jpg" />
        546<br />
        <a>我是一个链接</a>
      </q-text>
      <q-text
        self
        name=""
        avatar="http://q1.qlogo.cn/g?b=qq&nk=2595201156&s=640"
        tag-color="purple"
      >
        收到
      </q-text>
    </q-main>
  </section>
</template>

<style lang="scss" scoped>
.qq {
  width: 100%;
  background-color: var(--qq-background-03);
  border-radius: 10px;
  transition: background-color 0.2s;
}
</style>
