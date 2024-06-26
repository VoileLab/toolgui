export class Node {
  constructor(props) {
    this.props = props
    this.children = []
    this.removing = false
    this.parentID = ''
  }
}
