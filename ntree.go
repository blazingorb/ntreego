// Package ntree provides an idiomatic golang port of glib(Gnome)'s GNode N-ary tree https://developer.gnome.org/glib/stable/glib-N-ary-Trees.html
package ntree

import (
	"fmt"
)

const (
	TraverseInOrder TraverseType = iota
	TraversePreOrder
	TraversePostOrder
	TraverseLevelOrder
)

const (
	TraverseLeaves TraverseFlags = 1 << iota
	TraverseNonLeaves
	TraverseMask = 0x3
	TraverseAll  = TraverseLeaves | TraverseNonLeaves
)

type TraverseFunc func(*Node, interface{}) bool
type TraverseType int
type TraverseFlags int

type Node struct {
	Value    interface{}
	Next     *Node
	Previous *Node
	Parent   *Node
	Children *Node
}

type nodeVal struct {
	Value         interface{}
	NodeReference *Node
}

func New(v interface{}) *Node {
	return &Node{Value: v}
}

func Unlink(n *Node) {
	if n == nil {
		return
	}

	if n.Previous != nil {
		n.Previous.Next = n.Next
	} else if n.Parent != nil {
		n.Parent.Children = n.Next
	}

	n.Parent = nil
	if n.Next != nil {
		n.Next.Previous = n.Previous
		n.Next = nil
	}
	n.Previous = nil
}

func Depth(n *Node) int {
	depth := 0

	for n != nil {
		depth++
		n = n.Parent
	}

	return depth
}

func Insert(parent, n *Node) *Node {
	if parent == nil || n == nil || !IsRoot(n) {
		return nil
	}

	return AppendChild(parent, n)
}

func IsRoot(n *Node) bool {
	return n.Parent == nil && n.Previous == nil && n.Next == nil
}

func NodeCount(root *Node, flags TraverseFlags) int {
	if root == nil || flags > TraverseMask {
		return 0
	}

	n := 0
	nodeCountFunc(root, flags, &n)
	return n
}

func nodeCountFunc(n *Node, flags TraverseFlags, count *int) {
	if n.Children != nil {
		if flags&TraverseNonLeaves != 0 {
			(*count)++
		}

		child := n.Children
		for child != nil {
			nodeCountFunc(child, flags, count)
			child = child.Next
		}
	} else if flags&TraverseLeaves != 0 {
		(*count)++
	}
}

func AppendChild(parent, n *Node) *Node {
	if parent == nil || n == nil || !IsRoot(n) {
		return nil
	}

	n.Parent = parent
	if parent.Children != nil {
		sibling := parent.Children
		for sibling.Next != nil {
			sibling = sibling.Next
		}
		n.Previous = sibling
		sibling.Next = n
	} else {
		n.Parent.Children = n
	}

	return n
}

func GetRoot(n *Node) (*Node, int) {
	if n == nil {
		return nil, 0
	}

	depth := 0

	current := n
	for current.Parent != nil {
		depth++
		current = current.Parent
	}

	return current, depth
}

func FindNode(root *Node, order TraverseType, flags TraverseFlags, data interface{}) *Node {
	if root == nil {
		return nil
	}

	d := &nodeVal{Value: data}
	Traverse(root, order, flags, -1, nodeFindFunc, d)

	return d.NodeReference
}

func nodeFindFunc(n *Node, data interface{}) bool {
	nPtr := data.(*nodeVal)
	if n.Value != nPtr.Value {
		return false
	}

	nPtr.NodeReference = n
	return true
}

