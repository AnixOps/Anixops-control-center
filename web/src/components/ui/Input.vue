<template>
  <div class="relative">
    <input
      :id="id"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :class="inputClasses"
      @input="$emit('update:modelValue', $event.target.value)"
      @focus="$emit('focus', $event)"
      @blur="$emit('blur', $event)"
    />
    <label
      v-if="label"
      :for="id"
      class="absolute left-3 -top-2.5 px-1 text-xs bg-dark-800 text-dark-400 transition-all"
      :class="{ 'text-primary-400': focused }"
    >
      {{ label }}
    </label>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  id: {
    type: String,
    default: () => `input-${Math.random().toString(36).slice(2)}`
  },
  modelValue: {
    type: [String, Number],
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
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
  }
})

defineEmits(['update:modelValue', 'focus', 'blur'])

const focused = ref(false)

const inputClasses = computed(() => {
  const base = 'w-full px-4 py-2.5 bg-dark-700 border rounded-lg text-white placeholder-dark-500 focus:outline-none focus:ring-2 transition-colors'
  const borderColor = props.error
    ? 'border-red-500 focus:ring-red-500'
    : 'border-dark-600 focus:ring-primary-500'
  const disabledStyles = props.disabled ? 'opacity-50 cursor-not-allowed' : ''
  return `${base} ${borderColor} ${disabledStyles}`
})
</script>