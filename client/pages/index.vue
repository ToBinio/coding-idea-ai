<script setup lang="ts">

const question = ref("")
const thoughts = ref("")
const answers = ref<string[]>([])

const text = ref("")

const thinking = ref(false)
const hasStarted = ref(false)


const context = ref<number[]>([])

type response = {
  context: number[],
  thoughts: string,
  question: string,
  answers: string[],
}

async function onSubmit() {
  hasStarted.value = true
  thinking.value = true

  const data = await $fetch<{ response: response }>("http://localhost:8080/", {
    method: "post",
    body: {text: text.value, context: context.value}
  })

  thoughts.value = data.response.thoughts
  question.value = data.response.question
  answers.value = data.response.answers

  context.value = data.response.context
  thinking.value = false
}

async function onSelectAnswer(answer: string) {
  text.value = answer;
  await onSubmit()
}

</script>

<template>
  <div class="m-2 flex flex-col items-center">
    <div class="flex flex-col items-center w-1/2 gap-5">
      <h1 class="text-3xl">Idea Bot</h1>

      <div v-if="thinking">
        Thinking...
      </div>

      <div class="text-sm w-3/4 text-center">
        {{ thoughts }}
      </div>
      <div class="text-xl text-center">
        {{ question }}
      </div>

      <div class="flex gap-2" v-if="!hasStarted">
        <input class="border-b-2" placeholder="start the conversation" v-model="text" type="text" :disabled="thinking"/>
        <button class="bg-gray-200 px-2 rounded hover:bg-gray-300 transition" @click="onSubmit" :disabled="thinking">
          send
        </button>
      </div>

      <div class="grid grid-cols-3 gap-2">
        <button
            class="outline rounded p-0.5 px-2 outline-gray-200 outline-2 hover:bg-gray-200 hover:outline-gray-300 transition"
            :disabled="thinking"
            v-for="answer of answers" @click="onSelectAnswer(answer)">
          {{ answer }}
        </button>
      </div>
      <button class="text-xl rounded p-0.5 px-2 bg-gray-200 hover:bg-gray-300"
              :disabled="thinking"
              @click="onSelectAnswer('please give an project idea now!')">
        give me an idea
      </button>
    </div>
  </div>
</template>

<style scoped>

</style>