// Traverse traverses a node as the root node based on the passed in TraverseType, TraverseFlags, and Depth.  Each visited node will
// have TraverseFunc called, passing along Data to each node.
func Traverse(root *Node, order TraverseType, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data interface{}) {
	if root == nil || traverseFunc == nil || order > TraverseLevelOrder || flags > TraverseMask || (depth < -1 || depth == 0) {
		return
	}

	switch order {
	default:
		fallthrough
	case TraversePreOrder:
		if depth < 0 {
			traversePreOrder(root, flags, traverseFunc, data)
		} else {
			depthTraversePreOrder(root, flags, depth, traverseFunc, data)
		}
	case TraverseInOrder:
		if depth < 0 {
			traverseInOrder(root, flags, traverseFunc, data)
		} else {
			depthTraverseInOrder(root, flags, depth, traverseFunc, data)
		}
	case TraversePostOrder:
		if depth < 0 {
			traversePostOrder(root, flags, traverseFunc, data)
		} else {
			depthTraversePostOrder(root, flags, depth, traverseFunc, data)
		}

		// case Traverse_LevelOrder:
		// 	panic("Not Implemented")
		// 	// 	g_node_depth_traverse_level (root, flags, depth, func, data);
	}
}

func traversePreOrder(n *Node, flags TraverseFlags, traverseFunc TraverseFunc, data interface{}) bool {
	if n.Children != nil {
		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

		child := n.Children
		for child != nil {
			current := child
			child = current.Next
			if traversePreOrder(current, flags, traverseFunc, data) {
				return true
			}
		}
	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func depthTraversePreOrder(n *Node, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data interface{}) bool {
	if n.Children != nil {
		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

		depth--
		if depth == 0 {
			return false
		}

		child := n.Children
		for child != nil {
			current := child
			child = current.Next
			if depthTraversePreOrder(current, flags, depth, traverseFunc, data) {
				return true
			}
		}
	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func traverseInOrder(n *Node, flags TraverseFlags, traverseFunc TraverseFunc, data interface{}) bool {
	if n.Children != nil {
		child := n.Children
		current := child
		child = current.Next
		if traverseInOrder(current, flags, traverseFunc, data) {
			return true
		}

		if traverseFunc(n, data) { //flags & G_TRAVERSE_NON_LEAFS) &&
			return true
		}

		for child != nil {
			current = child
			child = current.Next
			if traverseInOrder(current, flags, traverseFunc, data) {
				return true
			}
		}
	} else if traverseFunc(n, data) { // (flags & G_TRAVERSE_LEAFS) &&
		return true
	}

	return false
}

func depthTraverseInOrder(n *Node, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data interface{}) bool {
	if n.Children != nil {
		depth--
		if depth > 0 {
			child := n.Children
			current := child
			child = current.Next

			if depthTraverseInOrder(current, flags, depth, traverseFunc, data) {
				return true
			}

			if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
				return true
			}

			for child != nil {
				current = child
				child = current.Next
				if depthTraverseInOrder(current, flags, depth, traverseFunc, data) {
					return true
				}
			}
		} else if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}
	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func traversePostOrder(n *Node, flags TraverseFlags, traverseFunc TraverseFunc, data interface{}) bool {
	if n.Children != nil {
		child := n.Children
		for child != nil {

			current := child
			child = current.Next
			if traversePostOrder(current, flags, traverseFunc, data) {
				return true
			}
		}

		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true

	}

	return false
}

func depthTraversePostOrder(n *Node, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data interface{}) bool {
	if n.Children != nil {
		depth--
		if depth > 0 {

			child := n.Children
			for child != nil {

				current := child
				child = current.Next
				if depthTraversePostOrder(current, flags, depth, traverseFunc, data) {
					return true
				}
			}
		}

		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func (n *Node) String() string {
	if n == nil {
		return "()"
	}

	currentLevel := 0
	lastNode := n
	levels := make([]string, 0, 10)
	levels = append(levels, "")
	tFunc := func(node *Node, value interface{}) bool {
		currentLevel = 0
		n := node.Parent
		for n != nil {
			currentLevel++
			if len(levels) <= currentLevel {
				levels = append(levels, "")
			}
			n = n.Parent
		}
		levels[currentLevel] += fmt.Sprintf("(%v)", node.Value) + "\t"
		lastNode = node
		return false
	}

	Traverse(n, TraversePreOrder, TraverseAll, -1, tFunc, nil)
	s := ""
	for _, v := range levels {
		s += v + "\n\n"
	}
	return s
}
