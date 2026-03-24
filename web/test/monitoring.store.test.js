import { describe, it, expect } from 'vitest'

// Mock monitoring data
const mockHealthChecks = [
  { name: 'API', status: 'healthy', latency: 12 },
  { name: 'Database', status: 'healthy', latency: 5 },
  { name: 'Cache', status: 'healthy', latency: 2 },
  { name: 'Queue', status: 'degraded', latency: 50 },
]

const mockMetrics = {
  requestRate: 1250,
  errorRate: 0.15,
  avgLatency: 45,
  p99Latency: 120,
}

const mockServices = [
  { name: 'api-gateway', health: 'healthy', requestRate: 500, errorRate: 0.1, latency: 25 },
  { name: 'auth-service', health: 'healthy', requestRate: 200, errorRate: 0.0, latency: 15 },
  { name: 'task-runner', health: 'healthy', requestRate: 100, errorRate: 0.5, latency: 50 },
]

const mockAlerts = [
  { id: '1', name: 'High Memory Usage', severity: 'warning', metric: 'memory_percent', value: 85, threshold: 80 },
  { id: '2', name: 'CPU Throttling', severity: 'critical', metric: 'cpu_throttle', value: 30, threshold: 20 },
]

describe('Health Checks', () => {
  it('filters healthy services', () => {
    const healthy = mockHealthChecks.filter(c => c.status === 'healthy')
    expect(healthy.length).toBe(3)
  })

  it('filters degraded services', () => {
    const degraded = mockHealthChecks.filter(c => c.status === 'degraded')
    expect(degraded.length).toBe(1)
  })

  it('calculates average latency', () => {
    const avgLatency = mockHealthChecks.reduce((sum, c) => sum + c.latency, 0) / mockHealthChecks.length
    expect(avgLatency).toBeCloseTo(17.25)
  })
})

describe('Metrics', () => {
  it('formats request rate correctly', () => {
    const formatNumber = (num) => {
      if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
      if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
      return num.toString()
    }

    expect(formatNumber(1250)).toBe('1.3K')
    expect(formatNumber(1500000)).toBe('1.5M')
    expect(formatNumber(500)).toBe('500')
  })

  it('calculates error percentage', () => {
    expect(mockMetrics.errorRate).toBeLessThan(1)
    expect(mockMetrics.errorRate).toBe(0.15)
  })

  it('compares latencies', () => {
    expect(mockMetrics.p99Latency).toBeGreaterThan(mockMetrics.avgLatency)
  })
})

describe('Services', () => {
  it('counts healthy services', () => {
    const healthyCount = mockServices.filter(s => s.health === 'healthy').length
    expect(healthyCount).toBe(3)
  })

  it('calculates total request rate', () => {
    const totalRequestRate = mockServices.reduce((sum, s) => sum + s.requestRate, 0)
    expect(totalRequestRate).toBe(800)
  })

  it('finds service with highest error rate', () => {
    const highestError = mockServices.reduce((max, s) =>
      s.errorRate > max.errorRate ? s : max
    )
    expect(highestError.name).toBe('task-runner')
  })

  it('filters services by name', () => {
    const filtered = mockServices.filter(s => s.name.includes('api'))
    expect(filtered.length).toBe(1)
    expect(filtered[0].name).toBe('api-gateway')
  })
})

describe('Alerts', () => {
  it('filters alerts by severity', () => {
    const warnings = mockAlerts.filter(a => a.severity === 'warning')
    const criticals = mockAlerts.filter(a => a.severity === 'critical')

    expect(warnings.length).toBe(1)
    expect(criticals.length).toBe(1)
  })

  it('checks if alert is firing', () => {
    const isFiring = (alert) => alert.value > alert.threshold

    expect(isFiring(mockAlerts[0])).toBe(true) // 85 > 80
    expect(isFiring(mockAlerts[1])).toBe(true) // 30 > 20
  })

  it('calculates alert duration', () => {
    const startedAt = new Date('2026-03-23T10:00:00Z')
    const now = new Date('2026-03-23T10:30:00Z')
    const durationMinutes = Math.floor((now - startedAt) / 60000)

    expect(durationMinutes).toBe(30)
  })
})

describe('Monitoring Dashboard', () => {
  it('computes overall health status', () => {
    const getOverallHealth = (checks) => {
      if (checks.every(c => c.status === 'healthy')) return 'healthy'
      if (checks.some(c => c.status === 'unhealthy')) return 'unhealthy'
      return 'degraded'
    }

    expect(getOverallHealth(mockHealthChecks)).toBe('degraded')
  })

  it('formats metrics display', () => {
    const formatMetrics = (metrics) => ({
      requestRate: `${(metrics.requestRate / 1000).toFixed(1)}K/s`,
      errorRate: `${metrics.errorRate.toFixed(2)}%`,
      latency: `${metrics.avgLatency}ms`,
    })

    const formatted = formatMetrics(mockMetrics)

    expect(formatted.requestRate).toBe('1.3K/s')
    expect(formatted.errorRate).toBe('0.15%')
    expect(formatted.latency).toBe('45ms')
  })
})