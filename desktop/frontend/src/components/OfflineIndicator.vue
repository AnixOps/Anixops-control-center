<template>
  <Transition
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0 translate-y-2"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition ease-in duration-150"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-2"
  >
    <div
      v-if="showIndicator"
      class="fixed bottom-4 right-4 z-50 flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg"
      :class="[
        isOnline
          ? 'bg-green-600/90 text-white'
          : 'bg-amber-600/90 text-white'
      ]"
    >
      <!-- Icon -->
      <div class="flex items-center justify-center w-8 h-8 rounded-full bg-white/20">
        <WifiIcon v-if="isOnline" class="w-5 h-5" />
        <WifiOffIcon v-else class="w-5 h-5" />
      </div>

      <!-- Message -->
      <div class="flex-1">
        <p class="font-medium text-sm">
          {{ isOnline ? 'Back Online' : 'Offline' }}
        </p>
        <p v-if="!isOnline && pendingCount > 0" class="text-xs text-white/80">
          {{ pendingCount }} pending action{{ pendingCount !== 1 ? 's' : '' }} will sync
        </p>
        <p v-if="isOnline && syncResult" class="text-xs text-white/80">
          Synced {{ syncResult.success }} action{{ syncResult.success !== 1 ? 's' : '' }}
        </p>
      </div>

      <!-- Close button -->
      <button
        @click="dismiss"
        class="p-1 hover:bg-white/20 rounded transition-colors"
      >
        <XMarkIcon class="w-4 h-4" />
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { WifiIcon, WifiOffIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useOfflineStore } from '../stores/offline'

const offlineStore = useOfflineStore()

const dismissed = ref(false)
const syncResult = ref<{ success: number; failed: number } | null>(null)
const showTimeout = ref<number | null>(null)

const isOnline = computed(() => offlineStore.isOnline.value)
const pendingCount = computed(() => offlineStore.pendingActionsCount.value)

const showIndicator = computed(() => {
  if (dismissed.value) return false
  return !isOnline.value || (isOnline.value && syncResult.value)
})

// Watch for online status changes
watch(isOnline, async (newVal, oldVal) => {
  // Reset dismissed when status changes
  dismissed.value = false
  syncResult.value = null

  // Clear any existing timeout
  if (showTimeout.value) {
    clearTimeout(showTimeout.value)
  }

  // If just came online, try to sync
  if (newVal && !oldVal) {
    syncResult.value = await offlineStore.syncPendingActions()

    // Auto-dismiss after showing sync result
    showTimeout.value = window.setTimeout(() => {
      dismissed.value = true
    }, 3000)
  }
})

function dismiss() {
  dismissed.value = true
  syncResult.value = null
}

onMounted(() => {
  offlineStore.initialize()
})

onUnmounted(() => {
  offlineStore.cleanup()
  if (showTimeout.value) {
    clearTimeout(showTimeout.value)
  }
})
</script>