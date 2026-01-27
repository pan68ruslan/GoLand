package documentStore

import (
	"encoding/json"
	"sync"
)

type TreeNode struct {
	Key   string
	ID    []string
	Left  *TreeNode
	Right *TreeNode
}

type BinaryTree struct {
	Root *TreeNode
	mu   sync.RWMutex
}

func BalanceTree(root *TreeNode) *TreeNode {
	var keys []string
	var ids [][]string

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

func (t *BinaryTree) Insert(key string, id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Root = insertNode(t.Root, key, id)
	t.Root = BalanceTree(t.Root)
}

func insertNode(node *TreeNode, key string, id string) *TreeNode {
	if node == nil {
		return &TreeNode{Key: key, ID: []string{id}}
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

func (t *BinaryTree) RangeSearch(min, max *string) []string {
	var result []string
	t.mu.RLock()
	defer t.mu.RUnlock()
	inOrderRange(t.Root, min, max, &result)
	return result
}

func inOrderRange(node *TreeNode, min, max *string, result *[]string) {
	if node == nil {
		return
	}
	if min == nil || node.Key >= *min {
		inOrderRange(node.Left, min, max, result)
	}
	if (min == nil || node.Key >= *min) && (max == nil || node.Key <= *max) {
		*result = append(*result, node.ID...)
	}
	if max == nil || node.Key <= *max {
		inOrderRange(node.Right, min, max, result)
	}
}

func (t *BinaryTree) RemoveFromIndex(key string, id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Root = removeIDFromNode(t.Root, key, id)
	t.Root = BalanceTree(t.Root)
}

func removeIDFromNode(node *TreeNode, key string, id string) *TreeNode {
	if node == nil {
		return nil
	}
	if key < node.Key {
		node.Left = removeIDFromNode(node.Left, key, id)
	} else if key > node.Key {
		node.Right = removeIDFromNode(node.Right, key, id)
	} else {
		newIDs := make([]string, 0, len(node.ID))
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

func removeNode(node *TreeNode, key string) *TreeNode {
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

func (t *BinaryTree) UpdateIndex(key string, id string, add bool) {
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
		Key   string    `json:"Key"`
		ID    []string  `json:"ID"`
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
