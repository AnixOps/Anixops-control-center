<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="variantClasses"
    class="inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-dark-800 disabled:opacity-50 disabled:cursor-not-allowed"
  >
    <svg
      v-if="loading"
      class="w-4 h-4 animate-spin"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      ></circle>
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      ></path>
    </svg>
    <slot v-if="!loading" name="icon"></slot>
    <slot></slot>
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'danger', 'ghost', 'outline'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  type: {
    type: String,
    default: 'button'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const variantClasses = computed(() => {
  const baseClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base'
  }

  const variants = {
    primary: `${baseClasses[props.size]} bg-primary-600 hover:bg-primary-700 text-white focus:ring-primary-500`,
    secondary: `${baseClasses[props.size]} bg-dark-700 hover:bg-dark-600 text-white focus:ring-dark-500`,
    danger: `${baseClasses[props.size]} bg-red-600 hover:bg-red-700 text-white focus:ring-red-500`,
    ghost: `${baseClasses[props.size]} bg-transparent hover:bg-dark-700 text-dark-300 hover:text-white focus:ring-dark-500`,
    outline: `${baseClasses[props.size]} bg-transparent border border-dark-600 hover:bg-dark-700 text-white focus:ring-dark-500`
  }

  return variants[props.variant]
})
</script>