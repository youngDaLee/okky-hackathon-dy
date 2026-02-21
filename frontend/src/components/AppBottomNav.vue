<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 px-4 z-40">
    <div class="max-w-4xl mx-auto flex items-center justify-around h-16">
      <RouterLink
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        :class="item.isSpecial ? 'flex flex-col items-center justify-center -mt-8' : navItemClass(item.path)"
      >
        <template v-if="item.isSpecial">
          <div class="bg-blue-600 text-white rounded-full p-4 shadow-lg hover:bg-blue-700 transition-all hover:scale-110 active:scale-95">
            <component :is="item.icon" class="size-6" />
          </div>
          <span class="text-xs text-gray-600 mt-2">{{ item.label }}</span>
        </template>
        <template v-else>
          <component :is="item.icon" class="size-6" />
          <span class="text-xs font-medium">{{ item.label }}</span>
        </template>
      </RouterLink>
    </div>
  </nav>
</template>

<script setup>
import { RouterLink, useRoute } from 'vue-router'
import { Home, Plus, ChefHat } from 'lucide-vue-next'

const route = useRoute()

const navItems = [
  { path: '/',               icon: Home,     label: '홈' },
  { path: '/add-ingredient', icon: Plus,     label: '추가', isSpecial: true },
  { path: '/recipes',        icon: ChefHat,  label: '레시피' },
]

const navItemClass = (path) => {
  const isActive = route.path === path
  return [
    'flex flex-col items-center justify-center gap-1 py-2 px-4 min-w-[64px] transition-colors',
    isActive ? 'text-blue-600' : 'text-gray-500 hover:text-gray-700',
  ]
}
</script>
