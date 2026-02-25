package main

func main() {
	CreateCollection("documents")
	CreateCollection("users")

	ListCollections()

	PutDocument("users", map[string]interface{}{"user_id": "123", "name": "Alex", "age": 43})
	PutDocument("users", map[string]interface{}{"user_id": "124", "name": "Bob", "age": 44})
	PutDocument("users", map[string]interface{}{"user_id": "125", "name": "John", "age": 45})

	GetDocument("users", "123")

	ListDocuments("users")

	DeleteDocument("users", "123")
	ListDocuments("users")

	CreateIndex("users", "user_id")
	CreateIndex("users", "name")

	DeleteIndex("users", "name_1")

	DeleteCollection("users")
	ListCollections()
}
