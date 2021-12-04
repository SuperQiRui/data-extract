<template>
  <router-view v-slot="{ Component }">
    <transition
      name="fade"
      mode="out-in"
    >
      <component :is="Component" />
    </transition>
  </router-view>
</template>
<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'App'
})
</script>

<style>
html, body {
  height: 100%;
  margin: 0;
}

#app {
  font-family: Helvetica, Tahoma, Arial;
  height: 100%;
}
</style>
