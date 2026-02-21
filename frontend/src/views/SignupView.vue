<template>
  <div class="container mx-auto px-4 py-8 max-w-md">
    <h1 class="text-3xl font-bold text-gray-900 mb-8 text-center">회원가입</h1>
    <form @submit.prevent="handleSignup" class="space-y-4">
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
        <label for="password" class="block text-sm font-medium text-gray-700 mb-1">
          비밀번호
        </label>
        <input
          id="password"
          v-model="password"
          type="password"
          required
          minlength="8"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          placeholder="최소 8자, 영문+숫자"
        />
        <p class="text-xs text-gray-400 mt-1">최소 8자, 영문과 숫자를 포함해야 해요.</p>
      </div>
      <div>
        <label for="nickname" class="block text-sm font-medium text-gray-700 mb-1">닉네임</label>
        <input
          id="nickname"
          v-model="nickname"
          type="text"
          required
          minlength="2"
          maxlength="20"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          placeholder="2~20자"
        />
      </div>

      <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

      <button
        type="submit"
        :disabled="submitting"
        class="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition font-medium"
      >
        {{ submitting ? '가입 중...' : '회원가입' }}
      </button>
    </form>
    <p class="mt-4 text-center text-sm text-gray-600">
      이미 계정이 있으신가요?
      <RouterLink to="/login" class="text-blue-600 hover:underline">로그인</RouterLink>
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
const nickname = ref('')
const submitting = ref(false)
const errorMsg = ref('')

const handleSignup = async () => {
  errorMsg.value = ''
  submitting.value = true
  try {
    await auth.signup(email.value, password.value, nickname.value)
    router.push('/')
  } catch (e) {
    const code = e?.response?.data?.code
    const field = e?.response?.data?.field
    if (code === 'DUPLICATE_EMAIL') {
      errorMsg.value = '이미 가입된 이메일이에요.'
    } else if (code === 'VALIDATION_ERROR' && field === 'password') {
      errorMsg.value = '비밀번호는 최소 8자, 영문과 숫자를 포함해야 해요.'
    } else if (code === 'VALIDATION_ERROR' && field === 'nickname') {
      errorMsg.value = '닉네임은 2~20자로 입력해주세요.'
    } else {
      errorMsg.value = '회원가입에 실패했어요. 다시 시도해주세요.'
    }
  } finally {
    submitting.value = false
  }
}
</script>
