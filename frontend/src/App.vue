<script setup lang="ts">
  import { ref } from 'vue';
  import apiUrl from'./services/api-url.ts'

  const url = ref("")
  const code = ref("")
  async function shorten() {
    try {
      const response = await apiUrl.createCode(JSON.stringify({
        url: url.value
      }))
      code.value = response.data.code
    } catch (error) {
      console.log(error)
    }
  }
</script>

<template>
  <h1>kurz Url Shortener</h1>
  <input v-model="url" placeholder="Enter URL here">
  <button @click="shorten">Shorten</button>
  <div>kurz.nhmk.de/{{ code }}</div>
</template>

<style scoped></style>
