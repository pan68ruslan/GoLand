package documentStore

import "encoding/json"

type TreeNode struct {
	Key   string
	ID    []string
	Left  *TreeNode
	Right *TreeNode
}

type BinaryTree struct {
	Root *TreeNode
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

func (t *BinaryTree) Insert(key string, id string) {
	t.Root = insertNode(t.Root, key, id)
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
