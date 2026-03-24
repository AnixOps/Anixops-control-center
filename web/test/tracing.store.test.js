import { describe, it, expect } from 'vitest'

// Mock tracing data
const mockTrace = {
  traceId: '0af7651916cd43dd8448eb211c80319c',
  status: 'ok',
  duration: 150,
  spanCount: 3,
  serviceCount: 1,
  rootSpan: { name: 'HTTP GET /api/users' },
  spans: [
    { spanId: '1', name: 'HTTP GET /api/users', duration: 150, kind: 'server', status: { code: 'ok' }, resource: { 'service.name': 'api-gateway' } },
    { spanId: '2', name: 'auth:validate_token', duration: 20, kind: 'client', status: { code: 'ok' }, resource: { 'service.name': 'api-gateway' } },
    { spanId: '3', name: 'db:query', duration: 50, kind: 'client', status: { code: 'ok' }, resource: { 'service.name': 'api-gateway' } },
  ],
}

const mockTraces = [
  mockTrace,
  {
    traceId: '1bf7651916cd43dd8448eb211c80319d',
    status: 'error',
    duration: 300,
    spanCount: 2,
    serviceCount: 2,
    rootSpan: { name: 'HTTP POST /api/tasks' },
    spans: [
      { spanId: '1', name: 'HTTP POST /api/tasks', duration: 300, kind: 'server', status: { code: 'error', message: 'Timeout' }, resource: { 'service.name': 'api-gateway' } },
      { spanId: '2', name: 'task:execute', duration: 280, kind: 'client', status: { code: 'error', message: 'Timeout' }, resource: { 'service.name': 'task-runner' } },
    ],
  },
]

const mockStats = {
  totalTraces: 1234,
  totalSpans: 5678,
  averageDuration: 145,
  errorRate: 0.02,
}

const mockServices = ['api-gateway', 'auth-service', 'task-runner', 'log-processor']

describe('Trace Data', () => {
  it('has valid trace ID format', () => {
    expect(mockTrace.traceId).toHaveLength(32)
    expect(/^[0-9a-f]{32}$/i.test(mockTrace.traceId)).toBe(true)
  })

  it('has valid status', () => {
    expect(['ok', 'error', 'unset']).toContain(mockTrace.status)
  })

  it('has spans array', () => {
    expect(Array.isArray(mockTrace.spans)).toBe(true)
    expect(mockTrace.spans.length).toBeGreaterThan(0)
  })

  it('calculates span count correctly', () => {
    expect(mockTrace.spanCount).toBe(mockTrace.spans.length)
  })
})

describe('Span Data', () => {
  it('has required span fields', () => {
    const span = mockTrace.spans[0]
    expect(span.spanId).toBeDefined()
    expect(span.name).toBeDefined()
    expect(span.duration).toBeDefined()
    expect(span.kind).toBeDefined()
    expect(span.status).toBeDefined()
  })

  it('has valid span kinds', () => {
    const validKinds = ['unspecified', 'internal', 'server', 'client', 'producer', 'consumer']
    mockTrace.spans.forEach(span => {
      expect(validKinds).toContain(span.kind)
    })
  })

  it('has valid status codes', () => {
    const validCodes = ['ok', 'error', 'unset']
    mockTrace.spans.forEach(span => {
      expect(validCodes).toContain(span.status.code)
    })
  })

  it('has service name in resource', () => {
    mockTrace.spans.forEach(span => {
      expect(span.resource['service.name']).toBeDefined()
    })
  })
})

describe('Trace Filtering', () => {
  it('filters by status', () => {
    const okTraces = mockTraces.filter(t => t.status === 'ok')
    const errorTraces = mockTraces.filter(t => t.status === 'error')

    expect(okTraces.length).toBe(1)
    expect(errorTraces.length).toBe(1)
  })

  it('filters by service', () => {
    const serviceFilter = 'task-runner'
    const filtered = mockTraces.filter(t =>
      t.spans.some(s => s.resource['service.name'] === serviceFilter)
    )

    expect(filtered.length).toBe(1)
    expect(filtered[0].traceId).toBe('1bf7651916cd43dd8448eb211c80319d')
  })

  it('filters by duration', () => {
    const minDuration = 200
    const filtered = mockTraces.filter(t => t.duration >= minDuration)

    expect(filtered.length).toBe(1)
    expect(filtered[0].duration).toBe(300)
  })
})

describe('Trace Statistics', () => {
  it('calculates average duration', () => {
    const avgDuration = mockTraces.reduce((sum, t) => sum + t.duration, 0) / mockTraces.length
    expect(avgDuration).toBe(225)
  })

  it('calculates error rate', () => {
    const errorCount = mockTraces.filter(t => t.status === 'error').length
    const errorRate = errorCount / mockTraces.length
    expect(errorRate).toBe(0.5)
  })

  it('calculates total spans', () => {
    const totalSpans = mockTraces.reduce((sum, t) => sum + t.spanCount, 0)
    expect(totalSpans).toBe(5)
  })

  it('calculates service count', () => {
    const allServices = new Set()
    mockTraces.forEach(t => {
      t.spans.forEach(s => {
        allServices.add(s.resource['service.name'])
      })
    })
    expect(allServices.size).toBe(2)
  })
})

describe('Service List', () => {
  it('has unique services', () => {
    const uniqueServices = [...new Set(mockServices)]
    expect(uniqueServices.length).toBe(mockServices.length)
  })

  it('filters services by name', () => {
    const filtered = mockServices.filter(s => s.includes('api'))
    expect(filtered.length).toBe(1)
    expect(filtered[0]).toBe('api-gateway')
  })
})

describe('Stats Display', () => {
  it('formats error rate as percentage', () => {
    const percentage = (mockStats.errorRate * 100).toFixed(1) + '%'
    expect(percentage).toBe('2.0%')
  })

  it('formats average duration', () => {
    const formatted = mockStats.averageDuration.toFixed(0) + 'ms'
    expect(formatted).toBe('145ms')
  })

  it('formats total traces', () => {
    const formatted = mockStats.totalTraces.toLocaleString()
    expect(formatted).toBe('1,234')
  })
})

describe('Error Span Analysis', () => {
  it('finds error spans', () => {
    const errorSpans = mockTraces[1].spans.filter(s => s.status.code === 'error')
    expect(errorSpans.length).toBe(2)
  })

  it('gets error messages', () => {
    const errorSpans = mockTraces[1].spans.filter(s => s.status.code === 'error')
    errorSpans.forEach(span => {
      expect(span.status.message).toBeDefined()
    })
  })

  it('identifies root cause span', () => {
    const errorTrace = mockTraces[1]
    const rootCause = errorTrace.spans.reduce((max, s) =>
      s.duration > max.duration ? s : max
    )
    expect(rootCause.name).toBe('HTTP POST /api/tasks')
  })
})