package documentStore

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"
)

func BenchmarkCollectionConcurrent(b *testing.B) {
	cfg := &CollectionConfig{
		PrimaryKey:    "id",
		IndexedFields: []string{"id", "name"},
	}
	coll := NewCollection("users", cfg, slog.Default())

	b.ResetTimer() // обнуляємо таймер перед основним циклом

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numGoroutines := 1000

		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()

				// створюємо документ
				doc := Document{
					Fields: map[string]DocumentField{
						"id":   {Type: DocumentFieldTypeString, Value: fmt.Sprintf("u%d-%d", i, j)},
						"name": {Type: DocumentFieldTypeString, Value: fmt.Sprintf("User%d-%d", i, j)},
					},
				}
				_ = coll.AddDocument(doc)

				// читаємо документ
				_, _ = coll.GetDocument(fmt.Sprintf("u%d-%d", i, j))

				// видаляємо документ
				coll.DeleteDocument(fmt.Sprintf("u%d-%d", i, j))
			}(j)
		}

		wg.Wait()
	}
}

// BenchmarkCollectionReadHeavy — сценарій з великою кількістю читань
func BenchmarkCollectionReadHeavy(b *testing.B) {
	cfg := &CollectionConfig{
		PrimaryKey:    "id",
		IndexedFields: []string{"id", "name"},
	}
	coll := NewCollection("users", cfg, slog.Default())

	// попередньо додаємо документи
	for i := 0; i < 1000; i++ {
		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: fmt.Sprintf("u%d", i)},
				"name": {Type: DocumentFieldTypeString, Value: fmt.Sprintf("User%d", i)},
			},
		}
		_ = coll.AddDocument(doc)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numReaders := 1000

		for j := 0; j < numReaders; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				_, _ = coll.GetDocument(fmt.Sprintf("u%d", j%1000))
			}(j)
		}

		wg.Wait()
	}
}

// BenchmarkCollectionWriteHeavy — сценарій з великою кількістю записів
func BenchmarkCollectionWriteHeavy(b *testing.B) {
	cfg := &CollectionConfig{
		PrimaryKey:    "id",
		IndexedFields: []string{"id", "name"},
	}
	coll := NewCollection("users", cfg, slog.Default())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numWriters := 1000

		for j := 0; j < numWriters; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				doc := Document{
					Fields: map[string]DocumentField{
						"id":   {Type: DocumentFieldTypeString, Value: fmt.Sprintf("w%d-%d", i, j)},
						"name": {Type: DocumentFieldTypeString, Value: fmt.Sprintf("Writer%d-%d", i, j)},
					},
				}
				_ = coll.AddDocument(doc)
				coll.DeleteDocument(fmt.Sprintf("w%d-%d", i, j))
			}(j)
		}

		wg.Wait()
	}
}
