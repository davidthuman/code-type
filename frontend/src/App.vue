<script setup>

import { ref, computed } from 'vue';
import { useMagicKeys, whenever, useEventListener } from '@vueuse/core';

const { current } = useMagicKeys();

const keys = ref([]);

useEventListener(window, 'keydown', (e) => {

  switch (e.key) {
    case "Tab":
      e.preventDefault()
      //keys.value.push("&nbsp;", "&nbsp;", "&nbsp;", "&nbsp;")
      keys.value.push("    ")
      break;
    case "Backspace":
      keys.value.pop()
      break;
    case "Enter":
      keys.value.push('\n')
      break;
    case "Shift":
      break;
    default:
      keys.value.push(e.key?.toLowerCase())
  }
})

</script>

<template>
  <div class="p-5">
    <div class="font-bold">Hello world</div>
    <code class="flex flex-col font-mono space-x-1">
      <div style="white-space: pre-wrap;">
        <letter v-for="key in keys" :key="key">{{ key }}</letter>
      </div>
    </code>
    <br />
    <div>{{ keys }}</div>
  </div>
</template>

<style scoped></style>
