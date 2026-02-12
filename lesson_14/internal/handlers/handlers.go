package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response struct {
	Ok bool `json:"ok"`
}

type Handler struct {
	DB *mongo.Client
}

func (h *Handler) PutDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string                 `json:"collection_name"`
		Document       map[string]interface{} `json:"document"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	collection := h.DB.Database("app").Collection(req.CollectionName)
	_, err := collection.InsertOne(context.Background(), req.Document)
	json.NewEncoder(w).Encode(Response{Ok: err == nil})
}

// /get_document
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string `json:"collection_name"`
		UserID         string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	collection := h.DB.Database("app").Collection(req.CollectionName)
	var result map[string]interface{}
	err := collection.FindOne(context.Background(), bson.M{"user_id": req.UserID}).Decode(&result)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Ok: false})
		return
	}
	json.NewEncoder(w).Encode(result)
}

// /list_documents
func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string `json:"collection_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	collection := h.DB.Database("app").Collection(req.CollectionName)
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(Response{Ok: false})
		return
	}
	defer cur.Close(context.Background())

	var results []map[string]interface{}
	for cur.Next(context.Background()) {
		var doc map[string]interface{}
		cur.Decode(&doc)
		results = append(results, doc)
	}
	json.NewEncoder(w).Encode(results)
}

// /delete_document
func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string `json:"collection_name"`
		UserID         string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	collection := h.DB.Database("app").Collection(req.CollectionName)
	_, err := collection.DeleteOne(context.Background(), bson.M{"user_id": req.UserID})
	json.NewEncoder(w).Encode(Response{Ok: err == nil})
}

// /create_collection
func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string `json:"collection_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.DB.Database("app").CreateCollection(context.Background(), req.CollectionName)
	json.NewEncoder(w).Encode(Response{Ok: err == nil})
}

// /list_collections
func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	names, err := h.DB.Database("app").ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(Response{Ok: false})
		return
	}
	json.NewEncoder(w).Encode(names)
}

// /delete_collection
func (h *Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string `json:"collection_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.DB.Database("app").Collection(req.CollectionName).Drop(context.Background())
	json.NewEncoder(w).Encode(Response{Ok: err == nil})
}

// /create_index
func (h *Handler) CreateIndex(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string `json:"collection_name"`
		Field          string `json:"field"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	collection := h.DB.Database("app").Collection(req.CollectionName)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: req.Field, Value: 1}},
	})
	json.NewEncoder(w).Encode(Response{Ok: err == nil})
}

// /delete_index
func (h *Handler) DeleteIndex(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CollectionName string `json:"collection_name"`
		IndexName      string `json:"index_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	collection := h.DB.Database("app").Collection(req.CollectionName)
	_, err := collection.Indexes().DropOne(context.Background(), req.IndexName)
	json.NewEncoder(w).Encode(Response{Ok: err == nil})
}
