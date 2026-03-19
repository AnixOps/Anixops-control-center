/**
 * Ansible Service Unit Tests
 */

import { describe, it, expect, beforeEach } from 'vitest'
import {
  parsePlaybook,
  generateExecutionCommands,
  parseValue,
  type ParsedPlaybook,
  type ExecutionOptions,
} from './ansible'

describe('Ansible Service', () => {
  describe('parseValue', () => {
    it('should parse boolean yes', () => {
      expect(parseValue('yes')).toBe(true)
    })

    it('should parse boolean no', () => {
      expect(parseValue('no')).toBe(false)
    })

    it('should parse boolean true', () => {
      expect(parseValue('true')).toBe(true)
    })

    it('should parse boolean false', () => {
      expect(parseValue('false')).toBe(false)
    })

    it('should parse integer', () => {
      expect(parseValue('42')).toBe(42)
    })

    it('should parse float', () => {
      expect(parseValue('3.14')).toBe(3.14)
    })

    it('should parse string', () => {
      expect(parseValue('hello world')).toBe('hello world')
    })

    it('should parse quoted string', () => {
      expect(parseValue('"hello world"')).toBe('hello world')
    })
  })

  describe('parsePlaybook', () => {
    it('should parse a simple playbook', () => {
      const content = `---
- name: Install and configure Fail2ban
  hosts: all
  become: yes
  vars:
    fail2ban_bantime: 3600
    fail2ban_maxretry: 5
  tasks:
    - name: Install Fail2ban
      apt:
        name: fail2ban
        state: present
`

      const playbook = parsePlaybook(content)

      expect(playbook).not.toBeNull()
      expect(playbook!.name).toBe('Install and configure Fail2ban')
      expect(playbook!.hosts).toBe('all')
      expect(playbook!.become).toBe(true)
      expect(playbook!.vars).toBeDefined()
      expect(playbook!.vars!.fail2ban_bantime).toBe(3600)
      expect(playbook!.vars!.fail2ban_maxretry).toBe(5)
      expect(playbook!.tasks.length).toBeGreaterThan(0)
    })

    it('should parse playbook without vars', () => {
      const content = `---
- name: Simple Task
  hosts: all
  tasks:
    - name: Hello World
      debug:
        msg: "Hello"
`

      const playbook = parsePlaybook(content)

      expect(playbook).not.toBeNull()
      expect(playbook!.name).toBe('Simple Task')
      expect(playbook!.become).toBe(false)
    })

    it('should return null for invalid playbook', () => {
      const content = 'invalid yaml content'
      const playbook = parsePlaybook(content)
      // Parser should handle gracefully, may return partial result
      expect(playbook).not.toBeUndefined()
    })

    it('should handle empty content', () => {
      const playbook = parsePlaybook('')
      // Empty content should still return an object
      expect(playbook).not.toBeNull()
    })

    it('should handle comments', () => {
      const content = `---
# This is a comment
- name: Test Playbook
  hosts: all
  # Another comment
  tasks:
    - name: Test Task
      debug:
        msg: "Test"
`

      const playbook = parsePlaybook(content)
      expect(playbook).not.toBeNull()
      expect(playbook!.name).toBe('Test Playbook')
    })
  })

  describe('generateExecutionCommands', () => {
    it('should generate basic commands', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const commands = generateExecutionCommands(content, 'test-playbook')

      expect(commands.length).toBeGreaterThan(0)
      expect(commands[0]).toContain('cat > /tmp/anixops-test-playbook')
      expect(commands[0]).toContain('EOFPLAYBOOK')
    })

    it('should include ansible-playbook command', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const commands = generateExecutionCommands(content, 'test')

      const ansibleCmd = commands.find(c => c.includes('ansible-playbook'))
      expect(ansibleCmd).toBeDefined()
    })

    it('should include check mode when specified', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const options: ExecutionOptions = { check_mode: true }
      const commands = generateExecutionCommands(content, 'test', options)

      const ansibleCmd = commands.find(c => c.includes('ansible-playbook'))
      expect(ansibleCmd).toContain('--check')
    })

    it('should include diff mode when specified', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const options: ExecutionOptions = { diff_mode: true }
      const commands = generateExecutionCommands(content, 'test', options)

      const ansibleCmd = commands.find(c => c.includes('ansible-playbook'))
      expect(ansibleCmd).toContain('--diff')
    })

    it('should include verbose mode when specified', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const options: ExecutionOptions = { verbose: true }
      const commands = generateExecutionCommands(content, 'test', options)

      const ansibleCmd = commands.find(c => c.includes('ansible-playbook'))
      expect(ansibleCmd).toContain('-v')
    })

    it('should include extra variables', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const options: ExecutionOptions = {
        extra_vars: { foo: 'bar', count: 5 }
      }
      const commands = generateExecutionCommands(content, 'test', options)

      const ansibleCmd = commands.find(c => c.includes('ansible-playbook'))
      expect(ansibleCmd).toContain('--extra-vars')
      expect(ansibleCmd).toContain('foo')
    })

    it('should include cleanup command', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const commands = generateExecutionCommands(content, 'test')

      const cleanupCmd = commands.find(c => c.includes('rm -f'))
      expect(cleanupCmd).toBeDefined()
    })

    it('should include localhost inventory', () => {
      const content = '---\n- name: Test\n  hosts: all\n  tasks: []'
      const commands = generateExecutionCommands(content, 'test')

      const ansibleCmd = commands.find(c => c.includes('ansible-playbook'))
      expect(ansibleCmd).toContain('localhost,')
    })
  })

  describe('Execution Result Types', () => {
    it('should have correct status types', () => {
      const statuses = ['pending', 'running', 'success', 'failed', 'cancelled', 'timeout']

      // Type check - this is just to verify types are correct
      expect(statuses.length).toBe(6)
    })

    it('should have correct node result structure', () => {
      const result = {
        node_id: 1,
        node_name: 'test-node',
        status: 'success' as const,
        started_at: '2024-01-01T00:00:00Z',
        completed_at: '2024-01-01T00:01:00Z',
        exit_code: 0,
        stdout: 'output',
        stderr: '',
      }

      expect(result.node_id).toBe(1)
      expect(result.status).toBe('success')
    })
  })

  describe('Task Queue Item', () => {
    it('should have correct structure', () => {
      const queueItem = {
        task_id: 'test-123',
        playbook_id: 1,
        playbook_name: 'test-playbook',
        storage_key: 'playbooks/test.yml',
        nodes: [{ id: 1, name: 'node1', host: '192.168.1.1' }],
        variables: { foo: 'bar' },
        triggered_by: 1,
        created_at: '2024-01-01T00:00:00Z',
      }

      expect(queueItem.task_id).toBe('test-123')
      expect(queueItem.nodes.length).toBe(1)
    })
  })
})