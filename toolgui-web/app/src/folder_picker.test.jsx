import React from 'react'
import { cleanup, fireEvent, waitFor } from '@testing-library/react'
import { afterEach, expect, test, vi } from 'vitest'
import { render } from './render'
import { Node } from '@toolgui-web/lib/src/app/Nodes'
import { FolderPickerContext } from '@toolgui-web/lib/src/app/FolderPicker'
import { TFolderPicker } from '@toolgui-web/lib/src/components/tcinput/folder_picker'

afterEach(cleanup)
const node = (disabled = false) => new Node('main/0', {
  name: 'folder_picker_component', id: 'folder', label: 'Output folder',
  default_dir: '/output', write: true, disabled,
})

test('web and wasm without a folder adapter show a disabled control and explanation', () => {
  const { getByRole, getByText } = render(<TFolderPicker node={node()} />)
  expect(getByRole('button')).toBeDisabled()
  expect(getByText('Folder selection is not available in this executor yet.')).toBeInTheDocument()
})

test('desktop calls the adapter with the component id and initial directory', async () => {
  let finish
  const pick = vi.fn().mockReturnValue(new Promise(resolve => { finish = resolve }))
  const update = vi.fn()
  const { getByRole } = render(
    <FolderPickerContext.Provider value={pick}>
      <TFolderPicker node={node()} update={update} />
    </FolderPickerContext.Provider>)
  const button = getByRole('button')
  fireEvent.click(button)
  expect(pick).toHaveBeenCalledWith('folder', '/output')
  expect(button).toBeDisabled()
  finish()
  await waitFor(() => expect(button).not.toBeDisabled())
  expect(update).not.toHaveBeenCalled()
})

test('disabled configuration prevents opening the desktop dialog', () => {
  const pick = vi.fn()
  const { getByRole } = render(
    <FolderPickerContext.Provider value={pick}>
      <TFolderPicker node={node(true)} />
    </FolderPickerContext.Provider>)
  fireEvent.click(getByRole('button'))
  expect(pick).not.toHaveBeenCalled()
})

test('dialog failures appear on the component without an unhandled rejection', async () => {
  const pick = vi.fn().mockRejectedValue(new Error('Dialog failed'))
  const { getByRole, getByText } = render(
    <FolderPickerContext.Provider value={pick}>
      <TFolderPicker node={node()} />
    </FolderPickerContext.Provider>)
  fireEvent.click(getByRole('button'))
  await waitFor(() => expect(getByRole('alert')).toHaveTextContent('Error: Dialog failed'))
  expect(getByRole('button')).not.toBeDisabled()
})
