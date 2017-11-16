package ntree_test

import (
	"fmt"
	"testing"

	ntree "github.com/blazingorb/ntreego"
)

type MockData struct {
	Id    string
	Value int
}

func GenerateTree() map[string]*ntree.Node {
	a_1 := &MockData{"a_1", 0}
	a_2 := &MockData{"a_2", 0}

	a_1_1 := &MockData{"a_1_1", 0}
	a_1_2 := &MockData{"a_1_2", 0}
	a_1_3 := &MockData{"a_1_3", 0}

	a_2_1 := &MockData{"a_2_1", 0}
	a_2_2 := &MockData{"a_2_2", 0}

	nodes := make(map[string]*ntree.Node)
	nodes["root"] = ntree.New(&MockData{"armor", 0})
	nodes["a_1"] = ntree.AppendChild(nodes["root"], ntree.New(a_1))
	nodes["a_2"] = ntree.AppendChild(nodes["root"], ntree.New(a_2))

	nodes["a_1_1"] = ntree.AppendChild(nodes["a_1"], ntree.New(a_1_1))
	nodes["a_1_2"] = ntree.AppendChild(nodes["a_1"], ntree.New(a_1_2))
	nodes["a_1_3"] = ntree.AppendChild(nodes["a_1"], ntree.New(a_1_3))

	nodes["a_2_1"] = ntree.AppendChild(nodes["a_2"], ntree.New(a_2_1))
	nodes["a_2_2"] = ntree.AppendChild(nodes["a_2"], ntree.New(a_2_2))

	return nodes
}

func TestCreateTree(t *testing.T) {
	var initialValue int = 1
	tree := ntree.New(initialValue)
	value, ok := tree.Value.(int)
	if !ok || value != initialValue {
		t.Errorf("Tree Creation Failed")
	}
}

func TestTreeSearch(t *testing.T) {
	nodes := GenerateTree()
	toSearch := nodes["a_2_2"]
	parentOfToSearch := nodes["a_2"]

	var resultNode *ntree.Node
	searchFunc := func(n *ntree.Node, value interface{}) bool {
		v := value.(string)
		if n.Value.(*MockData).Id == v {
			resultNode = n
			return true
		}
		return false
	}

	ntree.Traverse(nodes["root"], ntree.TraverseInOrder, ntree.TraverseAll, -1, searchFunc, toSearch.Value.(*MockData).Id)
	if toSearch != resultNode {
		t.Error("Search Failed", toSearch.Value, resultNode)
	}
	if toSearch.Parent != parentOfToSearch {
		t.Error("Unexpected Parent", toSearch, parentOfToSearch)
	}
}

func TestTreeModify(t *testing.T) {
	nodes := GenerateTree()
	addFunc := func(n *ntree.Node, value interface{}) bool {
		v := value.(int)
		n.Value.(*MockData).Value += v
		return false
	}

	ntree.Traverse(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, -1, addFunc, 5)

	for _, node := range nodes {
		if node.Value.(*MockData).Value != 5 {
			t.Error("Value of node is wrong!")
		}
	}

	str := nodes["root"].String()
	fmt.Println(str)
}

func TestNodeCount(t *testing.T) {
	nodes := GenerateTree()
	count := ntree.NodeCount(nodes["root"], ntree.TraverseAll)
	if count != len(nodes) {
		t.Error("Mismatched node count!")
	}

	if ntree.NodeCount(nil, ntree.TraverseAll) != 0 {
		t.Error("NodeCount should be zero when the tree assigned is nil!")
	}

	if ntree.NodeCount(nodes["root"], ntree.TraverseMask+1) != 0 {
		t.Error("NodeCount should be zero when flags is larger than TraverseMask")
	}
}

