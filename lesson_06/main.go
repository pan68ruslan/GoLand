package main

import (
	ds "github.com/pan68ruslan/GoLand/lesson_06/documentstore"
)

func main() {
	initialStore := ds.NewStore("InitialStore")
	cfg := &ds.CollectionConfig{
		PrimaryKey: "key",
	}
	if cfg == nil {
		ds.StoreLogger.Error("Initializing CollectionConfig failed", "PrimaryKey", cfg.PrimaryKey)
	} else {
		ds.StoreLogger.Info("The CollectionConfig initialized successfully", "PrimaryKey", cfg.PrimaryKey)
	}
	collectionName := "DocumentCollection"
	ok, collection := initialStore.CreateCollection(collectionName, cfg)
	if ok {
		ds.StoreLogger.Info("The empty collection was created in the store", "Name", collectionName)
	} else {
		ds.StoreLogger.Error("Collection failed to be created in the store", "Name", collectionName)
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
	ds.StoreLogger.Info("The documents were added to the collection", "Amount", len(collection.Documents), "Name", initialStore.Name)
	dumpDoc, e := initialStore.Dump()
	if e == nil {
		ds.StoreLogger.Info("The store dump was created", "Name", initialStore.Name)
		ds.StoreLogger.Debug(string(dumpDoc))
	} else {
		ds.StoreLogger.Error("The store dump was not created", "Name", initialStore.Name, "Error", e)
	}
	restoredStore, e := ds.NewStoreFromDump(dumpDoc)
	if e == nil {
		ds.StoreLogger.Info("The initial store was restored from dump", "Name", restoredStore.Name)
	} else {
		ds.StoreLogger.Error("The initial store wasn't restored from dump", "Error", e)
	}
	if e = initialStore.DumpToFile(initialStore.Name); e == nil {
		ds.StoreLogger.Info("The store dump to file was created", "NameOfFile", initialStore.Name)
	} else {
		ds.StoreLogger.Error("The store dump to file wasn't created", "NameOfFile", initialStore.Name, "Error", e)
	}
	if e = restoredStore.DumpToFile(restoredStore.Name + "Restored"); e == nil {
		ds.StoreLogger.Info("The dump of the restored store was created", "NameOfFile", restoredStore.Name)
	} else {
		ds.StoreLogger.Error("The dump of the restored store wasn't created", "NameOfFile", restoredStore.Name, "Error", e)
	}
	if fromFileStore, e := ds.NewStoreFromFile(initialStore.Name); e == nil {
		if fileStore, e := fromFileStore.Dump(); e == nil {
			ds.StoreLogger.Info("The dump of restored from the dump file store was created", "Name", initialStore.Name)
			ds.StoreLogger.Debug(string(fileStore))
		} else {
			ds.StoreLogger.Error("The store wasn't restored", "Name", initialStore.Name, "Error", e)
		}
	} else {
		ds.StoreLogger.Error("The dump of restored from the dump file store wasn't created ", "Name", initialStore.Name, "Error", e)
	}
}
