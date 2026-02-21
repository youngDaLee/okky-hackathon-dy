<template>
  <div class="bg-white rounded-xl shadow-sm p-4">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-base font-bold text-gray-900">나의 냉장고 현황</h2>
      <RouterLink to="/add-ingredient" class="text-sm text-blue-600 hover:text-blue-700">
        더보기 &gt;
      </RouterLink>
    </div>

    <div v-if="loading" class="text-sm text-gray-400 text-center py-4">불러오는 중...</div>

    <div v-else-if="items.length === 0" class="text-sm text-gray-400 text-center py-4">
      등록된 재료가 없어요. 재료를 추가해보세요!
    </div>

    <ul v-else class="space-y-2">
      <li
        v-for="item in items"
        :key="item.id"
        class="flex items-center justify-between"
      >
        <span class="text-sm text-gray-800">{{ categoryEmoji(item.category) }} {{ item.name }}</span>
        <span
          v-if="item.expiry_status !== 'NO_EXPIRY'"
          :class="expiryBadgeClass(item.expiry_status)"
          class="text-xs font-medium px-2 py-0.5 rounded-full"
        >
          {{ expiryLabel(item) }}
        </span>
        <span v-else class="text-xs text-gray-400">유통기한 없음</span>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'

const props = defineProps({
  items: {
    type: Array,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const CATEGORY_EMOJI = {
  채소: '🥕',
  과일: '🍎',
  육류: '🥩',
  해산물: '🐟',
  유제품: '🥛',
  기타: '🍱',
}

function categoryEmoji(category) {
  return CATEGORY_EMOJI[category] ?? '🍱'
}

function expiryBadgeClass(status) {
  if (status === 'URGENT') return 'bg-red-100 text-red-700'
  if (status === 'SOON') return 'bg-yellow-100 text-yellow-700'
  return 'bg-green-100 text-green-700'
}

function expiryLabel(item) {
  if (!item.expiry_date) return ''
  const diff = Math.ceil((new Date(item.expiry_date) - new Date()) / (1000 * 60 * 60 * 24))
  if (diff < 0) return '만료됨'
  if (diff === 0) return 'D-Day'
  return `D-${diff}`
}
</script>
