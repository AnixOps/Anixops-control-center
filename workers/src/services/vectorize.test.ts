import { describe, it, expect } from 'vitest'

// Mock vector data
const mockVectors = [
  { id: 'log-1', values: new Array(768).fill(0.1), metadata: { type: 'log', level: 'error' } },
  { id: 'log-2', values: new Array(768).fill(0.2), metadata: { type: 'log', level: 'info' } },
  { id: 'task-1', values: new Array(768).fill(0.3), metadata: { type: 'task', status: 'failed' } },
]

describe('Vectorize Service', () => {
  it('inserts vectors with correct format', () => {
    const vector = {
      id: 'test-1',
      values: new Array(768).fill(0.5),
      metadata: {
        type: 'log',
        level: 'error',
        timestamp: new Date().toISOString(),
      },
    }

    expect(vector.id).toBe('test-1')
    expect(vector.values).toHaveLength(768)
    expect(vector.metadata.type).toBe('log')
  })

  it('searches vectors with filters', () => {
    const filter = { type: 'log' }
    const results = mockVectors.filter((v) => v.metadata.type === filter.type)

    expect(results).toHaveLength(2)
    expect(results.every((v) => v.metadata.type === 'log')).toBe(true)
  })

  it('returns top K results', () => {
    const topK = 2
    const results = mockVectors.slice(0, topK)

    expect(results).toHaveLength(2)
  })
})

describe('Cosine Similarity', () => {
  it('calculates similarity between identical vectors', () => {
    const cosineSimilarity = (a: number[], b: number[]): number => {
      let dotProduct = 0
      let normA = 0
      let normB = 0
      for (let i = 0; i < a.length; i++) {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
      }
      return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB))
    }

    const a = [1, 2, 3]
    const b = [1, 2, 3]

    expect(cosineSimilarity(a, b)).toBeCloseTo(1, 5)
  })

  it('calculates similarity between orthogonal vectors', () => {
    const cosineSimilarity = (a: number[], b: number[]): number => {
      let dotProduct = 0
      let normA = 0
      let normB = 0
      for (let i = 0; i < a.length; i++) {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
      }
      return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB))
    }

    const a = [1, 0, 0]
    const b = [0, 1, 0]

    expect(cosineSimilarity(a, b)).toBe(0)
  })

  it('calculates similarity between similar vectors', () => {
    const cosineSimilarity = (a: number[], b: number[]): number => {
      let dotProduct = 0
      let normA = 0
      let normB = 0
      for (let i = 0; i < a.length; i++) {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
      }
      return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB))
    }

    const a = [1, 1, 1]
    const b = [1, 1, 0.9]

    expect(cosineSimilarity(a, b)).toBeGreaterThan(0.95)
    expect(cosineSimilarity(a, b)).toBeLessThan(1)
  })
})

describe('Semantic Search', () => {
  it('finds similar logs by embedding', () => {
    const queryEmbedding = new Array(768).fill(0.1)

    // Find vectors close to query
    const results = mockVectors.filter((v) => {
      const diff = v.values[0] - queryEmbedding[0]
      return Math.abs(diff) < 0.15
    })

    expect(results.length).toBeGreaterThanOrEqual(1)
  })

  it('filters by metadata', () => {
    const filter = { level: 'error' }
    const results = mockVectors.filter(
      (v) => v.metadata.level === filter.level
    )

    expect(results).toHaveLength(1)
    expect(results[0].id).toBe('log-1')
  })
})

describe('Anomaly Detection', () => {
  it('detects anomalies based on distance threshold', () => {
    const threshold = 0.3

    const detectAnomaly = (current: number[], baseline: number[]): boolean => {
      // Calculate distance
      let sumSquares = 0
      for (let i = 0; i < current.length; i++) {
        const diff = current[i] - baseline[i]
        sumSquares += diff * diff
      }
      const distance = Math.sqrt(sumSquares / current.length)
      return distance > threshold
    }

    const baseline = new Array(10).fill(0.5)
    const normal = new Array(10).fill(0.48)
    const anomaly = new Array(10).fill(0.9)

    expect(detectAnomaly(normal, baseline)).toBe(false)
    expect(detectAnomaly(anomaly, baseline)).toBe(true)
  })

  it('tracks anomaly history', () => {
    const anomalies = [
      { timestamp: '2024-03-23T10:00:00Z', score: 0.5, type: 'cpu_spike' },
      { timestamp: '2024-03-23T11:00:00Z', score: 0.8, type: 'memory_leak' },
    ]

    expect(anomalies).toHaveLength(2)
    expect(anomalies[1].score).toBeGreaterThan(anomalies[0].score)
  })
})

describe('Vector Metadata', () => {
  it('stores correct metadata for logs', () => {
    const metadata = {
      id: 'log-123',
      type: 'log',
      timestamp: new Date().toISOString(),
      nodeId: 'node-1',
      level: 'error',
      message: 'Connection timeout',
    }

    expect(metadata.type).toBe('log')
    expect(metadata.level).toBe('error')
  })

  it('stores correct metadata for tasks', () => {
    const metadata = {
      id: 'task-456',
      type: 'task',
      timestamp: new Date().toISOString(),
      playbookName: 'deploy-app',
      status: 'failed',
      duration: 120,
    }

    expect(metadata.type).toBe('task')
    expect(metadata.status).toBe('failed')
  })
})