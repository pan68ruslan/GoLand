package documentStore

import (
	"encoding/json"
	"sync"
)

type TreeNode struct {
	Key   int
	ID    []int
	Left  *TreeNode
	Right *TreeNode
}

type BinaryTree struct {
	Root *TreeNode
	mu   sync.RWMutex
}

func BalanceTree(root *TreeNode) *TreeNode {
	var keys []int
	var ids [][]int

	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil {
			return
		}
		inorder(node.Left)
		keys = append(keys, node.Key)
		ids = append(ids, node.ID)
		inorder(node.Right)
	}
	inorder(root)
	var build func(start, end int) *TreeNode
	build = func(start, end int) *TreeNode {
		if start > end {
			return nil
		}
		mid := (start + end) / 2
		node := &TreeNode{
			Key: keys[mid],
			ID:  ids[mid],
		}
		node.Left = build(start, mid-1)
		node.Right = build(mid+1, end)
		return node
	}
	return build(0, len(keys)-1)
}

func (t *BinaryTree) Insert(key, id int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Root = insertNode(t.Root, key, id)
	t.Root = BalanceTree(t.Root)
}

func insertNode(node *TreeNode, key, id int) *TreeNode {
	if node == nil {
		return &TreeNode{Key: key, ID: []int{id}}
	}
	if key < node.Key {
		node.Left = insertNode(node.Left, key, id)
	} else if key > node.Key {
		node.Right = insertNode(node.Right, key, id)
	} else {
		node.ID = append(node.ID, id)
	}
	return node
}

func (t *BinaryTree) RangeSearch(min, max int) []int {
	var result []int
	t.mu.RLock()
	defer t.mu.RUnlock()
	inOrderRange(t.Root, min, max, &result)
	return result
}

func inOrderRange(node *TreeNode, min, max int, result *[]int) {
	if node == nil {
		return
	}
	if node.Key >= min {
		inOrderRange(node.Left, min, max, result)
	}
	if node.Key >= min && node.Key <= max {
		*result = append(*result, node.ID...)
	}
	if node.Key <= max {
		inOrderRange(node.Right, min, max, result)
	}
}

func (t *BinaryTree) RemoveFromIndex(key, id int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Root = removeIDFromNode(t.Root, key, id)
	t.Root = BalanceTree(t.Root)
}

func removeIDFromNode(node *TreeNode, key, id int) *TreeNode {
	if node == nil {
		return nil
	}
	if key < node.Key {
		node.Left = removeIDFromNode(node.Left, key, id)
	} else if key > node.Key {
		node.Right = removeIDFromNode(node.Right, key, id)
	} else {
		newIDs := make([]int, 0, len(node.ID))
		for _, existingID := range node.ID {
			if existingID != id {
				newIDs = append(newIDs, existingID)
			}
		}
		node.ID = newIDs
		if len(node.ID) == 0 {
			if node.Left == nil {
				return node.Right
			} else if node.Right == nil {
				return node.Left
			} else {
				minNode := node.Right
				for minNode.Left != nil {
					minNode = minNode.Left
				}
				node.Key = minNode.Key
				node.ID = minNode.ID
				node.Right = removeNode(node.Right, minNode.Key)
			}
		}
	}
	return node
}

func removeNode(node *TreeNode, key int) *TreeNode {
	if node == nil {
		return nil
	}
	if key < node.Key {
		node.Left = removeNode(node.Left, key)
	} else if key > node.Key {
		node.Right = removeNode(node.Right, key)
	} else {
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		} else {
			minNode := node.Right
			for minNode.Left != nil {
				minNode = minNode.Left
			}
			node.Key = minNode.Key
			node.ID = minNode.ID
			node.Right = removeNode(node.Right, minNode.Key)
		}
	}
	return node
}

func (t *BinaryTree) UpdateIndex(key, id int, add bool) {
	if add {
		t.Insert(key, id)
	} else {
		t.RemoveFromIndex(key, id)
	}
}

func (n *TreeNode) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("null"), nil
	}
	fields := map[string]interface{}{
		"Key":   n.Key,
		"ID":    n.ID,
		"Left":  n.Left,
		"Right": n.Right,
	}
	return json.Marshal(fields)
}

func (n *TreeNode) UnmarshalJSON(data []byte) error {
	var out struct {
		Key   int       `json:"Key"`
		ID    []int     `json:"ID"`
		Left  *TreeNode `json:"Left"`
		Right *TreeNode `json:"Right"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	n.Key = out.Key
	n.ID = out.ID
	n.Left = out.Left
	n.Right = out.Right
	return nil
}
