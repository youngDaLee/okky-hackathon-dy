<template>
  <div
    v-if="count > 0"
    class="flex items-center gap-3 bg-red-50 border-l-4 border-red-500 px-4 py-3 rounded-r-lg"
  >
    <span class="text-lg">⚠️</span>
    <p class="text-sm text-red-700 font-medium">
      유통기한 임박 재료가 <strong>{{ count }}개</strong> 있어요! 빨리 요리하세요.
    </p>
  </div>
</template>

<script setup>
defineProps({
  count: {
    type: Number,
    required: true,
  },
})
</script>
