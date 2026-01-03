package main

import (
	"fmt"

	DocStore "github.com/pan68ruslan/GoLand/lesson_05/documentstore"
	Users "github.com/pan68ruslan/GoLand/lesson_05/users"
)

func main() {
	fmt.Println("\nCreate document store")
	storeName := "User_Store"
	serviceName := "User_Service"
	userId1 := "User1"
	store := DocStore.NewStore(storeName)
	if _, e := Users.NewService("", store); e != nil {
		fmt.Printf("\nCannot create the document store: %v", e)
	}
	userService, es := Users.NewService(serviceName, store)
	if es == nil {
		fmt.Printf("\nThe document store was created with name=%s", serviceName)
	}
	_, _ = userService.CreateUser(Users.User{ID: userId1, Name: "First_User"})
	_, _ = userService.CreateUser(Users.User{ID: "User2", Name: "Second_User"})
	_, _ = userService.CreateUser(Users.User{ID: "User3", Name: "Third_User"})

	fmt.Println("\nUser store contains:")
	if l, err := userService.ListUsers(); err == nil {
		for i, u := range l {
			fmt.Printf("%d.User %s with name %s\n", i+1, u.ID, u.Name)
		}
	}

	fmt.Println("\nFind the user with id=", userId1)
	u, ee := userService.GetUser(userId1)
	if ee == nil {
		fmt.Printf("User with Id=%s was found. It has name=%s\n", userId1, u.Name)
	} else {
		fmt.Printf("User with Id=%s was not found. Error: %v\n", userId1, ee)
	}

	fmt.Println("\nUnmarshal Document to User structure:")
	collection, e := store.GetCollection(serviceName)
	unmarshalledUser := Users.User{}
	if e == nil {
		if d, e := collection.Get(userId1); e == nil {
			if e := DocStore.UnmarshalDocumentToStructure(&d, &unmarshalledUser); e == nil {
				fmt.Printf("User structure was created with Id=%s\n", unmarshalledUser.ID)
			}
		}
	}
	fmt.Println("\nDelete the user with id=", userId1)
	if e := userService.DeleteUser(userId1); e == nil {
		fmt.Printf("User with Id=%s was deleted.\n", userId1)
	} else {
		fmt.Printf("User with Id=%s was not found\n", userId1)
	}

	fmt.Println("\nUser store contains:")
	if l, e := userService.ListUsers(); e == nil {
		for i, u := range l {
			fmt.Printf("%d.User %s with name %s\n", i+1, u.ID, u.Name)
		}
	}

	fmt.Println("\nFind the user with id=", userId1)
	if u, ok := userService.GetUser(userId1); ok == nil {
		fmt.Printf("User with Id=%s was found. It has name=%s:\n", userId1, u.Name)
	} else {
		fmt.Printf("User with Id=%s was not found\n", userId1)
	}

	fmt.Println("\nDelete the user with id=", userId1)
	if e := userService.DeleteUser(userId1); e == nil {
		fmt.Printf("User with Id=%s was deleted.\n", userId1)
	} else {
		fmt.Printf("User with Id=%s was not found\n", userId1)
	}

	fmt.Println("\nMarshal the User structure do Document:")
	unmarshalledUser.Name += "-Updated"
	marshalledDoc, er := DocStore.MarshalStructureToDocument(&unmarshalledUser)
	if er == nil && collection != nil {
		if e := collection.Put(marshalledDoc); e == nil {
			fmt.Printf("Marshalled Document was created with Id=%s\n", marshalledDoc.Fields["id"].Value)
		}
	}
	fmt.Println("\nUpdated User store contains:")
	if l, e := userService.ListUsers(); e == nil {
		for i, u := range l {
			fmt.Printf("%d.User %s with name %s\n", i+1, u.ID, u.Name)
		}
	}
}
