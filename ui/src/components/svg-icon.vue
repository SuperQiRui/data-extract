<template>
  <svg
    aria-hidden="true"
    preserveAspectRatio="xMidYMid meet"
    :with="size"
    :height="size"
    :viewBox="`0 0 ${size} ${size}`"
  >
    <use
      :xlink:href="symbolId"
      :stroke="color"
    />
  </svg>
</template>

<script>
import { defineComponent, computed } from 'vue'

export default defineComponent({
  name: 'SvgIcon',
  props: {
    size: {
      type: Number,
      default: 18
    },
    name: {
      type: String,
      required: true
    },
    color: {
      type: String,
      default: '#fff'
    }
  },
  setup(props) {
    const symbolId = computed(() => `#icon-${props.name}`)
    return { symbolId }
  }
})
</script>
<style scoped>
svg use {
  fill: v-bind(color);
}
</style>
