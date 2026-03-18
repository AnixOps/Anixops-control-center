<template>
  <div class="relative">
    <!-- Theme Toggle Button -->
    <button
      @click="toggleDropdown"
      class="p-2 rounded-lg hover:bg-dark-700 transition-colors flex items-center gap-2"
      :title="themeStore.isDark ? 'Switch to light mode' : 'Switch to dark mode'"
    >
      <SunIcon v-if="themeStore.isDark" class="w-5 h-5 text-yellow-400" />
      <MoonIcon v-else class="w-5 h-5 text-slate-400" />
    </button>

    <!-- Theme Dropdown -->
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
        class="absolute right-0 mt-2 w-64 bg-dark-800 border border-dark-700 rounded-xl shadow-xl z-50"
      >
        <!-- Theme Mode -->
        <div class="p-4 border-b border-dark-700">
          <h3 class="text-sm font-medium text-white mb-3">Theme Mode</h3>
          <div class="flex gap-2">
            <button
              v-for="theme in themes"
              :key="theme.value"
              @click="themeStore.setTheme(theme.value)"
              :class="[
                'flex-1 px-3 py-2 rounded-lg text-sm transition-colors',
                themeStore.config.theme === theme.value
                  ? 'bg-primary-600 text-white'
                  : 'bg-dark-700 text-dark-300 hover:bg-dark-600'
              ]"
            >
              {{ theme.label }}
            </button>
          </div>
        </div>

        <!-- Theme Color -->
        <div class="p-4 border-b border-dark-700">
          <h3 class="text-sm font-medium text-white mb-3">Accent Color</h3>
          <div class="flex gap-2">
            <button
              v-for="color in colors"
              :key="color.value"
              @click="themeStore.setColor(color.value)"
              :class="[
                'w-8 h-8 rounded-full transition-transform',
                themeStore.config.color === color.value ? 'ring-2 ring-white ring-offset-2 ring-offset-dark-800 scale-110' : ''
              ]"
              :style="{ backgroundColor: color.color }"
              :title="color.label"
            />
          </div>
        </div>

        <!-- Font Size -->
        <div class="p-4 border-b border-dark-700">
          <h3 class="text-sm font-medium text-white mb-3">Font Size</h3>
          <div class="flex gap-2">
            <button
              v-for="size in fontSizes"
              :key="size.value"
              @click="themeStore.setFontSize(size.value)"
              :class="[
                'flex-1 px-3 py-2 rounded-lg text-sm transition-colors',
                themeStore.config.fontSize === size.value
                  ? 'bg-primary-600 text-white'
                  : 'bg-dark-700 text-dark-300 hover:bg-dark-600'
              ]"
            >
              {{ size.label }}
            </button>
          </div>
        </div>

        <!-- Compact Mode -->
        <div class="p-4">
          <div class="flex items-center justify-between">
            <span class="text-sm text-white">Compact Mode</span>
            <button
              @click="themeStore.toggleCompactMode()"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                themeStore.config.compactMode ? 'bg-primary-600' : 'bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  themeStore.config.compactMode ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const showDropdown = ref(false)

const themes = [
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' },
  { value: 'system', label: 'System' }
]

const colors = [
  { value: 'blue', label: 'Blue', color: '#3B82F6' },
  { value: 'purple', label: 'Purple', color: '#8B5CF6' },
  { value: 'green', label: 'Green', color: '#22C55E' },
  { value: 'orange', label: 'Orange', color: '#F97316' },
  { value: 'red', label: 'Red', color: '#EF4444' }
]

const fontSizes = [
  { value: 'small', label: 'S' },
  { value: 'medium', label: 'M' },
  { value: 'large', label: 'L' }
]

function toggleDropdown() {
  showDropdown.value = !showDropdown.value
}

function handleClickOutside(event) {
  if (!event.target.closest('.relative')) {
    showDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<script>
// Icons
const SunIcon = {
  template: `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
  </svg>`
}

const MoonIcon = {
  template: `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
  </svg>`
}
</script>