<template>
  <div class="relative">
    <canvas ref="chartCanvas" :height="height"></canvas>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  height: {
    type: Number,
    default: 200
  },
  label: {
    type: String,
    default: 'Traffic'
  },
  color: {
    type: String,
    default: '#3b82f6'
  }
})

const chartCanvas = ref(null)
let chartInstance = null

function createChart() {
  if (!chartCanvas.value || !props.data.length) return

  const labels = props.data.map(d => d.label || d.time || '')
  const values = props.data.map(d => d.value || d.count || 0)

  if (chartInstance) {
    chartInstance.destroy()
  }

  chartInstance = new ChartJS(chartCanvas.value, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: props.label,
        data: values,
        borderColor: props.color,
        backgroundColor: `${props.color}20`,
        fill: true,
        tension: 0.4,
        pointRadius: 0,
        pointHoverRadius: 4
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        tooltip: {
          backgroundColor: '#1f2937',
          titleColor: '#fff',
          bodyColor: '#9ca3af',
          borderColor: '#374151',
          borderWidth: 1,
          padding: 10,
          displayColors: false
        }
      },
      scales: {
        x: {
          grid: {
            color: '#374151',
            drawBorder: false
          },
          ticks: {
            color: '#6b7280',
            maxTicksLimit: 8
          }
        },
        y: {
          grid: {
            color: '#374151',
            drawBorder: false
          },
          ticks: {
            color: '#6b7280'
          }
        }
      },
      interaction: {
        intersect: false,
        mode: 'index'
      }
    }
  })
}

watch(() => props.data, () => {
  createChart()
}, { deep: true })

onMounted(() => {
  createChart()
})

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.destroy()
  }
})
</script>