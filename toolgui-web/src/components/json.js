import { useState } from "react";

export function TJson({ node }) {
  return (
    <div class="block">
      <pre>
        <JsonValue val={JSON.parse(node.props.value)} pad={0} />
      </pre>
    </div>
  )
}

function JsonValue({ val, pad }) {
  if (val === null) {
    return JsonNull({ val })
  }

  if (Array.isArray(val)) {
    return JsonList({ val, pad })
  } else if (typeof val === "object") {
    return JsonDict({ val, pad })
  } else if (typeof val === 'string'){
    return JsonString({ val })
  } else if (typeof val === 'boolean'){
    return JsonBool({ val })
  } else {
    return JsonElement({ val })
  }
}

function JsonDict({ val, pad }) {
  const [open, setOpen] = useState(true)
  if (!open) {
    return <p style={{ display: 'inline' }}>
      <a onClick={(e) => {setOpen(true)}}>
        <b>{"{ ... }"}</b>
      </a>
    </p>
  }

  const kvs = []
  for (const [key, value] of Object.entries(val)) {
    kvs.push(<div>{' '.repeat(pad + 2)}"{key}": <JsonValue val={value} pad={pad + 2}/>,</div>)
  }
  return (
    <p style={{ display: 'inline' }}>
      <a onClick={(e) => {setOpen(false)}}>
        <b>{"{"}</b>
      </a>
      <div>{kvs}</div>
      <a onClick={(e) => {setOpen(false)}}>
        <b>{' '.repeat(pad)}{"}"}</b>
      </a>
    </p>
  )
}

function JsonList({ val, pad }) {
  const [open, setOpen] = useState(true)
  if (!open) {
    return <p style={{ display: 'inline' }}>
      <a onClick={(e) => {setOpen(true)}}>
        <b>{"[ ... ]"}</b>
      </a>
    </p>
  }

  return (
    <p style={{ display: 'inline' }}>
      <a onClick={(e) => {setOpen(false)}}>
        <b>{"["}</b>
      </a>
      <div>{
        val.map(value => (
          <div>{' '.repeat(pad + 2)}<JsonValue val={value} pad={pad + 2} />,</div>
        ))
      }</div>
      <a onClick={(e) => {setOpen(false)}}>
        <b>{' '.repeat(pad)}{"]"}</b>
      </a>
    </p>
  )
}

function JsonString({ val }) {
  return <p style={{ display: 'inline', color: 'green' }} >"{val}"</p>
}

function JsonBool({ val }) {
  if (val) {
    return <p style={{ display: 'inline', color: 'brown' }} >true</p>
  } else {
    return <p style={{ display: 'inline', color: 'brown' }} >false</p>
  }
}

function JsonNull({ val }) {
  return <p style={{ display: 'inline', color: 'brown' }} >null</p>
}

function JsonElement({ val }) {
  return <p style={{ display: 'inline', color: 'brown' }}>{val}</p>
}
