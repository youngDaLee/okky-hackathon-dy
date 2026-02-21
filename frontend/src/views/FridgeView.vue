<template>
  <div class="max-w-4xl mx-auto px-4 py-4">
    <!-- 헤더 -->
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-xl font-bold text-gray-900">🧊 냉장고 현황</h1>
      <RouterLink
        to="/add-ingredient"
        class="flex items-center gap-1 bg-blue-600 text-white text-sm font-medium px-3 py-2 rounded-lg hover:bg-blue-700 transition-colors"
      >
        <Plus class="size-4" /> 재료 추가
      </RouterLink>
    </div>

    <!-- 요약 카드 -->
    <div class="grid grid-cols-3 gap-3 mb-4">
      <div class="bg-white rounded-xl p-3 text-center shadow-sm">
        <p class="text-2xl font-bold text-gray-900">{{ store.summary.total }}</p>
        <p class="text-xs text-gray-500 mt-0.5">전체</p>
      </div>
      <div class="bg-red-50 rounded-xl p-3 text-center shadow-sm">
        <p class="text-2xl font-bold text-red-600">{{ store.summary.urgent }}</p>
        <p class="text-xs text-red-500 mt-0.5">임박</p>
      </div>
      <div class="bg-yellow-50 rounded-xl p-3 text-center shadow-sm">
        <p class="text-2xl font-bold text-yellow-600">{{ store.summary.soon }}</p>
        <p class="text-xs text-yellow-500 mt-0.5">주의</p>
      </div>
    </div>

    <!-- 검색 -->
    <div class="relative mb-3">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-gray-400" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="재료 검색..."
        class="w-full pl-9 pr-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:border-blue-500 bg-white"
      />
    </div>

    <!-- 카테고리 필터 -->
    <div class="flex gap-2 overflow-x-auto pb-2 mb-3 scrollbar-hide">
      <button
        v-for="cat in CATEGORY_FILTERS"
        :key="cat.value"
        :class="[
          'flex-shrink-0 px-3 py-1.5 rounded-full text-xs font-medium transition-colors',
          selectedCategory === cat.value
            ? 'bg-blue-600 text-white'
            : 'bg-white text-gray-600 border border-gray-200 hover:border-blue-300',
        ]"
        @click="selectedCategory = cat.value"
      >
        {{ cat.label }}
      </button>
    </div>

    <!-- 유통기한 필터 -->
    <div class="flex gap-2 mb-4">
      <button
        v-for="f in EXPIRY_FILTERS"
        :key="f.value"
        :class="[
          'px-3 py-1.5 rounded-full text-xs font-medium transition-colors',
          selectedExpiry === f.value
            ? 'bg-gray-800 text-white'
            : 'bg-white text-gray-600 border border-gray-200 hover:border-gray-400',
        ]"
        @click="selectedExpiry = f.value"
      >
        {{ f.label }}
      </button>
    </div>

    <!-- 로딩 -->
    <div v-if="store.loading" class="text-sm text-gray-400 text-center py-12">
      불러오는 중...
    </div>

    <!-- 빈 상태 -->
    <div v-else-if="filteredItems.length === 0" class="text-center py-12">
      <p class="text-gray-400 text-sm">등록된 재료가 없어요.</p>
      <RouterLink to="/add-ingredient" class="text-blue-600 text-sm mt-2 inline-block hover:underline">
        재료 추가하러 가기 →
      </RouterLink>
    </div>

    <!-- 재료 목록 -->
    <ul v-else class="space-y-2">
      <li
        v-for="item in filteredItems"
        :key="item.id"
        class="flex items-center gap-3 bg-white rounded-xl px-4 py-3 shadow-sm"
      >
        <span class="text-xl">{{ categoryEmoji(item.category) }}</span>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-900 truncate">{{ item.name }}</p>
          <p class="text-xs text-gray-400">{{ categoryLabel(item.category) }} · {{ item.quantity }}{{ item.unit }}</p>
        </div>
        <span
          :class="expiryBadgeClass(item.expiry_status)"
          class="text-xs font-medium px-2 py-0.5 rounded-full flex-shrink-0"
        >
          {{ expiryLabel(item) }}
        </span>
        <div class="flex items-center gap-1 flex-shrink-0">
          <button
            @click="handleEdit(item.id)"
            class="p-1.5 text-gray-300 hover:text-blue-500 transition-colors"
            aria-label="수정"
          >
            <Edit class="size-4" />
          </button>
          <button
            @click="handleDelete(item.id)"
            class="p-1.5 text-gray-300 hover:text-red-500 transition-colors"
            aria-label="삭제"
          >
            <Trash2 class="size-4" />
          </button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Plus, Search, Trash2, Edit } from 'lucide-vue-next'
