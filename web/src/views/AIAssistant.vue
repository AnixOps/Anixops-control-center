<template>
  <div class="h-full flex flex-col">
    <!-- Header -->
    <div class="bg-slate-800 border-b border-slate-700 px-6 py-4">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold text-white">AI Assistant</h1>
          <p class="text-slate-400 text-sm mt-1">Powered by Workers AI & Llama 3.1</p>
        </div>
        <div class="flex gap-2">
          <button
            @click="activeTab = 'chat'"
            :class="['px-4 py-2 rounded-lg text-sm font-medium transition-colors', activeTab === 'chat' ? 'bg-blue-600 text-white' : 'bg-slate-700 text-slate-300 hover:bg-slate-600']"
          >
            Chat
          </button>
          <button
            @click="activeTab = 'analyze'"
            :class="['px-4 py-2 rounded-lg text-sm font-medium transition-colors', activeTab === 'analyze' ? 'bg-blue-600 text-white' : 'bg-slate-700 text-slate-300 hover:bg-slate-600']"
          >
            Log Analysis
          </button>
          <button
            @click="activeTab = 'search'"
            :class="['px-4 py-2 rounded-lg text-sm font-medium transition-colors', activeTab === 'search' ? 'bg-blue-600 text-white' : 'bg-slate-700 text-slate-300 hover:bg-slate-600']"
          >
            Semantic Search
          </button>
        </div>
      </div>
    </div>

    <!-- Chat Tab -->
    <div v-if="activeTab === 'chat'" class="flex-1 flex flex-col overflow-hidden">
      <div ref="chatContainer" class="flex-1 overflow-y-auto p-6 space-y-4">
        <div v-if="aiStore.chatHistory.length === 0" class="text-center py-12">
          <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
            <svg class="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
            </svg>
          </div>
          <h3 class="text-lg font-medium text-white">AI DevOps Assistant</h3>
          <p class="text-slate-400 mt-2">Ask me anything about your infrastructure, logs, or operations.</p>
          <div class="mt-6 grid grid-cols-3 gap-3 max-w-lg mx-auto">
            <button @click="quickPrompt('Show node status')" class="p-3 bg-slate-700 rounded-lg text-sm text-slate-300 hover:bg-slate-600 transition-colors">
              Show node status
            </button>
            <button @click="quickPrompt('Analyze recent errors')" class="p-3 bg-slate-700 rounded-lg text-sm text-slate-300 hover:bg-slate-600 transition-colors">
              Analyze errors
            </button>
            <button @click="quickPrompt('Best practices for security')" class="p-3 bg-slate-700 rounded-lg text-sm text-slate-300 hover:bg-slate-600 transition-colors">
              Security tips
            </button>
          </div>
        </div>

        <div
          v-for="(msg, idx) in aiStore.chatHistory"
          :key="idx"
          :class="['flex', msg.role === 'user' ? 'justify-end' : 'justify-start']"
        >
          <div
            :class="[
              'max-w-2xl rounded-2xl px-4 py-3',
              msg.role === 'user'
                ? 'bg-blue-600 text-white'
                : msg.isError
                ? 'bg-red-900/50 text-red-200 border border-red-800'
                : 'bg-slate-700 text-slate-200'
            ]"
          >
            <p class="whitespace-pre-wrap">{{ msg.content }}</p>
            <p class="text-xs opacity-50 mt-1">{{ formatTime(msg.timestamp) }}</p>
          </div>
        </div>

        <div v-if="aiStore.isThinking" class="flex justify-start">
          <div class="bg-slate-700 rounded-2xl px-4 py-3">
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 bg-blue-500 rounded-full animate-bounce"></div>
              <div class="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
              <div class="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
              <span class="text-slate-400 text-sm ml-2">Thinking...</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Input -->
      <div class="border-t border-slate-700 p-4 bg-slate-800">
        <form @submit.prevent="sendMessage" class="flex gap-3">
          <input
            v-model="userInput"
            type="text"
            placeholder="Ask about your infrastructure..."
            class="flex-1 bg-slate-700 border border-slate-600 rounded-xl px-4 py-3 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            :disabled="aiStore.isThinking"
          />
          <button
            type="submit"
            :disabled="!userInput.trim() || aiStore.isThinking"
            class="px-6 py-3 bg-blue-600 text-white rounded-xl font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Send
          </button>
          <button
            type="button"
            @click="aiStore.clearHistory"
            class="px-4 py-3 bg-slate-700 text-slate-300 rounded-xl hover:bg-slate-600 transition-colors"
          >
            Clear
          </button>
        </form>
      </div>
    </div>

    <!-- Log Analysis Tab -->
    <div v-if="activeTab === 'analyze'" class="flex-1 overflow-y-auto p-6">
      <div class="max-w-4xl mx-auto">
        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <h2 class="text-lg font-medium text-white mb-4">Log Analysis</h2>
          <p class="text-slate-400 text-sm mb-4">Paste your logs for AI-powered analysis using Workers AI.</p>

          <textarea
            v-model="logContent"
            rows="10"
            placeholder="Paste log content here..."
            class="w-full bg-slate-700 border border-slate-600 rounded-lg px-4 py-3 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
          ></textarea>

          <div class="flex gap-3 mt-4">
            <button
              @click="analyzeLogs"
              :disabled="!logContent.trim() || aiStore.isThinking"
              class="px-6 py-3 bg-purple-600 text-white rounded-lg font-medium hover:bg-purple-700 disabled:opacity-50 transition-colors"
            >
              {{ aiStore.isThinking ? 'Analyzing...' : 'Analyze Logs' }}
            </button>
          </div>

          <div v-if="analysisResult" class="mt-6 p-4 bg-slate-700 rounded-lg">
            <h3 class="font-medium text-white mb-2">Analysis Result</h3>
            <pre class="text-slate-300 text-sm whitespace-pre-wrap">{{ analysisResult }}</pre>
          </div>
        </div>
      </div>
    </div>

    <!-- Semantic Search Tab -->
    <div v-if="activeTab === 'search'" class="flex-1 overflow-y-auto p-6">
      <div class="max-w-4xl mx-auto">
        <div class="bg-slate-800 rounded-xl border border-slate-700 p-6">
          <h2 class="text-lg font-medium text-white mb-4">Semantic Search</h2>
          <p class="text-slate-400 text-sm mb-4">Search logs and tasks using natural language with Vectorize embeddings.</p>

          <div class="flex gap-3">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="e.g., 'connection timeout errors' or 'failed deployment tasks'"
              class="flex-1 bg-slate-700 border border-slate-600 rounded-lg px-4 py-3 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              @keyup.enter="semanticSearch"
            />
            <button
              @click="semanticSearch"
              :disabled="!searchQuery.trim() || aiStore.isThinking"
              class="px-6 py-3 bg-green-600 text-white rounded-lg font-medium hover:bg-green-700 disabled:opacity-50 transition-colors"
            >
              Search
            </button>
          </div>

          <div v-if="searchResults.length > 0" class="mt-6 space-y-3">
            <h3 class="font-medium text-white">Results</h3>
            <div
              v-for="(result, idx) in searchResults"
              :key="idx"
              class="p-4 bg-slate-700 rounded-lg border border-slate-600"
            >
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm font-medium text-blue-400">{{ result.id }}</span>
                <span class="text-xs text-slate-400">Score: {{ result.score?.toFixed(3) || 'N/A' }}</span>
              </div>
              <p class="text-slate-300 text-sm">{{ JSON.stringify(result.metadata, null, 2) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import { useAiStore } from '@/stores/ai'

const aiStore = useAiStore()

const activeTab = ref('chat')
const userInput = ref('')
const logContent = ref('')
const analysisResult = ref(null)
const searchQuery = ref('')
const searchResults = ref([])
const chatContainer = ref(null)

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString()
}

const sendMessage = async () => {
  if (!userInput.value.trim() || aiStore.isThinking) return

  const message = userInput.value
  userInput.value = ''

  await aiStore.sendMessage(message)

  await nextTick()
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

const quickPrompt = (prompt) => {
  userInput.value = prompt
  sendMessage()
}

const analyzeLogs = async () => {
  if (!logContent.value.trim()) return
  analysisResult.value = await aiStore.analyzeLogs(logContent.value)
}

const semanticSearch = async () => {
  if (!searchQuery.value.trim()) return
  searchResults.value = await aiStore.semanticSearch(searchQuery.value) || []
}

watch(aiStore.chatHistory, async () => {
  await nextTick()
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}, { deep: true })
</script>