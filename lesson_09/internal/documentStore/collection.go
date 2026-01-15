package documentStore

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
)

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

type CollectionConfig struct {
	PrimaryKey    string   `json:"PrimaryKey"`
	IndexedFields []string `json:"IndexedFields"`
}

type Collection struct {
	Name      string                 `json:"Name"`
	Cfg       *CollectionConfig      `json:"Config"`
	Documents map[string]Document    `json:"Documents"`
	Indexes   map[string]*BinaryTree `json:"Indexes"`
	Logger    *slog.Logger
}

func NewCollection(name string, cfg *CollectionConfig, logger *slog.Logger) *Collection {
	return &Collection{
		Name:      name,
		Cfg:       cfg,
		Documents: make(map[string]Document),
		Indexes:   make(map[string]*BinaryTree),
		Logger:    logger,
	}
}

func (c *Collection) AddDocument(id string, doc Document) {
	if c.Documents == nil {
		c.Documents = make(map[string]Document)
	}
	c.Documents[id] = doc
	c.Logger.Info("Added document",
		slog.String("collection", c.Name),
		slog.String("id", id))
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

func (c *Collection) CreateIndex(fieldName string) error {
	if c.Indexes == nil {
		c.Indexes = make(map[string]*BinaryTree)
	}
	if _, exists := c.Indexes[fieldName]; exists {
		return fmt.Errorf("index already exists for field %s", fieldName)
	}
	tree := &BinaryTree{}
	for id, doc := range c.Documents {
		if field, ok := doc.Fields[fieldName]; ok {
			tree.Insert(fmt.Sprintf("%v", field.Value), id)
		}
	}
	tree.Root = BalanceTree(tree.Root)
	c.Indexes[fieldName] = tree
	c.Logger.Info("Created index",
		slog.String("collection", c.Name),
		slog.String("field", fieldName))
	return nil
}

func (c *Collection) DeleteIndex(fieldName string) error {
	if _, exists := c.Indexes[fieldName]; !exists {
		return fmt.Errorf("index does not exist for field %s", fieldName)
	}
	delete(c.Indexes, fieldName)
	c.Logger.Info("Deleted index",
		slog.String("collection", c.Name),
		slog.String("field", fieldName))
	for f, tree := range c.Indexes {
		if tree != nil && tree.Root != nil {
			tree.Root = BalanceTree(tree.Root)
			c.Logger.Info("Rebalanced index after deletion",
				slog.String("collection", c.Name),
				slog.String("field", f))
		}
	}
	return nil
}

func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	idx, exists := c.Indexes[fieldName]
	if !exists {
		return nil, fmt.Errorf("index does not exist for field %s", fieldName)
	}
	if params.MinValue == nil || params.MaxValue == nil || *params.MinValue > *params.MaxValue {
		return nil, fmt.Errorf("parametrs does not exist or wrong")
	}
	ids := idx.RangeSearch(params.MinValue, params.MaxValue)
	if params.Desc {
		sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] })
	} else {
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	}
	var result []Document
	for _, id := range ids {
		if doc, ok := c.Documents[id]; ok {
			result = append(result, doc)
		}
	}
	c.Logger.Info("Query executed",
		slog.String("collection", c.Name),
		slog.String("field", fieldName),
		slog.Int("found", len(result)))
	return result, nil
}

func (c *CollectionConfig) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{
		"PrimaryKey":    c.PrimaryKey,
		"IndexedFields": c.IndexedFields,
	}
	return json.Marshal(out)
}

func (c *CollectionConfig) UnmarshalJSON(data []byte) error {
	var tmp struct {
		PrimaryKey    string   `json:"PrimaryKey"`
		IndexedFields []string `json:"IndexedFields"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	c.PrimaryKey, c.IndexedFields = tmp.PrimaryKey, tmp.IndexedFields
	return nil
}

func (c *Collection) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}
	fields := map[string]interface{}{
		"Name":      c.Name,
		"Config":    c.Cfg,
		"Documents": c.Documents,
		"Indexes":   c.Indexes,
	}
	return json.Marshal(fields)
}

func (c *Collection) UnmarshalJSON(data []byte) error {
	var out struct {
		Name      string                 `json:"Name"`
		Config    *CollectionConfig      `json:"Config"`
		Documents map[string]Document    `json:"Documents"`
		Indexes   map[string]*BinaryTree `json:"Indexes"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	c.Name = out.Name
	c.Cfg = out.Config
	c.Documents = out.Documents
	c.Indexes = out.Indexes
	return nil
}

func (c *Collection) Put(doc Document) {
	if c.Cfg == nil || c.Cfg.PrimaryKey == "" {
		fmt.Println("[Collection]CollectionConfig is not configured")
		return
	}
	field, ok := doc.Fields[c.Cfg.PrimaryKey]
	if !ok || field.Type != DocumentFieldTypeString {
		fmt.Println("[Collection]Document is wrong")
		return
	}
	key, ok := field.Value.(string)
	if !ok {
		fmt.Println("[Collection]Key is not a string")
		return
	}
	fmt.Printf("[Collection]The document was added with key '%s'\n", key)
	c.Documents[key] = doc
}

func (c *Collection) Get(key string) (*Document, bool) {
	doc, ok := c.Documents[key]
	if !ok {
		fmt.Printf("[Collection]The document with key '%s' wasn't found\n", key)
		return nil, false
	}
	fmt.Printf("[Collection]The document with key '%s' was found\n", key)
	return &doc, true
}

func (c *Collection) Delete(key string) bool {
	if _, ok := c.Documents[key]; ok {
		fmt.Printf("[Collection]The document with key '%s' was deleted\n", key)
		delete(c.Documents, key)
		return true
	}
	fmt.Printf("[Collection]The document with key '%s' doesn't exist\n", key)
	return false
}

func (c *Collection) List() []Document {
	docs := make([]Document, 0)
	for _, doc := range c.Documents {
		docs = append(docs, doc)
	}
	if l := len(docs); l < 1 {
		fmt.Printf("[Collection]There are no documents in the collection '%s'\n", c.Name)
	} else {
		fmt.Printf("[Collection]There are %d documents in the collection '%s'\n", l, c.Name)
	}
	return docs
}
