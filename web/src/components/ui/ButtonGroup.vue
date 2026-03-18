<template>
  <div class="space-y-3">
    <label class="text-sm font-medium text-dark-300">{{ label }}</label>
    <div class="flex flex-wrap gap-2">
      <button
        v-for="option in options"
        :key="option[valueKey]"
        type="button"
        :class="[
          'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
          isSelected(option[valueKey])
            ? 'bg-primary-600 text-white'
            : 'bg-dark-700 text-dark-300 hover:bg-dark-600 hover:text-white'
        ]"
        @click="toggle(option[valueKey])"
      >
        {{ option[labelKey] }}
      </button>
    </div>
    <p v-if="error" class="text-sm text-red-400">{{ error }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: [Array, String, Number],
    default: () => []
  },
  label: {
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
  multiple: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const isSelected = computed(() => (value) => {
  if (props.multiple) {
    return Array.isArray(props.modelValue) && props.modelValue.includes(value)
  }
  return props.modelValue === value
})

function toggle(value) {
  if (props.multiple) {
    const current = Array.isArray(props.modelValue) ? [...props.modelValue] : []
    const index = current.indexOf(value)
    if (index === -1) {
      current.push(value)
    } else {
      current.splice(index, 1)
    }
    emit('update:modelValue', current)
  } else {
    emit('update:modelValue', value)
  }
}
</script>