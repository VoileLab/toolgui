export class Node {
  props: any
  children: Array<string>
  removing: boolean
  parentID: string

  constructor(props: any) {
    this.props = props
    this.children = []
    this.removing = false
    this.parentID = ''
  }
}

export class Forest {
  nodes: { [id: string]: Node }
  rootNodeID: string

  rootNode() {
    return new Node({
      name: 'container_component',
      id: this.rootNodeID,
    })
  }

  constructor(rootNodeID: string) {
    this.rootNodeID = rootNodeID
    this.nodes = {}
    this.nodes[rootNodeID] = this.rootNode()
  }

  swallowCopy(): Forest {
    const ret = new Forest(this.rootNodeID)
    ret.nodes = { ...this.nodes }
    return ret
  }

  setToRemoving() {
    for (const nodeID in this.nodes) {
      if (nodeID === this.rootNodeID) {
        continue
      }

      this.nodes[nodeID].removing = true
    }
  }

  createNode(props: any, parentID: string) {
    const nodeID: string = props.id

    if (nodeID in this.nodes && !this.nodes[nodeID].removing) {
      console.error('Depulicated component id:', nodeID)
      return
    }

    if (nodeID in this.nodes) {
      const parentID = this.nodes[nodeID].parentID
      const idx = this.nodes[parentID].children.indexOf(nodeID)
      this.nodes[parentID].children.splice(idx, 1)
    }

    const parentNode = this.nodes[parentID]
    var idx = 0
    for (var i = 0; i < parentNode.children.length; i++) {
      const prevNodeID = parentNode.children[i]
      if (!this.nodes[prevNodeID] || this.nodes[prevNodeID].removing) {
        break
      }
      idx = i + 1
    }
    parentNode.children.splice(idx, 0, nodeID)

    const oldNode = this.nodes[nodeID]
    if (oldNode) {
      oldNode.props = props
      oldNode.removing = false
    } else {
      this.nodes[nodeID] = new Node(props)
    }

    this.nodes[nodeID].parentID = parentID
  }

  updateNode(props: any) {
    const compID: string = props.id

    if (!(compID in this.nodes)) {
      console.error('Try to update a node that doesn\'t exist:', compID)
      return
    }

    this.nodes[compID].props = props
  }

  removeNode(nodeID: string) {
    if (!(nodeID in this.nodes)) {
      console.error('Try to remove a node that doesn\'t exist:', nodeID)
      return
    }

    const parentID = this.nodes[nodeID].parentID
    const idx = this.nodes[parentID].children.indexOf(nodeID)
    this.nodes[parentID].children.splice(idx, 1)
    delete this.nodes[nodeID]
  }

  removeNodeWithRemovingTag() {
    const removingID = new Set()
    for (const [key, node] of Object.entries(this.nodes)) {
      if (node.removing) {
        removingID.add(key)
      }
    }

    for (const id in removingID) {
      delete this.nodes[id]
    }

    for (const [key, node] of Object.entries(this.nodes)) {
      node.children = node.children.filter((nodeID) => {
        return !(removingID.has(nodeID))
      })
    }
  }
}