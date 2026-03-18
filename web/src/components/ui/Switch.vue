<template>
  <div class="flex items-center gap-3">
    <button
      type="button"
      role="switch"
      :aria-checked="modelValue"
      :disabled="disabled"
      :class="switchClasses"
      @click="toggle"
    >
      <span
        :class="thumbClasses"
        class="inline-block w-4 h-4 transform rounded-full bg-white transition-transform duration-200 ease-in-out"
      ></span>
    </button>
    <label v-if="label" :for="id" class="text-sm text-dark-300 cursor-pointer" @click="toggle">
      {{ label }}
    </label>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  id: {
    type: String,
    default: () => `switch-${Math.random().toString(36).slice(2)}`
  },
  modelValue: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const switchClasses = computed(() => {
  const base = 'relative inline-flex h-6 w-11 shrink-0 rounded-full transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 focus:ring-offset-dark-800'
  const activeColor = props.modelValue ? 'bg-primary-600' : 'bg-dark-600'
  const disabledStyles = props.disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
  return `${base} ${activeColor} ${disabledStyles}`
})

const thumbClasses = computed(() => {
  return props.modelValue ? 'translate-x-5' : 'translate-x-1'
})

function toggle() {
  if (!props.disabled) {
    emit('update:modelValue', !props.modelValue)
  }
}
</script>