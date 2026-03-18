<template>
  <div class="relative">
    <select
      :id="id"
      :value="modelValue"
      :disabled="disabled"
      :class="selectClasses"
      @change="$emit('update:modelValue', $event.target.value)"
    >
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option
        v-for="option in options"
        :key="option[valueKey]"
        :value="option[valueKey]"
      >
        {{ option[labelKey] }}
      </option>
    </select>
    <label
      v-if="label"
      :for="id"
      class="absolute left-3 -top-2.5 px-1 text-xs bg-dark-800 text-dark-400"
    >
      {{ label }}
    </label>
    <ChevronDownIcon class="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-400 pointer-events-none" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ChevronDownIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  id: {
    type: String,
    default: () => `select-${Math.random().toString(36).slice(2)}`
  },
  modelValue: {
    type: [String, Number],
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
  options: {
    type: Array,
    default: () => []
  },
  valueKey: {
    type: String,
    default: 'value'
  },
  labelKey: {
    type: String,
    default: 'label'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  }
})

defineEmits(['update:modelValue'])

const selectClasses = computed(() => {
  const base = 'w-full px-4 py-2.5 bg-dark-700 border rounded-lg text-white appearance-none focus:outline-none focus:ring-2 transition-colors cursor-pointer'
  const borderColor = props.error
    ? 'border-red-500 focus:ring-red-500'
    : 'border-dark-600 focus:ring-primary-500'
  const disabledStyles = props.disabled ? 'opacity-50 cursor-not-allowed' : ''
  return `${base} ${borderColor} ${disabledStyles}`
})
</script>