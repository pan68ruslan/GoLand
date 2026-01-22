package documentStore

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInsertAndRangeSearch(t *testing.T) {
	tree := &BinaryTree{}
	tree.Insert("b", "id1")
	tree.Insert("a", "id2")
	tree.Insert("c", "id3")
	tree.Insert("b", "id4")
	if tree.Root.Key != "b" {
		t.Errorf("expected root key 'b', got %s", tree.Root.Key)
	}
	if len(tree.Root.ID) != 2 {
		t.Errorf("expected 2 IDs at root, got %d", len(tree.Root.ID))
	}
	minVal := "a"
	maxVal := "c"
	result := tree.RangeSearch(&minVal, &maxVal)
	expected := []string{"id2", "id1", "id4", "id3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestMarshalUnmarshalTreeNode(t *testing.T) {
	node := &TreeNode{
		Key: "x",
		ID:  []string{"id1", "id2"},
		Left: &TreeNode{
			Key: "l",
			ID:  []string{"id3"},
		},
		Right: &TreeNode{
			Key: "r",
			ID:  []string{"id4"},
		},
	}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var restored TreeNode
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if restored.Key != "x" || len(restored.ID) != 2 {
		t.Errorf("expected key 'x' with 2 IDs, got %s with %d IDs", restored.Key, len(restored.ID))
	}
	if restored.Left.Key != "l" || restored.Right.Key != "r" {
		t.Errorf("expected left 'l' and right 'r', got %s and %s", restored.Left.Key, restored.Right.Key)
	}
}

func TestMarshalUnmarshalBinaryTree(t *testing.T) {
	tree := &BinaryTree{}
	tree.Insert("m", "id1")
	tree.Insert("a", "id2")
	tree.Insert("z", "id3")
	data, err := json.Marshal(tree)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var restored BinaryTree
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	result := restored.RangeSearch(nil, nil)
	expected := []string{"id2", "id1", "id3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestEmptyTreeMarshal(t *testing.T) {
	var tree *BinaryTree
	data, err := json.Marshal(tree)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected 'null', got %s", string(data))
	}
}
