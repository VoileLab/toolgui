import React, { useContext, useState } from "react"
import { Button, Text } from "@mantine/core"
import { FolderPickerContext } from "../../app/FolderPicker"
import { Props } from "../component_interface"

export function TFolderPicker({ node }: Props) {
  const pick = useContext(FolderPickerContext)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)
  return <div>
    <Button id={node.props.id} disabled={!pick || node.props.disabled || busy}
      onClick={async () => {
        if (!pick) return
        setBusy(true)
        setError(null)
        try { await pick(node.props.id, node.props.default_dir ?? "") }
        catch (e) { setError(String(e)) }
        finally { setBusy(false) }
      }}>{node.props.label}</Button>
    {!pick && <Text size="xs">Folder selection is not available in this executor yet.</Text>}
    {error && <Text size="xs" c="red">{error}</Text>}
  </div>
}
