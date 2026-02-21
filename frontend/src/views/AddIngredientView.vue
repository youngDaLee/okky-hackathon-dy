<template>
  <div class="max-w-4xl mx-auto px-4 py-4">
    <header class="mb-6">
      <div class="flex items-center justify-between mb-1">
        <h1 class="text-xl font-bold text-gray-900">
          {{ isEditMode ? '✏️ 재료 수정' : '🥕 재료 추가' }}
        </h1>
        <button
          type="button"
          @click="router.back()"
          class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          aria-label="닫기"
        >
          <X class="size-5" />
        </button>
      </div>
      <p class="text-sm text-gray-500">
        {{ isEditMode ? '재료 정보를 수정하세요' : '냉장고에 있는 재료를 등록하세요' }}
      </p>
    </header>

    <!-- 로딩 상태 -->
    <div v-if="loading" class="text-center py-8 text-gray-500">
      재료 정보를 불러오는 중...
    </div>

    <!-- 사진 빠른 추가 (추가 모드에서만 표시) -->
    <button
      v-if="!isEditMode"
      type="button"
      disabled
      class="w-full flex items-center justify-center gap-2 bg-gray-50 border border-dashed border-gray-300 rounded-xl py-3 mb-6 text-sm text-gray-400 cursor-not-allowed"
    >
      📷 영수증/재료 사진으로 빠른 추가 (준비 중)
    </button>

    <form v-if="!loading" @submit.prevent="handleSubmit" class="space-y-5">
      <!-- 재료 이름 -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">
          재료 이름 <span class="text-red-500">*</span>
        </label>
        <input
          v-model="form.name"
          type="text"
          placeholder="예: 감자, 양파, 우유"
          required
          class="w-full border border-gray-300 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
        />
      </div>

      <!-- 카테고리 -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">카테고리</label>
        <div class="flex flex-wrap gap-2">
          <CategoryButton
            v-for="cat in CATEGORIES"
            :key="cat.value"
            :label="cat.label"
            :selected="form.category === cat.value"
            @select="form.category = cat.value"
          />
        </div>
      </div>

      <!-- 수량 + 단위 -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">수량</label>
        <div class="flex gap-2">
          <input
            v-model.number="form.quantity"
            type="number"
            min="1"
            placeholder="1"
            class="w-28 border border-gray-300 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
          />
          <select
            v-model="form.unit"
            class="flex-1 border border-gray-300 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 bg-white"
          >
            <option v-for="u in UNITS" :key="u" :value="u">{{ u }}</option>
          </select>
        </div>
      </div>

      <!-- 유통기한 -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">
          유통기한 <span class="text-red-500">*</span>
        </label>
        <input
          v-model="form.expiryDate"
          type="date"
          required
          :min="today"
          class="w-full border border-gray-300 rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
        />
      </div>

      <!-- 에러 메시지 -->
      <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

      <!-- 버튼 -->
      <div class="flex gap-3 pt-2">
        <button
          type="button"
          @click="router.back()"
          class="flex-1 py-3 border border-gray-300 rounded-xl text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
        >
          취소
        </button>
        <button
          type="submit"
          :disabled="submitting"
          class="flex-1 py-3 bg-blue-600 rounded-xl text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50 transition-colors"
        >
          {{
            submitting
              ? (isEditMode ? '수정 중...' : '추가 중...')
              : (isEditMode ? '💾 수정하기' : '➕ 추가하기')
          }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { X } from 'lucide-vue-next'
import { useIngredientStore } from '@/stores/ingredient.js'
import CategoryButton from '@/components/CategoryButton.vue'

const CATEGORIES = [
  { label: '채소',   value: 'VEGETABLE' },
  { label: '과일',   value: 'FRUIT' },
  { label: '육류',   value: 'MEAT' },
  { label: '해산물', value: 'SEAFOOD' },
  { label: '유제품', value: 'DAIRY' },
  { label: '기타',   value: 'OTHER' },
]
const UNITS = ['개', 'g', 'kg', 'ml', 'L', '팩']

const router = useRouter()
const route = useRoute()
const store = useIngredientStore()

const ingredientId = computed(() => route.params.id)
const isEditMode = computed(() => !!ingredientId.value)
const loading = ref(false)

const today = computed(() => new Date().toISOString().split('T')[0])

const form = ref({
  name: '',
  category: '',
  quantity: 1,
  unit: '개',
  expiryDate: '',
})

const submitting = ref(false)
const errorMsg = ref('')

// 수정 모드일 때 기존 데이터 로드
onMounted(async () => {
  if (isEditMode.value) {
    loading.value = true
    try {
      const ingredient = await store.getIngredientById(ingredientId.value)
      form.value = {
        name: ingredient.name || '',
        category: ingredient.category || '',
        quantity: ingredient.quantity || 1,
        unit: ingredient.unit || '개',
        expiryDate: ingredient.expiry_date
          ? new Date(ingredient.expiry_date).toISOString().split('T')[0]
          : '',
      }
    } catch (e) {
      errorMsg.value = '재료 정보를 불러오는데 실패했어요.'
      console.error('Failed to load ingredient:', e)
    } finally {
      loading.value = false
    }
  }
})

async function handleSubmit() {
  errorMsg.value = ''
  submitting.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      category: form.value.category || 'OTHER',
      quantity: form.value.quantity ?? 1,
      unit: form.value.unit,
      expiryDate: form.value.expiryDate
        ? new Date(form.value.expiryDate).toISOString()
        : null,
    }

    if (isEditMode.value) {
      await store.updateIngredient(ingredientId.value, payload)
    } else {
      await store.addIngredient(payload)
    }
    router.push('/')
  } catch (e) {
    const code = e?.response?.data?.error || e?.response?.data?.code
    if (code === 'DUPLICATE_INGREDIENT') {
      errorMsg.value = '이미 등록된 재료예요.'
    } else if (code === 'FRIDGE_LIMIT_EXCEEDED') {
      errorMsg.value = '재료는 최대 200개까지 등록할 수 있어요.'
    } else if (code === 'INGREDIENT_NOT_FOUND') {
      errorMsg.value = '재료를 찾을 수 없어요.'
    } else if (code === 'VALIDATION_ERROR') {
      errorMsg.value = e?.response?.data?.message || '입력 정보를 확인해주세요.'
    } else {
      errorMsg.value = isEditMode.value
        ? '재료 수정에 실패했어요. 다시 시도해주세요.'
        : '재료 추가에 실패했어요. 다시 시도해주세요.'
    }
  } finally {
    submitting.value = false
  }
}
</script>
