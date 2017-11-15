# ntreego Golang N-ary Tree Implementation.  Roughly ported from glib's GNode with idiomatic go
Not Production Ready.  Pull Requests Welcomed!  (Need better tests!)


```go
root := ntree.New("A")
branchB := ntree.New("B")
leafC := ntree.New("C")
leafD := ntree.New("D")

root.AppendChild(branchB)
root.AppendChild(leafD)
branchB.AppendChild(leafC)

fmt.Println(root)
```


## ntree.Traverse is done through a few options, mainly the TraverseType, TraverseFlags, and Depth as per GNode's documentation
TraverseTypes:
- TraversePreOrder, visits a node, then its children.
- TraverseInOrder, vists a node's left child first, then the node itself then its right child. This is the one to use if you want the output sorted according to the compare function.
- TraversePostOrder, visits the node's children, then the node itself.
- TraverseLevelOrder (not implemented)

TraverseFlags: 
- TraverseAll
- TraverseLeaves
- TraverseNonLeaves

[Depth] is -1 to start at the root and 1->n for specified depths.  0 is root so is an invalid input

[TraverseFunc] is the function applied to each node that is traversed

[Data] is anything that should be passed to each node, allowing for lambda's to capture functionality outside of the library's traversal functions.


```go
root := ntree.New(1)
branchB := ntree.New(2)
leafC := ntree.New(3)
leafD := ntree.New(4)

root.AppendChild(branchB)
root.AppendChild(leafD)
branchB.AppendChild(leafC)

fmt.Println(root)

var matchedNode *ntree.Node
traverseFunc := func(node *Node, value interface{}) bool {
  v := value.(int)
  if v == 3 {
    matchedNode = node
  }
  node.Value.(*int).Value += v
}

data := 1
ntree.Traverse(root, ntree.TraversePreOrder, ntree.TraverseAll, -1, traverseFunc, data)

fmt.Println(matchedNode.Value) //should be 4, since it matched to the node with value 3 and incremented by 1
```
