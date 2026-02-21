<template>
  <button
    type="button"
    :class="[
      'px-4 py-2 rounded-full text-sm font-medium border transition-colors',
      selected
        ? 'bg-blue-600 border-blue-600 text-white'
        : 'bg-white border-gray-300 text-gray-600 hover:border-blue-400 hover:text-blue-600',
    ]"
    @click="$emit('select', label)"
  >
    {{ label }}
  </button>
</template>

<script setup>
defineProps({
  label: { type: String, required: true },
  selected: { type: Boolean, default: false },
})

defineEmits(['select'])
</script>
