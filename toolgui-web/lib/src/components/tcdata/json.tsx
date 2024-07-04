import React, { useState } from "react"

import { Props } from "../component_interface"

import '@toolgui-web/lib/src/assets/css/json.css'


export function TJson({ node }: Props) {
  return (
    <div className="block">
      <pre>
        <JsonValue val={JSON.parse(node.props.value)} pad={0} />
      </pre>
    </div>
  )
}

function JsonValue({ val, pad }: { val: any, pad: number }) {
  if (val === null) {
    return JsonNull({ val })
  }

  if (Array.isArray(val)) {
    return JsonList({ val, pad })
  } else if (typeof val === 'object') {
    return JsonDict({ val, pad })
  } else if (typeof val === 'string') {
    return JsonString({ val })
  } else if (typeof val === 'boolean') {
    return JsonBool({ val })
  } else {
    return JsonElement({ val })
  }
}

function JsonDict({ val, pad }: { val: any, pad: number }) {
  const [open, setOpen] = useState(true)
  if (!open) {
    return <p style={{ display: 'inline' }}>
      <span className="pseudolink" onClick={(e) => { setOpen(true) }}>
        <b>{"{ ... }"}</b>
      </span>
    </p>
  }

  const kvs = []
  for (const [key, value] of Object.entries(val)) {
    kvs.push(<div>{' '.repeat(pad + 2)}"{key}": <JsonValue val={value} pad={pad + 2} />,</div>)
  }
  return (
    <p style={{ display: 'inline' }}>
      <span className="pseudolink" onClick={(e) => { setOpen(false) }}>
        <b>{"{"}</b>
      </span>
      <div>{kvs}</div>
      <span className="pseudolink" onClick={(e) => { setOpen(false) }}>
        <b>{' '.repeat(pad)}{"}"}</b>
      </span>
    </p>
  )
}

function JsonList({ val, pad }: { val: { [key: string]: any }, pad: number }) {
  const [open, setOpen] = useState(true)
  if (!open) {
    return <p style={{ display: 'inline' }}>
      <span className="pseudolink" onClick={(e) => { setOpen(true) }}>
        <b>{"[ ... ]"}</b>
      </span>
    </p>
  }

  return (
    <p style={{ display: 'inline' }}>
      <span className="pseudolink" onClick={(e) => { setOpen(false) }}>
        <b>{"["}</b>
      </span>
      <div>{
        val.map((value: any) => (
          <div>{' '.repeat(pad + 2)}<JsonValue val={value} pad={pad + 2} />,</div>
        ))
      }</div>
      <span className="pseudolink" onClick={(e) => { setOpen(false) }}>
        <b>{' '.repeat(pad)}{"]"}</b>
      </span>
    </p>
  )
}

function JsonString({ val }: { val: string }) {
  return <p style={{ display: 'inline', color: 'green' }} >"{val}"</p>
}

function JsonBool({ val }: { val: boolean }) {
  if (val) {
    return <p style={{ display: 'inline', color: 'brown' }} >true</p>
  } else {
    return <p style={{ display: 'inline', color: 'brown' }} >false</p>
  }
}

function JsonNull({ val }: { val: null }) {
  return <p style={{ display: 'inline', color: 'brown' }} >null</p>
}

function JsonElement({ val }: { val: any }) {
  return <p style={{ display: 'inline', color: 'brown' }}>{val}</p>
}
