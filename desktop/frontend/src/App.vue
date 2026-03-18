<template>
  <div class="h-screen flex flex-col bg-slate-900">
    <!-- Custom Titlebar -->
    <div class="titlebar h-8 bg-slate-800 flex items-center justify-between px-4 border-b border-slate-700">
      <div class="flex items-center gap-2">
        <img src="/wails-logo.png" alt="Logo" class="w-4 h-4" v-if="false">
        <span class="text-xs text-slate-400">AnixOps Control Center</span>
      </div>
      <div class="flex items-center gap-1">
        <button @click="minimizeWindow" class="w-8 h-6 flex items-center justify-center hover:bg-slate-700 rounded">
          <svg class="w-3 h-3 text-slate-400" fill="currentColor" viewBox="0 0 24 24">
            <path d="M19 13H5v-2h14v2z"/>
          </svg>
        </button>
        <button @click="maximizeWindow" class="w-8 h-6 flex items-center justify-center hover:bg-slate-700 rounded">
          <svg class="w-3 h-3 text-slate-400" fill="currentColor" viewBox="0 0 24 24">
            <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V5h14v14z"/>
          </svg>
        </button>
        <button @click="closeWindow" class="w-8 h-6 flex items-center justify-center hover:bg-red-600 rounded">
          <svg class="w-3 h-3 text-slate-400 hover:text-white" fill="currentColor" viewBox="0 0 24 24">
            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 overflow-hidden">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { WindowMinimise, WindowMaximise, WindowClose } from '../../wailsjs/runtime/runtime'

function minimizeWindow() {
  WindowMinimise()
}

function maximizeWindow() {
  WindowMaximise()
}

function closeWindow() {
  WindowClose()
}
</script>