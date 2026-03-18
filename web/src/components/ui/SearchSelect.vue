<template>
  <div class="relative">
    <div
      class="flex items-center gap-2 px-3 py-2 bg-dark-700 border border-dark-600 rounded-lg"
      :class="{ 'ring-2 ring-primary-500': focused }"
    >
      <input
        ref="inputRef"
        v-model="search"
        type="text"
        :placeholder="placeholder"
        class="flex-1 bg-transparent text-white placeholder-dark-500 focus:outline-none"
        @focus="focused = true"
        @blur="handleBlur"
        @keydown="handleKeydown"
      />
      <ChevronDownIcon class="w-5 h-5 text-dark-400" />
    </div>

    <!-- Dropdown -->
    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="showDropdown"
        class="absolute z-10 mt-1 w-full bg-dark-700 border border-dark-600 rounded-lg shadow-lg max-h-60 overflow-auto"
      >
        <div v-if="loading" class="px-4 py-3 text-dark-400 text-center">
          Loading...
        </div>
        <div v-else-if="filteredOptions.length === 0" class="px-4 py-3 text-dark-400 text-center">
          No results found
        </div>
        <button
          v-else
          v-for="option in filteredOptions"
          :key="option[valueKey]"
          type="button"
          class="w-full px-4 py-2 text-left text-white hover:bg-dark-600 transition-colors"
          :class="{ 'bg-primary-600/20': isSelected(option[valueKey]) }"
          @mousedown.prevent="selectOption(option)"
        >
          <div class="flex items-center justify-between">
            <span>{{ option[labelKey] }}</span>
            <CheckIcon v-if="isSelected(option[valueKey])" class="w-4 h-4 text-primary-400" />
          </div>
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ChevronDownIcon, CheckIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  modelValue: {
    type: [String, Number, Object],
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Search...'
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
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const inputRef = ref(null)
const search = ref('')
const focused = ref(false)
const showDropdown = ref(false)
const highlightedIndex = ref(0)

watch(focused, (val) => {
  showDropdown.value = val
})

const filteredOptions = computed(() => {
  if (!search.value) return props.options
  return props.options.filter(option =>
    String(option[props.labelKey]).toLowerCase().includes(search.value.toLowerCase())
  )
})

const isSelected = computed(() => (value) => {
  if (typeof props.modelValue === 'object' && props.modelValue !== null) {
    return props.modelValue[props.valueKey] === value
  }
  return props.modelValue === value
})

function selectOption(option) {
  emit('update:modelValue', option[props.valueKey])
  search.value = ''
  focused.value = false
}

function handleBlur() {
  setTimeout(() => {
    focused.value = false
  }, 200)
}

function handleKeydown(e) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    highlightedIndex.value = Math.min(highlightedIndex.value + 1, filteredOptions.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    highlightedIndex.value = Math.max(highlightedIndex.value - 1, 0)
  } else if (e.key === 'Enter' && filteredOptions.value[highlightedIndex.value]) {
    e.preventDefault()
    selectOption(filteredOptions.value[highlightedIndex.value])
  } else if (e.key === 'Escape') {
    focused.value = false
  }
}
</script>