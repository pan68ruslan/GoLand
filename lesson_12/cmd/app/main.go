package main

import (
	"fmt"
	"log/slog"
	"sync"

	ds "lesson_12/internal/documentStore"
)

func main() {
	cfg := &ds.CollectionConfig{
		PrimaryKey:    "id",
		IndexedFields: []string{"id", "name"},
	}
	coll := ds.NewCollection("users", cfg, slog.Default())
	var wg sync.WaitGroup
	numGoroutines := 1000
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			doc := ds.Document{
				Fields: map[string]ds.DocumentField{
					"id":   {Type: ds.DocumentFieldTypeString, Value: fmt.Sprintf("u%d", i)},
					"name": {Type: ds.DocumentFieldTypeString, Value: fmt.Sprintf("User%d", i)},
				},
			}
			_ = coll.PutDocument(doc)
			if d, ok := coll.GetDocument(fmt.Sprintf("u%d", i)); ok {
				fmt.Printf("Found: %s\n", d.Fields["name"].Value)
			}
			coll.DeleteDocument(fmt.Sprintf("u%d", i))
		}(i)
	}
	wg.Wait()
	fmt.Println("All goroutines finished")
}
