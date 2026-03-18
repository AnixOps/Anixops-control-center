<template>
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <label v-if="label" :for="id" class="text-sm text-dark-400">
        {{ label }}
      </label>
      <span v-if="max" class="text-xs text-dark-500">
        {{ modelValue?.length || 0 }}/{{ max }}
      </span>
    </div>
    <textarea
      :id="id"
      :value="modelValue"
      :placeholder="placeholder"
      :rows="rows"
      :maxlength="max"
      :disabled="disabled"
      :readonly="readonly"
      :class="textareaClasses"
      @input="$emit('update:modelValue', $event.target.value)"
    ></textarea>
    <p v-if="error" class="text-sm text-red-400">{{ error }}</p>
    <p v-else-if="hint" class="text-sm text-dark-500">{{ hint }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  id: {
    type: String,
    default: () => `textarea-${Math.random().toString(36).slice(2)}`
  },
  modelValue: {
    type: String,
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  rows: {
    type: Number,
    default: 4
  },
  max: {
    type: Number,
    default: 0
  },
  disabled: {
    type: Boolean,
    default: false
  },
  readonly: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  hint: {
    type: String,
    default: ''
  }
})

defineEmits(['update:modelValue'])

const textareaClasses = computed(() => {
  const base = 'w-full px-4 py-3 bg-dark-700 border rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 transition-colors resize-none'
  const borderColor = props.error
    ? 'border-red-500 focus:ring-red-500'
    : 'border-dark-600 focus:ring-primary-500'
  const disabledStyles = props.disabled ? 'opacity-50 cursor-not-allowed' : ''
  return `${base} ${borderColor} ${disabledStyles}`
})
</script>