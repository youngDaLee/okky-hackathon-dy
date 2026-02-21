<template>
  <div class="bg-white rounded-xl shadow-sm p-4">
    <h2 class="text-base font-bold text-gray-900 mb-3">📖 저장된 요리책</h2>

    <div v-if="bookmarks.length === 0" class="text-sm text-gray-400 text-center py-4">
      저장된 레시피가 없어요.
    </div>

    <ul v-else class="space-y-3">
      <li
        v-for="item in bookmarks"
        :key="item.id"
        class="flex items-start gap-2"
      >
        <span class="text-yellow-500 mt-0.5">★</span>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-800 truncate">{{ item.title }}</p>
          <a
            :href="item.url"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-1 text-xs text-blue-600 hover:underline mt-0.5"
          >
            <span :class="item.type === 'youtube' ? 'bg-red-100 text-red-600' : 'bg-blue-100 text-blue-600'"
              class="px-1.5 py-0.5 rounded text-xs font-medium"
            >
              {{ item.type === 'youtube' ? 'YouTube' : 'Blog' }}
            </span>
            {{ item.source }}
            <span>↗</span>
          </a>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
defineProps({
  bookmarks: {
    type: Array,
    default: () => [],
  },
})
</script>
