<template>
  <div class="max-w-4xl mx-auto px-4 py-4 md:px-6 md:py-6">
    <header class="mb-6">
      <div class="flex items-center justify-between mb-2">
        <h1 class="text-xl md:text-2xl font-bold text-gray-900">🥕 재료 추가</h1>
        <button
          @click="router.back()"
          class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          aria-label="닫기"
        >
          <X class="size-5" />
        </button>
      </div>
      <p class="text-sm text-gray-600">냉장고에 있는 재료를 등록하세요</p>
    </header>
    <!-- Phase 3: AddIngredient 폼 구현 -->
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { X } from 'lucide-vue-next'

const router = useRouter()
</script>
