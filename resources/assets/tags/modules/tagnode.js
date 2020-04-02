
const pathSep = '/';

class TagNode {
  constructor(path, tag, parent) {
    this.path = path;
    this.tag = tag;
    this.parent = parent;
    this.children = [];
  }

  get tagID() {
    // Generate a random ID for synthetic nodes. This is used as the 'key' in the Vue template for loop so it needs to
    // be unique.
    if (this.tag === null) {
      let randInt = Math.floor((Math.random() * Number.MAX_SAFE_INTEGER));
      return `generated-${randInt}`;
    }
    return this.tag.tagID;
  }

  get depth() {
    return this.path.length - 1;
  }

  get name() {
    if (this.path.length === 0) {
      return 'root';
    }
    return this.path[this.path.length - 1]
  }

  get totalNumLineItems() {
    return this._getTotalCount(this, 'numLineItems');
  }

  get totalNumFingerprints() {
    return this._getTotalCount(this, 'numFingerprints');
  }

  _getTotalCount(tree, colName) {
    let sum = tree.tag === null ? 0 : tree.tag[colName];
    for (const child of tree.children) {
      sum += this._getTotalCount(child, colName)
    }
    return sum
  }

  insertTagNode(tag) {
    function getChild(node, tagPath) {
      for (let child of node.children) {
        if (child.path.join(pathSep) === tagPath.slice(0, node.depth + 2).join(pathSep)) {
          return child;
        }
      }
      return null;
    }

    let tagPath = `${pathSep}${tag.tag}`.split(pathSep);
    if (tagPath[this.depth] !== this.path[this.depth]) {
      return;
    }

    if (tagPath[tagPath.length - 2] === this.name) {
      let child = getChild(this, tagPath);
      if (child === null) {
        child = new TagNode(tagPath, tag, this);
        this.children.push(child);
      } else {
        child.tag_ = tag;
        if (child.path !== tagPath) {
          console.log("WARNING: paths not correct in node insertion")
        }
      }
      if (child.depth !== this.depth + 1) {
        console.log("WARNING: depth not correct in node insertion")
      }
      return
    } else {
      if (getChild(this, tagPath) === null) {
        let node = new TagNode(tagPath.slice(0, this.depth + 2), null, this);
        this.children.push(node);
      }
    }

    for (let child of this.children) {
      child.insertTagNode(tag)
    }
  }

  // Useful for debugging
  printTree() {
    this._printTree(this, 0);
  }

  _printTree(tree) {
    let indent = '';
    for (let i = 0; i < tree.depth; i++) {
      indent += '  ';
    }
    console.log(`${indent}${tree.name} ${tree.depth}`);
    for (const child of tree.children) {
      this._printTree(child);
    }
  }
}


export default TagNode;
