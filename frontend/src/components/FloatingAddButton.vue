<template>
  <button
    @click="router.push('/add-ingredient')"
    class="fixed bottom-24 right-4 z-30 bg-blue-600 text-white rounded-full p-4 shadow-lg hover:bg-blue-700 active:bg-blue-800 transition-all hover:scale-110 active:scale-95"
    aria-label="재료 추가"
  >
    <Plus class="size-6" />
  </button>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { Plus } from 'lucide-vue-next'

const router = useRouter()
</script>
