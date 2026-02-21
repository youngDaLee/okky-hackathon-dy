<template>
  <div v-if="recipes.length === 0" class="text-sm text-gray-400 text-center py-6">
    재료를 추가하면 맞춤 레시피를 추천해드려요 🍳
  </div>

  <div v-else>
    <div class="relative overflow-hidden rounded-lg">
        <div
          class="flex transition-transform duration-300"
          :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
        >
          <div
            v-for="recipe in recipes"
            :key="recipe.id"
            class="min-w-full cursor-pointer"
            @click="$router.push(`/recipes/${recipe.id}`)"
          >
            <div class="bg-gradient-to-br from-orange-50 to-yellow-50 rounded-lg p-4">
              <span v-if="recipe.urgencyBonus" class="inline-block bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded-full mb-2">
                🔥 긴급
              </span>
              <span v-else-if="recipe.matchRate === 100" class="inline-block bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded-full mb-2">
                완벽 매칭
              </span>
              <span v-else class="inline-block bg-blue-500 text-white text-xs font-bold px-2 py-0.5 rounded-full mb-2">
                Tier {{ recipe.tier }}
              </span>
              <h3 class="text-base font-bold text-gray-900 mb-1">{{ recipe.title }}</h3>
              <p class="text-sm text-gray-600 mb-3 line-clamp-2">{{ recipe.description }}</p>
              <div class="flex flex-wrap gap-1 mb-2">
                <span
                  v-for="ing in recipe.ingredients?.slice(0, 4)"
                  :key="ing"
                  class="text-xs bg-white text-gray-600 px-2 py-0.5 rounded-full border border-gray-200"
                >
                  {{ ing }}
                </span>
              </div>
              <div class="flex gap-3 text-xs text-gray-500">
                <span v-if="recipe.cookTime">⏱ {{ recipe.cookTime }}</span>
                <span v-if="recipe.servings">👥 {{ recipe.servings }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

    <div v-if="recipes.length > 1" class="flex items-center justify-center gap-2 mt-3">
      <button
        class="p-1 text-gray-400 hover:text-gray-600"
        @click="prev"
      >
        ←
      </button>
      <div class="flex gap-1">
        <span
          v-for="(_, i) in recipes"
          :key="i"
          :class="i === currentIndex ? 'bg-blue-600' : 'bg-gray-300'"
          class="w-2 h-2 rounded-full transition-colors"
        />
      </div>
      <button
        class="p-1 text-gray-400 hover:text-gray-600"
        @click="next"
      >
        →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  recipes: {
    type: Array,
    default: () => [],
  },
})

const router = useRouter()
const currentIndex = ref(0)

function prev() {
  currentIndex.value = (currentIndex.value - 1 + props.recipes.length) % props.recipes.length
}

function next() {
  currentIndex.value = (currentIndex.value + 1) % props.recipes.length
}
</script>
