package main

import (
	"log/slog"
	"os"

	ds "github.com/pan68ruslan/GoLand/lesson_07/internal/document_store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var initialStore = ds.NewStore("InitialStore")
	cfg := &ds.CollectionConfig{
		PrimaryKey: "key",
	}
	if cfg == nil {
		logger.Error("Initializing CollectionConfig failed", "PrimaryKey", cfg.PrimaryKey)
	} else {
		logger.Info("The CollectionConfig initialized successfully", "PrimaryKey", cfg.PrimaryKey)
	}
	collectionName := "DocumentCollection"
	ok, collection := initialStore.CreateCollection(collectionName, cfg)
	if ok {
		logger.Info("The empty collection was created in the store", "Name", collectionName)
	} else {
		logger.Error("Collection failed to be created in the store", "Name", collectionName)
	}
	collection.Put(ds.Document{
		Fields: map[string]ds.DocumentField{
			"key":        {Type: ds.DocumentFieldTypeString, Value: "doc1"},
			"title":      {Type: ds.DocumentFieldTypeString, Value: "firstDocument"},
			"isApproved": {Type: ds.DocumentFieldTypeBool, Value: false},
			"pages":      {Type: ds.DocumentFieldTypeNumber, Value: 42},
		},
	})
	collection.Put(ds.Document{
		Fields: map[string]ds.DocumentField{
			"key":        {Type: ds.DocumentFieldTypeString, Value: "doc2"},
			"title":      {Type: ds.DocumentFieldTypeString, Value: "secondDocument"},
			"isApproved": {Type: ds.DocumentFieldTypeBool, Value: true},
			"pages":      {Type: ds.DocumentFieldTypeNumber, Value: 100},
		},
	})
	logger.Info("The documents were added to the collection", "Amount", len(collection.Documents), "Name", initialStore.Name)
	jsonDoc, e := initialStore.Dump()
	if e == nil {
		logger.Info("The store dump was created:", "Name", initialStore.Name)
		logger.Debug(string(jsonDoc))
	} else {
		logger.Error("The store dump was not created:", "Error", e)
	}
	restoredStore, e := ds.NewStoreFromDump(jsonDoc)
	if e == nil {
		logger.Info("The initial store was restored\n", "Name", restoredStore.Name)
	} else {
		logger.Error("The initial store wasn't restored", "Error", e)
	}
	if e = initialStore.DumpToFile(initialStore.Name); e != nil {
		logger.Error("The store wasn't restored", "Name", initialStore.Name, "Error", e)
	}
	if e = restoredStore.DumpToFile(restoredStore.Name + "Restored"); e != nil {
		logger.Error("The store wasn't restored", "Error", e)
	}
	fromFileStore, e := ds.NewStoreFromFile(initialStore.Name)
	fileStore, e := fromFileStore.Dump()
	if e == nil {
		logger.Info("The store was created from the dump file:", "Name", initialStore.Name)
		logger.Debug(string(fileStore))
	} else {
		logger.Error("The store wasn't restored", "Error", e)
	}
}
