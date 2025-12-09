package main

import (
	"fmt"

	. "github.com/pan68ruslan/GoLand/DocStore"
)

func main() {
	doc1 := &Document{
		Fields: map[string]DocumentField{
			"key":        {Type: DocumentFieldTypeString, Value: "doc1"},
			"title":      {Type: DocumentFieldTypeString, Value: "First_Document"},
			"isApproved": {Type: DocumentFieldTypeBool, Value: false},
			"pages":      {Type: DocumentFieldTypeNumber, Value: 42},
		},
	}
	Put(doc1)

	doc2 := &Document{
		Fields: map[string]DocumentField{
			"key":        {Type: DocumentFieldTypeString, Value: "doc2"},
			"title":      {Type: DocumentFieldTypeString, Value: "Second_Document"},
			"isApproved": {Type: DocumentFieldTypeBool, Value: true},
			"pages":      {Type: DocumentFieldTypeNumber, Value: 100},
		},
	}
	Put(doc2)

	doc3 := &Document{
		Fields: map[string]DocumentField{
			"key":        {Type: DocumentFieldTypeString, Value: "doc3"},
			"title":      {Type: DocumentFieldTypeString, Value: "Third_Document"},
			"isApproved": {Type: DocumentFieldTypeBool, Value: true},
			"pages":      {Type: DocumentFieldTypeNumber, Value: 10},
		},
	}
	Put(doc3)

	fmt.Println("Content:")
	for _, d := range List() {
		fmt.Println("-", d.Fields["title"].Value)
	}

	docKey := "doc1"
	if d, ok := Get(docKey); ok {
		fmt.Printf("Document with key \"%s\" was found, its title is : %s\n", docKey, d.Fields["title"].Value)
	} else {
		fmt.Printf("Document with key \"%s\" wasn't found\n", docKey)
	}

	if Delete(docKey) {
		fmt.Printf("Document with key \"%s\" was removed\n", docKey)
	}

	fmt.Println("Content:")
	for _, d := range List() {
		fmt.Println("-", d.Fields["title"].Value)
	}

	if d, ok := Get(docKey); ok {
		fmt.Printf("Document with key \"%s\" was found, its title is : %s\n", docKey, d.Fields["title"].Value)
	} else {
		fmt.Printf("Document with key \"%s\" wasn't found\n", docKey)
	}
}