import { useIngredientStore } from '@/stores/ingredient.js'

const router = useRouter()

const store = useIngredientStore()

const searchQuery = ref('')
const selectedCategory = ref('')
const selectedExpiry = ref('')

const CATEGORY_FILTERS = [
  { label: '전체',   value: '' },
  { label: '채소',   value: 'VEGETABLE' },
  { label: '과일',   value: 'FRUIT' },
  { label: '육류',   value: 'MEAT' },
  { label: '해산물', value: 'SEAFOOD' },
  { label: '유제품', value: 'DAIRY' },
  { label: '기타',   value: 'OTHER' },
]

const EXPIRY_FILTERS = [
  { label: '전체',  value: '' },
  { label: '🔴 임박', value: 'URGENT' },
  { label: '🟡 주의', value: 'SOON' },
  { label: '🟢 여유', value: 'NORMAL' },
]

const CATEGORY_EMOJI_MAP = {
  VEGETABLE: '🥕',
  FRUIT: '🍎',
  MEAT: '🥩',
  SEAFOOD: '🐟',
  DAIRY: '🥛',
  GRAIN: '🌾',
  CONDIMENT: '🧂',
  FROZEN: '🧊',
  OTHER: '🍱',
}

const CATEGORY_LABEL_MAP = {
  VEGETABLE: '채소',
  FRUIT: '과일',
  MEAT: '육류',
  SEAFOOD: '해산물',
  DAIRY: '유제품',
  GRAIN: '곡류',
  CONDIMENT: '양념',
  FROZEN: '냉동',
  OTHER: '기타',
}

function categoryEmoji(cat) {
  return CATEGORY_EMOJI_MAP[cat] ?? '🍱'
}

function categoryLabel(cat) {
  return CATEGORY_LABEL_MAP[cat] ?? '기타'
}

function expiryBadgeClass(status) {
  if (status === 'URGENT') return 'bg-red-100 text-red-700'
  if (status === 'SOON') return 'bg-yellow-100 text-yellow-700'
  if (status === 'NORMAL') return 'bg-green-100 text-green-700'
  return 'bg-gray-100 text-gray-500'
}

function expiryLabel(item) {
  if (item.expiry_status === 'NO_EXPIRY') return '유통기한 없음'
  if (!item.expiry_date) return ''
  const diff = Math.ceil((new Date(item.expiry_date) - new Date()) / (1000 * 60 * 60 * 24))
  if (diff < 0) return '만료됨'
  if (diff === 0) return 'D-Day'
  return `D-${diff}`
}

const filteredItems = computed(() => {
  return store.ingredients.filter((item) => {
    const matchSearch = searchQuery.value
      ? item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
      : true
    const matchCategory = selectedCategory.value
      ? item.category === selectedCategory.value
      : true
    const matchExpiry = selectedExpiry.value
      ? item.expiry_status === selectedExpiry.value
      : true
    return matchSearch && matchCategory && matchExpiry
  })
})

function handleEdit(id) {
  router.push({ name: 'edit-ingredient', params: { id } })
}

async function handleDelete(id) {
  if (!confirm('재료를 삭제할까요?')) return
  await store.removeIngredient(id)
}

onMounted(async () => {
  await Promise.all([store.fetchIngredients(), store.fetchSummary()])
})
</script>
