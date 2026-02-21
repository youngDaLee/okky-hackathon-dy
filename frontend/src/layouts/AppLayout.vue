<template>
  <div id="app-layout" class="min-h-screen bg-gray-50 pb-20">
    <RouterView />
    <AppBottomNav />
  </div>
</template>

<script setup>
import { RouterView } from 'vue-router'
import AppBottomNav from '@/components/AppBottomNav.vue'
</script>
