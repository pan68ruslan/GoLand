package main

import (
	"fmt"
	"log/slog"
	"os"

	ds "lesson_09/internal/documentStore"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	store := ds.NewStore("InitialStore", logger)
	store.AddCollection("users", &ds.CollectionConfig{
		PrimaryKey:    "id",
		IndexedFields: []string{"id", "name"},
	})
	users := store.Collections["users"]
	doc1 := ds.Document{
		Fields: map[string]ds.DocumentField{
			"id":    {Type: ds.DocumentFieldTypeString, Value: "u1"},
			"name":  {Type: ds.DocumentFieldTypeString, Value: "Alice"},
			"email": {Type: ds.DocumentFieldTypeString, Value: "alice@mail.com"},
			"role":  {Type: ds.DocumentFieldTypeString, Value: "guest"},
		},
	}
	doc2 := ds.Document{
		Fields: map[string]ds.DocumentField{
			"id":    {Type: ds.DocumentFieldTypeString, Value: "u2"},
			"name":  {Type: ds.DocumentFieldTypeString, Value: "Bob"},
			"email": {Type: ds.DocumentFieldTypeString, Value: "bob@mail.com"},
			"role":  {Type: ds.DocumentFieldTypeString, Value: "admin"},
		},
	}
	doc3 := ds.Document{
		Fields: map[string]ds.DocumentField{
			"id":    {Type: ds.DocumentFieldTypeString, Value: "u3"},
			"name":  {Type: ds.DocumentFieldTypeString, Value: "Charlie"},
			"email": {Type: ds.DocumentFieldTypeString, Value: "charlie@mail.com"},
			"role":  {Type: ds.DocumentFieldTypeString, Value: "user"},
		},
	}
	users.AddDocument("u1", doc1)
	users.AddDocument("u2", doc2)
	users.AddDocument("u3", doc3)

	if err := users.CreateIndex("name"); err != nil {
		logger.Error("Index creation error", "Error", err)
	} else {
		logger.Info("Index Name created")
	}
	if err := users.CreateIndex("id"); err != nil {
		logger.Error("Index creation error", "Error", err)
	} else {
		logger.Info("Index ID created")
	}
	minVal := "A"
	maxVal := "D"
	params := ds.QueryParams{
		Desc:     true,
		MinValue: &minVal,
		MaxValue: &maxVal,
	}
	fmt.Println("Query by NAME index (descending)")
	if results, err := users.Query("name", params); err == nil {
		for _, doc := range results {
			fmt.Println("Found:", doc.Fields["id"].Value, "-", doc.Fields["name"].Value)
		}
	} else {
		logger.Error("Query by NAME error:", "Error", err)
	}

	fmt.Println("Query by ID index (ascending)")
	minVal = "u1"
	maxVal = "u2"
	params.Desc = false
	params.MinValue = &minVal
	params.MaxValue = &maxVal
	if results, err := users.Query("id", params); err == nil {
		for _, doc := range results {
			fmt.Println("Found:", doc.Fields["id"].Value, "-", doc.Fields["name"].Value)
		}
	} else {
		logger.Error("Query by ID error:", "Error", err)
	}

	fmt.Println("Dump the store")
	if jsonDoc, e := store.Dump(); e == nil {
		logger.Info("The store dump was created:", "Name", store.Name)
		logger.Debug(string(jsonDoc))
		if restored, e := ds.NewStoreFromDump(jsonDoc); e == nil {
			logger.Info("The initial store was restored", "Name", restored.Name)
			fmt.Println("Query by ID in restored store")
			maxVal = "u4"
			params.MaxValue = &maxVal
			if results, err := users.Query("id", params); err == nil {
				for _, doc := range results {
					fmt.Println("Found:", doc.Fields["name"].Value, "-", doc.Fields["role"].Value, ",", doc.Fields["email"].Value)
				}
			} else {
				logger.Error("Query by ID error:", "Error", err)
			}
		} else {
			logger.Error("The initial store wasn't restored", "Error", e)
		}
	} else {
		logger.Error("The store dump was not created:", "Error", e)
	}
}
