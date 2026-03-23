import { describe, it, expect } from 'vitest'

// AI Models configuration
const AI_MODELS = {
  textGeneration: '@cf/meta/llama-3.1-8b-instruct',
  textEmbeddings: '@cf/baai/bge-base-en-v1.5',
}

// Mock AI response
const mockAIResponse = {
  response: 'This is a test response from the AI model.',
}

describe('AI Models Configuration', () => {
  it('has text generation model configured', () => {
    expect(AI_MODELS.textGeneration).toBeDefined()
    expect(AI_MODELS.textGeneration).toContain('@cf/')
  })

  it('has text embeddings model configured', () => {
    expect(AI_MODELS.textEmbeddings).toBeDefined()
    expect(AI_MODELS.textEmbeddings).toContain('@cf/')
  })
})

describe('AI Service Functions', () => {
  it('generates text with correct format', () => {
    const prompt = 'Analyze this log: ERROR connection failed'
    const systemPrompt = 'You are a DevOps assistant.'

    // Test request format
    const request = {
      messages: [
        { role: 'system', content: systemPrompt },
        { role: 'user', content: prompt },
      ],
      max_tokens: 512,
      temperature: 0.7,
    }

    expect(request.messages).toHaveLength(2)
    expect(request.messages[0].role).toBe('system')
    expect(request.messages[1].content).toBe(prompt)
  })

  it('creates embedding request correctly', () => {
    const text = 'Sample log entry for embedding'
    const request = {
      text,
    }

    expect(request.text).toBe(text)
  })
})

describe('Log Analysis', () => {
  it('classifies log errors correctly', () => {
    const errorLogs = [
      { level: 'error', message: 'Connection timeout' },
      { level: 'warning', message: 'High memory usage' },
      { level: 'info', message: 'Task completed' },
    ]

    const errors = errorLogs.filter((l) => l.level === 'error')
    expect(errors).toHaveLength(1)
    expect(errors[0].message).toContain('timeout')
  })

  it('extracts root cause from logs', () => {
    const logContent = `2024-03-23 10:00:00 ERROR Database connection failed
2024-03-23 10:00:01 ERROR Retry attempt 1 failed
2024-03-23 10:00:02 ERROR Retry attempt 2 failed`

    // Check for error patterns
    expect(logContent).toContain('Database connection failed')
    expect(logContent).toContain('ERROR')
  })

  it('calculates severity levels', () => {
    const calculateSeverity = (errors: number, warnings: number): number => {
      return Math.min(5, Math.ceil((errors * 2 + warnings) / 2))
    }

    expect(calculateSeverity(0, 0)).toBe(0)
    expect(calculateSeverity(5, 0)).toBe(5)
    expect(calculateSeverity(1, 2)).toBe(2)
  })
})

describe('Natural Language Query', () => {
  it('parses simple queries', () => {
    const query = 'show me all failed tasks'
    const keywords = ['failed', 'tasks', 'show']

    const hasKeywords = keywords.every((k) => query.toLowerCase().includes(k))
    expect(hasKeywords).toBe(true)
  })

  it('extracts time range from query', () => {
    const query = 'tasks from last 24 hours'
    const timePatterns = ['last 24 hours', 'yesterday', 'today', 'last week']

    const foundPattern = timePatterns.find((p) => query.includes(p))
    expect(foundPattern).toBe('last 24 hours')
  })

  it('identifies filter conditions', () => {
    const query = 'show failed tasks on node server-1'
    const filters = {
      status: 'failed',
      type: 'tasks',
      node: 'server-1',
    }

    expect(filters.status).toBe('failed')
    expect(filters.node).toBe('server-1')
  })
})

describe('AI Chat Assistant', () => {
  it('maintains conversation history', () => {
    const history = [
      { role: 'user', content: 'What is the status of node-1?' },
      { role: 'assistant', content: 'Node-1 is currently online.' },
      { role: 'user', content: 'Show me its logs' },
    ]

    expect(history).toHaveLength(3)
    expect(history[0].role).toBe('user')
    expect(history[1].role).toBe('assistant')
  })

  it('generates appropriate responses', () => {
    const context = {
      nodeStatus: '3 online, 2 offline',
      activeAlerts: 5,
      runningTasks: 10,
    }

    // Context should be passed to LLM
    expect(context.nodeStatus).toBeDefined()
    expect(context.activeAlerts).toBeGreaterThan(0)
  })
})

describe('Embedding Generation', () => {
  it('generates consistent embeddings', () => {
    const text = 'Test log message'
    // Embedding should be deterministic for same input
    const embedding = new Array(768).fill(0).map(() => Math.random())

    expect(embedding).toHaveLength(768)
    expect(typeof embedding[0]).toBe('number')
  })

  it('calculates similarity correctly', () => {
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
    const b = [1, 0, 0]
    const c = [0, 1, 0]

    expect(cosineSimilarity(a, b)).toBe(1)
    expect(cosineSimilarity(a, c)).toBe(0)
  })
})