func TestUnlink(t *testing.T) {
	m := GenerateTree()

	ntree.Unlink(nil)

	if ntree.NodeCount(m["root"], ntree.TraverseAll) != len(m) {
		t.Error("NodeCount should remained the same when calling ntree.Unlink(nil)")
	}

	wrapperFunc := func(n *ntree.Node) {
		ntree.Unlink(n)
		if n.Parent != nil || n.Previous != nil || n.Next != nil {
			t.Error("Incomplete unlink!")
		}

		for _, node := range m {
			if node.Children == n || node.Previous == n || node.Next == n {
				t.Error("Incomplete unlink!")
			}
		}
	}

	wrapperFunc(m["a_1"])
	wrapperFunc(m["a_2_2"])

}

func TestGetRoot(t *testing.T) {
	m := GenerateTree()
	root, depth := ntree.GetRoot(nil)
	if root != nil || depth != 0 {
		t.Error("root should be nil and depth should be 0 when ntree.GetRoot(nil) is called!")
	}

	root, _ = ntree.GetRoot(m["a_1"])
	if root == nil {
		t.Error("root == nil")
	} else {
		if root.Value != m["root"].Value {
			t.Error("Mismatched root node!")
		}
	}
}

func TestDepth(t *testing.T) {
	m := GenerateTree()
	if ntree.Depth(m["root"]) != 1 {
		t.Error("Depth of root should be 1!")
	}

	if ntree.Depth(m["a_1"]) != 2 {
		t.Error("Depth of a_1 should be 2!")
	}

	if ntree.Depth(m["a_1_1"]) != 3 {
		t.Error("Depth of a_1_1 should be 3!")
	}

}

func TestInsert(t *testing.T) {
	nodes := GenerateTree()
	a_1 := nodes["a_1"]
	a_2 := nodes["a_2"]
	a_2_1 := nodes["a_2_1"]

	newNode := ntree.New(1)

	if ntree.Insert(nil, newNode) != nil || ntree.Insert(nodes["root"], nil) != nil {
		t.Error("Result should be nil when one of the arguments is nil")
	}

	if ntree.Insert(a_1, a_2) != nil {
		t.Error("Result should be nil when a non-root node is inserted as a child of other node")
	}

	ntree.Insert(nodes["root"], newNode)
	if a_2.Next != newNode || newNode.Previous != a_2 {
		t.Error("newNode should be a sibling of a_2")
	}

	newNode2 := ntree.New(2)
	ntree.Insert(a_2_1, newNode2)
	if a_2_1.Children != newNode2 || newNode2.Parent != a_2_1 {
		t.Error("newNode2 should be a child of a_2_1")
	}
}

func TestAppendChild(t *testing.T) {
	nodes := GenerateTree()
	a_1 := nodes["a_1"]
	a_2 := nodes["a_2"]
	a_2_1 := nodes["a_2_1"]

	newNode := ntree.New(1)

	if ntree.AppendChild(nil, newNode) != nil || ntree.AppendChild(nodes["root"], nil) != nil {
		t.Error("Result should be nil when one of the arguments is nil")
	}

	if ntree.AppendChild(a_1, a_2) != nil {
		t.Error("Result should be nil when a non-root node is appended as a child of other node")
	}

	ntree.AppendChild(nodes["root"], newNode)
	if a_2.Next != newNode || newNode.Previous != a_2 {
		t.Error("newNode should be a sibling of a_2")
	}

	newNode2 := ntree.New(2)
	ntree.AppendChild(a_2_1, newNode2)
	if a_2_1.Children != newNode2 || newNode2.Parent != a_2_1 {
		t.Error("newNode2 should be a child of a_2_1")
	}
}

func TestIsRoot(t *testing.T) {
	nodes := GenerateTree()
	if !ntree.IsRoot(nodes["root"]) {
		t.Error("Result should be true when root node is passed in")
	}

	if ntree.IsRoot(nodes["a_1"]) || ntree.IsRoot(nodes["a_2"]) {
		t.Error("Result should be false when non-root nodes are passed in")
	}
}

