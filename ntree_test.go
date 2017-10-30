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
	str := nodes["root"].String()
	fmt.Println(str)
}

func TestNodeCount(t *testing.T) {
	nodes := GenerateTree()
	count := ntree.NodeCount(nodes["root"], ntree.TraverseAll)
	if count != len(nodes) {
		t.Error("Mismatched node count!")
	}
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
