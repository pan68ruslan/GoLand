package main

import (
	"fmt"

	"github.com/pan68ruslan/GoLand/DocumentStore"
)

func main() {
	store := DocumentStore.NewStore()
	cfg := &DocumentStore.CollectionConfig{
		PrimaryKey: "key",
	}
	fmt.Println("Store was created")
	collectionName := "docCollection"
	fmt.Printf("\nCreate the collection '%s'\n", collectionName)
	if _, ok := store.GetCollection(collectionName); ok {
		fmt.Printf("Found the collection '%s'\n", collectionName)
	}
	if ok, _ := store.CreateCollection(collectionName, cfg); ok {
		fmt.Printf("Collection '%s' created\n", collectionName)
	}
	if ok, _ := store.CreateCollection(collectionName, cfg); ok {
		fmt.Printf("Collection '%s' created\n", collectionName)
	}
	docCollection, ok := store.GetCollection(collectionName)
	if ok {
		fmt.Printf("Found the collection with '%s' name\n", collectionName)
	}

	fmt.Println("\nCreate documents")
	docCollection.List()
	doc1 := DocumentStore.Document{
		Fields: map[string]DocumentStore.DocumentField{
			"key":        {Type: DocumentStore.DocumentFieldTypeString, Value: "doc1"},
			"title":      {Type: DocumentStore.DocumentFieldTypeString, Value: "First_Document"},
			"isApproved": {Type: DocumentStore.DocumentFieldTypeBool, Value: false},
			"pages":      {Type: DocumentStore.DocumentFieldTypeNumber, Value: 42},
		},
	}
	docCollection.Put(doc1)
	docCollection.Put(DocumentStore.Document{
		Fields: map[string]DocumentStore.DocumentField{
			"key":        {Type: DocumentStore.DocumentFieldTypeString, Value: "doc2"},
			"title":      {Type: DocumentStore.DocumentFieldTypeString, Value: "Second_Document"},
			"isApproved": {Type: DocumentStore.DocumentFieldTypeBool, Value: true},
			"pages":      {Type: DocumentStore.DocumentFieldTypeNumber, Value: 100},
		},
	})
	docCollection.Put(DocumentStore.Document{
		Fields: map[string]DocumentStore.DocumentField{
			"key":        {Type: DocumentStore.DocumentFieldTypeString, Value: "doc3"},
			"title":      {Type: DocumentStore.DocumentFieldTypeString, Value: "Third_Document"},
			"isApproved": {Type: DocumentStore.DocumentFieldTypeBool, Value: true},
			"pages":      {Type: DocumentStore.DocumentFieldTypeNumber, Value: 10},
		},
	})
	docCollection.List()
	fmt.Println("\nDelete document")
	docCollection.Delete("doc4")
	docCollection.Delete("doc3")
	docCollection.List()

	fmt.Println("\nDelete collection")
	store.DeleteCollection(collectionName)
	store.GetCollection(collectionName)
}
