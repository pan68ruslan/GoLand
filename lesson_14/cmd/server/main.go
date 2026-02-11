package main

import (
	"context"
	"fmt"
	"lesson_14/internal/handlers"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//const stringURL = "mongodb://root:root@localhost:27017/?authSource=admin"

func main() {
	ctx := context.Background()
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://root:root@localhost:27017/?authSource=admin"
	}
	clientDB, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongodb connect failed: %v", err)
	}
	defer func() {
		if err := clientDB.Disconnect(ctx); err != nil {
			log.Fatalf("mongodb disconnect failed: %v", err)
		}
	}()

	h := &handlers.Handler{DB: clientDB}

	http.HandleFunc("/put_document", h.PutDocument)
	http.HandleFunc("/get_document", h.GetDocument)
	http.HandleFunc("/list_documents", h.ListDocuments)
	http.HandleFunc("/delete_document", h.DeleteDocument)
	http.HandleFunc("/create_collection", h.CreateCollection)
	http.HandleFunc("/list_collections", h.ListCollections)
	http.HandleFunc("/delete_collection", h.DeleteCollection)
	http.HandleFunc("/create_index", h.CreateIndex)
	http.HandleFunc("/delete_index", h.DeleteIndex)

	fmt.Println("HTTP server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
