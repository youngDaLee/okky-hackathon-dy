<template>
  <div class="container mx-auto px-4 py-8">
    <div class="text-center">
      <h1 class="text-4xl font-bold text-gray-900 mb-4">
        냉털마스터
      </h1>
      <p class="text-lg text-gray-600 mb-8">
        남는 재료만 넣으세요, 최고의 레시피를 찾아드립니다
      </p>
      <div class="flex gap-4 justify-center">
        <RouterLink
          to="/login"
          class="px-6 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition"
        >
          로그인
        </RouterLink>
        <RouterLink
          to="/signup"
          class="px-6 py-3 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition"
        >
          회원가입
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>
