<template>
  <div class="relative" ref="dropdownRef">
    <button
      type="button"
      class="flex items-center gap-2 px-3 py-2 text-sm text-dark-300 hover:text-white bg-dark-700 hover:bg-dark-600 rounded-lg transition-colors"
      @click="isOpen = !isOpen"
    >
      <LanguageIcon class="w-5 h-5" />
      <span class="hidden sm:inline">{{ currentLocale?.nativeName }}</span>
      <ChevronDownIcon class="w-4 h-4" />
    </button>

    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute right-0 mt-2 w-48 bg-dark-700 border border-dark-600 rounded-lg shadow-lg overflow-hidden z-50"
      >
        <button
          v-for="locale in availableLocales"
          :key="locale.code"
          type="button"
          class="w-full flex items-center justify-between px-4 py-3 text-sm hover:bg-dark-600 transition-colors"
          :class="locale.code === currentLocaleCode ? 'text-primary-400 bg-dark-600/50' : 'text-dark-300 hover:text-white'"
          @click="changeLocale(locale.code)"
        >
          <div class="flex items-center gap-3">
            <span class="text-lg">{{ getLocaleFlag(locale.code) }}</span>
            <div class="text-left">
              <div class="font-medium">{{ locale.nativeName }}</div>
              <div class="text-xs text-dark-500">{{ locale.name }}</div>
            </div>
          </div>
          <CheckIcon v-if="locale.code === currentLocaleCode" class="w-4 h-4 text-primary-400" />
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { LanguageIcon, ChevronDownIcon, CheckIcon } from '@heroicons/vue/24/outline'
import { setLocale, getLocale, availableLocales } from '@/i18n'

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const currentLocaleCode = computed(() => getLocale())

const currentLocale = computed(() => {
  return availableLocales.find(l => l.code === currentLocaleCode.value)
})

function getLocaleFlag(code: string): string {
  const flags: Record<string, string> = {
    'en-US': '🇺🇸',
    'zh-CN': '🇨🇳',
    'zh-TW': '🇹🇼',
    'ja-JP': '🇯🇵',
    'ar-SA': '🇸🇦'
  }
  return flags[code] || '🌐'
}

function changeLocale(code: string) {
  setLocale(code)
  isOpen.value = false
}

// Close dropdown when clicking outside
function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>