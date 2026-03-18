import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Input from './Input.vue'

describe('Input', () => {
  it('renders properly with default props', () => {
    const wrapper = mount(Input)

    expect(wrapper.find('input').exists()).toBe(true)
    expect(wrapper.find('input').attributes('type')).toBe('text')
  })

  it('renders with label', () => {
    const wrapper = mount(Input, {
      props: {
        label: 'Username',
        id: 'username'
      }
    })

    expect(wrapper.find('label').exists()).toBe(true)
    expect(wrapper.find('label').text()).toBe('Username')
    expect(wrapper.find('label').attributes('for')).toBe('username')
  })

  it('binds modelValue correctly', async () => {
    const wrapper = mount(Input, {
      props: {
        modelValue: 'test value'
      }
    })

    expect(wrapper.find('input').element.value).toBe('test value')
  })

  it('emits update:modelValue on input', async () => {
    const wrapper = mount(Input)

    await wrapper.find('input').setValue('new value')

    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual(['new value'])
  })

  it('renders with different types', () => {
    const types = ['text', 'password', 'email', 'number']

    types.forEach(type => {
      const wrapper = mount(Input, {
        props: { type }
      })

      expect(wrapper.find('input').attributes('type')).toBe(type)
    })
  })

  it('shows error styles when error prop is set', () => {
    const wrapper = mount(Input, {
      props: {
        error: 'This field is required'
      }
    })

    expect(wrapper.find('input').classes().join(' ')).toContain('border-red-500')
  })

  it('is disabled when disabled prop is true', () => {
    const wrapper = mount(Input, {
      props: { disabled: true }
    })

    expect(wrapper.find('input').attributes('disabled')).toBeDefined()
    expect(wrapper.find('input').classes().join(' ')).toContain('cursor-not-allowed')
  })

  it('is readonly when readonly prop is true', () => {
    const wrapper = mount(Input, {
      props: { readonly: true }
    })

    expect(wrapper.find('input').attributes('readonly')).toBeDefined()
  })

  it('renders with placeholder', () => {
    const wrapper = mount(Input, {
      props: { placeholder: 'Enter your name' }
    })

    expect(wrapper.find('input').attributes('placeholder')).toBe('Enter your name')
  })

  it('emits focus event', async () => {
    const wrapper = mount(Input)

    await wrapper.find('input').trigger('focus')

    expect(wrapper.emitted('focus')).toBeTruthy()
  })

  it('emits blur event', async () => {
    const wrapper = mount(Input)

    await wrapper.find('input').trigger('blur')

    expect(wrapper.emitted('blur')).toBeTruthy()
  })

  it('generates unique id if not provided', () => {
    const wrapper1 = mount(Input)
    const wrapper2 = mount(Input)

    const id1 = wrapper1.find('input').attributes('id')
    const id2 = wrapper2.find('input').attributes('id')

    expect(id1).toBeDefined()
    expect(id2).toBeDefined()
    expect(id1).not.toBe(id2)
  })

  it('accepts numeric modelValue', () => {
    const wrapper = mount(Input, {
      props: {
        modelValue: 42,
        type: 'number'
      }
    })

    expect(wrapper.find('input').element.value).toBe('42')
  })
})