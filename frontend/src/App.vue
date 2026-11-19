<script setup lang="ts">
  import { ref } from 'vue';
  import apiUrl from'./services/api-url.ts'

  const url = ref("")
  const code = ref("")
  const origin = ref(window.location.origin)
  const shortUrl = ref("")
  async function shorten() {
    try {
      const response = await apiUrl.createCode(JSON.stringify({
        url: url.value
      }))
      code.value = response.data.code
      shortUrl.value = origin.value + '/' + code.value
    } catch (error) {
      console.log(error)
    }
  }
</script>

<template>
  <h1 class="mx-10 text-3xl text-white">kurz Url Shortener</h1>
  <div class="flex min-w-full m-10">
    <input v-model="url" placeholder="Enter URL here" class="p-2 min-w-4/5 text-2xl text-white
      outline-2 outline-lime-700 focus:outline-yellow-500 bg-stone-700">
    <button @click="shorten" class="mx-2 p-2 text-2xl font-bold text-white bg-stone-700
      outline-2 outline-lime-700">Shorten</button>
  </div>
  <div v-if="code" class="flex space-x-2 m-10 text-2xl text-white">
    <p>Short URL:</p>
    <a :href="shortUrl" class="text-blue-300 hover:underline">{{ shortUrl }}</a>
  </div>
</template>

<style scoped></style>
