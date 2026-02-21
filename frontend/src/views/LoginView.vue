<template>
  <div class="container mx-auto px-4 py-8 max-w-md">
    <h1 class="text-3xl font-bold text-gray-900 mb-8 text-center">로그인</h1>
    <form @submit.prevent="handleLogin" class="space-y-4">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 mb-1">이메일</label>
        <input
          id="email"
          v-model="email"
          type="email"
          required
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          placeholder="your@email.com"
        />
      </div>
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700 mb-1">비밀번호</label>
        <input
          id="password"
          v-model="password"
          type="password"
          required
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          placeholder="••••••••"
        />
      </div>

      <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

      <button
        type="submit"
        :disabled="submitting"
        class="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition font-medium"
      >
        {{ submitting ? '로그인 중...' : '로그인' }}
      </button>
    </form>
    <p class="mt-4 text-center text-sm text-gray-600">
      계정이 없으신가요?
      <RouterLink to="/signup" class="text-blue-600 hover:underline">회원가입</RouterLink>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const submitting = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  errorMsg.value = ''
  submitting.value = true
  try {
    await auth.login(email.value, password.value)
    router.push('/')
  } catch (e) {
    const code = e?.response?.data?.code
    if (code === 'INVALID_CREDENTIALS') {
      errorMsg.value = '이메일 또는 비밀번호가 올바르지 않아요.'
    } else {
      errorMsg.value = '로그인에 실패했어요. 다시 시도해주세요.'
    }
  } finally {
    submitting.value = false
  }
}
</script>
