<template>
  <div class="min-h-screen bg-dark-900">
    <div v-if="authStore.isAuthenticated">
      <Sidebar />
      <main class="ml-64 p-6">
        <router-view />
      </main>
    </div>
    <div v-else>
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import Sidebar from '@/components/Sidebar.vue'

const authStore = useAuthStore()
</script>