func TestFindNode(t *testing.T) {
	m := GenerateTree()

	if ntree.FindNode(nil, ntree.TraverseInOrder, ntree.TraverseAll, m["a_1_1"].Value) != nil {
		t.Error("Result should be nil when nil is passed as root argument!")
	}

	nodeFound := ntree.FindNode(m["root"], ntree.TraverseInOrder, ntree.TraverseAll, m["a_1_1"].Value)

	if nodeFound != m["a_1_1"] {
		t.Error("Wrong node has be found!")
	}
}

func TestTraverseFilter(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	depth_a_1_1 := ntree.Depth(nodes["a_1_1"])

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFilterFunc := func(root *ntree.Node, order ntree.TraverseType, flag ntree.TraverseFlags, depth int, traverseFunc func(n *ntree.Node, data interface{}) bool) {
		visitCount = 0
		ntree.Traverse(root, order, flag, depth, traverseFunc, 0)

		if visitCount != 0 {
			t.Error("Traverse Error! Visit count should be zero when argument meets the filter condition!")
		}
	}

	wrappedFilterFunc(nil, ntree.TraversePreOrder, ntree.TraverseAll, depth_a_1_1, traverseFunc)
	wrappedFilterFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, depth_a_1_1, nil)
	wrappedFilterFunc(nodes["root"], ntree.TraverseLevelOrder+1, ntree.TraverseAll, depth_a_1_1, traverseFunc)
	wrappedFilterFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseMask+1, depth_a_1_1, traverseFunc)
	wrappedFilterFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, -2, traverseFunc)
	wrappedFilterFunc(nodes["root"], ntree.TraversePreOrder, ntree.TraverseAll, 0, traverseFunc)
}
func TestTraverseAll(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseAll, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Error("Traverse Error! The expected visit count should be", expectedLength, "but return", visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return", lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be nil but return", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return nil")
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_2_2"], len(nodes))
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_2_2"], len(nodes))
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["root"], len(nodes))
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["root"], 3)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseAllwithConditions(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return true
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseAll, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Error("Traverse Error! The expected visit count should be", expectedLength, "but return", visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return", lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be nil but return", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return nil")
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_1_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_1_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["a_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, depth_a_1_1)
}

func TestTraverseLeaves(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Error("Traverse Error! The expected visit count should be", expectedLength, "but return", visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return", lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be nil but return", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return nil")
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_2_2"], 5)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_2_2"], 5)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_2_2"], 5)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nil, 0)
	wrappedFunc(ntree.TraverseInOrder, 2, nil, 0)
	wrappedFunc(ntree.TraversePostOrder, 2, nil, 0)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseLeavesWithConditions(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return true
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Error("Traverse Error! The expected visit count should be", expectedLength, "but return", visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return", lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be nil but return", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return nil")
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_1_1"], 1)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_1_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_1_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nil, 0)
	wrappedFunc(ntree.TraverseInOrder, 2, nil, 0)
	wrappedFunc(ntree.TraversePostOrder, 2, nil, 0)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseNonLeaves(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return false
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseNonLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Error("Traverse Error! The expected visit count should be", expectedLength, "but return", visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return", lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be nil but return", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return nil")
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["a_2"], 3)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_2"], 3)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["root"], 3)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_2"], 3)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["root"], 3)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

func TestTraverseNonLeavesWithCondistions(t *testing.T) {
	nodes := GenerateTree()
	visitCount := 0
	var lastVisitedNode *ntree.Node = nil

	traverseFunc := func(n *ntree.Node, data interface{}) bool {
		visitCount++
		lastVisitedNode = n
		//Return true will stop Traverse Function from traversing other remained nodes
		return true
	}

	wrappedFunc := func(order ntree.TraverseType, depth int, expectedNode *ntree.Node, expectedLength int) {
		visitCount = 0
		lastVisitedNode = nil
		ntree.Traverse(nodes["root"], order, ntree.TraverseNonLeaves, depth, traverseFunc, 0)

		if visitCount != expectedLength {
			t.Error("Traverse Error! The expected visit count should be", expectedLength, "but return", visitCount)
		}

		if lastVisitedNode != expectedNode {
			if lastVisitedNode != nil {
				if expectedNode != nil {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return", lastVisitedNode.Value.(*MockData).Id)
				} else {
					t.Error("Traverse Error! The expected node that traverseFunc last visited should be nil but return", lastVisitedNode.Value.(*MockData).Id)
				}
			} else {
				t.Error("Traverse Error! The expected node that traverseFunc last visited should be", expectedNode.Value.(*MockData).Id, "but return nil")
			}
		}
	}

	wrappedFunc(ntree.TraversePreOrder, -1, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, -1, nodes["a_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, -1, nodes["a_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, -1)

	wrappedFunc(ntree.TraversePreOrder, 2, nodes["root"], 1)
	wrappedFunc(ntree.TraverseInOrder, 2, nodes["a_1"], 1)
	wrappedFunc(ntree.TraversePostOrder, 2, nodes["a_1"], 1)
	//wrappedFunc(ntree.TraverseLevelOrder, 2)
}

// func TestFindRoot(t *testing.T) {
// 	_, node := GenerateTree()
// 	t.Log("\n", ntree.GetRoot(node).Value)
// }

// func TestTreeFindParentNode(t *testing.T) {
// 	type Mock struct {
// 		Id    string
// 		Value int
// 	}

// 	root := ntree.NewWithValue(&Mock{"armor", 0})
// 	child1 := ntree.NewWithValue(&Mock{"a_1", 1})
// 	child2 := ntree.NewWithValue(&Mock{"a_2", 1})

// 	child1_1 := ntree.NewWithValue(&Mock{"a_1_1", 2})
// 	child1_2 := ntree.NewWithValue(&Mock{"a_1_2", 2})
// 	child1_3 := ntree.NewWithValue(&Mock{"a_1_3", 2})

// 	child2_1 := ntree.NewWithValue(&Mock{"a_2_1", 2})
// 	child2_2 := ntree.NewWithValue(&Mock{"a_2_2", 2})

// 	ntree.InsertChild(root, child1)
// 	ntree.InsertChild(root, child2)

// 	ntree.InsertChild(child1, child1_1)
// 	ntree.InsertChild(child1, child1_2)
// 	ntree.InsertChild(child1, child1_3)

// 	ntree.InsertChild(child2, child2_1)
// 	ntree.InsertChild(child2, child2_2)

// 	// s := ntree.FindParent(root, child1_1)
// 	// if s == nil {
// 	// 	t.Error("Node Not Found: ", s)
// 	// 	return
// 	// }

// 	// t.Log("\n" + root.String())
// 	// t.Log("Found Val: ", s.Value)
// }

// func TestTreeHasPath(t *testing.T) {
// 	type Mock struct {
// 		Id    string
// 		Value int
// 	}

// 	root := ntree.NewWithValue(&Mock{"armor", 0})
// 	child1 := ntree.NewWithValue(&Mock{"a_1", 1})
// 	child2 := ntree.NewWithValue(&Mock{"a_2", 1})

// 	child1_1 := ntree.NewWithValue(&Mock{"a_1_1", 2})
// 	child1_2 := ntree.NewWithValue(&Mock{"a_1_2", 2})
// 	child1_3 := ntree.NewWithValue(&Mock{"a_1_3", 2})

// 	child2_1 := ntree.NewWithValue(&Mock{"a_2_1", 2})
// 	child2_2 := ntree.NewWithValue(&Mock{"a_2_2", 2})

// 	ntree.InsertChild(root, child1)
// 	ntree.InsertChild(root, child2)

// 	ntree.InsertChild(child1, child1_1)
// 	ntree.InsertChild(child1, child1_2)
// 	ntree.InsertChild(child1, child1_3)

// 	ntree.InsertChild(child2, child2_1)
// 	ntree.InsertChild(child2, child2_2)
// 	t.Log("\n" + root.String())

// 	result := make([]*ntree.NTree, 0)
// 	ntree.HasPath(root, child2_2, &result)

// 	t.Log("Path Length: ", len(result))
// 	for _, r := range result {
// 		t.Log(r.Value)
// 	}

// }
