import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Button from './Button.vue'

describe('Button', () => {
  it('renders properly with default props', () => {
    const wrapper = mount(Button, {
      slots: {
        default: 'Click me'
      }
    })

    expect(wrapper.text()).toContain('Click me')
    expect(wrapper.find('button').exists()).toBe(true)
    expect(wrapper.find('button').attributes('type')).toBe('button')
  })

  it('renders with different variants', () => {
    const variants = ['primary', 'secondary', 'danger', 'ghost', 'outline']

    variants.forEach(variant => {
      const wrapper = mount(Button, {
        props: { variant },
        slots: { default: 'Test' }
      })

      expect(wrapper.find('button').classes().some(c => c.includes(variant === 'primary' ? 'primary' : variant))).toBe(true)
    })
  })

  it('renders with different sizes', () => {
    const sizes = ['sm', 'md', 'lg']

    sizes.forEach(size => {
      const wrapper = mount(Button, {
        props: { size },
        slots: { default: 'Test' }
      })

      const sizeClasses = {
        sm: 'px-3 py-1.5',
        md: 'px-4 py-2',
        lg: 'px-6 py-3'
      }

      expect(wrapper.find('button').classes().join(' ')).toContain(sizeClasses[size])
    })
  })

  it('shows loading spinner when loading', () => {
    const wrapper = mount(Button, {
      props: { loading: true },
      slots: { default: 'Loading' }
    })

    expect(wrapper.find('svg').exists()).toBe(true)
    expect(wrapper.find('svg').classes()).toContain('animate-spin')
  })

  it('is disabled when disabled prop is true', () => {
    const wrapper = mount(Button, {
      props: { disabled: true },
      slots: { default: 'Disabled' }
    })

    expect(wrapper.find('button').attributes('disabled')).toBeDefined()
  })

  it('is disabled when loading', () => {
    const wrapper = mount(Button, {
      props: { loading: true },
      slots: { default: 'Loading' }
    })

    expect(wrapper.find('button').attributes('disabled')).toBeDefined()
  })

  it('emits click event when clicked', async () => {
    const wrapper = mount(Button, {
      slots: { default: 'Click me' }
    })

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })

  it('does not emit click when disabled', async () => {
    const wrapper = mount(Button, {
      props: { disabled: true },
      slots: { default: 'Disabled' }
    })

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('renders with custom type attribute', () => {
    const wrapper = mount(Button, {
      props: { type: 'submit' },
      slots: { default: 'Submit' }
    })

    expect(wrapper.find('button').attributes('type')).toBe('submit')
  })

  it('renders icon slot when not loading', () => {
    const wrapper = mount(Button, {
      slots: {
        default: 'With Icon',
        icon: '<span class="icon">+</span>'
      }
    })

    expect(wrapper.find('.icon').exists()).toBe(true)
  })

  it('hides icon slot when loading', () => {
    const wrapper = mount(Button, {
      props: { loading: true },
      slots: {
        default: 'Loading',
        icon: '<span class="icon">+</span>'
      }
    })

    expect(wrapper.find('.icon').exists()).toBe(false)
  })
})