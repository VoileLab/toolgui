import { createContext } from "react"

export type PickFolderFunc = (componentID: string, defaultDir: string, write: boolean) => Promise<void>
export const FolderPickerContext = createContext<PickFolderFunc | undefined>(undefined)
