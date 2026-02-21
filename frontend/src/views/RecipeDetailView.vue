<template>
  <div class="max-w-4xl mx-auto px-4 py-4 md:px-6 md:py-6">
    <div class="bg-white border-b sticky top-0 z-10 -mx-4 px-4 py-3 flex items-center gap-2">
      <button
        @click="router.back()"
        class="p-2 hover:bg-gray-100 rounded-lg transition-colors -ml-2"
        aria-label="뒤로 가기"
      >
        <ArrowLeft class="size-5" />
      </button>
      <h1 class="font-semibold text-base">레시피 상세</h1>
    </div>
    <!-- Phase 4: RecipeDetail 구현 -->
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'

const router = useRouter()
</script>
