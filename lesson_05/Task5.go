package main

import (
	"fmt"

	DocStore "github.com/pan68ruslan/GoLand/lesson_05/document_store"
	Users "github.com/pan68ruslan/GoLand/lesson_05/users"
)

func main() {
	fmt.Println("\nCreate document")
	initialDoc := DocStore.Document{
		Fields: map[string]DocStore.DocumentField{
			"key":        {Type: DocStore.DocumentFieldTypeString, Value: "doc1"},
			"title":      {Type: DocStore.DocumentFieldTypeString, Value: "Initial_Document"},
			"isApproved": {Type: DocStore.DocumentFieldTypeBool, Value: false},
			"pages":      {Type: DocStore.DocumentFieldTypeNumber, Value: 42},
		},
	}
	fmt.Println(initialDoc)
	docMap, _ := DocStore.MarshalDocument(initialDoc)
	docMap["title"] = "Updated_Document"
	docMap["isApproved"] = true
	fmt.Println("\nUpdate document")
	updatedDoc, err := DocStore.UnmarshalDocument(docMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(updatedDoc)

	fmt.Println("\nCreate user")
	coll := &Users.Collection{Title: "IT Store"}
	svc := Users.NewService(coll)
	if _, e := svc.CreateUser("u1", "Anna"); e != nil {
		fmt.Println("Error at user creation %w", e)
	}
	if _, e := svc.CreateUser("u2", "Bob"); e != nil {
		fmt.Println("Error at user creation %w", e)
	}
	if _, e := svc.CreateUser("u3", "Charlie"); e != nil {
		fmt.Println("Error at user creation %w", e)
	}
	users, _ := svc.ListUsers()
	fmt.Println(users)
	fmt.Println("\nGet user")
	user, _ := svc.GetUser("u1")
	fmt.Println(user)
	fmt.Println("\nDelete user")
	if e := svc.DeleteUser("u2"); e != nil {
		fmt.Println("Error at user deletion %w", e)
	}
	users, _ = svc.ListUsers()
	fmt.Println(users)
}
