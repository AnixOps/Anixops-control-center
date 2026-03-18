import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Modal from './Modal.vue'

// Mock Teleport
vi.mock('vue', async () => {
  const actual = await vi.importActual('vue')
  return {
    ...actual,
    Teleport: {
      name: 'Teleport',
      props: ['to'],
      setup(_, { slots }) {
        return () => slots.default?.()
      }
    }
  }
})

describe('Modal', () => {
  it('renders when modelValue is true', () => {
    const wrapper = mount(Modal, {
      props: { modelValue: true },
      slots: { default: 'Modal content' }
    })

    expect(wrapper.text()).toContain('Modal content')
  })

  it('does not render when modelValue is false', () => {
    const wrapper = mount(Modal, {
      props: { modelValue: false },
      slots: { default: 'Modal content' }
    })

    expect(wrapper.text()).not.toContain('Modal content')
  })

  it('renders with title', () => {
    const wrapper = mount(Modal, {
      props: {
        modelValue: true,
        title: 'Modal Title'
      }
    })

    expect(wrapper.text()).toContain('Modal Title')
  })

  it('renders with different sizes', () => {
    const sizes = ['sm', 'md', 'lg', 'xl', 'full']

    sizes.forEach(size => {
      const wrapper = mount(Modal, {
        props: {
          modelValue: true,
          size
        }
      })

      const sizeClasses = {
        sm: 'max-w-sm',
        md: 'max-w-md',
        lg: 'max-w-lg',
        xl: 'max-w-xl',
        full: 'max-w-4xl'
      }

      expect(wrapper.find('.relative.bg-dark-800').classes().join(' ')).toContain(sizeClasses[size])
    })
  })

  it('emits close event when close button is clicked', async () => {
    const wrapper = mount(Modal, {
      props: { modelValue: true }
    })

    await wrapper.find('button').trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
  })

  it('emits close event when backdrop is clicked', async () => {
    const wrapper = mount(Modal, {
      props: {
        modelValue: true,
        closeOnBackdrop: true
      }
    })

    await wrapper.find('.absolute.inset-0').trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('does not emit close when backdrop is clicked and closeOnBackdrop is false', async () => {
    const wrapper = mount(Modal, {
      props: {
        modelValue: true,
        closeOnBackdrop: false
      }
    })

    await wrapper.find('.absolute.inset-0').trigger('click')

    expect(wrapper.emitted('close')).toBeFalsy()
  })

  it('hides close button when showClose is false', () => {
    const wrapper = mount(Modal, {
      props: {
        modelValue: true,
        showClose: false
      }
    })

    expect(wrapper.find('button').exists()).toBe(false)
  })

  it('renders header slot', () => {
    const wrapper = mount(Modal, {
      props: { modelValue: true },
      slots: {
        header: '<div class="custom-header">Custom Header</div>'
      }
    })

    expect(wrapper.find('.custom-header').exists()).toBe(true)
  })

  it('renders footer slot', () => {
    const wrapper = mount(Modal, {
      props: { modelValue: true },
      slots: {
        footer: '<div class="custom-footer">Footer</div>'
      }
    })

    expect(wrapper.find('.custom-footer').exists()).toBe(true)
  })

  it('renders default slot', () => {
    const wrapper = mount(Modal, {
      props: { modelValue: true },
      slots: {
        default: '<div class="modal-body">Content</div>'
      }
    })

    expect(wrapper.find('.modal-body').exists()).toBe(true)
  })
})