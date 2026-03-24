import { defineStore } from 'pinia'
import { ref } from 'vue'
import { aiApi, vectorApi } from '@/api'

export const useAiStore = defineStore('ai', () => {
  const chatHistory = ref([])
  const isThinking = ref(false)
  const error = ref(null)

  async function sendMessage(message) {
    if (!message.trim()) return

    // Add user message
    chatHistory.value.push({
      role: 'user',
      content: message,
      timestamp: new Date().toISOString()
    })

    isThinking.value = true
    error.value = null

    try {
      const response = await aiApi.chat(message, chatHistory.value.slice(0, -1))
      const assistantMessage = response.data?.response || response.data?.data?.response || 'I apologize, but I could not generate a response.'

      chatHistory.value.push({
        role: 'assistant',
        content: assistantMessage,
        timestamp: new Date().toISOString()
      })

      return assistantMessage
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to get AI response'
      chatHistory.value.push({
        role: 'assistant',
        content: 'Sorry, I encountered an error. Please try again.',
        timestamp: new Date().toISOString(),
        isError: true
      })
      return null
    } finally {
      isThinking.value = false
    }
  }

  async function analyzeLogs(logContent) {
    isThinking.value = true
    error.value = null

    try {
      const response = await aiApi.analyzeLog(logContent)
      return response.data?.data || response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to analyze logs'
      return null
    } finally {
      isThinking.value = false
    }
  }

  async function getOpsAdvice(context) {
    isThinking.value = true
    error.value = null

    try {
      const response = await aiApi.opsAdvice(context)
      return response.data?.data || response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to get advice'
      return null
    } finally {
      isThinking.value = false
    }
  }

  async function generateEmbedding(text) {
    try {
      const response = await aiApi.embedding(text)
      return response.data?.embedding || response.data?.data?.embedding
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to generate embedding'
      return null
    }
  }

  async function semanticSearch(query) {
    isThinking.value = true
    error.value = null

    try {
      // First get embedding for query
      const embedding = await generateEmbedding(query)
      if (!embedding) return null

      // Then search vectors
      const response = await vectorApi.search(embedding, { topK: 10 })
      return response.data?.results || response.data?.data?.results || []
    } catch (err) {
      error.value = err.response?.data?.error || 'Search failed'
      return null
    } finally {
      isThinking.value = false
    }
  }

  function clearHistory() {
    chatHistory.value = []
  }

  return {
    chatHistory,
    isThinking,
    error,
    sendMessage,
    analyzeLogs,
    getOpsAdvice,
    generateEmbedding,
    semanticSearch,
    clearHistory
  }
})