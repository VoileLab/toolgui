import React from "react"
import { Props } from "../component_interface"
import { Editor } from "@monaco-editor/react"

export function TMonacoEditor({ node, update }: Props) {

  return (
    <div>
      <Editor
        height="100vh"
        language="javascript"
        value={"alert('Hello, world!');"}
      />
    </div>
  )
}

