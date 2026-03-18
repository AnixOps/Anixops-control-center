<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-lg font-semibold text-white">{{ title }}</h3>
        <p v-if="subtitle" class="text-dark-400 text-sm mt-1">{{ subtitle }}</p>
      </div>
      <div v-if="$slots.actions" class="flex items-center gap-2">
        <slot name="actions"></slot>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-dark-800 rounded-xl border border-dark-700 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-dark-700">
            <tr>
              <th
                v-for="column in columns"
                :key="column.key"
                :class="[
                  'px-6 py-3 text-xs font-medium text-dark-300 uppercase tracking-wider',
                  column.align === 'right' ? 'text-right' : 'text-left'
                ]"
              >
                {{ column.label }}
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-dark-700">
            <template v-if="loading">
              <tr>
                <td :colspan="columns.length" class="px-6 py-12 text-center">
                  <div class="flex flex-col items-center gap-3">
                    <svg class="w-8 h-8 text-primary-400 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <p class="text-dark-400">Loading...</p>
                  </div>
                </td>
              </tr>
            </template>
            <template v-else-if="data.length === 0">
              <tr>
                <td :colspan="columns.length" class="px-6 py-12 text-center">
                  <slot name="empty">
                    <div class="flex flex-col items-center gap-3">
                      <component :is="emptyIcon" class="w-12 h-12 text-dark-500" />
                      <p class="text-dark-400">{{ emptyText }}</p>
                    </div>
                  </slot>
                </td>
              </tr>
            </template>
            <template v-else>
              <tr
                v-for="(row, index) in data"
                :key="rowKey ? row[rowKey] : index"
                class="hover:bg-dark-700/50 transition-colors"
              >
                <td
                  v-for="column in columns"
                  :key="column.key"
                  :class="[
                    'px-6 py-4',
                    column.align === 'right' ? 'text-right' : 'text-left'
                  ]"
                >
                  <slot :name="`cell-${column.key}`" :row="row" :value="row[column.key]">
                    <span class="text-dark-300">{{ row[column.key] }}</span>
                  </slot>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="pagination && total > 0" class="px-6 py-4 border-t border-dark-700 flex items-center justify-between">
        <p class="text-dark-400 text-sm">
          Showing {{ startItem }} to {{ endItem }} of {{ total }} {{ itemLabel }}
        </p>
        <div class="flex items-center gap-2">
          <button
            @click="$emit('page-change', currentPage - 1)"
            :disabled="currentPage === 1"
            class="px-3 py-1 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span class="text-dark-400 text-sm">
            Page {{ currentPage }} of {{ totalPages }}
          </span>
          <button
            @click="$emit('page-change', currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="px-3 py-1 bg-dark-700 hover:bg-dark-600 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, h } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  columns: {
    type: Array,
    required: true
  },
  data: {
    type: Array,
    default: () => []
  },
  rowKey: {
    type: String,
    default: 'id'
  },
  loading: {
    type: Boolean,
    default: false
  },
  emptyText: {
    type: String,
    default: 'No data found'
  },
  emptyIcon: {
    type: Object,
    default: () => ({
      render: () => h('svg', { class: 'w-12 h-12', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
        h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4' })
      ])
    })
  },
  pagination: {
    type: Boolean,
    default: false
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  total: {
    type: Number,
    default: 0
  },
  itemLabel: {
    type: String,
    default: 'items'
  }
})

defineEmits(['page-change'])

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))
const startItem = computed(() => (props.currentPage - 1) * props.pageSize + 1)
const endItem = computed(() => Math.min(props.currentPage * props.pageSize, props.total))
</script>