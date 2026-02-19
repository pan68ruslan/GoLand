package documentStore

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInsertAndRangeSearch(t *testing.T) {
	tree := &BinaryTree{}
	tree.Insert(10, 1)
	tree.Insert(5, 2)
	tree.Insert(15, 3)
	tree.Insert(10, 4) // duplicate key
	result := tree.RangeSearch(5, 15)
	expected := []int{2, 1, 4, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RangeSearch got %v, want %v", result, expected)
	}
}

func TestRemoveFromIndex(t *testing.T) {
	tree := &BinaryTree{}
	tree.Insert(10, 1)
	tree.Insert(10, 2)
	tree.Insert(20, 3)
	tree.RemoveFromIndex(10, 1)
	result := tree.RangeSearch(0, 30)
	expected := []int{2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RemoveFromIndex got %v, want %v", result, expected)
	}
	tree.RemoveFromIndex(10, 2)
	result = tree.RangeSearch(0, 30)
	expected = []int{3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RemoveFromIndex got %v, want %v", result, expected)
	}
}

func TestUpdateIndex(t *testing.T) {
	tree := &BinaryTree{}
	tree.UpdateIndex(5, 100, true)  // додати
	tree.UpdateIndex(5, 200, true)  // додати
	tree.UpdateIndex(5, 100, false) // видалити
	result := tree.RangeSearch(0, 10)
	expected := []int{200}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("UpdateIndex got %v, want %v", result, expected)
	}
}

func TestJSONMarshalUnmarshal(t *testing.T) {
	tree := &BinaryTree{}
	tree.Insert(1, 1)
	tree.Insert(5, 2)
	tree.Insert(10, 3)
	data, err := json.Marshal(tree.Root)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	var node TreeNode
	if err := json.Unmarshal(data, &node); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if node.Key != 5 {
		t.Errorf("Expected root key 10, got %d", node.Key)
	}
}
