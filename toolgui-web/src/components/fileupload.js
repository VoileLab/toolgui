import { useState } from "react"
import { fileToBase64 } from "../util/base64";
import { sessionValues } from "./session"

export function TFileupload({ node, update }) {
  const [file, setFile] = useState(sessionValues[node.props.id]);

  const handleFileChange = async e => {
    e.preventDefault();

    if (!e.target.files) {
      return
    }

    const file = e.target.files[0]

    if (file.size >= 1024 * 100) {
      // TBD: use upload file method to make larger file work.
      console.error('File size:', file.size)
      return
    }

    const newFile = {
      name: file.name,
      type: file.type,
      size: file.size,
      body: await fileToBase64(file),
    }

    sessionValues[e.target.id] = newFile
    setFile(newFile)
    update({
      id: e.target.id,
      value: newFile,
    })
  };

  return (
    <div class="file has-name">
      <label class="file-label">
        <input class="file-input" type="file"
          id={node.props.id}
          name={node.props.id}
          onChange={handleFileChange} />
        <span class="file-cta">
          <span class="file-icon">
            <i class="fas fa-upload"></i>
          </span>
          <span class="file-label"> {node.props.label} </span>
        </span>
        <span class="file-name"> {file ? file.name : 'No file uploaded'} </span>
      </label>
    </div>
  )
}
