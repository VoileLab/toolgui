// @vitest-environment jsdom
import React from 'react'
import { cleanup, fireEvent, render } from '@testing-library/react'
import { afterEach, expect, test, vi } from 'vitest'
import { TSelect } from './select'
import { clearState } from '../state'
import { Node } from '../../app/Nodes'

// Use a native control to isolate TSelect's state lifecycle from Mantine's
// dropdown layout and browser APIs. The browser test covers the real dropdown.
vi.mock('@mantine/core', () => ({
  Select: ({ value, onChange, data }: any) => React.createElement('select', {
    value: value ?? '',
    onChange: (event: any) => onChange(event.target.value || null),
  }, React.createElement('option', { value: '' }, 'None'),
  ...data.map((item: any) => React.createElement('option', {
    key: item.value, value: item.value,
  }, item.label))),
}))

afterEach(() => { cleanup(); clearState() })

function props(defaultValue: number, update = vi.fn()) {
  return {
    node: new Node('select', { id: 'select', items: ['Value1', 'Value2'], default: defaultValue }),
    update, upload: vi.fn(), download: vi.fn(), theme: 'light' as any,
  }
}

test('an untouched mounted select follows changed defaults', () => {
  const view = render(React.createElement(TSelect, props(1)))
  const field = view.getByRole('combobox') as HTMLSelectElement
  expect(field.value).toBe('1')
  view.rerender(React.createElement(TSelect, props(2)))
  expect(field.value).toBe('2')
})

test('selection updates immediately and session reset restores current default', () => {
  const update = vi.fn()
  const view = render(React.createElement(TSelect, props(1, update)))
  const field = view.getByRole('combobox') as HTMLSelectElement
  fireEvent.change(field, { target: { value: '2' } })
  expect(field.value).toBe('2')
  expect(update).toHaveBeenLastCalledWith({ type: 'select', id: 'select', value: 2 })
  view.rerender(React.createElement(TSelect, props(1, update)))
  expect(field.value).toBe('2')
  clearState()
  view.rerender(React.createElement(TSelect, props(1, update)))
  expect(field.value).toBe('1')
})
