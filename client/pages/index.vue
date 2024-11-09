<script setup lang="ts">

const text = ref("")
const response = ref("")

const context = ref<number[]>([])

async function onSubmit() {
  const data = await $fetch<{ response: string, context: number[] }>("http://localhost:8080/", {
    method: "post",
    body: {text: text.value, context: context.value}
  })
  response.value = data.response
  context.value = data.context
}

</script>

<template>
  <h2>ChatBot</h2>

  {{ response }}

  <div>
    Text
    <input v-model="text" type="text">
    <button @click="onSubmit">send</button>
  </div>
</template>

<style scoped>

</style>