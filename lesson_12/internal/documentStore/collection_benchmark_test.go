package documentStore

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"
)

func BenchmarkCollectionConcurrent(b *testing.B) {
	coll := NewCollection("users", slog.Default())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numGoroutines := 1000
		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				doc := Document{
					Fields: map[string]DocumentField{
						"id":   {Type: DocumentFieldTypeNumber, Value: i - j}, //fmt.Sprintf("u%d-%d", i, j)},
						"name": {Type: DocumentFieldTypeString, Value: fmt.Sprintf("User%d-%d", i, j)},
					},
				}
				_ = coll.PutDocument(doc)
				_, _ = coll.GetDocument(i - j) //fmt.Sprintf("u%d-%d", i, j))
				coll.DeleteDocument(i - j)     //fmt.Sprintf("u%d-%d", i, j))
			}(j)
		}
		wg.Wait()
	}
}

func BenchmarkCollectionReadHeavy(b *testing.B) {
	coll := NewCollection("users", slog.Default())
	for i := 0; i < 1000; i++ {
		doc := Document{
			Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeNumber, Value: i}, //fmt.Sprintf("u%d", i)},
				"name": {Type: DocumentFieldTypeString, Value: fmt.Sprintf("User%d", i)},
			},
		}
		_ = coll.PutDocument(doc)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numReaders := 1000
		for j := 0; j < numReaders; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				_, _ = coll.GetDocument(j % 1000)
			}(j)
		}
		wg.Wait()
	}
}

func BenchmarkCollectionWriteHeavy(b *testing.B) {
	coll := NewCollection("users", slog.Default())
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
						"id":   {Type: DocumentFieldTypeNumber, Value: i - j}, //fmt.Sprintf("w%d-%d", i, j)},
						"name": {Type: DocumentFieldTypeString, Value: fmt.Sprintf("Writer%d-%d", i, j)},
					},
				}
				_ = coll.PutDocument(doc)
				coll.DeleteDocument(i - j) //fmt.Sprintf("w%d-%d", i, j))
			}(j)
		}
		wg.Wait()
	}
}
