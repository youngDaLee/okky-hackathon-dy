<template>
  <div class="container mx-auto px-4 py-8 max-w-md">
    <h1 class="text-3xl font-bold text-gray-900 mb-8 text-center">회원가입</h1>
    <form @submit.prevent="handleSignup" class="space-y-4">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
          이메일
        </label>
        <input
          id="email"
          v-model="email"
          type="email"
          required
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
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
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          placeholder="최소 8자, 영문+숫자"
        />
      </div>
      <div>
        <label for="nickname" class="block text-sm font-medium text-gray-700 mb-1">
          닉네임
        </label>
        <input
          id="nickname"
          v-model="nickname"
          type="text"
          required
          minlength="2"
          maxlength="20"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          placeholder="2-20자"
        />
      </div>
      <button
        type="submit"
        class="w-full px-4 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition font-medium"
      >
        회원가입
      </button>
    </form>
    <p class="mt-4 text-center text-sm text-gray-600">
      이미 계정이 있으신가요?
      <RouterLink to="/login" class="text-primary-600 hover:underline">
        로그인
      </RouterLink>
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'

const router = useRouter()
const email = ref('')
const password = ref('')
const nickname = ref('')

const handleSignup = async () => {
  // TODO: API 연동
  console.log('Signup:', {
    email: email.value,
    password: password.value,
    nickname: nickname.value
  })
  // router.push('/login')
}
</script>
