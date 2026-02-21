<template>
  <div id="app-layout" class="min-h-screen bg-gray-50 pb-20">
    <RouterView />
    <AppBottomNav />
    <FloatingAddButton v-if="!isAddIngredientPage" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import AppBottomNav from '@/components/AppBottomNav.vue'
import FloatingAddButton from '@/components/FloatingAddButton.vue'

const route = useRoute()

const isAddIngredientPage = computed(() => route.path === '/add-ingredient')
</script>
