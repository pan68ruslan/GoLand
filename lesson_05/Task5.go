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
	svc.CreateUser("u1", "Anna")
	svc.CreateUser("u2", "Bob")
	svc.CreateUser("u3", "Charlie")

	users, _ := svc.ListUsers()
	fmt.Println(users)
	fmt.Println("\nGet user")
	user, _ := svc.GetUser("u1")
	fmt.Println(user)
	fmt.Println("\nDelete user")
	svc.DeleteUser("u2")
	usrs, _ := svc.ListUsers()
	fmt.Println(usrs)
}
